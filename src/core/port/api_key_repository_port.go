package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"gorm.io/gorm"
)

type IApiKeyRepositoryPort interface {
	SaveNewApiKey(ctx context.Context, tx *gorm.DB, apiKeyEntity *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error)
	GetAllApiKeys(ctx context.Context, filter *request.GetApiKeyRequestFilter) ([]*entity.ApiKeyEntity, error)
	GetAPIKeyByHashKey(ctx context.Context, hashKey string) (*entity.ApiKeyEntity, error)
	GetAPIKeyByIDAndWorkspaceID(ctx context.Context, apiKeyID, workspaceID int64) (*entity.ApiKeyEntity, error)
	UpdateAPIKey(ctx context.Context, tx *gorm.DB, apiKeyEntity *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error)
}
