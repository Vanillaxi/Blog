package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateTagRequest struct {
	TagName string `json:"tag_name" binding:"required"`
	Slug    string `json:"slug" binding:"required"`
	Sort    uint32 `json:"sort" binding:"required"`
}

type UpdateTagStatusRequest struct {
	Status *int8 `json:"status" binding:"required"`
}

type UpdateTagSortRequest struct {
	Sort uint32 `json:"sort" binding:"required"`
}

// 创建标签
func CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误",
		})
		return
	}

	tag := models.Tag{
		TagName:      req.TagName,
		Slug:         req.Slug,
		ArticleCount: 0,
		Status:       1,
		Sort:         req.Sort,
	}

	if err := models.CreateTag(&tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建标签失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id": tag.ID,
		},
		"msg": "创建标签成功",
	})
}

// 更新标签状态
func UpdateTagStatus(c *gin.Context) {
	idParam := c.PostForm("id")

	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "标签ID不合法",
		})
		return
	}

	id := uint32(id64)
	var req UpdateTagStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误",
		})
		return
	}

	status := int8(*req.Status)
	if status != 0 && status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "状态值不合法",
		})
		return
	}

	if err := models.UpdateTagStatus(id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新标签状态失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新标签状态成功",
	})

}

// 前台显示
func GetTags(c *gin.Context) {
	categories, err := models.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取标签列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": categories,
		"msg":  "获取标签列表成功",
	})
}

// 后台展示
func GetAdminTags(c *gin.Context) {
	categories, err := models.GetAdminTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取标签列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": categories,
		"msg":  "获取标签列表成功",
	})
}

// 排序
func UpdateTagSort(c *gin.Context) {
	idParam := c.PostForm("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "标签ID不合法",
		})
		return
	}

	id := uint32(id64)
	var req UpdateTagSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误",
		})
		return
	}

	if err := models.UpdateTagSort(id, req.Sort); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新标签排序失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新标签排序成功",
	})

}
