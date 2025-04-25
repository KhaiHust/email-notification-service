package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"gorm.io/gorm"
)

type IEmailRequestRepositoryPort interface {
	SaveEmailRequestByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
	UpdateStatusByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
}
