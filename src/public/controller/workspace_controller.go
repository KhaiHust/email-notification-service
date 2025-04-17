package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/helper"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type WorkspaceController struct {
	BaseController
	workspaceService service.IWorkspaceService
}

func (w *WorkspaceController) CreateWorkspace(c *gin.Context) {
	var req request.CreateWorkspaceRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error(c, "Failed to bind request: %v", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err := w.validator.Struct(&req); err != nil {
		log.Error(c, "Failed to validate request: %v", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	userID, err := w.GetUserIDFromContext(c)
	if err != nil {
		log.Error(c, "Failed to get user ID from context: %v", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	workspace, err := w.workspaceService.CreateNewWorkspace(c, userID, &req)
	if err != nil {
		log.Error(c, "Failed to create workspace: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, workspace)

}
func NewWorkspaceController(workspaceService service.IWorkspaceService, validator *helper.CustomValidate) *WorkspaceController {
	return &WorkspaceController{
		BaseController:   *NewBaseController(validator),
		workspaceService: workspaceService,
	}
}
