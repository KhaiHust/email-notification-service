package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/service"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

type UserController struct {
	*BaseController
	userService service.IUserService
}

func (u *UserController) SignUp(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(c, "BindJSON error: %v", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	if err := u.validator.Struct(&req); err != nil {
		log.Error(c, "Validator error: %v", err)
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	_, err := u.userService.SignUp(c, &req)
	if err != nil {
		log.Error(c, "SignUp error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, nil)

}
func NewUserController(base *BaseController, userService service.IUserService) *UserController {
	return &UserController{
		BaseController: base,
		userService:    userService,
	}
}
