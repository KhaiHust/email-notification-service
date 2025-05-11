package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/utils"
)

func ToEmailRequestModel(emailRequestEntity *entity.EmailRequestEntity) *model.EmailRequestModel {
	emailRequestModel := &model.EmailRequestModel{
		BaseModel:     ToBaseModelMapper(&emailRequestEntity.BaseEntity),
		TemplateId:    emailRequestEntity.TemplateId,
		Recipient:     emailRequestEntity.Recipient,
		Data:          emailRequestEntity.Data,
		Status:        emailRequestEntity.Status,
		ErrorMessage:  emailRequestEntity.ErrorMessage,
		RetryCount:    emailRequestEntity.RetryCount,
		RequestID:     emailRequestEntity.RequestID,
		CorrelationID: emailRequestEntity.CorrelationID,
		WorkspaceID:   emailRequestEntity.WorkspaceID,
	}
	if emailRequestEntity.SentAt != nil {
		emailRequestModel.SentAt = utils.ToTimePointer(*emailRequestEntity.SentAt)
	}
	return emailRequestModel
}
func ToEmailRequestEntity(emailRequestModel *model.EmailRequestModel) *entity.EmailRequestEntity {
	emailRequestEntity := &entity.EmailRequestEntity{
		BaseEntity:    ToBaseEntityMapper(&emailRequestModel.BaseModel),
		TemplateId:    emailRequestModel.TemplateId,
		Recipient:     emailRequestModel.Recipient,
		Data:          emailRequestModel.Data,
		Status:        emailRequestModel.Status,
		ErrorMessage:  emailRequestModel.ErrorMessage,
		RetryCount:    emailRequestModel.RetryCount,
		RequestID:     emailRequestModel.RequestID,
		CorrelationID: emailRequestModel.CorrelationID,
		WorkspaceID:   emailRequestModel.WorkspaceID,
	}
	if emailRequestModel.SentAt != nil {
		emailRequestEntity.SentAt = utils.ToUnixTimeToPointer(emailRequestModel.SentAt)
	}
	return emailRequestEntity
}
func ToListEmailRequestModel(emailRequestEntities []*entity.EmailRequestEntity) []*model.EmailRequestModel {
	emailRequestModels := make([]*model.EmailRequestModel, len(emailRequestEntities))
	for i, emailRequestEntity := range emailRequestEntities {
		emailRequestModels[i] = ToEmailRequestModel(emailRequestEntity)
	}
	return emailRequestModels
}
func ToListEmailRequestEntity(emailRequestModels []*model.EmailRequestModel) []*entity.EmailRequestEntity {
	emailRequestEntities := make([]*entity.EmailRequestEntity, len(emailRequestModels))
	for i, emailRequestModel := range emailRequestModels {
		emailRequestEntities[i] = ToEmailRequestEntity(emailRequestModel)
	}
	return emailRequestEntities
}
func ToEmailStatusCountEntity(emailStatusCountModel *model.EmailRequestStatusCountModel) *entity.EmailRequestStatusCountEntity {
	return &entity.EmailRequestStatusCountEntity{
		EmailTemplateId: emailStatusCountModel.EmailTemplateID,
		Status:          emailStatusCountModel.Status,
		Total:           emailStatusCountModel.Total,
	}
}
func ToListEmailStatusCountEntity(emailStatusCountModels []*model.EmailRequestStatusCountModel) []*entity.EmailRequestStatusCountEntity {
	emailStatusCountEntities := make([]*entity.EmailRequestStatusCountEntity, len(emailStatusCountModels))
	for i, emailStatusCountModel := range emailStatusCountModels {
		emailStatusCountEntities[i] = ToEmailStatusCountEntity(emailStatusCountModel)
	}
	return emailStatusCountEntities
}
