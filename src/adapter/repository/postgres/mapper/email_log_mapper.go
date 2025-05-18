package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"time"
)

func ToEmailLogModel(emailLogEntity *entity.EmailLogsEntity) *model.EmailLogsModel {
	if emailLogEntity == nil {
		return nil
	}
	return &model.EmailLogsModel{
		BaseModel:       ToBaseModelMapper(&emailLogEntity.BaseEntity),
		EmailRequestID:  emailLogEntity.EmailRequestID,
		TemplateId:      emailLogEntity.TemplateId,
		Recipient:       emailLogEntity.Recipient,
		Status:          emailLogEntity.Status,
		ErrorMessage:    emailLogEntity.ErrorMessage,
		RetryCount:      emailLogEntity.RetryCount,
		RequestID:       emailLogEntity.RequestID,
		WorkspaceID:     emailLogEntity.WorkspaceID,
		EmailProviderID: emailLogEntity.EmailProviderID,
		LoggedAt:        time.Unix(emailLogEntity.LoggedAt, 0),
	}
}
func ToEmailLogEntity(emailLogModel *model.EmailLogsModel) *entity.EmailLogsEntity {
	if emailLogModel == nil {
		return nil
	}
	return &entity.EmailLogsEntity{
		BaseEntity:      ToBaseEntityMapper(&emailLogModel.BaseModel),
		EmailRequestID:  emailLogModel.EmailRequestID,
		TemplateId:      emailLogModel.TemplateId,
		Recipient:       emailLogModel.Recipient,
		Status:          emailLogModel.Status,
		ErrorMessage:    emailLogModel.ErrorMessage,
		RetryCount:      emailLogModel.RetryCount,
		RequestID:       emailLogModel.RequestID,
		WorkspaceID:     emailLogModel.WorkspaceID,
		EmailProviderID: emailLogModel.EmailProviderID,
		LoggedAt:        emailLogModel.LoggedAt.Unix(),
	}
}
func ToListEmailLogEntity(emailLogModels []*model.EmailLogsModel) []*entity.EmailLogsEntity {
	if emailLogModels == nil {
		return nil
	}
	emailLogEntities := make([]*entity.EmailLogsEntity, len(emailLogModels))
	for i, emailLogModel := range emailLogModels {
		emailLogEntities[i] = ToEmailLogEntity(emailLogModel)
	}
	return emailLogEntities
}
