package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IApiKeyRepositoryPort interface {
	SaveNewApiKey(ctx context.Context, tx *gorm.DB, apiKeyEntity *entity.ApiKeyEntity) (*entity.ApiKeyEntity, error)
}
