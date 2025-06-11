package model

type WebhookModel struct {
	BaseModel
	WorkspaceID int64  `gorm:"column:workspace_id"`
	URL         string `gorm:"column:url"`
	Type        string `gorm:"column:type"`
	Enabled     bool   `gorm:"column:enabled"`
	Name        string `gorm:"column:name"`
}

func (WebhookModel) TableName() string {
	return "webhooks"
}
