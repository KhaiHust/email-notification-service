package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity/dto"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
	"time"
)

type IAnalyticService interface {
	GetSendVolumes(ctx context.Context, filter *request.SendVolumeFilter) (interface{}, error)
	GetTemplateMetrics(ctx context.Context, filter *request.TemplateMetricFilter) (*dto.TemplateMetricDTO, error)
}
type AnalyticService struct {
	analyticUsecase usecase.IAnalyticUsecase
}

func (a AnalyticService) GetTemplateMetrics(ctx context.Context, filter *request.TemplateMetricFilter) (*dto.TemplateMetricDTO, error) {
	// Call the usecase to get template metrics
	metrics, err := a.analyticUsecase.GetTemplateMetrics(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Return the template metrics
	return metrics, nil
}

func (a AnalyticService) GetSendVolumes(ctx context.Context, filter *request.SendVolumeFilter) (interface{}, error) {
	if filter.StartDate == nil {
		if filter.EndDate != nil {
			startDate := *filter.EndDate - int64(14*24*time.Hour/time.Second)
			filter.StartDate = &startDate
		} else {
			now := time.Now()
			startDate := now.Unix() - int64(14*24*time.Hour/time.Second)
			filter.StartDate = &startDate
			filter.EndDate = utils.ToUnixTimeToPointer(&now)
		}
	}

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
