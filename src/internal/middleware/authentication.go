package middleware

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/core/usecase"
	"github.com/KhaiHust/email-notification-service/internal/apihelper"
	"github.com/gin-gonic/gin"
)

type APIKeyMiddleware struct {
	validateApiKeyUsecase usecase.IValidateApiKeyUsecase
}

func (a *APIKeyMiddleware) AuthenticationMiddlewareHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawKey := c.GetHeader("Authorization")
		if rawKey == "" {
			apihelper.AbortErrorHandle(c, common.ErrUnauthorized)
			return
		}

		if len(rawKey) > 7 && rawKey[:7] == "Bearer " {
			rawKey = rawKey[7:]
		} else {
			apihelper.AbortErrorHandle(c, common.ErrUnauthorized)
			return
		}
		valid, err := a.validateApiKeyUsecase.ValidateKey(c, rawKey)
		if err != nil || !valid {
			apihelper.AbortErrorHandle(c, common.ErrUnauthorized)
			return
		}
		c.Next()
	}
}
func NewAPIKeyMiddleware(
	validateApiKeyUsecase usecase.IValidateApiKeyUsecase,
) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		validateApiKeyUsecase: validateApiKeyUsecase,
	}
}
