package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required" example:"Go 后端"`
	Slug         string `json:"slug" example:"go-backend"`
	Sort         uint32 `json:"sort" example:"10"`
}

type UpdateCategoryStatusRequest struct {
	Status *int8 `json:"status" binding:"required" example:"1"`
}

type UpdateCategorySortRequest struct {
	Sort uint32 `json:"sort" binding:"required" example:"20"`
}

// CreateCategory godoc
// @Summary 创建分类
// @Description 管理员创建文章分类。
// @Tags admin/category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body CreateCategoryRequest true "创建分类请求体。category_name 示例：Go 后端；slug 示例：go-backend；sort 示例：10"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"id\":1},\"msg\":\"创建分类成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数格式错误\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"创建分类失败\"}"
// @Router /api/admin/categories/create [post]
func CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误",
		})
		return
	}

	category := models.Category{
		CategoryName: req.CategoryName,
		Slug:         req.Slug,
		Sort:         req.Sort,
		Status:       1,
		ArticleCount: 0,
	}

	if err := models.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建分类失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id": category.ID,
		},
		"msg": "创建分类成功",
	})

}

// UpdateCategoryStatus godoc
// @Summary 更新分类状态
// @Description 管理员更新分类启用状态。
// @Tags admin/category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "分类 ID"
// @Param body body UpdateCategoryStatusRequest true "分类状态请求体。status 示例：1 表示启用，0 表示停用"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"更新分类状态成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"分类ID不合法\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"更新分类状态失败\"}"
// @Router /api/admin/categories/{id}/status [put]
func UpdateCategoryStatus(c *gin.Context) {
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "分类ID不合法",
		})
		return
	}

	id := uint32(id64)
	var req UpdateCategoryStatusRequest
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

	if err := models.UpdateCategoryStatus(id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新分类状态失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新分类状态成功",
	})

}

// GetCategories godoc
// @Summary 获取前台分类
// @Description 获取启用状态的分类列表。
// @Tags category
// @Accept json
// @Produce json
// @Param unused query string false "无业务参数"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":[],\"msg\":\"获取分类列表成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"获取分类列表失败\"}"
// @Router /api/categories [get]
func GetCategories(c *gin.Context) {
	categories, err := models.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取分类列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": categories,
		"msg":  "获取分类列表成功",
	})
}

// GetAdminCategories godoc
// @Summary 获取后台分类
// @Description 管理员获取全部分类列表。
// @Tags admin/category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":[],\"msg\":\"获取分类列表成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"获取分类列表失败\"}"
// @Router /api/admin/categories [post]
func GetAdminCategories(c *gin.Context) {
	categories, err := models.GetAdminCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取分类列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": categories,
		"msg":  "获取分类列表成功",
	})
}

// UpdateCategorySort godoc
// @Summary 更新分类排序
// @Description 管理员更新分类排序值。
// @Tags admin/category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "分类 ID"
// @Param body body UpdateCategorySortRequest true "分类排序请求体。sort 示例：20"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"更新分类排序成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"分类ID不合法\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"更新分类排序失败\"}"
// @Router /api/admin/categories/{id}/sort [put]
func UpdateCategorySort(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "分类ID不合法",
		})
		return
	}

	id := uint32(id64)
	var req UpdateCategorySortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误",
		})
		return
	}

	if err := models.UpdateCategorySort(id, req.Sort); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新分类排序失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新分类排序成功",
	})

}
