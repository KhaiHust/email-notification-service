package usecase

import "github.com/KhaiHust/email-notification-service/core/port"

type ICreateUserUseCase interface {
}
type CreateUserUseCase struct {
	getUserUseCase             IGetUserUseCase
	databaseTransactionUseCase IDatabaseTransactionUseCase
	userRepo                   port.IUserRepositoryPort
}

func NewCreateUserUseCase(
	getUserUseCase IGetUserUseCase,
	databaseTransactionUseCase IDatabaseTransactionUseCase,
	userRepo port.IUserRepositoryPort,
) ICreateUserUseCase {
	return &CreateUserUseCase{
		getUserUseCase:             getUserUseCase,
		databaseTransactionUseCase: databaseTransactionUseCase,
		userRepo:                   userRepo,
	}
}
