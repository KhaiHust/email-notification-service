package postgres

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailTemplateRepositoryAdapter struct {
	base
}

func (e EmailTemplateRepositoryAdapter) SaveNewTemplate(ctx context.Context, tx *gorm.DB, template *entity.EmailTemplateEntity) (*entity.EmailTemplateEntity, error) {
	emailTemplateModel := mapper.ToEmailTemplateModel(template)
	if err := tx.WithContext(ctx).Model(&model.EmailTemplateModel{}).Create(emailTemplateModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailTemplateEntity(emailTemplateModel), nil
}

func NewEmailTemplateRepositoryAdapter(db *gorm.DB) port.IEmailTemplateRepositoryPort {
	return &EmailTemplateRepositoryAdapter{
		base: base{db: db},
	}
}
