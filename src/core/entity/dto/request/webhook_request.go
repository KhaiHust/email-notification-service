package request

type CreateWebhookRequest struct {
	WorkspaceID int64
	URL         string
	Type        string
	Name        string
	Enabled     bool
}
