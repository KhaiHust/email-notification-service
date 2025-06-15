package request

type EmailSendingRequestDto struct {
	TemplateID  int64
	WorkspaceID int64
	ProviderID  int64
	Datas       []*EmailSendingData
}
type EmailSendingData struct {
	EmailRequestID int64
	TrackingID     string
	To             string
	SendAt         *int64
	IsRetry        bool
	Subject        map[string]string
	Body           map[string]string
}
