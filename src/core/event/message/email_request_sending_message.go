package message

type EmailRequestSendingMessage struct {
	SendData      []*EmailSendData `json:"send_data"`
	TemplateId    int64            `json:"template_id"`
	WorkspaceID   int64            `json:"workspace_id"`
	IntegrationID int64            `json:"integration_id"`
}
type EmailSendData struct {
	To      string            `json:"to"`
	Subject map[string]string `json:"subject"`
	Body    map[string]string `json:"body"`
}
