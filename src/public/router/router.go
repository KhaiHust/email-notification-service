package router

import (
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
}

func RegisterGinRouters(p RegisterRoutersIn) {
	group := p.Engine.Group(p.App.Path())
	group.GET("/actuator/health", gin.WrapF(p.Actuator.Health))
	group.GET("/actuator/info", gin.WrapF(p.Actuator.Info))
	v1Auth := group.Group("/v1/auth")
	{
		v1Auth.POST("/signup", p.UserController.SignUp)
		v1Auth.POST("/login", p.UserController.Login)
	}
	v1Integration := group.Group("/v1/integrations")
	{
		v1Integration.GET("/:emailProvider/oauth", p.EmailProviderController.GetOAuthUrl)
	}
	V1Workspace := group.Group("/v1/workspaces")
	{
		V1Workspace.POST("", p.WorkspaceController.CreateWorkspace)
		V1Workspace.POST("/:workspaceCode/providers/:emailProvider/oauth",
			p.WorkspaceAccessMiddleware.WorkspaceAccessMiddlewareHandle(),
			p.EmailProviderController.CreateEmailProvider)
	}
}
