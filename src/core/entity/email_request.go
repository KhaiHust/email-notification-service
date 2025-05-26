package entity

import (
	"encoding/json"
)

type EmailRequestEntity struct {
	BaseEntity
	TemplateId          int64
	Recipient           string
	Data                json.RawMessage
	Status              string
	ErrorMessage        string
	RetryCount          int64
	SentAt              *int64
	RequestID           string
	CorrelationID       string
	WorkspaceID         int64
	EmailProviderID     int64
	TrackingID          string
	OpenedAt            *int64
	OpenedCount         int64
	EmailTemplateEntity *EmailTemplateEntity
	EmailProviderEntity *EmailProviderEntity
	SendAt              *int64 // For scheduling emails
}
