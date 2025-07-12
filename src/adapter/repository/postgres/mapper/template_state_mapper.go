package mapper

import (
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres/model"
	"github.com/KhaiHust/email-notification-service/core/entity/dto"
)

func ToChartStatDto(chartModel *model.ChartStatModel) *dto.ChartStatDto {
	if chartModel == nil {
		return nil
	}
	return &dto.ChartStatDto{
		Period:    chartModel.Period,
		Sent:      chartModel.Sent,
		Error:     chartModel.Error,
		Open:      chartModel.Open,
		Scheduled: chartModel.Scheduled,
	}
}
func ToChartStatDtos(chartModels []*model.ChartStatModel) []*dto.ChartStatDto {
	if chartModels == nil {
		return nil
	}
	chartStats := make([]*dto.ChartStatDto, len(chartModels))
	for i, chartModel := range chartModels {
		chartStats[i] = ToChartStatDto(chartModel)
	}
	return chartStats
}
func ToTemplateStatDto(templateStatModel *model.TemplateStatModel) *dto.TemplateStat {
	if templateStatModel == nil {
		return nil
	}
	return &dto.TemplateStat{
		Sent:          templateStatModel.Sent,
		Error:         templateStatModel.Error,
		Open:          templateStatModel.Open,
		ProviderStats: ToProviderStats(templateStatModel.ProviderStats),
	}
}
func ToProviderStats(providerStatModels []*model.ProviderStatModel) []*dto.ProviderStat {
	if providerStatModels == nil {
		return nil
	}
	providerStats := make([]*dto.ProviderStat, len(providerStatModels))
	for i, providerStatModel := range providerStatModels {
		providerStats[i] = &dto.ProviderStat{
			ProviderID:   providerStatModel.ProviderID,
			ProviderName: providerStatModel.ProviderName,
			Sent:         providerStatModel.Sent,
			Error:        providerStatModel.Error,
			Open:         providerStatModel.Open,
		}
	}
	return providerStats
}
func ToProviderStatDto(providerStatModel *model.ProviderStatModel) *dto.ProviderStat {
	if providerStatModel == nil {
		return nil
	}
	return &dto.ProviderStat{
		ProviderID:   providerStatModel.ProviderID,
		ProviderName: providerStatModel.ProviderName,
		Sent:         providerStatModel.Sent,
		Error:        providerStatModel.Error,
		Open:         providerStatModel.Open,
	}
}
