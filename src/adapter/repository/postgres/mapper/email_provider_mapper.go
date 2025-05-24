package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"time"
)

func ToEmailProviderModel(emailProviderEntity *entity.EmailProviderEntity) *model.EmailProviderModel {
	if emailProviderEntity == nil {
		return nil
	}
	return &model.EmailProviderModel{
		BaseModel:         ToBaseModelMapper(&emailProviderEntity.BaseEntity),
		WorkspaceId:       emailProviderEntity.WorkspaceId,
		Provider:          emailProviderEntity.Provider,
		SmtpHost:          emailProviderEntity.SmtpHost,
		SmtpPort:          emailProviderEntity.SmtpPort,
		OAuthToken:        emailProviderEntity.OAuthToken,
		OAuthRefreshToken: emailProviderEntity.OAuthRefreshToken,
		OAuthExpiredAt:    time.Unix(emailProviderEntity.OAuthExpiredAt, 0),
		UseTLS:            emailProviderEntity.UseTLS,
		Email:             emailProviderEntity.Email,
		FromName:          emailProviderEntity.FromName,
		Environment:       emailProviderEntity.Environment,
	}
}
func ToEmailProviderEntity(emailProviderModel *model.EmailProviderModel) *entity.EmailProviderEntity {
	if emailProviderModel == nil {
		return nil
	}
	return &entity.EmailProviderEntity{
		BaseEntity: entity.BaseEntity{
			ID:        emailProviderModel.ID,
			CreatedAt: emailProviderModel.CreatedAt.Unix(),
			UpdatedAt: emailProviderModel.UpdatedAt.Unix(),
		},
		WorkspaceId:       emailProviderModel.WorkspaceId,
		Provider:          emailProviderModel.Provider,
		SmtpHost:          emailProviderModel.SmtpHost,
		SmtpPort:          emailProviderModel.SmtpPort,
		OAuthToken:        emailProviderModel.OAuthToken,
		OAuthRefreshToken: emailProviderModel.OAuthRefreshToken,
		OAuthExpiredAt:    emailProviderModel.OAuthExpiredAt.Unix(),
		UseTLS:            emailProviderModel.UseTLS,
		Email:             emailProviderModel.Email,
		FromName:          emailProviderModel.FromName,
		Environment:       emailProviderModel.Environment,
	}
}
func ToListEmailProviderEntity(emailProviderModels []*model.EmailProviderModel) []*entity.EmailProviderEntity {
	emailProviderEntities := make([]*entity.EmailProviderEntity, len(emailProviderModels))
	for i, emailProviderModel := range emailProviderModels {
		emailProviderEntities[i] = ToEmailProviderEntity(emailProviderModel)
	}
	return emailProviderEntities
}
