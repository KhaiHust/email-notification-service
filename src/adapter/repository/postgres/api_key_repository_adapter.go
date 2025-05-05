package postgres

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type ApiKeyRepositoryAdapter struct {
	base
}

func (a ApiKeyRepositoryAdapter) SaveNewApiKey(ctx context.Context, tx *gorm.DB, apiKeyEntity *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error) {
	apiKeyModel := mapper.ToApiKeyModel(apiKeyEntity)
	if err := tx.WithContext(ctx).Model(&model.ApiKeyModel{}).Create(apiKeyModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToApiKeyEntity(apiKeyModel), nil
}

func NewApiKeyRepositoryAdapter(db *gorm.DB) port.IApiKeyRepositoryPort {
	return &ApiKeyRepositoryAdapter{
		base: base{db: db},
	}
}
