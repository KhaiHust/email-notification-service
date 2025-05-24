package response

import "github.com/KhaiHust/email-notification-service/core/entity/dto"

type TemplateMetricResponse struct {
	ChartStats   []*ChartStatResponse  `json:"chart_stats,omitempty"`
	TemplateStat *TemplateStatResponse `json:"template_stat,omitempty"`
}
type ChartStatResponse struct {
	Period    string  `json:"period"`
	Sent      int64   `json:"sent"`
	Error     int64   `json:"error"`
	Open      float64 `json:"open_rate"`
	Scheduled int64   `json:"scheduled"`
}
type TemplateStatResponse struct {
	Sent          int64                   `json:"sent"`
	Error         int64                   `json:"error"`
	Open          float64                 `json:"open_rate"`
	Scheduled     int64                   `json:"scheduled"`
	ProviderStats []*ProviderStatResponse `json:"provider_stats"`
}
type ProviderStatResponse struct {
	ProviderID   int64   `json:"provider_id"`
	ProviderName string  `json:"provider_name"`
	Sent         int64   `json:"sent"`
	Error        int64   `json:"error"`
	Open         float64 `json:"open"`
	Scheduled    int64   `json:"scheduled"`
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
			Period:    chartStat.Period.Format("2006-01-02"),
			Sent:      chartStat.Sent,
			Error:     chartStat.Error,
			Open:      calculateOpenRate(chartStat.Open, chartStat.Sent),
			Scheduled: chartStat.Scheduled,
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
		Open:          calculateOpenRate(templateStat.Open, templateStat.Sent),
		Scheduled:     templateStat.Scheduled,
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
			Open:         calculateOpenRate(providerStat.Open, providerStat.Sent),
			Scheduled:    providerStat.Scheduled,
		}
	}
	return providerStatResponses
}
func calculateOpenRate(open int64, sent int64) float64 {
	if sent == 0 {
		return 0
	}
	return float64(open) / float64(sent+open) * 100
}
