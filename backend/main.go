package main

import (
	"MyBlog/config"
	_ "MyBlog/docs"
	"MyBlog/global"
	"MyBlog/initialize"
	"MyBlog/router"
	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title MyBlog API
// @version 1.0
// @description Go Gin personal blog backend API.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	//读取配置文件
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("读取配置失败：", err)
	}

	//保存到全局变量
	global.Config = cfg

	//初始化数据库
	err = initialize.InitDB()
	if err != nil {
		log.Fatal("初始化数据库失败：", err)
	}

	//初始化路由
	r := router.InitRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//启动服务
	err = r.Run(":" + global.Config.App.Port)
	if err != nil {
		log.Fatal("启动服务失败：", err)
	}

}
