package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type EmailLogResponse struct {
	ID              int64  `json:"id"`
	EmailRequestID  int64  `json:"email_request_id"`
	Status          string `json:"status"`
	ErrorMessage    string `json:"error_message"`
	LoggedAt        int64  `json:"logged_at"`
	RetryCount      int64  `json:"retry_count"`
	RequestID       string `json:"request_id"`
	WorkspaceID     int64  `json:"workspace_id"`
	EmailProviderID int64  `json:"email_provider_id"`
	TemplateId      int64  `json:"template_id"`
	Recipient       string `json:"recipient"`
	CreatedAt       int64  `json:"create_at"`
	UpdatedAt       int64  `json:"update_at"`
}

func ToEmailLogResource(emailLogEntity *entity.EmailLogsEntity) *EmailLogResponse {
	if emailLogEntity == nil {
		return nil
	}
	return &EmailLogResponse{
		ID:              emailLogEntity.ID,
		EmailRequestID:  emailLogEntity.EmailRequestID,
		Status:          emailLogEntity.Status,
		ErrorMessage:    emailLogEntity.ErrorMessage,
		LoggedAt:        emailLogEntity.LoggedAt,
		RetryCount:      emailLogEntity.RetryCount,
		RequestID:       emailLogEntity.RequestID,
		WorkspaceID:     emailLogEntity.WorkspaceID,
		EmailProviderID: emailLogEntity.EmailProviderID,
		TemplateId:      emailLogEntity.TemplateId,
		Recipient:       emailLogEntity.Recipient,
		CreatedAt:       emailLogEntity.CreatedAt,
		UpdatedAt:       emailLogEntity.CreatedAt,
	}
}
func ToListEmailLogResource(emailLogEntities []*entity.EmailLogsEntity) []*EmailLogResponse {
	if emailLogEntities == nil {
		return nil
	}
	emailLogResources := make([]*EmailLogResponse, len(emailLogEntities))
	for i, emailLogEntity := range emailLogEntities {
		emailLogResources[i] = ToEmailLogResource(emailLogEntity)
	}
	return emailLogResources
}
