package entity

import (
	"encoding/json"
	"github.com/KhaiHust/email-notification-service/core/entity/dto"
)

type EmailTemplateEntity struct {
	BaseEntity
	Name          string
	Subject       string
	Body          string
	Variables     json.RawMessage
	WorkspaceId   int64
	CreatedBy     int64
	LastUpdatedBy int64
	Version       string
	Metric        *dto.EmailTemplateMetric
}
