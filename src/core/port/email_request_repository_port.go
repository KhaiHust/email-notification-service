package port

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"gorm.io/gorm"
)

type IEmailRequestRepositoryPort interface {
	SaveEmailRequestByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
	UpdateStatusByBatches(ctx context.Context, tx *gorm.DB, emailRequests []*entity.EmailRequestEntity) ([]*entity.EmailRequestEntity, error)
	UpdateEmailRequestByID(ctx context.Context, tx *gorm.DB, emailRequest *entity.EmailRequestEntity) (*entity.EmailRequestEntity, error)
	GetEmailRequestByID(ctx context.Context, emailRequestID int64) (*entity.EmailRequestEntity, error)
	CountEmailRequestStatuses(ctx context.Context, filter *request.EmailRequestFilter) ([]*entity.EmailRequestStatusCountEntity, error)
	GetAllEmailRequest(ctx context.Context, filter *request.EmailRequestFilter) ([]*entity.EmailRequestEntity, error)
	CountAllEmailRequest(ctx context.Context, filter *request.EmailRequestFilter) (int64, error)
	GetEmailRequestForUpdateByIDOrTrackingID(ctx context.Context, tx *gorm.DB, emailRequestID int64, trackingID string) (*entity.EmailRequestEntity, error)
	GetTotalSendVolumeByDate(ctx context.Context, filter *request.SendVolumeFilter) (map[string]int64, error)
	GetTotalSendVolumeByProvider(ctx context.Context, filter *request.SendVolumeFilter) (map[string]interface{}, error)
}
