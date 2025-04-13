package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
)

type IUserService interface {
	SignUp(ctx context.Context, req *request.CreateUserRequest) (*entity.UserEntity, error)
	Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error)
}
type UserService struct {
	createUserUseCase usecase.ICreateUserUseCase
	loginUsecase      usecase.ILoginUsecase
}

func (u UserService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	result, err := u.loginUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return response.ToLoginResponseResource(result), nil
}

func (u UserService) SignUp(ctx context.Context, req *request.CreateUserRequest) (*entity.UserEntity, error) {
	userEntity := request.ToUserEntity(req)
	user, err := u.createUserUseCase.CreateNewUser(ctx, userEntity)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserService(createUserUseCase usecase.ICreateUserUseCase, loginUsecase usecase.ILoginUsecase) IUserService {
	return &UserService{
		createUserUseCase: createUserUseCase,
		loginUsecase:      loginUsecase,
	}
}
