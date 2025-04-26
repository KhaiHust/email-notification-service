package request

type EmailDataDto struct {
	Tos     []string
	Subject string
	Body    string
}
