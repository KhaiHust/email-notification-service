package usecase

import "github.com/KhaiHust/email-notification-service/core/port"

type ICreateEmailProviderUseCase interface {
}
type CreateEmailProviderUseCase struct {
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
	getWorkspaceUseCase         IGetWorkspaceUseCase
	databaseTransactionUseCase  IDatabaseTransactionUseCase
}

func NewCreateEmailProviderUseCase(emailProviderRepositoryPort port.IEmailProviderRepositoryPort, getWorkspaceUseCase IGetWorkspaceUseCase, databaseTransactionUseCase IDatabaseTransactionUseCase) ICreateEmailProviderUseCase {
	return &CreateEmailProviderUseCase{
		emailProviderRepositoryPort: emailProviderRepositoryPort,
		getWorkspaceUseCase:         getWorkspaceUseCase,
		databaseTransactionUseCase:  databaseTransactionUseCase,
	}
}
