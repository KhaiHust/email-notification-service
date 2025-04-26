package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IUpdateEmailProviderUseCase interface {
	UpdateOAuthInfoByRefreshToken(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error)
}
type UpdateEmailProviderUseCase struct {
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
	getWorkspaceUseCase         IGetWorkspaceUseCase
	databaseTransactionUseCase  IDatabaseTransactionUseCase
	emailProviderPort           port.IEmailProviderPort
}

func (u UpdateEmailProviderUseCase) UpdateOAuthInfoByRefreshToken(ctx context.Context, emailProviderEntity *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error) {
	oauthResponse, err := u.emailProviderPort.GetOAuthByRefreshToken(ctx, emailProviderEntity)
	if err != nil {
		log.Error(ctx, "GetOAuthByRefreshToken error: %v", err)
		return nil, err
	}
	emailProviderEntity.OAuthToken = oauthResponse.AccessToken
	emailProviderEntity.OAuthRefreshToken = oauthResponse.RefreshToken
	emailProviderEntity.OAuthExpiredAt = oauthResponse.ExpiredAt

	//save to database
	tx := u.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Error(ctx, "Rollback error: %v", errRollback)
		} else {
			log.Info(ctx, "Rollback successfully")
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
	//todo: update cache
	return emailProviderEntity, nil

}

func NewUpdateEmailProviderUseCase(
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort,
	getWorkspaceUseCase IGetWorkspaceUseCase,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	emailProviderPort port.IEmailProviderPort,
) IUpdateEmailProviderUseCase {
	return &UpdateEmailProviderUseCase{
		emailProviderRepositoryPort: emailProviderRepositoryPort,
		getWorkspaceUseCase:         getWorkspaceUseCase,
		databaseTransactionUseCase:  databaseTransactionUseCase,
		emailProviderPort:           emailProviderPort,
	}
}
