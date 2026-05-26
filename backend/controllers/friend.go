package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateFriendLinkRequest struct {
	Name        string `json:"name" binding:"required" example:"Example Blog"`
	URL         string `json:"url" binding:"required" example:"https://example.com"`
	Logo        string `json:"logo" example:"https://example.com/logo.png"`
	Description string `json:"description" example:"这个站点的一句话介绍"`
	Sort        int    `json:"sort" example:"10"`
	Status      int    `json:"status" example:"1"`
}

type UpdateFriendLinkRequest struct {
	ID          uint64 `json:"id" binding:"required" example:"1"`
	Name        string `json:"name" binding:"required" example:"Example Blog Updated"`
	URL         string `json:"url" binding:"required" example:"https://example.com"`
	Logo        string `json:"logo" example:"https://example.com/logo-updated.png"`
	Description string `json:"description" example:"这个站点的一句话介绍"`
	Sort        int    `json:"sort" example:"20"`
	Status      int    `json:"status" example:"1"`
}

// GetFriendLinks godoc
// @Summary 获取前台友链
// @Description 分页获取启用状态的友链。
// @Tags friendlink
// @Accept json
// @Produce json
// @Param page query int false "页码，默认 1"
// @Param pageSize query int false "每页数量，默认 10"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"links\":[],\"total\":0,\"page\":1,\"size\":10},\"msg\":\"获取成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"获取友链失败\"}"
// @Router /api/friendlinks [get]
func GetFriendLinks(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	links, total, err := models.GetActiveFriendLinks(pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取友链失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"links": links,
			"total": total,
			"page":  pageNum,
			"size":  pageSize,
		},
		"msg": "获取成功",
	})
}

// GetAdminFriendLinks godoc
// @Summary 获取后台友链
// @Description 管理员分页获取全部友链。
// @Tags admin/friendlink
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码，默认 1"
// @Param pageSize query int false "每页数量，默认 10"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"links\":[],\"total\":0,\"page\":1,\"size\":10},\"msg\":\"获取成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"获取友链失败\"}"
// @Router /api/admin/friendlinks [get]
func GetAdminFriendLinks(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	links, total, err := models.GetAdminFriendLinks(pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取友链失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"links": links,
			"total": total,
			"page":  pageNum,
			"size":  pageSize,
		},
		"msg": "获取成功",
	})
}

// CreateFriendLink godoc
// @Summary 创建友链
// @Description 管理员创建友链。
// @Tags admin/friendlink
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body CreateFriendLinkRequest true "创建友链请求体。name 示例：Example Blog；url 示例：https://example.com；status 示例：1"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"创建成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数错误\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"创建失败\"}"
// @Router /api/admin/friendlinks/add [post]
func CreateFriendLink(c *gin.Context) {
	var req CreateFriendLinkRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	link := models.FriendLink{
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	if err := models.CreateFriendLink(&link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// UpdateFriendLink godoc
// @Summary 更新友链
// @Description 管理员更新友链。
// @Tags admin/friendlink
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "友链 ID"
// @Param body body UpdateFriendLinkRequest true "更新友链请求体。id 需要和 path 中的友链 ID 保持一致"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"更新成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数错误\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"更新失败\"}"
// @Router /api/admin/friendlinks/{id}/update [put]
func UpdateFriendLink(c *gin.Context) {
	var req UpdateFriendLinkRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	link := models.FriendLink{
		ID:          req.ID,
		Name:        req.Name,
		URL:         req.URL,
		Logo:        req.Logo,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	if err := models.UpdateFriendLink(&link); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})

}

// DeleteFriendLink godoc
// @Summary 删除友链
// @Description 管理员按友链 ID 删除友链。
// @Tags admin/friendlink
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "友链 ID"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"删除成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"删除失败\"}"
// @Router /api/admin/friendlinks/{id}/delete [delete]
func DeleteFriendLink(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := models.DeleteFriendLink(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
