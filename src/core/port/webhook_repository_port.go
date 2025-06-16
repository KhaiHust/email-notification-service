package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IWebhookRepositoryPort interface {
	CreateNewWebhook(ctx context.Context, tx *gorm.DB, webhookEntity *entity.WebhookEntity) (*entity.WebhookEntity, error)
	GetWebhooksByWorkspaceID(ctx context.Context, workspaceID int64) ([]*entity.WebhookEntity, error)
	GetActiveWebhooksByWorkspaceIDs(ctx context.Context, workspaceIDs []int64) ([]*entity.WebhookEntity, error)
	GetWebhookByWorkspaceIDAndWebhookID(ctx context.Context, workspaceID, webhookID int64) (*entity.WebhookEntity, error)
	UpdateWebhook(ctx context.Context, tx *gorm.DB, webhookEntity *entity.WebhookEntity) (*entity.WebhookEntity, error)
	DeleteWebhook(ctx context.Context, tx *gorm.DB, workspaceID, webhookID int64) error
}
