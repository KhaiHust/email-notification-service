package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IEmailLogRepositoryPort interface {
	SaveNewEmailLog(ctx context.Context, tx *gorm.DB, emailLog *entity.EmailLogsEntity) (*entity.EmailLogsEntity, error)
}
