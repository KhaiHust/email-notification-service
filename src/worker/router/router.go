package router

import (
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib"
	"github.com/golibs-starter/golib/web/actuator"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/fx"
)

type RegisterRoutersIn struct {
	fx.In
	App      *golib.App
	Engine   *gin.Engine
	Newrelic *newrelic.Application
	Actuator *actuator.Endpoint
}

func RegisterGinRouters(p RegisterRoutersIn) {
	p.Engine.Use(nrgin.Middleware(p.Newrelic))
	group := p.Engine.Group(p.App.Path())
	group.GET("/actuator/health", gin.WrapF(p.Actuator.Health))
	group.GET("/actuator/info", gin.WrapF(p.Actuator.Info))
}
