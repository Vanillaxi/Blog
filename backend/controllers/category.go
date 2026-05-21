package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required"`
	Slug         string `json:"slug"`
	Sort         uint32 `json:"sort"`
}

type UpdateCategoryStatusRequest struct {
	Status *int8 `json:"status" binding:"required"`
}

type UpdateCategorySortRequest struct {
	Sort uint32 `json:"sort" binding:"required"`
}

// 创建分类
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

// 更新分类状态
func UpdateCategoryStatus(c *gin.Context) {
	idParam := c.PostForm("id")

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

// 前台显示
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

// 后台展示
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

// 排序
func UpdateCategorySort(c *gin.Context) {
	idParam := c.PostForm("id")
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
