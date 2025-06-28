package request

type EmailSendTaskRequest struct {
	EmailRequestID int64 `json:"email_request_id" validate:"required"`
}
