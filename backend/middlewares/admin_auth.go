package middlewares

import (
	"MyBlog/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 鉴权
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		log.Printf("[auth-debug] AdminAuthMiddleware path=%s method=%s Authorization=%q", c.Request.URL.Path, c.Request.Method, authHeader)

		if authHeader == "" {
			log.Printf("[auth-debug] AdminAuthMiddleware missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少 Authorization 请求头",
			})
			c.Abort() //停止执行后续的Handler
			return
		}

		//Authorization:Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("[auth-debug] AdminAuthMiddleware invalid Authorization format parts_len=%d prefix_ok=%v", len(parts), len(parts) == 2 && parts[0] == "Bearer")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization 格式错误，应为Bearer token",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		log.Printf("[auth-debug] AdminAuthMiddleware bearer token_len=%d", len(tokenString))

		claims, err := utils.ParseAdminToken(tokenString)
		if err != nil {
			log.Printf("[auth-debug] AdminAuthMiddleware ParseAdminToken error=%v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		log.Printf("[auth-debug] AdminAuthMiddleware claims admin_id=%d username=%s subject=%s", claims.AdminID, claims.Username, claims.Subject)
		log.Printf("[auth-debug] AdminAuthMiddleware login_state_check redis=false db=false result=skipped")

		//把管理员信息放进上下文，后面的 controller 可以拿
		c.Set("admin_id", claims.AdminID)
		c.Set("username", claims.Username)
		c.Set("admin_claims", claims)
		log.Printf("[auth-debug] AdminAuthMiddleware context set admin_id=%d username=%s admin_claims_set=true", claims.AdminID, claims.Username)

		c.Next()
	}
}
