package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
)

type IAnalyticService interface {
	GetSendVolumes(ctx context.Context, filter *request.SendVolumeFilter) (interface{}, error)
}
type AnalyticService struct {
	analyticUsecase usecase.IAnalyticUsecase
}

func (a AnalyticService) GetSendVolumes(ctx context.Context, filter *request.SendVolumeFilter) (interface{}, error) {
	// Call the usecase to get send volumes
	volumes, err := a.analyticUsecase.GetSendVolumes(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Return the send volumes
	return response.ToSendVolumeResponse(volumes), nil
}

func NewAnalyticService(analyticUsecase usecase.IAnalyticUsecase) IAnalyticService {
	return &AnalyticService{
		analyticUsecase: analyticUsecase,
	}
}
