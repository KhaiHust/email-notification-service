package model

import (
	"time"
)

type EmailRequestModel struct {
	BaseModel
	TemplateId         int64               `gorm:"column:template_id"`
	Recipient          string              `gorm:"column:recipient"`
	Data               string              `gorm:"column:data"`
	Status             string              `gorm:"column:status"`
	ErrorMessage       string              `gorm:"column:error_message"`
	RetryCount         int64               `gorm:"column:retry_count"`
	SentAt             *time.Time          `gorm:"column:sent_at"`
	RequestID          string              `gorm:"column:request_id"`
	CorrelationID      string              `gorm:"column:correlation_id"`
	WorkspaceID        int64               `gorm:"column:workspace_id"`
	EmailProviderID    int64               `gorm:"column:email_provider_id"`
	TrackingID         string              `gorm:"column:tracking_id"`
	OpenedAt           *time.Time          `gorm:"column:opened_at"`
	OpenedCount        int64               `gorm:"column:opened_count"`
	EmailTemplateModel *EmailTemplateModel `gorm:"foreignKey:template_id;references:id"`
	EmailProviderModel *EmailProviderModel `gorm:"foreignKey:email_provider_id;references:id"`
	SendAt             *time.Time          `gorm:"column:send_at"`
}

func (EmailRequestModel) TableName() string {
	return "email_requests"
}
