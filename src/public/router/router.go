package router

import (
	"fmt"
	"github.com/KhaiHust/email-notification-service/core/constant"
	"github.com/KhaiHust/email-notification-service/public/controller"
	"github.com/KhaiHust/email-notification-service/public/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib"
	"github.com/golibs-starter/golib/web/actuator"
	"go.uber.org/fx"
)

type RegisterRoutersIn struct {
	fx.In
	App                       *golib.App
	Engine                    *gin.Engine
	Actuator                  *actuator.Endpoint
	WorkspaceAccessMiddleware *middleware.WorkspaceAccessMiddleware
	*controller.BaseController
	EmailProviderController *controller.EmailProviderController
	UserController          *controller.UserController
	WorkspaceController     *controller.WorkspaceController
	EmailTemplateController *controller.EmailTemplateController
	EmailSendingController  *controller.EmailSendingController
	ApiKeyController        *controller.ApiKeyController
	EmailRequestController  *controller.EmailRequestController
	EmailTrackingController *controller.EmailTrackingController
	EmailLogController      *controller.EmailLogController
	AnalyticController      *controller.AnalyticController
}

func RegisterGinRouters(p RegisterRoutersIn) {
	group := p.Engine.Group(p.App.Path())
	group.GET("/actuator/health", gin.WrapF(p.Actuator.Health))
	group.GET("/actuator/info", gin.WrapF(p.Actuator.Info))
	v1Auth := group.Group("/v1/auth")
	{
		v1Auth.POST("/signup", p.UserController.SignUp)
		v1Auth.POST("/login", p.UserController.Login)
		v1Auth.POST("/refresh-token", p.UserController.GenerateTokenFromRefreshToken)
	}
	v1Integration := group.Group("/v1/integrations")
	{
		v1Integration.GET("/:emailProvider/oauth", p.EmailProviderController.GetOAuthUrl)
	}
	v1Workspace := group.Group("/v1/workspaces")
	{
		v1Workspace.POST("", p.WorkspaceController.CreateWorkspace)
		v1Workspace.GET("", p.WorkspaceController.GetWorkspaces)
		v1Workspace.POST("/:workspaceCode/send", p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailSendingController.SendEmailRequest)
	}
	v1EmailProvider := v1Workspace.Group("/:workspaceCode/providers")
	{
		v1EmailProvider.POST("/:emailProvider/oauth",
			p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailProviderController.CreateEmailProvider)
		v1EmailProvider.GET("/:emailProvider",
			p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailProviderController.GetEmailProvider)
	}
	v1Template := v1Workspace.Group("/:workspaceCode/templates")
	{
		v1Template.POST("",
			p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailTemplateController.CreateTemplate)
		v1Template.GET("", p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailTemplateController.GetAllEmailTemplate)
		v1Template.GET(fmt.Sprintf("/:%s", constant.ParamTemplateId),
			p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailTemplateController.GetTemplateDetail)
		v1Template.PATCH(fmt.Sprintf("/:%s", constant.ParamTemplateId),
			p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailTemplateController.UpdateTemplate)
		v1Template.GET("/:templateId/metrics",
			p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.AnalyticController.GetTemplateMetrics)
	}
	v1ApiKey := v1Workspace.Group("/:workspaceCode/api-keys")
	{
		v1ApiKey.POST("", p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.ApiKeyController.CreateNewApiKey)
		v1ApiKey.GET("", p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.ApiKeyController.GetListApiKey)
	}
	v1EmailRequest := v1Workspace.Group("/:workspaceCode/logs")
	{
		v1EmailRequest.GET("", p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailRequestController.GetListEmailRequest)
		v1EmailRequest.GET(fmt.Sprintf("/:%s", constant.ParamEmailRequestId),
			p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailLogController.GetLogs)
	}
	v1Tracking := group.Group("/v1/tracking")
	{
		v1Tracking.GET("/open", p.EmailTrackingController.OpenEmailTracking)
		v1Tracking.GET("/open.png", p.EmailTrackingController.OpenEmailTracking)
	}
	v1Analytic := v1Workspace.Group("/:workspaceCode/analytics")
	{
		v1Analytic.GET("/send-volumes", p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.AnalyticController.GetSendVolumes)
	}
}
