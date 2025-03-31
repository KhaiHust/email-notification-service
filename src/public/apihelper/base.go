package apihelper

import (
	"errors"
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AbortErrorHandle(ctx *gin.Context, err error) {
	var e common.Errs
	if errors.As(err, &e) {
		ctx.JSON(e.HttpStatusCode, Response{
			Data: nil,
			Meta: ResponseMeta{
				Code:    e.Code,
				Message: e.Error(),
			},
		})
		return
	}
	internalErr := common.ErrInternalServer
	ctx.JSON(internalErr.HttpStatusCode, Response{
		Data: nil,
		Meta: ResponseMeta{
			Code:    internalErr.Code,
			Message: internalErr.Error(),
		},
	})
}
func AbortErrorHandleWithData(ctx *gin.Context, err error, data interface{}) {
	var e common.Errs
	if errors.As(err, &e) {
		ctx.JSON(e.HttpStatusCode, Response{
			Data: data,
			Meta: ResponseMeta{
				Code:    e.Code,
				Message: e.Error(),
			},
		})
		return
	}
	internalErr := common.ErrInternalServer
	ctx.JSON(internalErr.HttpStatusCode, Response{
		Data: data,
		Meta: ResponseMeta{
			Code:    internalErr.Code,
			Message: internalErr.Error(),
		},
	})
}
func SuccessfulHandle(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Data: data,
		Meta: ResponseMeta{
			Code:    http.StatusOK,
			Message: "OK",
		},
	})
}

type Response struct {
	Data interface{} `json:"data,omitempty"`
	Meta interface{} `json:"meta,omitempty"`
}
type ResponseMeta struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}
