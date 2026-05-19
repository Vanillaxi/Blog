package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	//测试接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//前台公开接口
	api := r.Group("/api")
	{
		api.GET("/articles", controllers.GetArticles)

		api.GET("/category", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "分类列表接口",
			})
		})
	}

	//后台接口
	admin := r.Group("/admin")
	{
		admin.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "管理员登录接口",
			})
		})

		//TODO:加jwt
		admin.POST("/articles", controllers.CreateArticle)
	}

	return r
}
