package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/event"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
	"time"
)

const (
	BatchSize     = int64(100)
	MaxRetryCount = 5
)

type IEmailSendRetryUsecase interface {
	ProcessBatches(ctx context.Context) error
}
type EmailSendRetryUsecase struct {
	getEmailRequestUsecase    IGetEmailRequestUsecase
	eventPublisher            port.IEventPublisher
	webhookUsecase            IWebhookUsecase
	updateEmailRequestUsecase IUpdateEmailRequestUsecase
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
		RetryCount: utils.ToInt64Pointer(MaxRetryCount),
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
		emailRequestRetries := make([]*entity.EmailRequestEntity, 0, len(emailRequests))
		emailRequestMaxRetries := make([]*entity.EmailRequestEntity, 0, len(emailRequests))
		for _, emailRequest := range emailRequests {
			emailRequest.IsRetry = true
			if emailRequest.RetryCount >= MaxRetryCount {
				emailRequestMaxRetries = append(emailRequestMaxRetries, emailRequest)
			} else {
				emailRequestRetries = append(emailRequestRetries, emailRequest)
			}
		}
		if len(emailRequestMaxRetries) > 0 {
			go func() {
				if err := e.HandleMaxRetryCountReached(ctx, emailRequestMaxRetries); err != nil {
					log.Error(ctx, "Failed to handle max retry count reached: %v", err)
				} else {
					log.Info(ctx, "Handled max retry count reached for %d email requests", len(emailRequestMaxRetries))
				}
			}()
		}
		ev := event.NewEventRequestSendingEmail(ctx, emailRequestRetries)
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
	webhookUsecase IWebhookUsecase,
	updateEmailRequestUsecase IUpdateEmailRequestUsecase,
) IEmailSendRetryUsecase {
	return &EmailSendRetryUsecase{
		getEmailRequestUsecase:    getEmailRequestUsecase,
		eventPublisher:            eventPublisher,
		webhookUsecase:            webhookUsecase,
		updateEmailRequestUsecase: updateEmailRequestUsecase,
	}
}

// todo: implement when max retry count is reached
func (e EmailSendRetryUsecase) HandleMaxRetryCountReached(ctx context.Context, emailRequests []*entity.EmailRequestEntity) error {

	// For example, you can update the status of the email request to "failed permanently"
	// or send a notification to the user.
	for _, emailRequest := range emailRequests {
		emailRequest.Status = constant.EmailSendingStatusFailed
		emailRequest.ErrorMessage = "Max retry count reached"
		emailRequest.RetryCount = MaxRetryCount + 1
	}
	if _, err := e.updateEmailRequestUsecase.UpdateStatusByBatches(ctx, emailRequests); err != nil {
		log.Error(ctx, "Failed to update email requests status: %v", err)
		return err
	}
	e.webhookUsecase.SendNotifyMaxRetry(ctx, emailRequests)
	return nil
}
