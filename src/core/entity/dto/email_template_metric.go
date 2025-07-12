package dto

type EmailTemplateMetric struct {
	TotalSent   int64 `json:"total_sent"`
	TotalErrors int64 `json:"total_errors"`
}
