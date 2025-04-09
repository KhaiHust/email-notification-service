package request

import "github.com/KhaiHust/email-notification-service/core/entity"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	FullName string `json:"full_name" validate:"required,min=6,max=20"`
}

func ToUserEntity(req *CreateUserRequest) *entity.UserEntity {
	return &entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	}
}
