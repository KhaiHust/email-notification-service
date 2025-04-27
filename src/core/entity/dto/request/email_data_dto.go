package request

type EmailDataDto struct {
	EmailRequestID int64
	Tos            []string
	Subject        string
	Body           string
}
