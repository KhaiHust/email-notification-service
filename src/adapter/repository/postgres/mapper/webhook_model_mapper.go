package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
)

func ToWebhookEntity(webhookModel *model.WebhookModel) *entity.WebhookEntity {
	if webhookModel == nil {
		return nil
	}
	return &entity.WebhookEntity{
		BaseEntity:  ToBaseEntityMapper(&webhookModel.BaseModel),
		WorkspaceID: webhookModel.WorkspaceID,
		URL:         webhookModel.URL,
		Type:        webhookModel.Type,
		Enabled:     webhookModel.Enabled,
	}
}
func ToWebhookModel(webhookEntity *entity.WebhookEntity) *model.WebhookModel {
	if webhookEntity == nil {
		return nil
	}
	return &model.WebhookModel{
		BaseModel:   ToBaseModelMapper(&webhookEntity.BaseEntity),
		WorkspaceID: webhookEntity.WorkspaceID,
		URL:         webhookEntity.URL,
		Type:        webhookEntity.Type,
		Enabled:     webhookEntity.Enabled,
	}
}
