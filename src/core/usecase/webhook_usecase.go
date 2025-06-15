package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IWebhookUsecase interface {
	SendNotifyMaxRetry(ctx context.Context, emailRequests []*entity.EmailRequestEntity)
}
type WebhookUsecase struct {
	webhookRepositoryPort port.IWebhookRepositoryPort
	webhookServicePort    port.IWebhookServicePort
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

	templateMessage := "Email delivery failed to %s\\nError: %s\\nRequest ID: %s"
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
) IWebhookUsecase {
	return &WebhookUsecase{
		webhookRepositoryPort: webhookRepositoryPort,
		webhookServicePort:    webhookServicePort,
	}
}
