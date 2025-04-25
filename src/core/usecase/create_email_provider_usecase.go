package usecase

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type ICreateEmailProviderUseCase interface {
	CreateEmailProvider(ctx context.Context, userId int64, workspaceCode, provider, code string) (*entity.EmailProviderEntity, error)
}
type CreateEmailProviderUseCase struct {
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
	getWorkspaceUseCase         IGetWorkspaceUseCase
	databaseTransactionUseCase  IDatabaseTransactionUseCase
	emailProviderPort           port.IEmailProviderPort
}

func (c CreateEmailProviderUseCase) CreateEmailProvider(ctx context.Context, userId int64, workspaceCode, provider, code string) (*entity.EmailProviderEntity, error) {
	workspace, err := c.getWorkspaceUseCase.GetWorkspaceByCode(ctx, userId, workspaceCode)
	if err != nil {
		log.Error(ctx, "GetWorkspaceByCode error: %v", err)
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, common.ErrForbidden
		}
		return nil, err
	}

	oauthResponse, err := c.emailProviderPort.GetOAuthInfo(ctx, provider, code)
	if err != nil {
		log.Error(ctx, "GetOAuthInfoByCode error: %v", err)
		return nil, err
	}
	providerEntity := entity.EmailProviderEntity{
		WorkspaceId:       workspace.ID,
		Provider:          provider,
		SmtpHost:          oauthResponse.SmtpHost,
		SmtpPort:          oauthResponse.SmtpPort,
		OAuthToken:        oauthResponse.AccessToken,
		OAuthRefreshToken: oauthResponse.RefreshToken,
		OAuthExpiredAt:    oauthResponse.ExpiredAt,
		UseTLS:            oauthResponse.UseTLS,
		Email:             oauthResponse.Email,
	}
	tx := c.databaseTransactionUseCase.StartTx()
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
	emailProvider, err := c.emailProviderRepositoryPort.SaveEmailProvider(ctx, tx, &providerEntity)
	if err != nil {
		log.Error(ctx, "SaveEmailProvider error: %v", err)
		return nil, err
	}
	if errCommit := tx.Commit().Error; errCommit != nil {
		log.Error(ctx, "Commit error: %v", errCommit)
		return nil, errCommit
	}
	return emailProvider, nil
}

func NewCreateEmailProviderUseCase(
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort,
	getWorkspaceUseCase IGetWorkspaceUseCase,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	emailProviderPort port.IEmailProviderPort,
) ICreateEmailProviderUseCase {
	return &CreateEmailProviderUseCase{
		emailProviderRepositoryPort: emailProviderRepositoryPort,
		getWorkspaceUseCase:         getWorkspaceUseCase,
		databaseTransactionUseCase:  databaseTransactionUseCase,
		emailProviderPort:           emailProviderPort,
	}
}
