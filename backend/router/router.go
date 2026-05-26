package router

import (
	"MyBlog/controllers"
	"MyBlog/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		panic(err)
	}
	r.Use(middlewares.CorsMiddleware())
	r.Use(middlewares.GlobalRateLimiter())

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
		v1.GET("/articles/search", controllers.SearchArticles)
		v1.GET("/articles/:id", controllers.GetArtileDetail)

		v1.GET("/visitor-identity", controllers.GetVisitorIdentity)
		v1.POST("/comments/add", controllers.AddComment)
		v1.GET("/cmments/get", controllers.GetComment)

		v1.GET("/friendlinks", controllers.GetFriendLinks)

	}

	//admin接口
	admin := r.Group("/api/admin")
	{
		//TODO:后期需要隐藏，并且取消register
		admin.POST("/login", controllers.Login)
		admin.POST("/register", controllers.Register)

		admin.Use(middlewares.AdminAuthMiddleware())
		{
			admin.POST("/username", controllers.UpdateUsername)
			admin.POST("/password", controllers.UpdatePassword)
			admin.GET("/dashboard", controllers.GetAdminDashboard)

			admin.POST("/articles/create", controllers.CreateArticle)
			admin.PUT("/articles/update/:id", controllers.UpdateArticle)
			admin.PUT("/articles/updateStatus/:id", controllers.UpdateArticleStatus)
			admin.DELETE("/articles/delete/:id", controllers.DeleteArticle)
			admin.GET("/articles/:id", controllers.GetAdminArticleDetail)

			admin.GET("/tags", controllers.GetAdminTags)
			admin.POST("/tag/create", controllers.CreateTag)
			admin.PUT("/tag/:id/status", controllers.UpdateTagStatus)
			admin.PUT("/tag/:id/sort", controllers.UpdateTagSort)

			admin.POST("/categories", controllers.GetAdminCategories)
			admin.POST("/categories/create", controllers.CreateCategory)
			admin.PUT("/categories/:id/status", controllers.UpdateCategoryStatus)
			admin.PUT("/categories/:id/sort", controllers.UpdateCategorySort)

			admin.GET("/articles", controllers.GetAdminArticles)

			//评论/留言
			admin.GET("/comments", controllers.GetAdminComments)
			admin.DELETE("/comments/delete/:id", controllers.DeleteComment)
			admin.PUT("/comments/:id/restore", controllers.RestoreComment)

			//友链
			admin.GET("/friendlinks", controllers.GetAdminFriendLinks)
			admin.POST("/friendlinks/add", controllers.CreateFriendLink)
			admin.PUT("/friendlinks/:id/update", controllers.UpdateFriendLink)
			admin.DELETE("/friendlinks/:id/delete", controllers.DeleteFriendLink)
		}

	}

	return r
}
