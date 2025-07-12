package middleware

import (
	"context"
	"github.com/golibs-starter/golib/web/constant"
	webContext "github.com/golibs-starter/golib/web/context"
	"github.com/golibs-starter/golib/web/event"
)

func InitContextWorker(ctx context.Context) context.Context {
	attributes := ctx.Value(constant.ContextEventAttributes).(*event.Attributes)

	return context.WithValue(context.Background(), constant.ContextReqAttribute, &webContext.RequestAttributes{
		CorrelationId: attributes.CorrelationId,
		SecurityAttributes: webContext.SecurityAttributes{
			UserId:            attributes.UserId,
			TechnicalUsername: attributes.TechnicalUsername,
		},
	})
}
