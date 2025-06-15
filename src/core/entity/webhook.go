package entity

type WebhookEntity struct {
	BaseEntity
	WorkspaceID int64
	URL         string
	Type        string
	Enabled     bool
	Name        string
}
