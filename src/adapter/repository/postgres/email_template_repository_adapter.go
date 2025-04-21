package postgres

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailTemplateRepositoryAdapter struct {
	base
}

func (e EmailTemplateRepositoryAdapter) GetTemplateByID(ctx context.Context, ID int64) (*entity.EmailTemplateEntity, error) {
	var emailTemplateModel model.EmailTemplateModel
	if err := e.db.WithContext(ctx).Where("id = ?", ID).First(&emailTemplateModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return mapper.ToEmailTemplateEntity(&emailTemplateModel), nil
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
