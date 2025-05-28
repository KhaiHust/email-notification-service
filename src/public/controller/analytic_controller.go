package controller

import (
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/core/entity/dto/request"
	"github.com/KhaiHust/email-notification-service/core/utils"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
	"strconv"
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
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
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
func (a AnalyticController) GetTemplateMetrics(c *gin.Context) {
	workspaceID := a.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	templateID, err := strconv.ParseInt(c.Param(constant.ParamTemplateId), 10, 64)
	if err != nil {
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	filter, err := a.buildTemplateMetricFilter(c)
	if err != nil {
		log.Error(c, "Error when build template metric filter", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	filter.WorkspaceID = workspaceID
	filter.TemplateID = templateID
	templateMetrics, err := a.analyticService.GetTemplateMetrics(c, filter)
	if err != nil {
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, response.ToTemplateMertricResponse(templateMetrics))
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

func (a AnalyticController) buildTemplateMetricFilter(c *gin.Context) (*request.TemplateMetricFilter, error) {
	filter := &request.TemplateMetricFilter{}
	values := c.Request.URL.Query()
	isChart, err := strconv.ParseBool(values.Get(constant.QueryParamIsChart))
	if err != nil {
		return nil, fmt.Errorf("isChart param is invalid")
	}
	var internal *string
	if isChart {
		internal = utils.GetQueryStringPointer(values, constant.QueryParamInterval)
		mapInterval := map[string]bool{
			constant.IntervalDay:   true,
			constant.IntervalMonth: true,
			constant.IntervalWeek:  true,
		}
		if internal == nil || !mapInterval[*internal] {
			return nil, fmt.Errorf("invalid interval param")
		}
		duration, err := strconv.ParseInt(values.Get(constant.QueryParamDuration), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("duration param is invalid")
		}
		if duration <= 0 {
			return nil, fmt.Errorf("duration param must be greater than 0")
		}
		filter.Duration = int(duration)
		filter.Interval = *internal
		filter.IsChart = isChart
	} else {
		startDate, err := utils.GetQueryInt64Pointer(values, constant.QueryParamStartDate)
		if startDate == nil || err != nil {
			return nil, fmt.Errorf("start date is required")
		}
		endDate, err := utils.GetQueryInt64Pointer(values, constant.QueryParamEndDate)
		if endDate == nil || err != nil {
			return nil, fmt.Errorf("end date is required")
		}
		if *startDate > *endDate {
			return nil, fmt.Errorf("start date must be less than end date")
		}
		filter.StartDate = startDate
		filter.EndDate = endDate
	}

	return filter, nil

}
func NewAnalyticController(analyticService service.IAnalyticService, base *BaseController) *AnalyticController {
	return &AnalyticController{
		BaseController:  base,
		analyticService: analyticService,
	}
}
