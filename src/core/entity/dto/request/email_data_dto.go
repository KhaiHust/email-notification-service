package request

type EmailDataDto struct {
	EmailRequestID int64
	TrackingID     string
	Tos            []string
	Subject        string
	Body           string
}
