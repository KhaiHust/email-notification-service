package response

type LoginResponseDto struct {
	AccessToken  string
	RefreshToken string
	UserInfo     UserInfoDto
}
type UserInfoDto struct {
	FullName string
	Email    string
}
