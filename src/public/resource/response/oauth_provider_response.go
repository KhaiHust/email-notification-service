package response

import "github.com/KhaiHust/email-notification-service/core/entity/dto/response"

type OAuthProviderResponse struct {
	Url string `json:"url"`
}

func ToOAuthProviderResponse(rspDto *response.OAuthUrlResponseDto) *OAuthProviderResponse {
	return &OAuthProviderResponse{
		Url: rspDto.Url,
	}
}
