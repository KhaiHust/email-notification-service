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
)

type EmailTemplateController struct {
	*BaseController
	emailTemplateService service.IEmailTemplateService
}

func (e *EmailTemplateController) CreateTemplate(c *gin.Context) {
	userId, err := e.GetUserIDFromContext(c)
	if err != nil {
		log.Error(c, "Error when get user id from context", err)
		apihelper.AbortErrorHandle(c, common.ErrForbidden)
		return
	}
	workspaceID := c.Param(constant.ParamWorkspaceCode)
	if workspaceID == "" {
		log.Error(c, "Error when get workspace id from context", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	var req request.CreateEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(c, "Error when bind json", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err := e.validator.Struct(&req); err != nil {
		log.Error(c, "Error when validate request", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	emailTemplate, err := e.emailTemplateService.CreateTemplate(c, userId, workspaceID, &req)
	if err != nil {
		log.Error(c, "Error when create email template", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, response.ToEmailTemplateResponse(emailTemplate))

}
func (e *EmailTemplateController) GetAllEmailTemplate(c *gin.Context) {
	workspaceId := e.GetWorkspaceIDFromContext(c)
	if workspaceId == 0 {
		log.Error(c, "Error when get workspace id from context", common.ErrBadRequest)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	queryParams, err := e.buildGetListTemplateQueryParams(c)
	if err != nil {
		log.Error(c, "Error when build query params", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	emailTemplates, paging, err := e.emailTemplateService.GetAllEmailTemplateWithMetrics(c, workspaceId, queryParams)
	if err != nil {
		log.Error(c, "Error when get all email template", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandleWithPaging(c, response.ToListEmailTemplateResponse(emailTemplates), paging)
}
func NewEmailTemplateController(
	emailTemplateService service.IEmailTemplateService,
	base *BaseController,
) *EmailTemplateController {
	return &EmailTemplateController{
		BaseController:       base,
		emailTemplateService: emailTemplateService,
	}
}
func (e *EmailTemplateController) buildGetListTemplateQueryParams(ctx *gin.Context) (*request.GetEmailTemplateParams, error) {
	values := ctx.Request.URL.Query()
	queryParams := &request.GetEmailTemplateParams{
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
	if queryParams.ErCreatedAtFrom, err = utils.GetQueryInt64Pointer(values, constant.QueryParamErCreatedAtFrom); err != nil {
		return nil, err
	}
	if queryParams.ErCreatedAtTo, err = utils.GetQueryInt64Pointer(values, constant.QueryParamErCreatedAtTo); err != nil {
		return nil, err
	}
	if queryParams.ErSentAtFrom, err = utils.GetQueryInt64Pointer(values, constant.QueryParamErSentAtFrom); err != nil {
		return nil, err
	}
	if queryParams.ErSentAtTo, err = utils.GetQueryInt64Pointer(values, constant.QueryParamErSentAtTo); err != nil {
		return nil, err
	}

	// String and array fields
	queryParams.Name = utils.GetQueryStringPointer(values, constant.QueryParamName)
	queryParams.ErStatuses = utils.GetQueryStringArray(values, constant.QueryParamErStatuses)

	// SortOrder with validation
	if sortOrder := values.Get(constant.QueryParamSortOrder); sortOrder != "" {
		if sortOrder != constant.ASC && sortOrder != constant.DESC {
			return nil, common.ErrBadRequest
		}
		queryParams.SortOrder = sortOrder
	}

	return queryParams, nil
}
