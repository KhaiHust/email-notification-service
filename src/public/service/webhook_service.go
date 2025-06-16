package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
)

type IWebhookService interface {
	CreateWebhook(ctx context.Context, workspaceID int64, req *request.CreateWebhookRequest) (*response.WebhookResponse, error)
	UpdateWebhook(ctx context.Context, workspaceID, webID int64, req *request.UpdateWebhookRequest) (*response.WebhookResponse, error)
	GetAllWebhooksByWorkspaceID(ctx context.Context, workspaceID int64) ([]*response.WebhookResponse, error)
	DeleteWebhook(ctx context.Context, workspaceID, webID int64) error
	GetWebhookDetail(ctx context.Context, workspaceID, webID int64) (*response.WebhookResponse, error)
	TestWebhook(ctx context.Context, workspaceID, webID int64) error
}
type WebhookService struct {
	createWebhookUseCase usecase.ICreateWebhookUseCase
	webhookUsecase       usecase.IWebhookUsecase
}

func (w WebhookService) TestWebhook(ctx context.Context, workspaceID, webID int64) error {
	err := w.webhookUsecase.TestWebhook(ctx, workspaceID, webID)
	if err != nil {
		return err
	}
	return nil
}

func (w WebhookService) GetWebhookDetail(ctx context.Context, workspaceID, webID int64) (*response.WebhookResponse, error) {
	webhookEntity, err := w.webhookUsecase.GetWebhookDetail(ctx, workspaceID, webID)
	if err != nil {
		return nil, err
	}
	return response.ToWebhookResponse(webhookEntity), nil
}

func (w WebhookService) DeleteWebhook(ctx context.Context, workspaceID, webID int64) error {
	err := w.webhookUsecase.DeleteWebhook(ctx, workspaceID, webID)
	if err != nil {
		return err
	}
	return nil
}

func (w WebhookService) GetAllWebhooksByWorkspaceID(ctx context.Context, workspaceID int64) ([]*response.WebhookResponse, error) {
	webhooks, err := w.webhookUsecase.GetAllWebhookByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	webhookResponses := make([]*response.WebhookResponse, len(webhooks))
	for i, webhook := range webhooks {
		webhookResponses[i] = response.ToWebhookResponse(webhook)
	}
	return webhookResponses, nil
}

func (w WebhookService) UpdateWebhook(ctx context.Context, workspaceID, webID int64, req *request.UpdateWebhookRequest) (*response.WebhookResponse, error) {
	webhookDto := request.ToUpdateWebhookRequestDto(req)
	webhookEntity, err := w.webhookUsecase.UpdateWebhook(ctx, workspaceID, webID, webhookDto)
	if err != nil {
		return nil, err
	}
	return response.ToWebhookResponse(webhookEntity), nil
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

func NewWebhookService(
	createWebhookUseCase usecase.ICreateWebhookUseCase,
	webhookUsecase usecase.IWebhookUsecase,
) IWebhookService {
	return &WebhookService{
		createWebhookUseCase: createWebhookUseCase,
		webhookUsecase:       webhookUsecase,
	}
}
