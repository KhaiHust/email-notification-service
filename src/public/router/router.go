package router

import (
	"github.com/KhaiHust/email-notification-service/public/controller"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib"
	"github.com/golibs-starter/golib/web/actuator"
	"go.uber.org/fx"
)

type RegisterRoutersIn struct {
	fx.In
	App                     *golib.App
	Engine                  *gin.Engine
	Actuator                *actuator.Endpoint
	EmailProviderController *controller.EmailProviderController
}

func RegisterGinRouters(p RegisterRoutersIn) {
	group := p.Engine.Group(p.App.Path())
	group.GET("/actuator/health", gin.WrapF(p.Actuator.Health))
	group.GET("/actuator/info", gin.WrapF(p.Actuator.Info))
	v1Integration := group.Group("/v1/integrations")
	{
		v1Integration.GET("/:emailProvider/oauth", p.EmailProviderController.GetOAuthUrl)
		v1Integration.POST("/:emailProvider/oauth", p.EmailProviderController.CreateEmailProvider)
	}
}
