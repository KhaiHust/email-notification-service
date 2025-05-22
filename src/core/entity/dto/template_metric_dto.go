package dto

type TemplateMetricDTO struct {
	Overview *OverviewMetricDTO
}
type OverviewMetricDTO struct {
	TotalSend  int64
	TotalSent  int64
	TotalError int64
	TotalOpen  int64
}
