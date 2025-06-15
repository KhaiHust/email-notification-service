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
