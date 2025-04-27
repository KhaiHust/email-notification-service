package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IEmailRequestRepositoryPort interface {
	SaveEmailRequestByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
	UpdateStatusByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
	UpdateEmailRequestByID(ctx context.Context, tx *gorm.DB, emailRequest *entity.EmailRequestEntity) (*entity.EmailRequestEntity, error)
	GetEmailRequestByID(ctx context.Context, emailRequestID int64) (*entity.EmailRequestEntity, error)
}
