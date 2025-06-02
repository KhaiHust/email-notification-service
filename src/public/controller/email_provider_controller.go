package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	coreRequestDto "github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type EmailProviderController struct {
	*BaseController
	emailProviderService service.IEmailProviderService
}

func (e EmailProviderController) GetAllEmailProviders(c *gin.Context) {
	workspaceID := e.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "workspaceID is 0")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}

	filter := e.buildQueryParamGetAllProviders(c)
	if filter == nil {
		log.Error(c, "filter is nil")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	filter.WorkspaceID = &workspaceID
	results, err := e.emailProviderService.GetAllEmailProviders(c, filter)
	if err != nil {
		log.Error(c, "GetAllEmailProviders error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, results)

}
func (e EmailProviderController) GetOAuthUrl(c *gin.Context) {
	provider := c.Param(constant.ParamEmailProvider)
	if provider == "" {
		log.Error(c, "provider is empty")
		apihelper.AbortErrorHandle(c, common.ErrEmailProviderParamNotFound)
		return
	}
	result, err := e.emailProviderService.GetOAuthUrl(c, provider)
	if err != nil {
		log.Error(c, "GetOAuthUrl error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, result)
}
func (e EmailProviderController) CreateEmailProvider(c *gin.Context) {
	workspaceCode := c.Param(constant.ParamWorkspaceCode)
	if workspaceCode == "" {
		log.Error(c, "workspaceCode is empty")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	provider := c.Param(constant.ParamEmailProvider)
	if provider == "" {
		log.Error(c, "provider is empty")
		apihelper.AbortErrorHandle(c, common.ErrEmailProviderParamNotFound)
		return
	}
	var req request.CreateEmailProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(c, "Failed to bind the request's body to create email provider")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err := e.validator.Struct(&req); err != nil {
		log.Error(c, "Failed to validate the request's body to create email provider")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	userID, err := e.GetUserIDFromContext(c)
	if err != nil {
		log.Error(c, "Failed to get user ID from context: %v", err)
		apihelper.AbortErrorHandle(c, common.ErrForbidden)
		return
	}

	_, err = e.emailProviderService.CreateEmailProvider(c, provider, userID, workspaceCode, &req)
	if err != nil {
		log.Error(c, "CreateEmailProvider error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, nil)
}
func (e EmailProviderController) GetEmailProvider(c *gin.Context) {
	workspaceCode := c.Param(constant.ParamWorkspaceCode)
	if workspaceCode == "" {
		log.Error(c, "workspaceCode is empty")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	provider := c.Param(constant.ParamEmailProvider)
	if provider == "" {
		log.Error(c, "provider is empty")
		apihelper.AbortErrorHandle(c, common.ErrEmailProviderParamNotFound)
		return
	}
	result, err := e.emailProviderService.GetEmailProviderByWorkspaceCodeAndProvider(c, workspaceCode, provider)
	if err != nil {
		log.Error(c, "GetEmailProvider error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, result)
}
func NewEmailProviderController(base *BaseController, emailProviderService service.IEmailProviderService) *EmailProviderController {
	return &EmailProviderController{
		BaseController:       base,
		emailProviderService: emailProviderService,
	}
}
func (e EmailProviderController) buildQueryParamGetAllProviders(c *gin.Context) *coreRequestDto.GetEmailProviderRequestFilter {
	query := c.Request.URL.Query()
	filter := &coreRequestDto.GetEmailProviderRequestFilter{}
	if provider := query.Get(constant.QueryParamProvider); provider != "" {
		filter.Provider = &provider
	}
	if env := query.Get(constant.QueryParamEnvironment); env != "" {
		filter.Environment = &env
	}
	return filter
}
