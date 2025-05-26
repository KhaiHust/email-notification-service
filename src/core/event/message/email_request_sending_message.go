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
}
