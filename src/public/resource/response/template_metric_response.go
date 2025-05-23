package response

import "github.com/KhaiHust/email-notification-service/core/entity/dto"

type TemplateMetricResponse struct {
	ChartStats   []*ChartStatResponse  `json:"chart_stats"`
	TemplateStat *TemplateStatResponse `json:"template_stat"`
}
type ChartStatResponse struct {
	Period string `json:"period"`
	Sent   int64  `json:"sent"`
	Error  int64  `json:"error"`
	Open   int64  `json:"open"`
}
type TemplateStatResponse struct {
	Sent          int64                   `json:"sent"`
	Error         int64                   `json:"error"`
	Open          int64                   `json:"open"`
	ProviderStats []*ProviderStatResponse `json:"provider_stats"`
}
type ProviderStatResponse struct {
	ProviderID   int64  `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	Sent         int64  `json:"sent"`
	Error        int64  `json:"error"`
	Open         int64  `json:"open"`
}

func ToTemplateMertricResponse(metricDto *dto.TemplateMetricDTO) *TemplateMetricResponse {
	if metricDto == nil {
		return nil
	}
	return &TemplateMetricResponse{
		ChartStats:   ToChartStatResponse(metricDto.ChartStats),
		TemplateStat: ToTemplateStatResponse(metricDto.TemplateStat),
	}
}

func ToChartStatResponse(chartStats []*dto.ChartStatDto) []*ChartStatResponse {
	if chartStats == nil {
		return nil
	}
	chartStatResponses := make([]*ChartStatResponse, len(chartStats))
	for i, chartStat := range chartStats {
		chartStatResponses[i] = &ChartStatResponse{
			Period: chartStat.Period.Format("2006-01-02"),
			Sent:   chartStat.Sent,
			Error:  chartStat.Error,
			Open:   chartStat.Open,
		}
	}
	return chartStatResponses
}
func ToTemplateStatResponse(templateStat *dto.TemplateStat) *TemplateStatResponse {
	if templateStat == nil {
		return nil
	}
	return &TemplateStatResponse{
		Sent:          templateStat.Sent,
		Error:         templateStat.Error,
		Open:          templateStat.Open,
		ProviderStats: ToProviderStatResponse(templateStat.ProviderStats),
	}
}
func ToProviderStatResponse(providerStats []*dto.ProviderStat) []*ProviderStatResponse {
	if providerStats == nil {
		return nil
	}
	providerStatResponses := make([]*ProviderStatResponse, len(providerStats))
	for i, providerStat := range providerStats {
		providerStatResponses[i] = &ProviderStatResponse{
			ProviderID:   providerStat.ProviderID,
			ProviderName: providerStat.ProviderName,
			Sent:         providerStat.Sent,
			Error:        providerStat.Error,
			Open:         providerStat.Open,
		}
	}
	return providerStatResponses
}
