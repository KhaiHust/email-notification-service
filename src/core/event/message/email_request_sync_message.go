package message

type EmailRequestSyncMessage struct {
	EmailRequestID int64  `json:"email_request_id"`
	Status         string `json:"status"`
	ErrorMessage   string `json:"error_message"`
	SentAt         *int64 `json:"sent_at"`
}
