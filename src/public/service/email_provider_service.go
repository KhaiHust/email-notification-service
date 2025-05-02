package service

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
)

type IEmailProviderService interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthProviderResponse, error)
	CreateEmailProvider(ctx context.Context, provider string, userId int64, workspaceCode string, req *request.CreateEmailProviderRequest) (*entity.EmailProviderEntity, error)
	GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*response.EmailProviderResponse, error)
}
type EmailProviderService struct {
	getEmailProviderUseCase    usecase.IGetEmailProviderUseCase
	createEmailProviderUseCase usecase.ICreateEmailProviderUseCase
}

func (e EmailProviderService) GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*response.EmailProviderResponse, error) {
	result, err := e.getEmailProviderUseCase.GetEmailProviderByWorkspaceCodeAndProvider(ctx, workspaceCode, provider)
	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		return nil, err
	}
	return []*response.EmailProviderResponse{response.ToEmailProviderResponse(result)}, nil
}

func (e EmailProviderService) CreateEmailProvider(ctx context.Context, provider string, userId int64, workspaceCode string, req *request.CreateEmailProviderRequest) (*entity.EmailProviderEntity, error) {
	result, err := e.createEmailProviderUseCase.CreateEmailProvider(ctx, userId, workspaceCode, provider, request.ToEmailProviderDto(req))
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
