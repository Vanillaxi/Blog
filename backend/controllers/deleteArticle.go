package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteArticle(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章ID不合法",
		})
		return
	}

	id := uint32(id64)

	if err := models.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除文章失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除文章成果",
	})
}
