package message

type EmailRequestSyncMessage struct {
	EmailRequestID int64  `json:"email_request_id"`
	Status         string `json:"status"`
	ErrorMessage   string `json:"error_message"`
	SentAt         *int64 `json:"sent_at"`
	TrackingID     string `json:"tracking_id"`
	OpenedAt       *int64 `json:"opened_at"`
	OpenedCount    int64  `json:"opened_count"`
}
