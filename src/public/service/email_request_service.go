package service

import (
	"context"
	"github.com/KhaiHust/email-notification-service/core/entity"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
)

type IEmailRequestService interface {
	GetAllEmailRequest(ctx context.Context, workspaceID int64, params *request.GetListEmailRequestParams) ([]*entity.EmailRequestEntity, *apihelper.PagingMetadata, error)
}
type EmailRequestService struct {
	getEmailRequestUsecase usecase.IGetEmailRequestUsecase
}

func (e EmailRequestService) GetAllEmailRequest(ctx context.Context, workspaceID int64, params *request.GetListEmailRequestParams) ([]*entity.EmailRequestEntity, *apihelper.PagingMetadata, error) {
	filter := request.ToGetEmailRequestFilter(params)
	filter.WorkspaceIDs = []int64{workspaceID}
	emailRequestEntities, total, err := e.getEmailRequestUsecase.GetAllEmailRequest(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	getIDOfEmailRequestFunc := func(emailRequest *entity.EmailRequestEntity) int64 {
		if emailRequest == nil {
			return 0
		}
		return emailRequest.ID
	}
	pagingMetadata := apihelper.BuildIDPaginatedResponse(emailRequestEntities, params.Since, params.Until, params.Limit, &total, getIDOfEmailRequestFunc, params.SortOrder)
	return emailRequestEntities, &pagingMetadata, nil
}

func NewEmailRequestService(getEmailRequestUsecase usecase.IGetEmailRequestUsecase) IEmailRequestService {
	return &EmailRequestService{
		getEmailRequestUsecase: getEmailRequestUsecase,
	}
}
