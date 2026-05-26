package controllers

import (
	"MyBlog/models"
	"MyBlog/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type AdminRegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=20" example:"admin"`
	Password string `json:"password" binding:"required,min=6,max=20" example:"123456"`
}

type UpdateUsernameRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"123456"`
	NewUsername string `json:"new_username" binding:"required,min=2,max=20" example:"newadmin"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"123456"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20" example:"654321"`
}

// Login godoc
// @Summary 管理员登录
// @Description 管理员使用用户名和密码登录，成功后返回 JWT token。
// @Tags admin
// @Accept json
// @Produce json
// @Param body body AdminLoginRequest true "管理员登录请求体。username 示例：admin；password 示例：123456"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"登录成功\",\"data\":{\"token\":\"xxx\",\"username\":\"admin\",\"userID\":1}}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数格式错误\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"Token生成失败\"}"
// @Router /api/admin/login [post]
func Login(c *gin.Context) {
	var input AdminLoginRequest

	//绑定前端传来的JSON参数
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数格式错误"})
		return
	}

	user, err := models.GetAdminByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "用户或密码错误",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 401,
			"msg":  "用户或密码密码错误",
		})
		return
	}

	//生成token
	token, err := utils.GenerateAdminToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Token生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token":    token,
			"username": user.Username,
			"userID":   user.ID,
		},
	})

	utils.Logger.WithFields(logrus.Fields{
		"userID":   user.ID,
		"username": user.Username,
		"ip":       c.ClientIP(),
		"action":   "login",
	}).Info("User logged in")

}

// Register godoc
// @Summary 管理员注册
// @Description 注册管理员账号，当前路由未加鉴权。
// @Tags admin
// @Accept json
// @Produce json
// @Param body body AdminRegisterRequest true "管理员注册请求体。username 示例：admin；password 示例：123456"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"Welcome!\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数校验失败\"}"
// @Failure 409 {object} object "示例：{\"code\":409,\"msg\":\"用户名已被注册\"}"
// @Router /api/admin/register [post]
func Register(c *gin.Context) {
	//1.定义带有校验标签的结构体
	var input AdminRegisterRequest

	//2.参数绑定与校验
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数校验失败：" + err.Error()})
		return
	}

	//3.密码哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码加密失败"})
		return
	}

	user := models.Admin{
		Username:   input.Username,
		Password:   string(hashedPassword),
		Nickname:   input.Username,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	//4.写入数据库
	if err := models.CreateAdmin(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "用户名已被注册"})
		return
	}

	//6.成功返回
	msg := fmt.Sprintf("Welcome!")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  msg,
	})

	utils.Logger.WithFields(logrus.Fields{
		"userID":   user.ID,
		"username": user.Username,
		"ip":       c.ClientIP(),
		"action":   "register",
	}).Info("User registred")

}

// UpdateUsername godoc
// @Summary 修改管理员用户名
// @Description 已登录管理员修改用户名。
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body UpdateUsernameRequest true "修改用户名请求体。old_password 示例：123456；new_username 示例：newadmin"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"用户名修改成功，请重新登录\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数校验失败\"}"
// @Failure 401 {object} object "示例：{\"code\":401,\"msg\":\"用户未登录\"}"
// @Router /api/admin/username [post]
func UpdateUsername(c *gin.Context) {
	adminIDVal, exists := c.Get("admin_id")
	if !exists {
		username, usernameExists := c.Get("username")
		log.Printf("[auth-debug] UpdateUsername admin_id missing username_exists=%v username=%v", usernameExists, username)
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	adminID, ok := adminIDVal.(uint32)
	if !ok {
		log.Printf("[auth-debug] UpdateUsername admin_id type invalid value=%v", adminIDVal)
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	currentUser, err := models.GetAdminByID(uint64(adminID))
	if err != nil {
		log.Printf("[auth-debug] UpdateUsername admin lookup failed admin_id=%d err=%v", adminID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	var input UpdateUsernameRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数校验失败"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "用户名已被占用"})
		return
	}

	if err := models.UpdateAdminName(currentUser.ID, input.NewUsername); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "用户名更新失败！"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "用户名修改成功，请重新登录",
	})

}

// UpdatePassword godoc
// @Summary 修改管理员密码
// @Description 已登录管理员修改密码。
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body UpdatePasswordRequest true "修改密码请求体。old_password 示例：123456；new_password 示例：654321"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"密码修改成功，请重新登录\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数校验失败\"}"
// @Failure 401 {object} object "示例：{\"code\":401,\"msg\":\"旧密码错误\"}"
// @Router /api/admin/password [post]
func UpdatePassword(c *gin.Context) {
	adminIDVal, exists := c.Get("admin_id")
	if !exists {
		username, usernameExists := c.Get("username")
		log.Printf("[auth-debug] UpdatePassword admin_id missing username_exists=%v username=%v", usernameExists, username)
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	adminID, ok := adminIDVal.(uint32)
	if !ok {
		log.Printf("[auth-debug] UpdatePassword admin_id type invalid value=%v", adminIDVal)
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	currentUser, err := models.GetAdminByID(uint64(adminID))
	if err != nil {
		log.Printf("[auth-debug] UpdatePassword admin lookup failed admin_id=%d err=%v", adminID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	var input UpdatePasswordRequest

	//参数绑定
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数校验失败：" + err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "旧密码错误"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码更新失败"})
		return
	}

	if err := models.UpdatePassword(currentUser.ID, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "密码修改成功，请重新登录",
	})

}
