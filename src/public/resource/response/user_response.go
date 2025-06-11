package response

import "github.com/KhaiHust/email-notification-service/core/entity"

type UserResponse struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func ToUserResponseResource(userEntity *entity.UserEntity) *UserResponse {
	if userEntity == nil {
		return nil
	}
	user := &UserResponse{
		FullName: userEntity.FullName,
		Email:    userEntity.Email,
	}
	if userEntity.WorkspaceUserEntity != nil {
		user.Role = userEntity.WorkspaceUserEntity.Role
	}
	return user
}
func ToListUserResponseResource(userEntities []*entity.UserEntity) []*UserResponse {
	if userEntities == nil {
		return nil
	}
	var users []*UserResponse
	for _, user := range userEntities {
		users = append(users, ToUserResponseResource(user))
	}
	return users
}
