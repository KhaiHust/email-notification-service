package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
	"github.com/samber/lo"
)

type IGetEmailTemplateUseCase interface {
	GetTemplateByID(ctx context.Context, ID int64) (*entity.EmailTemplateEntity, error)
	GetAllTemplates(ctx context.Context, filter *request.GetListEmailTemplateFilter) ([]*entity.EmailTemplateEntity, error)
	GetAllTemplatesWithMetrics(ctx context.Context, filter *request.GetListEmailTemplateFilter) ([]*entity.EmailTemplateEntity, int64, error)
	GetTemplateByIDAndWorkspaceID(ctx context.Context, ID int64, workspaceID int64) (*entity.EmailTemplateEntity, error)
	GetTemplatesByIDs(ctx context.Context, IDs []int64) ([]*entity.EmailTemplateEntity, error)
}
type GetEmailTemplateUseCase struct {
	emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort
	getEmailRequestUsecase      IGetEmailRequestUsecase
	encryptUseCase              IEncryptUseCase
}

func (e GetEmailTemplateUseCase) GetTemplatesByIDs(ctx context.Context, IDs []int64) ([]*entity.EmailTemplateEntity, error) {
	return e.emailTemplateRepositoryPort.GetTemplatesByIDs(ctx, IDs)
}

func (e GetEmailTemplateUseCase) GetTemplateByIDAndWorkspaceID(ctx context.Context, ID int64, workspaceID int64) (*entity.EmailTemplateEntity, error) {
	result, err := e.emailTemplateRepositoryPort.GetTemplateByIDAndWorkspaceID(ctx, ID, workspaceID)
	if err != nil {
		log.Error(ctx, "Get email template by ID and workspace ID error: %v", err)
		return nil, err
	}
	version := result.Version
	version, err = e.encryptUseCase.EncryptVersionTemplate(ctx, version)
	if err != nil {
		log.Error(ctx, "Encrypt email template version error: %v", err)
		return nil, err
	}
	result.Version = version
	return result, nil
}

func (e GetEmailTemplateUseCase) GetAllTemplatesWithMetrics(ctx context.Context, filter *request.GetListEmailTemplateFilter) ([]*entity.EmailTemplateEntity, int64, error) {
	emailTemplates, err := e.emailTemplateRepositoryPort.GetAllTemplates(ctx, filter)
	if err != nil || len(emailTemplates) == 0 {
		return nil, 0, err
	}
	countAllTemplates, err := e.emailTemplateRepositoryPort.CountAllTemplates(ctx, filter)
	if err != nil {
		log.Error(ctx, "Count email templates error: %v", err)
		return nil, 0, err
	}
	if len(emailTemplates) == 0 {
		return make([]*entity.EmailTemplateEntity, 0), countAllTemplates, nil
	}
	templateIDs := lo.Map(emailTemplates, func(item *entity.EmailTemplateEntity, _ int) int64 {
		return item.ID
	})
	emailRequestFilter := filter.EmailRequestFilter
	emailRequestFilter.EmailTemplateIDs = templateIDs
	emailRequestStatusCounts, err := e.getEmailRequestUsecase.CountEmailRequestStatuses(ctx, emailRequestFilter)
	if err != nil {
		log.Error(ctx, "Get email request status counts error: %v", err)
		return nil, 0, err
	}
	countMap := make(map[int64]map[string]int64)
	for _, r := range emailRequestStatusCounts {
		if _, ok := countMap[r.EmailTemplateId]; !ok {
			countMap[r.EmailTemplateId] = make(map[string]int64)
		}
		countMap[r.EmailTemplateId][r.Status] = r.Total
	}
	emailTemplates = lo.Map(emailTemplates, func(item *entity.EmailTemplateEntity, _ int) *entity.EmailTemplateEntity {
		count, ok := countMap[item.ID]
		if ok {
			item.Metric = &dto.EmailTemplateMetric{
				TotalSent:   count[constant.EmailSendingStatusSent],
				TotalErrors: count[constant.EmailSendingStatusFailed],
			}
		}
		return item
	})

	return emailTemplates, countAllTemplates, nil
}

func (e GetEmailTemplateUseCase) GetAllTemplates(ctx context.Context, filter *request.GetListEmailTemplateFilter) ([]*entity.EmailTemplateEntity, error) {
	return e.emailTemplateRepositoryPort.GetAllTemplates(ctx, filter)
}

func (e GetEmailTemplateUseCase) GetTemplateByID(ctx context.Context, ID int64) (*entity.EmailTemplateEntity, error) {
	return e.emailTemplateRepositoryPort.GetTemplateByID(ctx, ID)
}

func NewGetEmailTemplateUseCase(
	emailTemplateRepositoryPort port.IEmailTemplateRepositoryPort,
	getEmailRequestUsecase IGetEmailRequestUsecase,
	encryptUseCase IEncryptUseCase,
) IGetEmailTemplateUseCase {
	return &GetEmailTemplateUseCase{
		emailTemplateRepositoryPort: emailTemplateRepositoryPort,
		getEmailRequestUsecase:      getEmailRequestUsecase,
		encryptUseCase:              encryptUseCase,
	}
}
