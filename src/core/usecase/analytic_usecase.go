package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IAnalyticUsecase interface {
	GetSendVolumes(ctx context.Context, filter *request.SendVolumeFilter) (map[string]*dto.SendVolumeDTO, error)
	GetTemplateMetrics(ctx context.Context, filter *request.TemplateMetricFilter) (*dto.TemplateMetricDTO, error)
}
type AnalyticUsecase struct {
	emailRequestRepositoryPort  port.IEmailRequestRepositoryPort
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
}

func (a AnalyticUsecase) GetTemplateMetrics(ctx context.Context, filter *request.TemplateMetricFilter) (*dto.TemplateMetricDTO, error) {
	var response dto.TemplateMetricDTO
	//get chart
	chartStats, err := a.emailRequestRepositoryPort.GetChartStats(ctx, filter)
	if err != nil {
		log.Error(ctx, "Error when get chart stats", err)
		return nil, err
	}
	response.ChartStats = chartStats
	//get template stats
	templateStat, err := a.emailRequestRepositoryPort.GetTemplateStats(ctx, filter)
	if err != nil {
		log.Error(ctx, "Error when get template stats", err)
		return nil, err
	}
	response.TemplateStat = templateStat
	//get provider stats
	providerStats, err := a.emailRequestRepositoryPort.GetTemplateStatsByProvider(ctx, filter)
	if err != nil {
		log.Error(ctx, "Error when get provider stats", err)
		return nil, err
	}
	response.TemplateStat.ProviderStats = providerStats
	return &response, nil
}
func (a AnalyticUsecase) GetSendVolumes(ctx context.Context, filter *request.SendVolumeFilter) (map[string]*dto.SendVolumeDTO, error) {
	// Get total send volume by date
	volumesByDate, err := a.GetSendVolumeByDate(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Get total send volume by provider
	volumesByProvider, err := a.GetSendVolumeByProvider(ctx, filter)
	if err != nil {
		return nil, err
	}

	//get provider id
	providerIDs := make([]int64, 0)
	for _, v := range volumesByProvider {
		if m, ok := v.(map[int64]int64); ok {
			for providerID := range m {
				providerIDs = append(providerIDs, providerID)
			}
		}
	}
	providerMap, err := a.buildProviderMap(ctx, providerIDs)
	if err != nil {
		log.Error(ctx, "Error when build provider map", err)
		return nil, err
	}
	//build dto
	sendVolumeDtoMap := make(map[string]*dto.SendVolumeDTO)
	for date, total := range volumesByDate {
		sendVolumeDto := &dto.SendVolumeDTO{
			TotalSend: total,
		}
		volumesProvider := volumesByProvider[date]
		if m, ok := volumesProvider.(map[int64]int64); ok {
			for providerID, totalSend := range m {
				if provider, ok := providerMap[providerID]; ok {
					volumesProviderMap := make(map[string]int64)
					volumesProviderMap[provider.Provider] = totalSend
					sendVolumeDto.TotalSendByProvider = volumesProviderMap
				}
			}
		}
		sendVolumeDtoMap[date] = sendVolumeDto
	}
	return sendVolumeDtoMap, nil
}

func (a AnalyticUsecase) GetSendVolumeByDate(ctx context.Context, filter *request.SendVolumeFilter) (map[string]int64, error) {
	volumesByDate, err := a.emailRequestRepositoryPort.GetTotalSendVolumeByDate(ctx, filter)
	if err != nil {
		return nil, err
	}
	return volumesByDate, nil
}
func (a AnalyticUsecase) GetSendVolumeByProvider(ctx context.Context, filter *request.SendVolumeFilter) (map[string]interface{}, error) {
	volumesByProvider, err := a.emailRequestRepositoryPort.GetTotalSendVolumeByProvider(ctx, filter)
	if err != nil {
		return nil, err
	}
	return volumesByProvider, nil
}
func (a AnalyticUsecase) buildProviderMap(ctx context.Context, providerIDs []int64) (map[int64]*entity.EmailProviderEntity, error) {
	providerMap := make(map[int64]*entity.EmailProviderEntity)
	providers, err := a.emailProviderRepositoryPort.GetEmailProviderByIds(ctx, providerIDs)
	if err != nil {
		return nil, err
	}
	for _, provider := range providers {
		providerMap[provider.ID] = provider
	}
	return providerMap, nil
}
func NewAnalyticUsecase(
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort,
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort,
) IAnalyticUsecase {
	return &AnalyticUsecase{
		emailRequestRepositoryPort:  emailRequestRepositoryPort,
		emailProviderRepositoryPort: emailProviderRepositoryPort,
	}
}
