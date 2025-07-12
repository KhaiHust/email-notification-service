package postgres

import (
	"context"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/mapper"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/port"
	"gorm.io/gorm"
)

type EmailLogRepositoryAdapter struct {
	base
}

func (e EmailLogRepositoryAdapter) SaveEmailLogsByBatches(ctx context.Context, tx *gorm.DB, emailLogs []*entity.EmailLogsEntity) ([]*entity.EmailLogsEntity, error) {
	emailLogsModels := mapper.ToListEmailLogModel(emailLogs)
	if err := tx.WithContext(ctx).Model(&model.EmailLogsModel{}).Save(emailLogsModels).Error; err != nil {
		return nil, err
	}
	return mapper.ToListEmailLogEntity(emailLogsModels), nil
}

func (e EmailLogRepositoryAdapter) GetLogsByEmailRequestIDAndWorkspaceID(ctx context.Context, emailRequestID int64, workspaceID int64) ([]*entity.EmailLogsEntity, error) {
	var emailLogs []*model.EmailLogsModel
	if err := e.db.WithContext(ctx).Where("email_request_id = ? AND workspace_id = ?", emailRequestID, workspaceID).Find(&emailLogs).Error; err != nil {
		return nil, err
	}
	return mapper.ToListEmailLogEntity(emailLogs), nil
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
