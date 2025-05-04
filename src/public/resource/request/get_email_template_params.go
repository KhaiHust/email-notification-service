package request

import "github.com/KhaiHust/email-notification-service/core/entity/dto/request"

type GetEmailTemplateParams struct {
	WorkspaceID     *int64
	Name            *string
	Limit           *int64
	Since           *int64
	Until           *int64
	SortOrder       string
	CreatedAtFrom   *int64
	CreatedAtTo     *int64
	UpdatedAtFrom   *int64
	UpdatedAtTo     *int64
	ErStatuses      []string
	ErCreatedAtFrom *int64
	ErCreatedAtTo   *int64
	ErSentAtFrom    *int64
	ErSentAtTo      *int64
}

func ToGetEmailTemplateFilter(req *GetEmailTemplateParams) *request.GetListEmailTemplateFilter {
	return &request.GetListEmailTemplateFilter{
		WorkspaceID:        req.WorkspaceID,
		Name:               req.Name,
		Limit:              req.Limit,
		Since:              req.Since,
		Until:              req.Until,
		CreatedAtFrom:      req.CreatedAtFrom,
		CreatedAtTo:        req.CreatedAtTo,
		UpdatedAtFrom:      req.UpdatedAtFrom,
		UpdatedAtTo:        req.UpdatedAtTo,
		SortOrder:          req.SortOrder,
		EmailRequestFilter: &request.EmailRequestFilter{Statuses: req.ErStatuses, CreatedAtFrom: req.ErCreatedAtFrom, CreatedAtTo: req.ErCreatedAtTo, SentAtFrom: req.ErSentAtFrom, SentAtTo: req.ErSentAtTo},
	}
}
