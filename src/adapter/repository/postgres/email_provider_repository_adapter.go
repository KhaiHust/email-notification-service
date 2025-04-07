package postgres

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailProviderRepositoryAdapter struct {
	base
}

func (e EmailProviderRepositoryAdapter) SaveEmailProvider(ctx context.Context, tx *gorm.DB, emailProvider *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error) {
	//TODO implement me
	panic("implement me")
}

func NewEmailProviderRepositoryAdapter(db *gorm.DB) port.IEmailProviderRepositoryPort {
	return &EmailProviderRepositoryAdapter{
		base: base{db},
	}
}
