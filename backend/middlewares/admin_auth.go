package middlewares

import (
	"MyBlog/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "缺少 Authorization 请求头",
			})
			c.Abort()
			return
		}

		//Authorization:Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization 格式错误，应为Bearer token",
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseAdminToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		//把管理员信息放进上下文，后面的 controller 可以拿
		c.Set("claims", claims.AdminID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
