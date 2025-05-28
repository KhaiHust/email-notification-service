package request

type EmailDataDto struct {
	EmailRequestID int64
	TrackingID     string
	To             string
	Subject        string
	Body           string
}
