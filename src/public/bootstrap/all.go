package bootstrap

import (
	"github.com/KhaiHust/email-notification-service/adapter/http/client"
	strategyAdapterImpl "github.com/KhaiHust/email-notification-service/adapter/http/strategy/impl"
	"github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/public/controller"
	"github.com/KhaiHust/email-notification-service/public/router"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/golibs-starter/golib"
	golibdata "github.com/golibs-starter/golib-data"
	golibgin "github.com/golibs-starter/golib-gin"
	golibsec "github.com/golibs-starter/golib-security"
	"go.uber.org/fx"
)

func All() fx.Option {
	return fx.Options(
		golib.AppOpt(),
		golib.PropertiesOpt(),
		golib.LoggingOpt(),
		golib.EventOpt(),
		golib.BuildInfoOpt(Version, CommitHash, BuildTime),
		golib.ActuatorEndpointOpt(),
		golib.HttpRequestLogOpt(),

		// Http security auto config and authentication filters
		//golibsec.HttpSecurityOpt(),

		golib.HttpClientOpt(),
		golibsec.SecuredHttpClientOpt(),

		// Provide datasource auto config
		golibdata.RedisOpt(),
		golibdata.DatasourceOpt(),

		// Provide all application properties
		golib.ProvideProps(properties.NewGmailProviderProperties),

		// Provide port's implements
		fx.Provide(client.NewGmailProviderAdapter),
		fx.Provide(strategyAdapterImpl.NewEmailProviderAdapter),
		fx.Provide(postgres.NewWorkspaceRepositoryAdapter),
		fx.Provide(postgres.NewEmailProviderRepositoryAdapter),
		fx.Provide(postgres.NewDatabaseTransactionAdapter),

		// Provide use cases
		fx.Provide(usecase.NewDatabaseTransactionUseCase),
		fx.Provide(usecase.NewGetEmailProviderUseCase),
		fx.Provide(usecase.NewCreateEmailProviderUseCase),
		fx.Provide(usecase.NewGetWorkspaceUseCase),

		// Provide services
		fx.Provide(service.NewEmailProviderService),

		fx.Provide(controller.NewEmailProviderController),

		golibgin.GinHttpServerOpt(),
		fx.Invoke(router.RegisterGinRouters),

		//Graceful shutdown
		golibgin.OnStopHttpServerOpt(),
	)
}
