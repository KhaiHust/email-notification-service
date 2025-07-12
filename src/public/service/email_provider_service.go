package service

import (
	"context"
	"errors"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/entity"
	coreRequest "github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
)

type IEmailProviderService interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthProviderResponse, error)
	CreateEmailProvider(ctx context.Context, provider string, userId int64, workspaceCode string, req *request.CreateEmailProviderRequest) (*entity.EmailProviderEntity, error)
	GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*response.EmailProviderResponse, error)
	GetAllEmailProviders(ctx context.Context, filter *coreRequest.GetEmailProviderRequestFilter) ([]*response.EmailProviderResponse, error)
	UpdateEmailProviderRequest(ctx context.Context, workspaceID, providerID int64, req *request.UpdateEmailProviderRequest) (*entity.EmailProviderEntity, error)
	DeactivateEmailProvider(ctx context.Context, workspaceID, providerID int64) (*entity.EmailProviderEntity, error)
}
type EmailProviderService struct {
	getEmailProviderUseCase    usecase.IGetEmailProviderUseCase
	createEmailProviderUseCase usecase.ICreateEmailProviderUseCase
	updateEmailProviderUseCase usecase.IUpdateEmailProviderUseCase
}

func (e EmailProviderService) DeactivateEmailProvider(ctx context.Context, workspaceID, providerID int64) (*entity.EmailProviderEntity, error) {
	result, err := e.updateEmailProviderUseCase.DeactivateEmailProvider(ctx, workspaceID, providerID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e EmailProviderService) UpdateEmailProviderRequest(ctx context.Context, workspaceID, providerID int64, req *request.UpdateEmailProviderRequest) (*entity.EmailProviderEntity, error) {
	result, err := e.updateEmailProviderUseCase.UpdateInfoProvider(ctx, workspaceID, providerID, request.ToUpdateEmailProviderDto(req))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e EmailProviderService) GetAllEmailProviders(ctx context.Context, filter *coreRequest.GetEmailProviderRequestFilter) ([]*response.EmailProviderResponse, error) {
	results, err := e.getEmailProviderUseCase.GetAllEmailProviders(ctx, filter)
	if err != nil {
		return nil, err
	}
	return response.ToEmailProviderResponseList(results), nil
}

func (e EmailProviderService) GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*response.EmailProviderResponse, error) {
	result, err := e.getEmailProviderUseCase.GetEmailProviderByWorkspaceCodeAndProvider(ctx, workspaceCode, provider)
	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		return nil, err
	}
	return response.ToEmailProviderResponseList(result), nil
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
	updateEmailProviderUseCase usecase.IUpdateEmailProviderUseCase,
) IEmailProviderService {
	return &EmailProviderService{
		getEmailProviderUseCase:    getEmailProviderUseCase,
		createEmailProviderUseCase: createEmailProviderUseCase,
		updateEmailProviderUseCase: updateEmailProviderUseCase,
	}
}
