package postgres

import (
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailProviderRepositoryAdapter struct {
	base
}

func NewEmailProviderRepositoryAdapter(db *gorm.DB) port.IEmailProviderRepositoryPort {
	return &EmailProviderRepositoryAdapter{
		base: base{db},
	}
}
