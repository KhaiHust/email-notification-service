package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/exception"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IWebhookUsecase interface {
	SendNotifyMaxRetry(ctx context.Context, emailRequests []*entity.EmailRequestEntity)
	GetAllWebhookByWorkspaceID(ctx context.Context, workspaceID int64) ([]*entity.WebhookEntity, error)
	UpdateWebhook(ctx context.Context, workspaceID, webID int64, req *request.UpdateWebhookRequest) (*entity.WebhookEntity, error)
	DeleteWebhook(ctx context.Context, workspaceID, webID int64) error
	GetWebhookDetail(ctx context.Context, workspaceID, webID int64) (*entity.WebhookEntity, error)
	TestWebhook(ctx context.Context, workspaceID, webID int64) error
}
type WebhookUsecase struct {
	webhookRepositoryPort      port.IWebhookRepositoryPort
	webhookServicePort         port.IWebhookServicePort
	databaseTransactionUseCase IDatabaseTransactionUseCase
}

func (w WebhookUsecase) TestWebhook(ctx context.Context, workspaceID, webID int64) error {
	webhook, err := w.webhookRepositoryPort.GetWebhookByWorkspaceIDAndWebhookID(ctx, workspaceID, webID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error when get webhook by workspace id %d and webhook id %d", workspaceID, webID), err)
		return err
	}
	message := fmt.Sprintf("Test webhook from workspace %d with webhook id %d", workspaceID, webID)
	if err := w.webhookServicePort.Send(ctx, webhook, message); err != nil {
		log.Error(ctx, fmt.Sprintf("Error when send test webhook to %s", webhook.URL), err)
		return err
	}
	return nil
}

func (w WebhookUsecase) GetWebhookDetail(ctx context.Context, workspaceID, webID int64) (*entity.WebhookEntity, error) {
	webhook, err := w.webhookRepositoryPort.GetWebhookByWorkspaceIDAndWebhookID(ctx, workspaceID, webID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error when get webhook by workspace id %d and webhook id %d", workspaceID, webID), err)
		return nil, err
	}
	return webhook, nil
}

func (w WebhookUsecase) DeleteWebhook(ctx context.Context, workspaceID, webID int64) error {
	webhook, err := w.webhookRepositoryPort.GetWebhookByWorkspaceIDAndWebhookID(ctx, workspaceID, webID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error when get webhook by workspace id %d and webhook id %d", workspaceID, webID), err)
		return err
	}
	tx := w.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := w.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "Failed to rollback transaction", "error", errRollback)
			} else {
				log.Info(ctx, "Transaction rolled back successfully")
			}
		}
	}()
	err = w.webhookRepositoryPort.DeleteWebhook(ctx, tx, webhook.WorkspaceID, webhook.ID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error when delete webhook with id %d", webID), err)
		return err
	}
	if err = w.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Failed to commit transaction", "error", err)
		return err
	}
	return nil
}

func (w WebhookUsecase) UpdateWebhook(ctx context.Context, workspaceID, webID int64, req *request.UpdateWebhookRequest) (*entity.WebhookEntity, error) {
	webhook, err := w.webhookRepositoryPort.GetWebhookByWorkspaceIDAndWebhookID(ctx, workspaceID, webID)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error when get webhook by workspace id %d and webhook id %d", workspaceID, webID), err)
		return nil, err
	}
	if req.Name != nil {
		webhook.Name = *req.Name
	}
	if req.URL != nil {
		webhook.URL = *req.URL
	}
	if req.Enabled != nil {
		webhook.Enabled = *req.Enabled
	}
	tx := w.databaseTransactionUseCase.StartTx()
	defer func() {
		if r := recover(); r != nil {
			err = exception.InternalServerException
		}
		if err != nil {
			if errRollback := w.databaseTransactionUseCase.RollbackTx(tx); errRollback != nil {
				log.Error(ctx, "Failed to rollback transaction", "error", errRollback)
			} else {
				log.Info(ctx, "Transaction rolled back successfully")
			}
		}
	}()
	webhook, err = w.webhookRepositoryPort.UpdateWebhook(ctx, tx, webhook)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("Error when update webhook with id %d", webID), err)
		return nil, err
	}
	if err = w.databaseTransactionUseCase.CommitTx(tx); err != nil {
		log.Error(ctx, "Failed to commit transaction", "error", err)
		return nil, err
	}
	return webhook, nil
}

func (w WebhookUsecase) GetAllWebhookByWorkspaceID(ctx context.Context, workspaceID int64) ([]*entity.WebhookEntity, error) {
	return w.webhookRepositoryPort.GetWebhooksByWorkspaceID(ctx, workspaceID)
}

func (w WebhookUsecase) SendNotifyMaxRetry(ctx context.Context, emailRequests []*entity.EmailRequestEntity) {
	workspaceIDs := make([]int64, 0, len(emailRequests))
	for _, emailRequest := range emailRequests {
		if emailRequest.WorkspaceID > 0 {
			workspaceIDs = append(workspaceIDs, emailRequest.WorkspaceID)
		}
	}
	//build message
	webhookEntities, err := w.webhookRepositoryPort.GetActiveWebhooksByWorkspaceIDs(ctx, workspaceIDs)
	if err != nil || len(webhookEntities) == 0 {
		log.Error(ctx, fmt.Sprintf("Error when get active webhooks by workspace ids : %v", workspaceIDs), err)
		return
	}
	mapWebhooksByWorkspaceID := make(map[int64][]*entity.WebhookEntity)
	for _, webhookEntity := range webhookEntities {
		if exist, ok := mapWebhooksByWorkspaceID[webhookEntity.WorkspaceID]; ok {
			mapWebhooksByWorkspaceID[webhookEntity.WorkspaceID] = append(exist, webhookEntity)
		} else {
			mapWebhooksByWorkspaceID[webhookEntity.WorkspaceID] = []*entity.WebhookEntity{webhookEntity}
		}
	}

	templateMessage := "Email delivery failed to %s\nError: %s\nRequest ID: %s"
	for _, emailRequest := range emailRequests {
		webhooks, ok := mapWebhooksByWorkspaceID[emailRequest.WorkspaceID]
		if !ok || len(webhooks) == 0 {
			log.Warn(ctx, fmt.Sprintf("No active webhooks found for workspace ID %d", emailRequest.WorkspaceID))
			continue
		}
		message := fmt.Sprintf(templateMessage, emailRequest.Recipient, emailRequest.ErrorMessage, emailRequest.RequestID)
		for _, webhook := range webhooks {
			if err := w.webhookServicePort.Send(ctx, webhook, message); err != nil {
				log.Error(ctx, fmt.Sprintf("Error when send webhook to %s", webhook.URL), err)
			}
		}
	}
}

func NewWebhookUsecase(
	webhookRepositoryPort port.IWebhookRepositoryPort,
	webhookServicePort port.IWebhookServicePort,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
) IWebhookUsecase {
	return &WebhookUsecase{
		webhookRepositoryPort:      webhookRepositoryPort,
		webhookServicePort:         webhookServicePort,
		databaseTransactionUseCase: databaseTransactionUseCase,
	}
}
