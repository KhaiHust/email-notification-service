package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IEmailProviderRepositoryPort interface {
	SaveEmailProvider(ctx context.Context, tx *gorm.DB, emailProvider *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error)
}
