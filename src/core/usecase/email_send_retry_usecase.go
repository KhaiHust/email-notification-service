package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/event"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
	"time"
)

const (
	BatchSize = int64(100)
)

type IEmailSendRetryUsecase interface {
	ProcessBatches(ctx context.Context) error
}
type EmailSendRetryUsecase struct {
	getEmailRequestUsecase IGetEmailRequestUsecase
	eventPublisher         port.IEventPublisher
}

func (e EmailSendRetryUsecase) ProcessBatches(ctx context.Context) error {
	now := time.Now().Unix()
	filter := &request.EmailRequestFilter{
		Statuses: []string{constant.EmailSendingStatusFailed},
		BaseFilter: &request.BaseFilter{
			Limit:       utils.ToInt64Pointer(BatchSize),
			SortOrder:   constant.ASC,
			UpdatedAtTo: utils.ToInt64Pointer(now),
		},
	}
	for {
		emailRequests, total, err := e.getEmailRequestUsecase.GetAllEmailRequest(ctx, filter)
		if err != nil {
			log.Error(ctx, "Failed to get email requests: %v", err)
			return err
		}
		if total == 0 || len(emailRequests) == 0 {
			break
		}
		ev := event.NewEventRequestSendingEmail(ctx, emailRequests)
		if err = e.eventPublisher.SyncPublish(ctx, ev); err != nil {
			log.Error(ctx, "Failed to publish email sending event: %v", err)
			return err
		}
		log.Info(ctx, "Processing batch of email requests: %d", len(emailRequests))
		filter.BaseFilter.Since = &emailRequests[len(emailRequests)-1].ID
	}
	log.Info(ctx, "Finished processing all batches of email requests")
	return nil
}

func NewEmailSendRetryUsecase(
	getEmailRequestUsecase IGetEmailRequestUsecase,
	eventPublisher port.IEventPublisher,
) IEmailSendRetryUsecase {
	return &EmailSendRetryUsecase{
		getEmailRequestUsecase: getEmailRequestUsecase,
		eventPublisher:         eventPublisher,
	}
}
