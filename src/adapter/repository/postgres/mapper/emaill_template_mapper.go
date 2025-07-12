package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
)

func ToEmailTemplateModel(emailTemplateEntity *entity.EmailTemplateEntity) *model.EmailTemplateModel {
	return &model.EmailTemplateModel{
		BaseModel:     ToBaseModelMapper(&emailTemplateEntity.BaseEntity),
		Name:          emailTemplateEntity.Name,
		Subject:       emailTemplateEntity.Subject,
		Body:          emailTemplateEntity.Body,
		Variables:     emailTemplateEntity.Variables,
		WorkspaceId:   emailTemplateEntity.WorkspaceId,
		CreatedBy:     emailTemplateEntity.CreatedBy,
		LastUpdatedBy: emailTemplateEntity.LastUpdatedBy,
		Version:       emailTemplateEntity.Version,
	}
}
func ToEmailTemplateEntity(emailTemplateModel *model.EmailTemplateModel) *entity.EmailTemplateEntity {
	if emailTemplateModel == nil {
		return nil
	}
	return &entity.EmailTemplateEntity{
		BaseEntity:    ToBaseEntityMapper(&emailTemplateModel.BaseModel),
		Name:          emailTemplateModel.Name,
		Subject:       emailTemplateModel.Subject,
		Body:          emailTemplateModel.Body,
		Variables:     emailTemplateModel.Variables,
		WorkspaceId:   emailTemplateModel.WorkspaceId,
		CreatedBy:     emailTemplateModel.CreatedBy,
		LastUpdatedBy: emailTemplateModel.LastUpdatedBy,
		Version:       emailTemplateModel.Version,
	}
}
func ToEmailTemplateEntities(emailTemplateModels []*model.EmailTemplateModel) []*entity.EmailTemplateEntity {
	var emailTemplateEntities []*entity.EmailTemplateEntity
	for _, emailTemplateModel := range emailTemplateModels {
		emailTemplateEntities = append(emailTemplateEntities, ToEmailTemplateEntity(emailTemplateModel))
	}
	return emailTemplateEntities
}
