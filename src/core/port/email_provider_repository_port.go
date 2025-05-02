package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IEmailProviderRepositoryPort interface {
	SaveEmailProvider(ctx context.Context, tx *gorm.DB, emailProvider *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error)
	GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error)
	UpdateEmailProvider(ctx context.Context, tx *gorm.DB, emailProvider *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error)
	GetEmailProviderByWorkspaceID(ctx context.Context, workspaceID int64) (*entity.EmailProviderEntity, error)
	GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) (*entity.EmailProviderEntity, error)
}
