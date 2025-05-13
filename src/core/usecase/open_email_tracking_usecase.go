package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IEmailTrackingUsecase interface {
	OpenEmailTracking(ctx context.Context, encryptTrackingID string) error
}
type EmailTrackingUsecase struct {
	encryptUseCase         IEncryptUseCase
	getEmailRequestUsecase IGetEmailRequestUsecase
	eventPublisher         port.IEventPublisher
}

func (e EmailTrackingUsecase) OpenEmailTracking(ctx context.Context, encryptTrackingID string) error {
	trackingID, err := e.encryptUseCase.DecryptTrackingID(ctx, encryptTrackingID)
	if err != nil {
		log.Error(ctx, "[EmailTrackingUsecase] Error when decrypt tracking id", err)
		return common.ErrInvalidEmailTrackingID
	}

}

func NewEmailTrackingUsecase() IEmailTrackingUsecase {
	return &EmailTrackingUsecase{}
}
