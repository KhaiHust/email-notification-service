package model

import (
	"encoding/json"
	"time"
)

type EmailRequestModel struct {
	BaseModel
	TemplateId         int64               `gorm:"column:template_id"`
	Recipient          string              `gorm:"column:recipient"`
	Data               json.RawMessage     `gorm:"column:data"`
	Status             string              `gorm:"column:status"`
	ErrorMessage       string              `gorm:"column:error_message"`
	RetryCount         int64               `gorm:"column:retry_count"`
	SentAt             *time.Time          `gorm:"column:sent_at"`
	RequestID          string              `gorm:"column:request_id"`
	CorrelationID      string              `gorm:"column:correlation_id"`
	WorkspaceID        int64               `gorm:"column:workspace_id"`
	ProviderID         int64               `gorm:"column:provider_id"`
	EmailTemplateModel *EmailTemplateModel `gorm:"foreignKey:template_id;references:id"`
	EmailProviderModel *EmailProviderModel `gorm:"foreignKey:email_provider_id"`
}

func (EmailRequestModel) TableName() string {
	return "email_requests"
}
