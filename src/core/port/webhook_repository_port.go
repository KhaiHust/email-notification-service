package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IWebhookRepositoryPort interface {
	CreateNewWebhook(ctx context.Context, tx *gorm.DB, webhookEntity *entity.WebhookEntity) (*entity.WebhookEntity, error)
	GetActiveWebhooksByWorkspaceID(ctx context.Context, workspaceID int64) ([]*entity.WebhookEntity, error)
	GetActiveWebhooksByWorkspaceIDs(ctx context.Context, workspaceIDs []int64) ([]*entity.WebhookEntity, error)
}
