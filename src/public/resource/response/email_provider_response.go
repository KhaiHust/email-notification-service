package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type EmailProviderResponse struct {
	ID       int64  `json:"id,omitempty"`
	Provider string `json:"provider,omitempty"`
	Email    string `json:"email,omitempty"`
	FromName string `json:"from_name,omitempty"`
	Env      string `json:"env,omitempty"`
}

func ToEmailProviderResponse(emailProvider *entity.EmailProviderEntity) *EmailProviderResponse {
	if emailProvider == nil {
		return nil
	}
	return &EmailProviderResponse{
		ID:       emailProvider.ID,
		Provider: emailProvider.Provider,
		Email:    emailProvider.Email,
		FromName: emailProvider.FromName,
		Env:      emailProvider.Environment,
	}
}
func ToEmailProviderResponseList(emailProviders []*entity.EmailProviderEntity) []*EmailProviderResponse {
	if emailProviders == nil {
		return nil
	}
	emailProviderResponses := make([]*EmailProviderResponse, len(emailProviders))
	for i, emailProvider := range emailProviders {
		emailProviderResponses[i] = ToEmailProviderResponse(emailProvider)
	}
	return emailProviderResponses
}
