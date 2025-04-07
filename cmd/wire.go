//go:build wireinject
// +build wireinject

package main

import (
	"ddd/controller"
	"ddd/infrastructure/ioc"
	"ddd/repo"
	"ddd/repo/cache"
	"ddd/repo/dao"
	"ddd/service"
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
		//初始化gin
		ioc.InitMiddlewares, //加载全局中间件
		ioc.InitGin,
	))
}
