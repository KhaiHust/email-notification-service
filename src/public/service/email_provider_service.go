package service

import "github.com/KhaiHust/email-notification-service/core/usecase"

type IEmailProviderService interface {
}
type EmailProviderService struct {
	getEmailProviderUseCase usecase.IGetEmailProviderUseCase
}

func NewEmailProviderService(getEmailProviderUseCase usecase.IGetEmailProviderUseCase) IEmailProviderService {
	return &EmailProviderService{
		getEmailProviderUseCase: getEmailProviderUseCase,
	}
}
