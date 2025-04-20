package postgres

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailProviderRepositoryAdapter struct {
	base
}

func (e EmailProviderRepositoryAdapter) SaveEmailProvider(ctx context.Context, tx *gorm.DB, emailProvider *entity.EmailProviderEntity) (*entity.EmailProviderEntity, error) {
	emailProviderModel := mapper.ToEmailProviderModel(emailProvider)

	if err := tx.WithContext(ctx).Create(emailProviderModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailProviderEntity(emailProviderModel), nil
}

func NewEmailProviderRepositoryAdapter(db *gorm.DB) port.IEmailProviderRepositoryPort {
	return &EmailProviderRepositoryAdapter{
		base: base{db},
	}
}
