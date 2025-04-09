package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
)

type IUserService interface {
	SignUp(ctx context.Context, req *request.CreateUserRequest) (*entity.UserEntity, error)
}
type UserService struct {
	createUserUseCase usecase.ICreateUserUseCase
}

func (u UserService) SignUp(ctx context.Context, req *request.CreateUserRequest) (*entity.UserEntity, error) {
	userEntity := request.ToUserEntity(req)
	user, err := u.createUserUseCase.CreateNewUser(ctx, userEntity)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserService(createUserUseCase usecase.ICreateUserUseCase) IUserService {
	return &UserService{
		createUserUseCase: createUserUseCase,
	}
}
