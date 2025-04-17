package controller

import (
	"framework/domain"
	"framework/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	zap         *zap.SugaredLogger
	userService *service.UserService
}

func NewUserController(zap *zap.SugaredLogger, userService *service.UserService) *UserController {
	return &UserController{
		zap:         zap,
		userService: userService,
	}
}

// RegisterRouter
//
//	@Description: 对于每一个controller都提供一个方法去注册路由
//	@receiver uc
//	@param r
func (uc *UserController) RegisterRouter(r *gin.Engine) {
	rg := r.Group("user")
	rg.GET("", uc.GetUser)
}

func (uc *UserController) GetUser(c *gin.Context) {
	// 绑定参数
	uc.zap.Info(uc.userService.GetUser(c.Request.Context(), domain.User{}))
}
