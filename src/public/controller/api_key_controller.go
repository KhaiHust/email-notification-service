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

type ApiKeyController struct {
	*BaseController
	ApiKeyService service.IApiKeyService
}

func (a *ApiKeyController) GetListApiKey(c *gin.Context) {
	workspaceID := a.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "Failed to get workspace ID from context")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	envs := utils.GetQueryStringArray(c.Request.URL.Query(), constant.QueryParamEnvironment)
	params := &request.GetListApiKeyRequest{
		Environment: envs,
		WorkspaceID: workspaceID,
	}
	result, err := a.ApiKeyService.GetAll(c, params)
	if err != nil {
		log.Error(c, "GetListApiKey error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, response.ToListApiKeyResponse(result))

}
func (a *ApiKeyController) CreateNewApiKey(c *gin.Context) {
	var req request.CreateApiKeyRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error(c, "Failed to bind the request's body to create new api key", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err := a.validator.Struct(&req); err != nil {
		log.Error(c, "Failed to validate the request's body to create new api key", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	workspaceID := a.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "Failed to get workspace ID from context")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	result, err := a.ApiKeyService.CreateNewApiKey(c, workspaceID, &req)
	if err != nil {
		log.Error(c, "CreateNewApiKey error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, response.ToApiKeyResponse(result))
}
func NewApiKeyController(base *BaseController, apiKeyService service.IApiKeyService) *ApiKeyController {
	return &ApiKeyController{
		BaseController: base,
		ApiKeyService:  apiKeyService,
	}
}
