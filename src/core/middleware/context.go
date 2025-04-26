package middleware

import (
	"context"
	"github.com/golibs-starter/golib/web/constant"
	webContext "github.com/golibs-starter/golib/web/context"
	"github.com/google/uuid"
)

func InitContextWorker() context.Context {
	ctx := context.WithValue(context.Background(), constant.ContextReqAttribute, &webContext.RequestAttributes{
		CorrelationId: uuid.New().String(),
	})
	return ctx
}
