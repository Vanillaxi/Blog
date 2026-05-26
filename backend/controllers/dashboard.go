package controllers

import (
	"MyBlog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAdminDashboard godoc
// @Summary 获取后台 Dashboard 统计
// @Description 获取文章、分类、标签、评论、留言、友链统计和最近文章。
// @Tags admin/dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"article_count\":1},\"msg\":\"获取 Dashboard 成功\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"获取 Dashboard 失败\"}"
// @Router /api/admin/dashboard [get]
func GetAdminDashboard(c *gin.Context) {
	var stats models.DashboardStats
	var err error

	if stats.ArticleCount, err = models.CountArticles(-1, -1); err != nil {
		failDashboard(c)
		return
	}
	if stats.PublishedCount, err = models.CountArticles(1, 0); err != nil {
		failDashboard(c)
		return
	}
	if stats.DraftCount, err = models.CountArticles(0, 0); err != nil {
		failDashboard(c)
		return
	}
	if stats.OfflineCount, err = models.CountArticles(2, 0); err != nil {
		failDashboard(c)
		return
	}
	if stats.DeletedArticleCount, err = models.CountArticles(-1, 1); err != nil {
		failDashboard(c)
		return
	}
	if stats.CategoryCount, err = models.CountCategories(); err != nil {
		failDashboard(c)
		return
	}
	if stats.TagCount, err = models.CountTags(); err != nil {
		failDashboard(c)
		return
	}
	if stats.CommentCount, err = models.CountCommentsByTargetType(1); err != nil {
		failDashboard(c)
		return
	}
	if stats.GuestbookCount, err = models.CountCommentsByTargetType(2); err != nil {
		failDashboard(c)
		return
	}
	if stats.FriendlinkCount, err = models.CountFriendLinks(); err != nil {
		failDashboard(c)
		return
	}
	if stats.RecentArticles, err = models.GetRecentArticles(5); err != nil {
		failDashboard(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": stats,
		"msg":  "获取 Dashboard 成功",
	})
}

func failDashboard(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": 500,
		"msg":  "获取 Dashboard 失败",
	})
}
