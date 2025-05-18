package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/event"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
	"time"
)

type IEmailTrackingUsecase interface {
	OpenEmailTracking(ctx context.Context, encryptTrackingID string) error
}
type EmailTrackingUsecase struct {
	encryptUseCase IEncryptUseCase
	eventPublisher port.IEventPublisher
}

func (e EmailTrackingUsecase) OpenEmailTracking(ctx context.Context, encryptTrackingID string) error {
	trackingID, err := e.encryptUseCase.DecryptTrackingID(ctx, encryptTrackingID)
	if err != nil {
		log.Error(ctx, "[EmailTrackingUsecase] Error when decrypt tracking id", err)
		return common.ErrInvalidEmailTrackingID
	}
	openedAt := time.Now()
	emailRequestEntity := &entity.EmailRequestEntity{
		TrackingID: trackingID,
		Status:     constant.EmailSendingStatusOpened,
		OpenedAt:   utils.ToUnixTimeToPointer(&openedAt),
	}
	ev := event.NewEventEmailRequestSync(ctx, emailRequestEntity)
	e.eventPublisher.Publish(ev)
	return nil
}

func NewEmailTrackingUsecase(
	encryptUseCase IEncryptUseCase,
	eventPublisher port.IEventPublisher,
) IEmailTrackingUsecase {
	return &EmailTrackingUsecase{
		encryptUseCase: encryptUseCase,
		eventPublisher: eventPublisher,
	}
}
