package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateTagRequest struct {
	TagName string `json:"tag_name" binding:"required" example:"Gin"`
	Slug    string `json:"slug" binding:"required" example:"gin"`
	Sort    uint32 `json:"sort" binding:"required" example:"10"`
}

type UpdateTagStatusRequest struct {
	Status *int8 `json:"status" binding:"required" example:"1"`
}

type UpdateTagSortRequest struct {
	Sort uint32 `json:"sort" binding:"required" example:"20"`
}

// CreateTag godoc
// @Summary 创建标签
// @Description 管理员创建文章标签。
// @Tags admin/tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body CreateTagRequest true "创建标签请求体。tag_name 示例：Gin；slug 示例：gin；sort 示例：10"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"id\":1},\"msg\":\"创建标签成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数格式错误\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"创建标签失败\"}"
// @Router /api/admin/tag/create [post]
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

// UpdateTagStatus godoc
// @Summary 更新标签状态
// @Description 管理员更新标签启用状态。
// @Tags admin/tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "标签 ID"
// @Param body body UpdateTagStatusRequest true "标签状态请求体。status 示例：1 表示启用，0 表示停用"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"更新标签状态成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"标签ID不合法\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"更新标签状态失败\"}"
// @Router /api/admin/tag/{id}/status [put]
func UpdateTagStatus(c *gin.Context) {
	idParam := c.Param("id")

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

// GetTags godoc
// @Summary 获取前台标签
// @Description 获取启用状态的标签列表。
// @Tags tag
// @Accept json
// @Produce json
// @Param unused query string false "无业务参数"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":[],\"msg\":\"获取标签列表成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"获取标签列表失败\"}"
// @Router /api/tags [get]
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

// GetAdminTags godoc
// @Summary 获取后台标签
// @Description 管理员获取全部标签列表。
// @Tags admin/tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":[],\"msg\":\"获取标签列表成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"获取标签列表失败\"}"
// @Router /api/admin/tags [get]
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

// UpdateTagSort godoc
// @Summary 更新标签排序
// @Description 管理员更新标签排序值。
// @Tags admin/tag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "标签 ID"
// @Param body body UpdateTagSortRequest true "标签排序请求体。sort 示例：20"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"更新标签排序成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"标签ID不合法\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"更新标签排序失败\"}"
// @Router /api/admin/tag/{id}/sort [put]
func UpdateTagSort(c *gin.Context) {
	idParam := c.Param("id")
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
