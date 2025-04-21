package bootstrap

import (
	"github.com/KhaiHust/email-notification-service/adapter/http/client"
	strategyAdapterImpl "github.com/KhaiHust/email-notification-service/adapter/http/strategy/impl"
	"github.com/KhaiHust/email-notification-service/adapter/properties"
	"github.com/KhaiHust/email-notification-service/adapter/repository/postgres"
	"github.com/KhaiHust/email-notification-service/adapter/service/thirdparty"
	coreProperties "github.com/KhaiHust/email-notification-service/core/properties"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/worker/handler"
	"github.com/golibs-starter/golib"
	golibdata "github.com/golibs-starter/golib-data"
	golibmsg "github.com/golibs-starter/golib-message-bus"
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

		golib.HttpClientOpt(),
		golibsec.SecuredHttpClientOpt(),

		golibmsg.KafkaCommonOpt(),
		golibmsg.KafkaConsumerOpt(),

		// Provide datasource auto config
		golibdata.RedisOpt(),
		golibdata.DatasourceOpt(),

		// Provide all application properties
		golib.ProvideProps(properties.NewGmailProviderProperties),
		golib.ProvideProps(coreProperties.NewAuthProperties),
		// Provide port's implements
		fx.Provide(client.NewGmailProviderAdapter),
		fx.Provide(strategyAdapterImpl.NewEmailProviderAdapter),
		fx.Provide(postgres.NewWorkspaceRepositoryAdapter),
		fx.Provide(postgres.NewEmailProviderRepositoryAdapter),
		fx.Provide(postgres.NewDatabaseTransactionAdapter),
		fx.Provide(postgres.NewUserRepositoryAdapter),
		fx.Provide(postgres.NewWorkspaceUserRepositoryAdapter),
		fx.Provide(postgres.NewEmailTemplateRepositoryAdapter),

		// Provide third-party services
		fx.Provide(thirdparty.NewRedisService),

		// Provide use cases
		fx.Provide(usecase.NewDatabaseTransactionUseCase),
		fx.Provide(usecase.NewGetEmailProviderUseCase),
		fx.Provide(usecase.NewCreateEmailProviderUseCase),
		fx.Provide(usecase.NewGetWorkspaceUseCase),
		fx.Provide(usecase.NewCreateUserUseCase),
		fx.Provide(usecase.NewGetUserUseCase),
		fx.Provide(usecase.NewHashPasswordUseCase),
		fx.Provide(usecase.NewLoginUsecase),
		fx.Provide(usecase.NewCreateWorkspaceUseCase),
		fx.Provide(usecase.NewValidateAccessWorkspaceUsecase),
		fx.Provide(usecase.NewCreateTemplateUseCase),
		fx.Provide(usecase.NewEmailSendingUsecase),
		fx.Provide(usecase.NewGetEmailTemplateUseCase),

		//provider handler
		golibmsg.ProvideConsumer(handler.NewEmailSendingRequestHandler),

		// Graceful shutdown.
		// OnStop hooks will run in reverse order.
		//golibgin.OnStopHttpServerOpt(),
		//golibmsg.OnStopProducerOpt(),
		golibmsg.OnStopConsumerOpt(),
	)
}
