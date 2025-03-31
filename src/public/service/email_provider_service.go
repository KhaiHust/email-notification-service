package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
)

type IEmailProviderService interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthProviderResponse, error)
	CreateEmailProvider(ctx context.Context, provider string, userId int64, workspaceCode, code string) (*entity.EmailProviderEntity, error)
}
type EmailProviderService struct {
	getEmailProviderUseCase    usecase.IGetEmailProviderUseCase
	createEmailProviderUseCase usecase.ICreateEmailProviderUseCase
}

func (e EmailProviderService) CreateEmailProvider(ctx context.Context, provider string, userId int64, workspaceCode, code string) (*entity.EmailProviderEntity, error) {
	result, err := e.createEmailProviderUseCase.CreateEmailProvider(ctx, userId, workspaceCode, provider, code)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e EmailProviderService) GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthProviderResponse, error) {
	result, err := e.getEmailProviderUseCase.GetOAuthUrl(ctx, provider)
	if err != nil {
		return nil, err
	}
	return response.ToOAuthProviderResponse(result), nil

}

func NewEmailProviderService(
	getEmailProviderUseCase usecase.IGetEmailProviderUseCase,
	createEmailProviderUseCase usecase.ICreateEmailProviderUseCase,
) IEmailProviderService {
	return &EmailProviderService{
		getEmailProviderUseCase:    getEmailProviderUseCase,
		createEmailProviderUseCase: createEmailProviderUseCase,
	}
}
