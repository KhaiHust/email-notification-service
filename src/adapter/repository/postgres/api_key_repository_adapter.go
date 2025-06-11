package postgres

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/specification"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type ApiKeyRepositoryAdapter struct {
	base
}

func (a ApiKeyRepositoryAdapter) UpdateAPIKey(ctx context.Context, tx *gorm.DB, apiKeyEntity *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error) {
	apiKeyModel := mapper.ToApiKeyModel(apiKeyEntity)
	if err := tx.WithContext(ctx).Model(&model.ApiKeyModel{}).
		Where("id = ?", apiKeyModel.ID).Updates(apiKeyModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToApiKeyEntity(apiKeyModel), nil
}

func (a ApiKeyRepositoryAdapter) GetAPIKeyByIDAndWorkspaceID(ctx context.Context, apiKeyID, workspaceID int64) (*entity.ApiKeyEntity, error) {
	var apiKeyModel model.ApiKeyModel
	if err := a.db.WithContext(ctx).Model(&model.ApiKeyModel{}).
		Where("id = ? AND workspace_id = ?", apiKeyID, workspaceID).First(&apiKeyModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToApiKeyEntity(&apiKeyModel), nil
}

func (a ApiKeyRepositoryAdapter) GetAPIKeyByHashKey(ctx context.Context, hashKey string) (*entity.ApiKeyEntity, error) {
	var apiKeyModel model.ApiKeyModel
	if err := a.db.WithContext(ctx).Model(&model.ApiKeyModel{}).Where("key_hash = ?", hashKey).First(&apiKeyModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToApiKeyEntity(&apiKeyModel), nil
}

func (a ApiKeyRepositoryAdapter) GetAllApiKeys(ctx context.Context, filter *request.GetApiKeyRequestFilter) ([]*entity.ApiKeyEntity, error) {
	spec := specification.ToApiKeySpecification(filter)
	query, args, err := specification.NewApiKeySpecification(spec)
	if err != nil {
		return nil, err
	}
	var apiKeyModels []*model.ApiKeyModel
	if err := a.db.WithContext(ctx).Raw(query, args...).Scan(&apiKeyModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListApiKeyEntity(apiKeyModels), nil
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
