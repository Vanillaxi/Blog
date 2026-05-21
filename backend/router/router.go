package router

import (
	"MyBlog/controllers"
	"MyBlog/middlewares"

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

	//publish接口
	v1 := r.Group("/api")
	{
		v1.GET("/categories", controllers.GetCategories)
		v1.GET("/tags", controllers.GetTags)

		v1.GET("/categories/:id/articles", controllers.GetArticlesByCategoryID)
		v1.GET("/tags/:id/articles", controllers.GetArticlesByTagID)
		v1.GET("/articles/timeline", controllers.GetArticlesTimeline)
		v1.GET("/articles/:id", controllers.GetArtileDetail)

	}

	//admin接口
	admin := r.Group("/api/admin")
	{
		//TODO:后期需要隐藏，并且取消register
		admin.POST("/login", controllers.Login)
		admin.POST("/register", controllers.Register)
		admin.POST("/username", controllers.UpdateUsername)
		admin.POST("/password", controllers.UpdatePassword)

		admin.Use(middlewares.AdminAuthMiddleware())
		{
			admin.POST("/articles/create", controllers.CreateArticle)
			admin.PUT("/articles/update/:id", controllers.UpdateArticle)
			admin.PUT("/articles/updateStatus/:id", controllers.UpdateArticleStatus)
			admin.DELETE("/articles/delete/:id", controllers.DeleteArticle)

			admin.GET("/tags", controllers.GetAdminTags)
			admin.POST("/tag/create", controllers.CreateTag)
			admin.PUT("/tag/:id/status", controllers.UpdateTagStatus)
			admin.PUT("/tag/:id/sort", controllers.UpdateTagSort)

			admin.POST("/categories", controllers.GetAdminCategories)
			admin.POST("/categories/create", controllers.CreateCategory)
			admin.PUT("/categories/:id/status", controllers.UpdateCategoryStatus)
			admin.PUT("/categories/:id/sort", controllers.UpdateCategorySort)

			admin.GET("/articles", controllers.GetAdminArticles)
		}

	}

	return r
}
