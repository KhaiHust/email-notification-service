package postgres

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type WebhookRepositoryAdapter struct {
	*base
}

func (w WebhookRepositoryAdapter) DeleteWebhook(ctx context.Context, tx *gorm.DB, workspaceID, webhookID int64) error {
	webhookModel := &model.WebhookModel{}
	if err := tx.WithContext(ctx).Where("id = ? AND workspace_id = ?", webhookID, workspaceID).
		Delete(webhookModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrRecordNotFound // No record found
		}
		return err
	}
	return nil
}

func (w WebhookRepositoryAdapter) UpdateWebhook(ctx context.Context, tx *gorm.DB, webhookEntity *entity.WebhookEntity) (*entity.WebhookEntity, error) {
	webhookModel := mapper.ToWebhookModel(webhookEntity)
	// build map of fields to update
	mapForUpdate := map[string]interface{}{
		"enabled": webhookModel.Enabled,
		"name":    webhookModel.Name,
		"url":     webhookModel.URL,
	}
	if err := tx.WithContext(ctx).Model(&model.WebhookModel{}).
		Where("id = ? AND workspace_id = ?", webhookModel.ID, webhookModel.WorkspaceID).
		Updates(mapForUpdate).Error; err != nil {
		return nil, err
	}
	return mapper.ToWebhookEntity(webhookModel), nil
}

func (w WebhookRepositoryAdapter) GetWebhookByWorkspaceIDAndWebhookID(ctx context.Context, workspaceID, webhookID int64) (*entity.WebhookEntity, error) {
	var webhookModel model.WebhookModel
	if err := w.db.WithContext(ctx).
		Where("workspace_id = ? AND id = ?", workspaceID, webhookID).
		First(&webhookModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound // No record found
		}
		return nil, err // Other errors
	}
	return mapper.ToWebhookEntity(&webhookModel), nil
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

func (w WebhookRepositoryAdapter) GetWebhooksByWorkspaceID(ctx context.Context, workspaceID int64) ([]*entity.WebhookEntity, error) {
	var webhookModels []*model.WebhookModel
	if err := w.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
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
