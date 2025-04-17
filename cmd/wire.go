//go:build wireinject
// +build wireinject

package main

import (
	"framework/controller"
	"framework/infrastructure/ioc"
	"framework/repo"
	"framework/repo/cache"
	"framework/repo/dao"
	"framework/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

//go:generate wire
func InitWeb() (*gin.Engine, error) {
	panic(wire.Build(
		//公共资源
		ioc.InitMysql,
		ioc.InitRedis,
		ioc.InitZap,
		ioc.InitQiNiuOSS,
		//UserController
		//dao
		dao.NewUserDao,
		//cache
		cache.NewUserCache,
		//repo
		repo.NewUserRepo,
		//service
		service.NewUserService,
		//controller
		controller.NewUserController,
		//CommonController
		//oss
		//controller
		controller.NewCommonController,
		//初始化gin
		ioc.InitMiddlewares, //加载全局中间件
		ioc.InitGin,
	))
}
