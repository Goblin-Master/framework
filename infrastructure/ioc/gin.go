package ioc

import (
	"ddd/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// InitGin
//
//	@Description: 初始化gin
//	@param middlewares
//	@param userController
//	@return *gin.Engine
func InitGin(middlewares []gin.HandlerFunc, userController *controller.UserController) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares...)
	userController.RegisterRouter(r)
	return r
}

// InitMiddlewares
//
//	@Description: 这里放全局中间件
//	@return []gin.HandlerFunc
func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		Cors(),
	}
}

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:  []string{"Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		//是否允许你带cookie之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://127.0.0.1") {
				return true
			}
			return strings.HasPrefix(origin, "http://www.example.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
