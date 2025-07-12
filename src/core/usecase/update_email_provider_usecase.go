package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"time"
)

const (
	PrefixEmailProviderKeyLease   = "email_provider_lock:id:%d"
	DefaultEmailProviderLeaseTime = 60 // seconds
)

type IUpdateEmailProviderUseCase interface {
	UpdateOAuthInfoByRefreshToken(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error)
	UpdateInfoProvider(ctx context.Context, workspaceID, providerID int64, req *request.UpdateEmailProviderDto) (*entity.EmailProviderEntity, error)
	DeactivateEmailProvider(ctx context.Context, workspaceID, providerID int64) (*entity.EmailProviderEntity, error)
	HandleTokenExpired(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error)
}
type UpdateEmailProviderUseCase struct {
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
	getWorkspaceUseCase         IGetWorkspaceUseCase
	databaseTransactionUseCase  IDatabaseTransactionUseCase
	emailProviderPort           port.IEmailProviderPort
	encryptUseCase              IEncryptUseCase
	redisPort                   port.IRedisPort
	getEmailProviderUseCase     IGetEmailProviderUseCase
}

func (u UpdateEmailProviderUseCase) HandleTokenExpired(
	ctx context.Context,
	emailProviderEntity *entity.EmailProviderEntity,
) (*entity.EmailProviderEntity, error) {
	acquireLockKey := fmt.Sprintf(PrefixEmailProviderKeyLease, emailProviderEntity.ID)

	// Try to acquire the distributed lock
	acquireLock, err := u.redisPort.SetLock(ctx, acquireLockKey, "1", DefaultEmailProviderLeaseTime)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("SetLock error: %v", err))
		return nil, err // Return early on Redis error
	}

	if acquireLock {
		// Only the lock owner should delete the lock!
		defer func() {
			if err := u.redisPort.DeleteKey(ctx, acquireLockKey); err != nil {
				log.Error(ctx, fmt.Sprintf("DeleteKey error: %v", err))
			}
		}()
		// Perform the refresh
		return u.UpdateOAuthInfoByRefreshToken(ctx, emailProviderEntity)
	}

	// Could not acquire the lock: someone else is refreshing
	const maxRetries = 3
	retryDelay := time.Second
	for i := 0; i < maxRetries; i++ {
		// Wait a bit for the lock to be released
		log.Info(ctx, fmt.Sprintf("Lock is already acquired, retrying in %v... (attempt %d)", retryDelay, i+1))
		time.Sleep(retryDelay)

		// Check if lock has been released (no need to check value, just check existence)
		lockExists, err := u.redisPort.GetFromRedis(ctx, acquireLockKey)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("Exists error: %v", err))
			return nil, err
		}
		if lockExists == nil {
			// Lock released, fetch the updated provider info
			emailProviderEntity, err = u.getEmailProviderUseCase.GetEmailProviderByID(ctx, emailProviderEntity.ID)
			if err != nil {
				log.Error(ctx, fmt.Sprintf("GetEmailProviderByID error: %v", err))
				return nil, err
			}
			return emailProviderEntity, nil
		}
	}

	emailProviderEntity, err = u.emailProviderRepositoryPort.GetEmailProviderByID(ctx, emailProviderEntity.ID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("GetEmailProviderByID error after retries: %v", err))
		return nil, err
	}
	return emailProviderEntity, nil
}

func (u UpdateEmailProviderUseCase) DeactivateEmailProvider(ctx context.Context, workspaceID, providerID int64) (*entity.EmailProviderEntity, error) {
	emailProviderEntity, err := u.emailProviderRepositoryPort.GetEmailProviderByWorkspaceIDAndID(ctx, workspaceID, providerID)
	if err != nil {
		log.Error(ctx, "GetEmailProviderByWorkspaceIDAndID error: %v", err)
		return nil, err
	}

	tx := u.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Error(ctx, "Rollback error: %v", errRollback)
			} else {
				log.Info(ctx, "Rollback successfully")
			}
		}
	}()
	deactivateEmail := fmt.Sprintf("%d_deactive_%s", time.Now().UnixMicro(), emailProviderEntity.Email)
	active := false
	emailProviderEntity.OAuthToken = ""
	emailProviderEntity.OAuthRefreshToken = ""
	emailProviderEntity.Active = &active
	emailProviderEntity.Email = deactivateEmail
	emailProviderEntity, err = u.emailProviderRepositoryPort.UpdateEmailProvider(ctx, tx, emailProviderEntity)
	if err != nil {
		log.Error(ctx, "UpdateEmailProvider error: %v", err)
		return nil, err
	}

	if err = u.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Commit error: %v", err)
		return nil, err
	}
	return emailProviderEntity, nil
}

func (u UpdateEmailProviderUseCase) UpdateInfoProvider(ctx context.Context, workspaceID, providerID int64, req *request.UpdateEmailProviderDto) (*entity.EmailProviderEntity, error) {
	emailProviderEntity, err := u.emailProviderRepositoryPort.GetEmailProviderByWorkspaceIDAndID(ctx, workspaceID, providerID)
	if err != nil {
		log.Error(ctx, "GetEmailProviderByWorkspaceIDAndID error: %v", err)
		return nil, err
	}
	var oauthResponse *response.OAuthInfoResponseDto
	if req.Code != nil && *req.Code != "" {
		oauthResponse, err = u.emailProviderPort.GetOAuthInfo(ctx, emailProviderEntity.Provider, *req.Code)
		if err != nil {
			log.Error(ctx, "GetOAuthInfoByCode error: %v", err)
			return nil, err
		}
		// Encrypt token before save
		accessToken, err := u.encryptUseCase.EncryptProviderToken(ctx, oauthResponse.AccessToken)
		if err != nil {
			log.Error(ctx, "EncryptProviderToken error: %v", err)
			return nil, err
		}
		refreshToken, err := u.encryptUseCase.EncryptProviderToken(ctx, oauthResponse.RefreshToken)
		if err != nil {
			log.Error(ctx, "EncryptProviderToken error: %v", err)
			return nil, err
		}
		emailProviderEntity.OAuthToken = accessToken
		emailProviderEntity.OAuthRefreshToken = refreshToken
		emailProviderEntity.OAuthExpiredAt = oauthResponse.ExpiredAt
		emailProviderEntity.SmtpHost = oauthResponse.SmtpHost
		emailProviderEntity.SmtpPort = oauthResponse.SmtpPort
		emailProviderEntity.UseTLS = oauthResponse.UseTLS
		emailProviderEntity.Email = oauthResponse.Email
	}
	if req.FromName != nil {
		emailProviderEntity.FromName = *req.FromName
	}
	tx := u.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Error(ctx, "Rollback error: %v", errRollback)
			} else {
				log.Info(ctx, "Rollback successfully")
			}
		}
	}()
	emailProviderEntity, err = u.emailProviderRepositoryPort.UpdateEmailProvider(ctx, tx, emailProviderEntity)
	if err != nil {
		log.Error(ctx, "UpdateEmailProvider error: %v", err)
		return nil, err
	}
	if err = u.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Commit error: %v", err)
		return nil, err
	}
	emailProviderEntity.OAuthToken = oauthResponse.AccessToken
	emailProviderEntity.OAuthRefreshToken = oauthResponse.RefreshToken
	return emailProviderEntity, nil
}

func (u UpdateEmailProviderUseCase) UpdateOAuthInfoByRefreshToken(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error) {
	cacheKey := fmt.Sprintf(PrefixEmailProviderKeyCache, emailProviderEntity.ID)
	if err := u.redisPort.DeleteKey(ctx, cacheKey); err != nil {
		log.Error(ctx, fmt.Sprintf("DeleteKey error: %v", err))
		return nil, err
	}
	oauthResponse, err := u.emailProviderPort.GetOAuthByRefreshToken(ctx, emailProviderEntity)
	if err != nil {
		log.Error(ctx, "GetOAuthByRefreshToken error: %v", err)
		return nil, err
	}
	// Encrypt token before save
	accessToken, err := u.encryptUseCase.EncryptProviderToken(ctx, oauthResponse.AccessToken)
	if err != nil {
		log.Error(ctx, "EncryptProviderToken error: %v", err)
		return nil, err
	}
	refreshToken, err := u.encryptUseCase.EncryptProviderToken(ctx, oauthResponse.RefreshToken)
	if err != nil {
		log.Error(ctx, "EncryptProviderToken error: %v", err)
		return nil, err
	}
	emailProviderEntity.OAuthToken = accessToken
	emailProviderEntity.OAuthRefreshToken = refreshToken
	emailProviderEntity.OAuthExpiredAt = oauthResponse.ExpiredAt

	//save to database
	tx := u.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Error(ctx, "Rollback error: %v", errRollback)
			} else {
				log.Info(ctx, "Rollback successfully")
			}
		}
	}()
	emailProviderEntity, err = u.emailProviderRepositoryPort.UpdateEmailProvider(ctx, tx, emailProviderEntity)
	if err != nil {
		log.Error(ctx, "UpdateEmailProvider error: %v", err)
		return nil, err
	}
	if err = u.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Commit error: %v", err)
		return nil, err
	}

	emailProviderEntity.OAuthToken = oauthResponse.AccessToken
	emailProviderEntity.OAuthRefreshToken = oauthResponse.RefreshToken
	return emailProviderEntity, nil

}

func NewUpdateEmailProviderUseCase(
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort,
	getWorkspaceUseCase IGetWorkspaceUseCase,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	emailProviderPort port.IEmailProviderPort,
	encryptUseCase IEncryptUseCase,
	redisPort port.IRedisPort,
	getEmailProviderUseCase IGetEmailProviderUseCase,
) IUpdateEmailProviderUseCase {
	return &UpdateEmailProviderUseCase{
		emailProviderRepositoryPort: emailProviderRepositoryPort,
		getWorkspaceUseCase:         getWorkspaceUseCase,
		databaseTransactionUseCase:  databaseTransactionUseCase,
		emailProviderPort:           emailProviderPort,
		encryptUseCase:              encryptUseCase,
		redisPort:                   redisPort,
		getEmailProviderUseCase:     getEmailProviderUseCase,
	}
}
