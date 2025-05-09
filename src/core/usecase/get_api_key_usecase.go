package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
)

type IGetApiKeyUseCase interface {
	GetAll(ctx context.Context, filter *request.GetApiKeyRequestFilter) ([]*entity.ApiKeyEntity, error)
}
type GetApiKeyUseCase struct {
	apiKeyRepositoryPort port.IApiKeyRepositoryPort
}

func (g GetApiKeyUseCase) GetAll(ctx context.Context, filter *request.GetApiKeyRequestFilter) ([]*entity.ApiKeyEntity, error) {
	return g.apiKeyRepositoryPort.GetAllApiKeys(ctx, filter)
}

func NewGetApiKeyUseCase(apiKeyRepositoryPort port.IApiKeyRepositoryPort) IGetApiKeyUseCase {
	return &GetApiKeyUseCase{
		apiKeyRepositoryPort: apiKeyRepositoryPort,
	}
}
