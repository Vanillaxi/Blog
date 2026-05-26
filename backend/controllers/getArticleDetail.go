package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetArtileDetail godoc
// @Summary 获取文章详情
// @Description 根据文章 ID 获取已发布且未删除的文章详情。
// @Tags article
// @Accept json
// @Produce json
// @Param id path int true "文章 ID"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"id\":1,\"title\":\"标题\"},\"msg\":\"查询成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"ID格式错误，请输入数字\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"查询失败或文章不存在\"}"
// @Router /api/articles/{id} [get]
func GetArtileDetail(c *gin.Context) {
	//1. 获取路径参数 id （字符串）
	idStr := c.Param("id")

	//2.将字符串转换为数字(uint32)
	//Atoi 是转为int,如果是uint32建议用ParseUint
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		//如果转化失败（比如输入了非数字），返回400
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "ID格式错误，请输入数字",
		})
		return
	}

	//3. 转换为uint32并调用模型层方法
	id := uint(idUint64)
	article, err := models.GetArticleDetail(uint32(id))
	if err != nil {
		//处理数据库查询错误或文章不存在
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败或文章不存在",
		})
		return
	}

	//4.返回成功数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": article,
		"msg":  "查询成功",
	})

}

func GetAdminArticleDetail(c *gin.Context) {
	idUint64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章ID不合法",
		})
		return
	}

	article, err := models.GetAdminArticleDetail(uint32(idUint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询文章详情失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": article,
		"msg":  "查询成功",
	})
}
