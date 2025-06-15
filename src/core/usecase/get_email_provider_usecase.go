package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/KhaiHust/email-notification-service/core/port"
)

type IGetEmailProviderUseCase interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error)
	GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error)
	GetEmailProviderByIDAndWorkspaceID(ctx context.Context, providerID int64, workspaceID int64) (*entity.EmailProviderEntity, error)
	GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*entity.EmailProviderEntity, error)
	GetAllEmailProviders(ctx context.Context, filter *request.GetEmailProviderRequestFilter) ([]*entity.EmailProviderEntity, error)
	GetProvidersByIds(ctx context.Context, ids []int64) ([]*entity.EmailProviderEntity, error)
}
type GetEmailProviderUseCase struct {
	emailProviderPort           port.IEmailProviderPort
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
}

func (g GetEmailProviderUseCase) GetProvidersByIds(ctx context.Context, ids []int64) ([]*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetProvidersByIds(ctx, ids)
}

func (g GetEmailProviderUseCase) GetAllEmailProviders(ctx context.Context, filter *request.GetEmailProviderRequestFilter) ([]*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetAllEmailProviders(ctx, filter)
}

func (g GetEmailProviderUseCase) GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetEmailProviderByWorkspaceCodeAndProvider(ctx, workspaceCode, provider)
}

func (g GetEmailProviderUseCase) GetEmailProviderByIDAndWorkspaceID(ctx context.Context, providerID, workspaceID int64) (*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetEmailProviderByIDAndWorkspaceID(ctx, providerID, workspaceID)
}

func (g GetEmailProviderUseCase) GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetEmailProviderByID(ctx, ID)
}

func (g GetEmailProviderUseCase) GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error) {
	return g.emailProviderPort.GetOAuthUrl(ctx, provider)
}

func NewGetEmailProviderUseCase(
	emailProviderPort port.IEmailProviderPort,
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort,
) IGetEmailProviderUseCase {
	return &GetEmailProviderUseCase{
		emailProviderPort:           emailProviderPort,
		emailProviderRepositoryPort: emailProviderRepositoryPort,
	}
}
