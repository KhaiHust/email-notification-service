package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
)

type IEmailTemplateService interface {
	CreateTemplate(ctx context.Context, userId int64, workspaceCode string, req *request.CreateEmailTemplateRequest) (*entity.EmailTemplateEntity, error)
}
type EmailTemplateService struct {
	createTemplateUseCase usecase.ICreateTemplateUseCase
	getWorkspaceUseCase   usecase.IGetWorkspaceUseCase
}

func (e EmailTemplateService) CreateTemplate(ctx context.Context, userId int64, workspaceCode string, req *request.CreateEmailTemplateRequest) (*entity.EmailTemplateEntity, error) {
	emailTemplateEntity := request.ToEmailTemplateEntity(req)
	workspace, err := e.getWorkspaceUseCase.GetWorkspaceByCode(ctx, userId, workspaceCode)
	if err != nil {
		return nil, err
	}
	emailTemplateEntity.CreatedBy = userId
	emailTemplateEntity.WorkspaceId = workspace.ID
	emailTemplateEntity.LastUpdatedBy = userId
	return e.createTemplateUseCase.CreateTemplate(ctx, emailTemplateEntity)
}

func NewEmailTemplateService(
	createTemplateUseCase usecase.ICreateTemplateUseCase,
	getWorkspaceUseCase usecase.IGetWorkspaceUseCase,
) IEmailTemplateService {
	return &EmailTemplateService{
		createTemplateUseCase: createTemplateUseCase,
		getWorkspaceUseCase:   getWorkspaceUseCase,
	}
}
