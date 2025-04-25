package entity

import (
	"encoding/json"
)

type EmailRequestEntity struct {
	BaseEntity
	TemplateId   int64
	Recipient    string
	Data         json.RawMessage
	Status       string
	ErrorMessage string
	RetryCount   int64
	SentAt       *int64
}
