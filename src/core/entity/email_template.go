package entity

import "encoding/json"

type EmailTemplateEntity struct {
	BaseEntity
	Name          string
	Subject       string
	Body          string
	Variables     json.RawMessage
	WorkspaceId   int64
	CreatedBy     int64
	LastUpdatedBy int64
}
