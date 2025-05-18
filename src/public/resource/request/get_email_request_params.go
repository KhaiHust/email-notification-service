package request

import "github.com/KhaiHust/email-notification-service/core/entity/dto/request"

type GetListEmailRequestParams struct {
	RequestID        *string
	Statuses         []string
	Limit            *int64
	Since            *int64
	Until            *int64
	SortOrder        string
	CreatedAtFrom    *int64
	CreatedAtTo      *int64
	UpdatedAtFrom    *int64
	UpdatedAtTo      *int64
	Providers        []string
	Email            *string
	EmailTemplateIDs []int64
}

func ToGetEmailRequestFilter(param *GetListEmailRequestParams) *request.EmailRequestFilter {
	if param == nil {
		return nil
	}
	return &request.EmailRequestFilter{
		BaseFilter: &request.BaseFilter{
			SortOrder:     param.SortOrder,
			CreatedAtFrom: param.CreatedAtFrom,
			CreatedAtTo:   param.CreatedAtTo,
			UpdatedAtFrom: param.UpdatedAtFrom,
			UpdatedAtTo:   param.UpdatedAtTo,
			Limit:         param.Limit,
			Since:         param.Since,
			Until:         param.Until,
		},
		RequestID:        param.RequestID,
		Statuses:         param.Statuses,
		EmailTemplateIDs: param.EmailTemplateIDs,
	}
}
