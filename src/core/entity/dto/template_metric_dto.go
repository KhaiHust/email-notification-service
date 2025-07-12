package dto

import "time"

type TemplateMetricDTO struct {
	ChartStats   []*ChartStatDto
	TemplateStat *TemplateStat
}
type ChartStatDto struct {
	Period    time.Time
	Sent      int64
	Error     int64
	Scheduled int64
	Open      int64
}
type TemplateStat struct {
	Sent          int64
	Error         int64
	Open          int64
	Scheduled     int64
	ProviderStats []*ProviderStat
}
type ProviderStat struct {
	ProviderID   int64
	ProviderName string
	Sent         int64
	Error        int64
	Open         int64
	Scheduled    int64
}
