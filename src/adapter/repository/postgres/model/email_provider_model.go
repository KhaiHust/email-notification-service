package model

import "time"

type EmailProviderModel struct {
	BaseModel
	WorkspaceId       int64     `gorm:"column:workspace_id"`
	Provider          string    `gorm:"column:provider"`
	SmtpHost          string    `gorm:"column:smtp_host"`
	SmtpPort          int       `gorm:"column:smtp_port"`
	OAuthToken        string    `gorm:"column:oauth_token"`
	OAuthRefreshToken string    `gorm:"column:oauth_refresh_token"`
	OAuthExpiredAt    time.Time `gorm:"column:oauth_expires_at"`
	UseTLS            bool      `gorm:"column:use_tls"`
	Email             string    `gorm:"column:email"`
	FromName          string    `gorm:"column:from_name"`
	Environment       string    `gorm:"column:environment"`
}

func (EmailProviderModel) TableName() string {
	return "email_providers"
}
