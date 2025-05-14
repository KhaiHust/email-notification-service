package postgres

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailLogRepositoryAdapter struct {
	base
}

func (e EmailLogRepositoryAdapter) SaveNewEmailLog(ctx context.Context, tx *gorm.DB, emailLog *entity.EmailLogsEntity) (*entity.EmailLogsEntity, error) {
	logModel := mapper.ToEmailLogModel(emailLog)
	if err := tx.WithContext(ctx).Create(logModel).Error; err != nil {
		return nil, err
	}
	return mapper.ToEmailLogEntity(logModel), nil
}

func NewEmailLogRepositoryAdapter(db *gorm.DB) port.IEmailLogRepositoryPort {
	return &EmailLogRepositoryAdapter{base{db: db}}
}
