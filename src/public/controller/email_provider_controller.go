package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type EmailProviderController struct {
	BaseController
	emailProviderService service.IEmailProviderService
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

	userId := int64(1)
	workspaceCode := c.Param(constant.ParamWorkspaceCode)
	result, err := e.emailProviderService.CreateEmailProvider(c, provider, userId, workspaceCode, req.Code)
	if err != nil {
		log.Error(c, "CreateEmailProvider error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, result)
}
func NewEmailProviderController(emailProviderService service.IEmailProviderService) *EmailProviderController {
	return &EmailProviderController{
		emailProviderService: emailProviderService,
	}
}
