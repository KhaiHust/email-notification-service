package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
)

type AnalyticController struct {
	*BaseController
	analyticService service.IAnalyticService
}

func (a AnalyticController) GetSendVolumes(c *gin.Context) {
	workspaceID := a.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	filter, err := a.buildSendVolumeFilter(c)
	if err != nil {
		apihelper.AbortErrorHandle(c, err)
		return
	}
	filter.WorkspaceId = workspaceID

	sendVolumes, err := a.analyticService.GetSendVolumes(c, filter)
	if err != nil {
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, sendVolumes)
}

func (a AnalyticController) buildSendVolumeFilter(c *gin.Context) (*request.SendVolumeFilter, error) {
	filter := &request.SendVolumeFilter{}
	values := c.Request.URL.Query()
	var err error
	if filter.StartDate, err = utils.GetQueryInt64Pointer(values, constant.QueryParamStartDate); err != nil {
		return nil, err
	}
	if filter.EndDate, err = utils.GetQueryInt64Pointer(values, constant.QueryParamEndDate); err != nil {
		return nil, err
	}
	return filter, nil
}
func NewAnalyticController(analyticService service.IAnalyticService, base *BaseController) *AnalyticController {
	return &AnalyticController{
		BaseController:  base,
		analyticService: analyticService,
	}
}
