package usecase

import (
	"context"
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/response"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IGetEmailProviderUseCase interface {
	GetOAuthUrl(ctx context.Context, provider string) (*response.OAuthUrlResponseDto, error)
	GetEmailProviderByID(ctx context.Context, ID int64) (*entity.EmailProviderEntity, error)
	GetEmailProviderByIDAndWorkspaceID(ctx context.Context, providerID int64, workspaceID int64) (*entity.EmailProviderEntity, error)
	GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*entity.EmailProviderEntity, error)
	GetAllEmailProviders(ctx context.Context, filter *request.GetEmailProviderRequestFilter) ([]*entity.EmailProviderEntity, error)
	GetProvidersByIds(ctx context.Context, ids []int64) ([]*entity.EmailProviderEntity, error)
	GetProviderByProviderAndWorkspaceIDAndEnvironment(ctx context.Context, provider string, workspaceID int64, environment string) (*entity.EmailProviderEntity, error)
}
type GetEmailProviderUseCase struct {
	emailProviderPort           port.IEmailProviderPort
	emailProviderRepositoryPort port.IEmailProviderRepositoryPort
	encryptUseCase              IEncryptUseCase
}

func (g GetEmailProviderUseCase) GetProviderByProviderAndWorkspaceIDAndEnvironment(ctx context.Context, provider string, workspaceID int64, environment string) (*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetProviderByProviderAndWorkspaceIDAndEnvironment(ctx, provider, workspaceID, environment)
}

func (g GetEmailProviderUseCase) GetProvidersByIds(ctx context.Context, ids []int64) ([]*entity.EmailProviderEntity, error) {
	providers, err := g.emailProviderRepositoryPort.GetProvidersByIds(ctx, ids)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("GetProvidersByIds error: %v", err))
		return nil, err
	}
	for _, provider := range providers {
		if provider.OAuthToken != "" {
			decryptedToken, err := g.encryptUseCase.DecryptProviderToken(ctx, provider.OAuthToken)
			if err != nil {
				log.Error(ctx, fmt.Sprintf("DecryptProviderToken error: %v", err))
				return nil, err
			}
			provider.OAuthToken = decryptedToken
		}
		if provider.OAuthRefreshToken != "" {
			decryptedRefreshToken, err := g.encryptUseCase.DecryptProviderToken(ctx, provider.OAuthRefreshToken)
			if err != nil {
				log.Error(ctx, fmt.Sprintf("DecryptProviderToken error: %v", err))
				return nil, err
			}
			provider.OAuthRefreshToken = decryptedRefreshToken
		}
	}
	return providers, nil
}

func (g GetEmailProviderUseCase) GetAllEmailProviders(ctx context.Context, filter *request.GetEmailProviderRequestFilter) ([]*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetAllEmailProviders(ctx, filter)
}

func (g GetEmailProviderUseCase) GetEmailProviderByWorkspaceCodeAndProvider(ctx context.Context, workspaceCode string, provider string) ([]*entity.EmailProviderEntity, error) {
	return g.emailProviderRepositoryPort.GetEmailProviderByWorkspaceCodeAndProvider(ctx, workspaceCode, provider)
}

func (g GetEmailProviderUseCase) GetEmailProviderByIDAndWorkspaceID(ctx context.Context, providerID, workspaceID int64) (*entity.EmailProviderEntity, error) {

	provider, err := g.emailProviderRepositoryPort.GetEmailProviderByIDAndWorkspaceID(ctx, providerID, workspaceID)
	if err != nil {
		log.Error(ctx, "GetEmailProviderByIDAndWorkspaceID error: %v", err)
		return nil, err
	}
	//decrypt OAuth tokens if they exist
	if provider.OAuthToken != "" {
		decryptedToken, err := g.encryptUseCase.DecryptProviderToken(ctx, provider.OAuthToken)
		if err != nil {
			log.Error(ctx, "DecryptProviderToken error: %v", err)
			return nil, err
		}
		provider.OAuthToken = decryptedToken
	}
	if provider.OAuthRefreshToken != "" {
		decryptedRefreshToken, err := g.encryptUseCase.DecryptProviderToken(ctx, provider.OAuthRefreshToken)
		if err != nil {
			log.Error(ctx, "DecryptProviderToken error: %v", err)
			return nil, err
		}
		provider.OAuthRefreshToken = decryptedRefreshToken
	}
	return provider, nil
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
	encryptUseCase IEncryptUseCase,
) IGetEmailProviderUseCase {
	return &GetEmailProviderUseCase{
		emailProviderPort:           emailProviderPort,
		emailProviderRepositoryPort: emailProviderRepositoryPort,
		encryptUseCase:              encryptUseCase,
	}
}
