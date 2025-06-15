package postgres

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type WebhookRepositoryAdapter struct {
	*base
}

func (w WebhookRepositoryAdapter) GetActiveWebhooksByWorkspaceIDs(ctx context.Context, workspaceIDs []int64) ([]*entity.WebhookEntity, error) {
	var webhookModels []*model.WebhookModel
	if err := w.db.WithContext(ctx).
		Where("workspace_id IN ? AND enabled = ?", workspaceIDs, true).
		Find(&webhookModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListWebhookEntity(webhookModels), nil
}

func (w WebhookRepositoryAdapter) GetActiveWebhooksByWorkspaceID(ctx context.Context, workspaceID int64) ([]*entity.WebhookEntity, error) {
	var webhookModels []*model.WebhookModel
	if err := w.db.WithContext(ctx).
		Where("workspace_id = ? AND enabled = ?", workspaceID, true).
		Find(&webhookModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListWebhookEntity(webhookModels), nil
}

func (w WebhookRepositoryAdapter) CreateNewWebhook(ctx context.Context, tx *gorm.DB, webhookEntity *entity.WebhookEntity) (*entity.WebhookEntity, error) {
	webhookModel := mapper.ToWebhookModel(webhookEntity)
	if err := w.db.WithContext(ctx).Model(&model.WebhookModel{}).Create(webhookModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToWebhookEntity(webhookModel), nil
}

func NewWebhookRepositoryAdapter(db *gorm.DB) port.IWebhookRepositoryPort {
	return &WebhookRepositoryAdapter{base: &base{db: db}}
}
