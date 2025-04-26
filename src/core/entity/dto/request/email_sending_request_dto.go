package request

type EmailSendingRequestDto struct {
	TemplateId int64
	Data       []*EmailSendingData
}
type EmailSendingData struct {
	To      string
	Subject map[string]string
	Body    map[string]string
}
