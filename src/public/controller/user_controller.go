package controller

import (
	"github.com/KhaiHust/email-notification-service/core/common"
	"github.com/KhaiHust/email-notification-service/public/apihelper"
	"github.com/KhaiHust/email-notification-service/public/resource/request"
	"github.com/KhaiHust/email-notification-service/public/resource/response"
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
func (u *UserController) Login(c *gin.Context) {
	var req request.LoginRequest
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
	result, err := u.userService.Login(c, &req)
	if err != nil {
		log.Error(c, "Login error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, result)
}
func (u *UserController) GenerateTokenFromRefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
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
	result, err := u.userService.GenerateTokenFromRefreshToken(c, req.RefreshToken)
	if err != nil {
		log.Error(c, "GenerateTokenFromRefreshToken error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, result)
}
func (u *UserController) GetListMembers(c *gin.Context) {
	workspaceID := u.GetWorkspaceIDFromContext(c)
	if workspaceID == 0 {
		log.Error(c, "GetListMembers error: workspaceID is 0")
		apihelper.AbortErrorHandle(c, common.ErrBadRequest)
		return
	}
	users, err := u.userService.GetListMembers(c, workspaceID)
	if err != nil {
		log.Error(c, "GetListMembers error: %v", err)
		apihelper.AbortErrorHandle(c, err)
		return
	}
	apihelper.SuccessfulHandle(c, response.ToListUserResponseResource(users))
}
func NewUserController(base *BaseController, userService service.IUserService) *UserController {
	return &UserController{
		BaseController: base,
		userService:    userService,
	}
}
