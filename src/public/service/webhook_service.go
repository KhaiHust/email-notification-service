package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
)

type IWebhookService interface {
	CreateWebhook(ctx context.Context, workspaceID int64, req *request.CreateWebhookRequest) (*response.WebhookResponse, error)
}
type WebhookService struct {
	createWebhookUseCase usecase.ICreateWebhookUseCase
}

func (w WebhookService) CreateWebhook(ctx context.Context, workspaceID int64, req *request.CreateWebhookRequest) (*response.WebhookResponse, error) {
	webhookDto := request.ToCreateWebhookRequestDto(req)
	webhookDto.WorkspaceID = workspaceID
	webhookEntity, err := w.createWebhookUseCase.CreateNewWebhook(ctx, webhookDto)
	if err != nil {
		return nil, err
	}

	return response.ToWebhookResponse(webhookEntity), nil
}

func NewWebhookService(createWebhookUseCase usecase.ICreateWebhookUseCase) IWebhookService {
	return &WebhookService{
		createWebhookUseCase: createWebhookUseCase,
	}
}
