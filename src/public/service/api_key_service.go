package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/golibs-starter/golib/log"
)

type IApiKeyService interface {
	CreateNewApiKey(ctx context.Context, workspaceId int64, req *request.CreateApiKeyRequest) (*entity.ApiKeyEntity, error)
	GetAll(ctx context.Context, param *request.GetListApiKeyRequest) ([]*entity.ApiKeyEntity, error)
}
type ApiKeyService struct {
	createApiKeyUseCase usecase.ICreateApiKeyUseCase
	getApiKeyUseCase    usecase.IGetApiKeyUseCase
}

func (a ApiKeyService) GetAll(ctx context.Context, param *request.GetListApiKeyRequest) ([]*entity.ApiKeyEntity, error) {
	filter := request.NewGetListApiKeyFilter(param)
	apiKeys, err := a.getApiKeyUseCase.GetAll(ctx, filter)
	if err != nil {
		log.Error(ctx, "GetAllApiKey error: %v", err)
		return nil, err
	}
	return apiKeys, nil
}

func (a ApiKeyService) CreateNewApiKey(ctx context.Context, workspaceId int64, req *request.CreateApiKeyRequest) (*entity.ApiKeyEntity, error) {
	apiKeyEntity := &entity.ApiKeyEntity{
		Name:        req.Name,
		WorkspaceID: workspaceId,
		Environment: req.Environment,
	}
	apiKeyEntity, err := a.createApiKeyUseCase.CreateNewApiKey(ctx, apiKeyEntity)
	if err != nil {
		log.Error(ctx, "CreateNewApiKey error: %v", err)
		return nil, err
	}
	return apiKeyEntity, nil
}

func NewApiKeyService(
	createApiKeyUseCase usecase.ICreateApiKeyUseCase,
	getApiKeyUseCase usecase.IGetApiKeyUseCase,
) IApiKeyService {
	return &ApiKeyService{
		createApiKeyUseCase: createApiKeyUseCase,
		getApiKeyUseCase:    getApiKeyUseCase,
	}
}
