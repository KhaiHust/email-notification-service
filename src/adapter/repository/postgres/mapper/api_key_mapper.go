package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/utils"
)

func ToApiKeyEntity(apiKeyModel *model.ApiKeyModel) *entity.ApiKeyEntity {
	if apiKeyModel == nil {
		return nil
	}
	return &entity.ApiKeyEntity{
		BaseEntity:  ToBaseEntityMapper(&apiKeyModel.BaseModel),
		WorkspaceID: apiKeyModel.WorkspaceID,
		Name:        apiKeyModel.Name,
		KeyHash:     apiKeyModel.KeyHash,
		RawPrefix:   apiKeyModel.RawPrefix,
		Environment: apiKeyModel.Environment,
		ExpiresAt:   utils.ToUnixTimeToPointer(apiKeyModel.ExpiresAt),
		Revoked:     apiKeyModel.Revoked,
	}
}
func ToApiKeyModel(apiKeyEntity *entity.ApiKeyEntity) *model.ApiKeyModel {
	if apiKeyEntity == nil {
		return nil
	}
	return &model.ApiKeyModel{
		BaseModel:   ToBaseModelMapper(&apiKeyEntity.BaseEntity),
		WorkspaceID: apiKeyEntity.WorkspaceID,
		Name:        apiKeyEntity.Name,
		KeyHash:     apiKeyEntity.KeyHash,
		RawPrefix:   apiKeyEntity.RawPrefix,
		Environment: apiKeyEntity.Environment,
		ExpiresAt:   utils.FromUnixPointerToTime(apiKeyEntity.ExpiresAt),
		Revoked:     apiKeyEntity.Revoked,
	}
}
