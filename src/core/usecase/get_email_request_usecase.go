package usecase

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/port"
	"github.com/golibs-starter/golib/log"
)

type IGetEmailRequestUsecase interface {
	GetEmailRequestByID(ctx context.Context, emailRequestID int64) (*entity.EmailRequestEntity, error)
	CountEmailRequestStatuses(ctx context.Context, filter *request.EmailRequestFilter) ([]*entity.EmailRequestStatusCountEntity, error)
	GetAllEmailRequest(ctx context.Context, filter *request.EmailRequestFilter) ([]*entity.EmailRequestEntity, int64, error)
	CountAllEmailRequest(ctx context.Context, filter *request.EmailRequestFilter) (int64, error)
}
type GetEmailRequestUsecase struct {
	emailRequestRepositoryPort port.IEmailRequestRepositoryPort
}

func (g GetEmailRequestUsecase) CountAllEmailRequest(ctx context.Context, filter *request.EmailRequestFilter) (int64, error) {
	emailRequestCount, err := g.emailRequestRepositoryPort.CountAllEmailRequest(ctx, filter)
	if err != nil {
		return 0, err
	}
	return emailRequestCount, nil
}

func (g GetEmailRequestUsecase) GetAllEmailRequest(ctx context.Context, filter *request.EmailRequestFilter) ([]*entity.EmailRequestEntity, int64, error) {
	emailRequestEntities, err := g.emailRequestRepositoryPort.GetAllEmailRequest(ctx, filter)
	if err != nil {
		log.Error(ctx, "Error when get all email request", err)
		return nil, 0, err
	}
	total, err := g.CountAllEmailRequest(ctx, filter)
	if err != nil {
		log.Error(ctx, "Error when count all email request", err)
		return nil, 0, err
	}
	return emailRequestEntities, total, nil
}

func (g GetEmailRequestUsecase) CountEmailRequestStatuses(ctx context.Context, filter *request.EmailRequestFilter) ([]*entity.EmailRequestStatusCountEntity, error) {
	emailRequestStatusCounts, err := g.emailRequestRepositoryPort.CountEmailRequestStatuses(ctx, filter)
	if err != nil {
		return nil, err
	}
	return emailRequestStatusCounts, nil
}

func (g GetEmailRequestUsecase) GetEmailRequestByID(ctx context.Context, emailRequestID int64) (*entity.EmailRequestEntity, error) {
	return g.emailRequestRepositoryPort.GetEmailRequestByID(ctx, emailRequestID)
}

func NewGetEmailRequestUsecase(emailRequestRepositoryPort port.IEmailRequestRepositoryPort) IGetEmailRequestUsecase {
	return &GetEmailRequestUsecase{
		emailRequestRepositoryPort: emailRequestRepositoryPort,
	}
}
