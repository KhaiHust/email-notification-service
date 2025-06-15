package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type ICreateWebhookUseCase interface {
	CreateNewWebhook(ctx context.Context, webhookReq *request.CreateWebhookRequest) (*entity.WebhookEntity, error)
}
type CreateWebhookUseCase struct {
	webhookRepositoryPort      port.IWebhookRepositoryPort
	databaseTransactionUseCase IDatabaseTransactionUseCase
}

func (c CreateWebhookUseCase) CreateNewWebhook(ctx context.Context, webhookReq *request.CreateWebhookRequest) (*entity.WebhookEntity, error) {
	var err error
	tx := c.databaseTransactionUseCase.StartTx()
	commit := false
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil || !commit {
			if errRollback := c.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "Failed to rollback transaction", "error", errRollback)
			} else {
				log.Info(ctx, "Transaction rolled back successfully")
			}

		}
	}()
	webhookEntity := &entity.WebhookEntity{
		WorkspaceID: webhookReq.WorkspaceID,
		URL:         webhookReq.URL,
		Type:        webhookReq.Type,
		Name:        webhookReq.Name,
		Enabled:     webhookReq.Enabled,
	}
	webhookEntity, err = c.webhookRepositoryPort.CreateNewWebhook(ctx, tx, webhookEntity)
	if err != nil {
		log.Error(ctx, "Failed to create new webhook", "error", err)
		return nil, err
	}
	if err = c.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Failed to commit transaction", "error", err)
		return nil, err
	}
	commit = true
	return webhookEntity, nil
}

func NewCreateWebhookUseCase(webhookRepositoryPort port.IWebhookRepositoryPort, databaseTransactionUseCase IDatabaseTransactionUseCase) ICreateWebhookUseCase {
	return &CreateWebhookUseCase{
		webhookRepositoryPort:      webhookRepositoryPort,
		databaseTransactionUseCase: databaseTransactionUseCase,
	}
}
