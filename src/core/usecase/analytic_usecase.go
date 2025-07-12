package usecase

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/golibs-starter/golib/log"
	"time"
)

type IAnalyticUsecase interface {
	GetSendVolumes(ctx context.Context, filter *request.SendVolumeFilter) (map[string]*dto.SendVolumeDTO, error)
	GetTemplateMetrics(ctx context.Context, filter *request.TemplateMetricFilter) (*dto.TemplateMetricDTO, error)
	GetSendVolumeByProvider(ctx context.Context, filter *request.SendVolumeFilter) ([]*dto.SendVolumeByProviderDto, error)
}
type AnalyticUsecase struct {
	emailRequestRepositoryPort  port.IEmailRequestRepositoryPort
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
}

func (a AnalyticUsecase) GetSendVolumeByProvider(ctx context.Context, filter *request.SendVolumeFilter) ([]*dto.SendVolumeByProviderDto, error) {
	volumesByProvider, err := a.emailRequestRepositoryPort.GetVolumeProvider(ctx, filter)
	if err != nil {
		log.Error(ctx, "Error when get send volume by provider", err)
		return nil, err
	}
	if len(volumesByProvider) == 0 {
		return []*dto.SendVolumeByProviderDto{}, nil
	}
	providerIDs := make([]int64, 0)
	for _, v := range volumesByProvider {
		providerIDs = append(providerIDs, v.ProviderID)
	}
	providerMap, err := a.buildProviderMap(ctx, providerIDs)
	if err != nil {
		log.Error(ctx, "Error when build provider map", err)
		return nil, err
	}
	for _, v := range volumesByProvider {
		if provider, ok := providerMap[v.ProviderID]; ok {
			v.Provider = provider.Provider
		}
	}

	volumesByProviderMap := make(map[string]*dto.SendVolumeByProviderDto)
	for _, v := range volumesByProvider {
		if existing, ok := volumesByProviderMap[v.Provider]; ok {
			existing.Total += v.Total
			existing.TotalSent += v.TotalSent
			existing.TotalError += v.TotalError
		} else {
			volumesByProviderMap[v.Provider] = v
		}
	}
	//convert map to slice
	volumesByProviderSlice := make([]*dto.SendVolumeByProviderDto, 0, len(volumesByProviderMap))
	for _, v := range volumesByProviderMap {
		volumesByProviderSlice = append(volumesByProviderSlice, v)
	}
	return volumesByProviderSlice, nil
}

func (a AnalyticUsecase) GetTemplateMetrics(ctx context.Context, filter *request.TemplateMetricFilter) (*dto.TemplateMetricDTO, error) {
	var response dto.TemplateMetricDTO
	if filter.IsChart {
		//get chart
		chartStats, err := a.getChartStats(ctx, filter)
		if err != nil {
			log.Error(ctx, "Error when get chart stats", err)
			return nil, err
		}
		response.ChartStats = chartStats
		return &response, nil
	}

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
	// combine provider stats
	providerStatsMap := make(map[string]*dto.ProviderStat)
	for _, providerStat := range providerStats {
		if existing, ok := providerStatsMap[providerStat.ProviderName]; ok {
			existing.Sent += providerStat.Sent
			existing.Error += providerStat.Error
			existing.Scheduled += providerStat.Scheduled
			existing.Open += providerStat.Open
		} else {
			providerStatsMap[providerStat.ProviderName] = providerStat
		}
	}
	// convert map to slice
	providerStats = make([]*dto.ProviderStat, 0, len(providerStatsMap))
	for _, providerStat := range providerStatsMap {
		providerStats = append(providerStats, providerStat)
	}
	response.TemplateStat.ProviderStats = providerStats
	return &response, nil
}
func (a AnalyticUsecase) getChartStats(ctx context.Context, filter *request.TemplateMetricFilter) ([]*dto.ChartStatDto, error) {
	//build start date and end date
	// End date is the current date and time is 23:59:59
	endDate := time.Now().Truncate(time.Hour * 24).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
	filter.EndDate = utils.ToUnixTimeToPointer(&endDate)

	internval := time.Hour
	switch filter.Interval {
	case constant.IntervalDay:
		internval = time.Hour * 24
	case constant.IntervalWeek:
		internval = time.Hour * 24 * 7
	case constant.IntervalMonth:
		internval = time.Hour * 24 * 30
	default:
		return nil, errors.New("Invalid interval")
	}
	// StartDate is endDate - filter.Interval * filter.Duration and time is 00:00:00
	startDate := endDate.Add(-internval * time.Duration(filter.Duration)).Truncate(time.Hour * 24)
	filter.StartDate = utils.ToUnixTimeToPointer(&startDate)
	chartStats, err := a.emailRequestRepositoryPort.GetChartStats(ctx, filter)
	if err != nil {
		log.Error(ctx, "Error when get chart stats", err)
		return nil, err
	}
	return chartStats, nil
}
func (a AnalyticUsecase) GetSendVolumes(ctx context.Context, filter *request.SendVolumeFilter) (map[string]*dto.SendVolumeDTO, error) {
	volumesByDate, err := a.GetSendVolumeByDate(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Get total send volume by provider
	volumesByProvider, err := a.GetSendVolumeByProviderByDate(ctx, filter)
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
			volumesProviderMap := make(map[string]int64)
			// Iterate over the map to build the provider map
			for providerID, totalSend := range m {
				if provider, ok := providerMap[providerID]; ok {
					if _, exists := volumesProviderMap[provider.Provider]; !exists {
						volumesProviderMap[provider.Provider] = totalSend
					} else {
						// If the provider already exists, add to the total
						volumesProviderMap[provider.Provider] += totalSend
					}
				}
			}
			sendVolumeDto.TotalSendByProvider = volumesProviderMap
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
func (a AnalyticUsecase) GetSendVolumeByProviderByDate(ctx context.Context, filter *request.SendVolumeFilter) (map[string]interface{}, error) {
	volumesByProvider, err := a.emailRequestRepositoryPort.GetTotalSendVolumeProviderByDate(ctx, filter)
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
