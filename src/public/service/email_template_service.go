package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
)

type IEmailTemplateService interface {
	CreateTemplate(ctx context.Context, userId int64, workspaceCode string, req *request.CreateEmailTemplateRequest) (*entity.EmailTemplateEntity, error)
	GetAllEmailTemplateWithMetrics(ctx context.Context, workspaceId int64, filter *request.GetEmailTemplateParams) ([]*entity.EmailTemplateEntity, *apihelper.PagingMetadata, error)
	GetTemplateDetail(ctx context.Context, workspaceID, templateID int64) (*entity.EmailTemplateEntity, error)
	UpdateEmailTemplate(ctx context.Context, userId int64, workspaceID, templateID int64, req *request.CreateEmailTemplateRequest) (*entity.EmailTemplateEntity, error)
	DeleteEmailTemplate(ctx context.Context, userId, workspaceID, templateID int64) error
}
type EmailTemplateService struct {
	createTemplateUseCase      usecase.ICreateTemplateUseCase
	getEmailTemplateUseCase    usecase.IGetEmailTemplateUseCase
	getWorkspaceUseCase        usecase.IGetWorkspaceUseCase
	updateEmailTemplateUseCase usecase.IUpdateEmailTemplateUseCase
	deleteTemplateUseCase      usecase.IDeleteTemplateUseCase
}

func (e EmailTemplateService) DeleteEmailTemplate(ctx context.Context, userId, workspaceID, templateID int64) error {
	err := e.deleteTemplateUseCase.DeleteTemplate(ctx, workspaceID, userId, templateID)
	if err != nil {
		return err
	}
	return nil
}

func (e EmailTemplateService) UpdateEmailTemplate(ctx context.Context, userId int64, workspaceID, templateID int64, req *request.CreateEmailTemplateRequest) (*entity.EmailTemplateEntity, error) {
	templateEntity := request.ToEmailTemplateEntity(req)
	templateEntity.WorkspaceId = workspaceID
	templateEntity.ID = templateID
	templateEntity.LastUpdatedBy = userId
	templateEntity, err := e.updateEmailTemplateUseCase.UpdateEmailTemplate(ctx, templateEntity)
	if err != nil {
		return nil, err
	}
	return templateEntity, nil
}

func (e EmailTemplateService) GetTemplateDetail(ctx context.Context, workspaceID, templateID int64) (*entity.EmailTemplateEntity, error) {
	return e.getEmailTemplateUseCase.GetTemplateByIDAndWorkspaceID(ctx, templateID, workspaceID)
}

func (e EmailTemplateService) GetAllEmailTemplateWithMetrics(ctx context.Context, workspaceId int64, filter *request.GetEmailTemplateParams) ([]*entity.EmailTemplateEntity, *apihelper.PagingMetadata, error) {
	emailTemplateFilter := request.ToGetEmailTemplateFilter(filter)
	emailTemplateFilter.WorkspaceID = &workspaceId
	emailTemplates, total, err := e.getEmailTemplateUseCase.GetAllTemplatesWithMetrics(ctx, emailTemplateFilter)
	if err != nil {
		return nil, nil, err
	}
	getIDOfEmailTemplateFunc := func(emailTemplate *entity.EmailTemplateEntity) int64 {
		if emailTemplate == nil {
			return 0
		}
		return emailTemplate.ID
	}
	pagingMetadata := apihelper.BuildIDPaginatedResponse(emailTemplates, filter.Since, filter.Until, filter.Limit, &total, getIDOfEmailTemplateFunc, filter.SortOrder)
	return emailTemplates, &pagingMetadata, nil
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
	getEmailTemplateUseCase usecase.IGetEmailTemplateUseCase,
	getWorkspaceUseCase usecase.IGetWorkspaceUseCase,
	updateEmailTemplateUseCase usecase.IUpdateEmailTemplateUseCase,
	deleteTemplateUseCase usecase.IDeleteTemplateUseCase,
) IEmailTemplateService {
	return &EmailTemplateService{
		createTemplateUseCase:      createTemplateUseCase,
		getEmailTemplateUseCase:    getEmailTemplateUseCase,
		getWorkspaceUseCase:        getWorkspaceUseCase,
		updateEmailTemplateUseCase: updateEmailTemplateUseCase,
		deleteTemplateUseCase:      deleteTemplateUseCase,
	}
}
