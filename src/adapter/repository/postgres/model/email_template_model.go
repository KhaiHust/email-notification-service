package model

import "encoding/json"

type EmailTemplateModel struct {
	BaseModel
	Name          string          `gorm:"column:name"`
	Subject       string          `gorm:"column:subject"`
	Body          string          `gorm:"column:body"`
	Variables     json.RawMessage `gorm:"column:variables"`
	WorkspaceId   int64           `gorm:"column:workspace_id"`
	CreatedBy     int64           `gorm:"column:created_by"`
	LastUpdatedBy int64           `gorm:"column:last_updated_by"`
	Version       string          `gorm:"column:version"`
}

func (EmailTemplateModel) TableName() string {
	return "email_templates"
}
