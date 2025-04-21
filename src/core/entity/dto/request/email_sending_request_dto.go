package request

type EmailSendingRequestDto struct {
	TemplateID string
	Data       []EmailSendingData
}
type EmailSendingData struct {
	To      string
	Subject map[string]string
	Body    map[string]string
}
