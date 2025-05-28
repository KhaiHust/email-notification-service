package request

type EmailSendingRequestDto struct {
	TemplateID    int64
	WorkspaceID   int64
	IntegrationID int64
	Datas         []*EmailSendingData
}
type EmailSendingData struct {
	EmailRequestID int64
	TrackingID     string
	To             string
	SendAt         *int64
	Subject        map[string]string
	Body           map[string]string
}
