package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteArticle godoc
// @Summary 删除文章
// @Description 管理员按文章 ID 删除文章。
// @Tags admin/article
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "文章 ID"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"删除文章成果\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"文章ID不合法\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"删除文章失败\"}"
// @Router /api/admin/articles/delete/{id} [delete]
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
