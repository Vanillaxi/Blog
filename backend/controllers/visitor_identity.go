package controllers

import (
	"MyBlog/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetVisitorIdentity godoc
// @Summary 查询邮箱绑定昵称
// @Description 根据邮箱查询是否已有绑定昵称。只返回是否存在和昵称，不返回邮箱、IP 或 UA。
// @Tags visitor
// @Accept json
// @Produce json
// @Param email query string true "邮箱"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"exists\":true,\"nickname\":\"Vanillaxi\"},\"msg\":\"查询成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"邮箱格式不正确\"}"
// @Router /api/visitor-identity [get]
func GetVisitorIdentity(c *gin.Context) {
	email := strings.TrimSpace(c.Query("email"))
	nickname, exists, err := models.FindVisitorNicknameByEmail(email)
	if err != nil {
		if models.IsCommentValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询访客身份失败",
		})
		return
	}

	data := gin.H{"exists": exists}
	if exists {
		data["nickname"] = nickname
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
		"msg":  "查询成功",
	})
}
