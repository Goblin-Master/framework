package middleware

import (
	"framework/infrastructure/utils/jwt"
	"github.com/gin-gonic/gin"
)

type AuthenticationMiddlewareBuilder struct {
	//这里以后可以加入一些参数，如 redis用于存储token
}

// 思考：通过依赖注入redis，但是其它就没有办法简单使用中间件，只能通过加一个ignore方法忽略路径
func NewAuthenticationMiddlewareBuilder() *AuthenticationMiddlewareBuilder {
	return &AuthenticationMiddlewareBuilder{}
}

func (l *AuthenticationMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			//还没包装响应
			c.Abort()
			return
		}
		//解析token是否有效，并取出上一次的值
		data, err := jwt.IdentifyToken(token)
		if err != nil {
			//对应token无效，直接让他返回
			c.Abort()
			return
		}
		//判断其是否为atoken
		if data.Class != jwt.AUTH_ENUMS_ATOKEN {
			// 类型错误
			c.Abort()
			return
		}
		//将token内部数据传下去
		c.Set(jwt.TOKEN_USER_ID, data.UserID)
		c.Next()
	}
}
