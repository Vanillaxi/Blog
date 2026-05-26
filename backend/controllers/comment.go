package controllers

import (
	"MyBlog/models"
	"MyBlog/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateCommentRequest struct {
	TargetType int8   `json:"target_type" binding:"required" example:"1"`
	TargetID   uint64 `json:"target_id" example:"1"`
	ParentID   uint64 `json:"parent_id" example:"0"`
	Nickname   string `json:"nickname" example:"Apifox Tester"`
	Email      string `json:"email" example:"tester@example.com"`
	Website    string `json:"website" example:"https://example.com"`
	Content    string `json:"content" example:"这是一条 Apifox 测试评论"`
}

// GetAdminComments godoc
// @Summary 获取后台评论/留言列表
// @Description 管理员分页获取评论或留言，支持目标类型、文章 ID 和删除状态筛选。
// @Tags admin/comment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param target_type query int false "目标类型，1 文章评论，2 留言"
// @Param target_id query int false "目标 ID"
// @Param include_deleted query bool false "是否包含已删除，默认 false"
// @Param page query int false "页码，默认 1"
// @Param page_size query int false "每页数量，默认 10"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"list\":[],\"total\":0},\"msg\":\"获取评论成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"获取评论失败\"}"
// @Router /api/admin/comments [get]
func GetAdminComments(c *gin.Context) {
	targetType64, _ := strconv.Atoi(c.DefaultQuery("target_type", "0"))
	targetID, _ := strconv.ParseUint(c.DefaultQuery("target_id", "0"), 10, 64)
	includeDeleted, _ := strconv.ParseBool(c.DefaultQuery("include_deleted", "false"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", c.DefaultQuery("pageSize", "10")))

	comments, total, err := models.GetAdminComments(int8(targetType64), targetID, includeDeleted, pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      comments,
			"total":     total,
			"page":      pageNum,
			"page_size": pageSize,
		},
		"msg": "获取评论成功",
	})
}

type CreateMessageRequest struct {
	TargetType int8   `json:"target_type" binding:"required" example:"2"`
	TargetID   uint64 `json:"target_id" example:"0"`
	ParentID   uint64 `json:"parent_id" example:"0"`
	Nickname   string `json:"nickname" example:"留言用户"`
	Email      string `json:"email" example:"message@example.com"`
	Website    string `json:"website" example:"https://example.com"`
	Content    string `json:"content" example:"这是一条 Apifox 测试留言"`
}

// AddComment godoc
// @Summary 添加评论
// @Description 添加文章评论或留言板评论。
// @Tags comment
// @Accept json
// @Produce json
// @Param body body CreateCommentRequest true "评论请求体。文章评论 target_type 示例：1；留言 target_type 示例：2，留言示例可参考 CreateMessageRequest 字段"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"id\":1},\"msg\":\"评论成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数错误\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"添加评论失败\"}"
// @Router /api/comments/add [post]
func AddComment(c *gin.Context) {
	var req CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	clientIP := utils.GetClientIP(c)
	ipLocation := utils.ResolveIPLocation(clientIP)
	log.Printf(
		"[comment] ip debug x_real_ip=%q x_forwarded_for=%q gin_client_ip=%q resolved_client_ip=%q remote_addr=%q location=%q",
		c.GetHeader("X-Real-IP"),
		c.GetHeader("X-Forwarded-For"),
		c.ClientIP(),
		clientIP,
		c.Request.RemoteAddr,
		ipLocation,
	)
	comment := &models.Comment{
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		ParentID:   req.ParentID,
		Nickname:   req.Nickname,
		Email:      req.Email,
		Website:    req.Website,
		Content:    req.Content,
		IP:         clientIP,
		IPLocation: ipLocation,
		UserAgent:  c.Request.UserAgent(),
	}

	if err := models.AddComment(comment); err != nil {
		if models.IsCommentValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": comment,
		"msg":  "评论成功",
	})

}

// GetComment godoc
// @Summary 获取评论
// @Description 根据评论目标类型和目标 ID 分页获取评论列表。
// @Tags comment
// @Accept json
// @Produce json
// @Param target_type query int true "目标类型，1 文章，2 留言板"
// @Param target_id query int true "目标 ID"
// @Param page query int false "页码，默认 1"
// @Param pageSize query int false "每页数量，默认 10"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"list\":[],\"total\":0,\"page\":1,\"page_size\":10},\"msg\":\"获取评论成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"获取评论失败\"}"
// @Router /api/cmments/get [get]
func GetComment(c *gin.Context) {
	targetType64, _ := strconv.Atoi(c.Query("target_type"))
	targetID, _ := strconv.ParseUint(c.Query("target_id"), 10, 64)
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	comments, total, err := models.GetComments(int8(targetType64), targetID, pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      comments,
			"total":     total,
			"page":      pageNum,
			"page_size": pageSize,
		},
		"msg": "获取评论成功",
	})

}

// DeleteComment godoc
// @Summary 删除评论
// @Description 管理员按评论 ID 删除评论。
// @Tags admin/comment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "评论 ID"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"删除评论成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"删除评论失败\"}"
// @Router /api/admin/comments/delete/{id} [delete]
func DeleteComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := models.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除评论失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除评论成功",
	})
}

// RestoreComment godoc
// @Summary 恢复评论
// @Description 管理员按评论 ID 恢复已删除评论或留言。
// @Tags admin/comment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论 ID"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"恢复评论成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"恢复评论失败\"}"
// @Router /api/admin/comments/{id}/restore [put]
func RestoreComment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := models.RestoreComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "恢复评论失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "恢复评论成功",
	})
}
