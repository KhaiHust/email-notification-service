package response

import "github.com/KhaiHust/email-notification-service/core/entity/dto/response"

type LoginResponse struct {
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	UserInfo     UserInfoResponse `json:"user_info"`
}
type UserInfoResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func ToLoginResponseResource(loginDto *response.LoginResponseDto) *LoginResponse {
	return &LoginResponse{
		AccessToken:  loginDto.AccessToken,
		RefreshToken: loginDto.RefreshToken,
		UserInfo: UserInfoResponse{
			FullName: loginDto.UserInfo.FullName,
			Email:    loginDto.UserInfo.Email,
		},
	}
}
