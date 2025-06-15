package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
	"strconv"
)

type EmailRequestController struct {
	*BaseController
	emailRequestService service.IEmailRequestService
}

func (e *EmailRequestController) GetListEmailRequest(ctx *gin.Context) {
	workspaceID := e.GetWorkspaceIDFromContext(ctx)
	if workspaceID == 0 {
		log.Error(ctx, "Error when get workspace id from context", common.ErrBadRequest)
		apihelper.AbortErrorHandle(ctx, common.ErrBadRequest)
		return
	}
	params, err := e.buildParamsGetListEmailRequests(ctx)
	if err != nil {
		log.Error(ctx, "Error when build params", err)
		apihelper.AbortErrorHandle(ctx, err)
		return
	}
	emailRequestEntities, paging, err := e.emailRequestService.GetAllEmailRequest(ctx, workspaceID, params)
	if err != nil {
		log.Error(ctx, "Error when get email request", err)
		apihelper.AbortErrorHandle(ctx, err)
		return
	}
	apihelper.SuccessfulHandleWithPaging(ctx, response.ToListEmailRequestResponse(emailRequestEntities), paging)
}
func (e *EmailRequestController) buildParamsGetListEmailRequests(ctx *gin.Context) (*request.GetListEmailRequestParams, error) {
	values := ctx.Request.URL.Query()
	queryParams := &request.GetListEmailRequestParams{
		Limit:     utils.ToInt64Pointer(constant.DefaultLimit),
		SortOrder: constant.DefaultSortOrderAsc,
	}
	var err error
	// Numeric fields
	queryParams.Limit, err = utils.GetInt64PointerWithDefault(values, constant.QueryParamLimit, constant.DefaultLimit)
	if err != nil {
		return nil, err
	}
	queryParams.Since, err = utils.GetQueryInt64Pointer(values, constant.QueryParamSince)
	if err != nil {
		return nil, err
	}
	if queryParams.Until, err = utils.GetQueryInt64Pointer(values, constant.QueryParamUntil); err != nil {
		return nil, err
	}
	if queryParams.CreatedAtFrom, err = utils.GetQueryInt64Pointer(values, constant.QueryParamCreatedAtFrom); err != nil {
		return nil, err
	}
	if queryParams.CreatedAtTo, err = utils.GetQueryInt64Pointer(values, constant.QueryParamCreatedAtTo); err != nil {
		return nil, err
	}
	if queryParams.UpdatedAtFrom, err = utils.GetQueryInt64Pointer(values, constant.QueryParamUpdatedAtFrom); err != nil {
		return nil, err
	}
	if queryParams.UpdatedAtTo, err = utils.GetQueryInt64Pointer(values, constant.QueryParamUpdatedAtTo); err != nil {
		return nil, err
	}
	requestID := utils.GetQueryStringPointer(values, constant.QueryParamRequestID)
	if requestID != nil && *requestID != "" {
		queryParams.RequestID = requestID
	}

	queryParams.Statuses = utils.GetQueryStringArray(values, constant.QueryParamErStatuses)

	recipientEmail := utils.GetQueryStringPointer(values, constant.QueryRecipientEmail)
	if recipientEmail != nil && *recipientEmail != "" {
		queryParams.Email = recipientEmail
	}
	templateIDs := utils.GetQueryStringArray(values, constant.QueryParamEmailTemplateIDs)
	if len(templateIDs) > 0 {
		queryParams.EmailTemplateIDs = make([]int64, len(templateIDs))
		for i, idStr := range templateIDs {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return nil, common.ErrBadRequest
			}
			queryParams.EmailTemplateIDs[i] = id
		}
	}
	// SortOrder with validation
	if sortOrder := values.Get(constant.QueryParamSortOrder); sortOrder != "" {
		if sortOrder != constant.ASC && sortOrder != constant.DESC {
			return nil, common.ErrBadRequest
		}
		queryParams.SortOrder = sortOrder
	}
	return queryParams, nil
}
func NewEmailRequestController(
	base *BaseController,
	emailRequestService service.IEmailRequestService,
) *EmailRequestController {
	return &EmailRequestController{
		BaseController:      base,
		emailRequestService: emailRequestService,
	}
}
