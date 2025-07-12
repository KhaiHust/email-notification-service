package model

import "time"

type EmailLogsModel struct {
	BaseModel
	EmailRequestID  int64     `gorm:"column:email_request_id"`
	TemplateId      int64     `gorm:"column:template_id"`
	Recipient       string    `gorm:"column:recipient"`
	Status          string    `gorm:"column:status"`
	ErrorMessage    string    `gorm:"column:error_message"`
	RetryCount      int64     `gorm:"column:retry_count"`
	RequestID       string    `gorm:"column:request_id"`
	WorkspaceID     int64     `gorm:"column:workspace_id"`
	EmailProviderID int64     `gorm:"column:email_provider_id"`
	LoggedAt        time.Time `gorm:"column:logged_at"`
}

func (EmailLogsModel) TableName() string {
	return "email_logs"
}
