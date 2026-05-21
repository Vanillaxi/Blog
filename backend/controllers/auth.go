package controllers

import (
	"MyBlog/models"
	"MyBlog/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// 只有Login有效，其余仅供测试
func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

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

func Register(c *gin.Context) {
	//1.定义带有校验标签的结构体
	var input struct {
		UserName string `json:"username" binding:"required,min=2,max=20"`
		Password string `json:"password" binding:"required,min=6,max=20"`
	}

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
		Username: input.UserName,
		Password: string(hashedPassword),
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

func UpdateUsername(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	currentUser := user.(*models.Admin)

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewUsername string `json:"new_username" binding:"required,min=2,max=20"`
	}

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

func UpdatePassword(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	currentUser := user.(*models.Admin)
	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
	}

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
