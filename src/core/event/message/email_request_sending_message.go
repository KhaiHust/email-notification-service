package message

type EmailRequestSendingMessage struct {
	SendData      []*EmailSendData `json:"send_data"`
	TemplateId    int64            `json:"template_id"`
	WorkspaceID   int64            `json:"workspace_id"`
	IntegrationID int64            `json:"integration_id"`
}
type EmailSendData struct {
	EmailRequestID int64  `json:"email_request_id"`
	TrackingID     string `json:"tracking_id"`
	SendAt         *int64 `json:"send_at,omitempty"` // SendAt is optional, it can be nil if not scheduled
}
