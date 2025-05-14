package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/usecase"
)

type IEmailTrackingService interface {
	OpenEmailTracking(ctx context.Context, encryptTrackingID string) error
}
type EmailTrackingService struct {
	emailTrackingUsecase usecase.IEmailTrackingUsecase
}

func (e EmailTrackingService) OpenEmailTracking(ctx context.Context, encryptTrackingID string) error {
	return e.emailTrackingUsecase.OpenEmailTracking(ctx, encryptTrackingID)
}

func NewEmailTrackingService(
	emailTrackingUsecase usecase.IEmailTrackingUsecase,
) IEmailTrackingService {
	return &EmailTrackingService{
		emailTrackingUsecase: emailTrackingUsecase,
	}
}
