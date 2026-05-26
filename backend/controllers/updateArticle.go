package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateArticleRequest struct {
	Title      string   `json:"title" binding:"required" example:"我的第一篇博客 - 已更新"`
	CategoryID uint32   `json:"category_id" binding:"required" example:"1"`
	Summary    string   `json:"summary" example:"这是一篇更新后的 Go 后端开发文章摘要"`
	Content    string   `json:"content" binding:"required" example:"这里是更新后的文章正文内容。"`
	IsTop      int8     `json:"is_top" example:"1"`
	CoverURL   string   `json:"cover_url" example:"https://example.com/cover-updated.jpg"`
	TagIDs     []uint32 `json:"tag_ids" example:"1,2"`
}

type UpdateArticleStatusRequest struct {
	Status *int8 `json:"status" binding:"required" example:"1"` //要用指针，不然传0可能会认为没传，导致格式错误
}

// UpdateArticle godoc
// @Summary 更新文章
// @Description 管理员按文章 ID 更新文章内容和标签。
// @Tags admin/article
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "文章 ID"
// @Param body body UpdateArticleRequest true "更新文章请求体。需要先准备 article_id、category_id 和 tag_ids"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"id\":1},\"msg\":\"修改文章成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"文章ID不合法\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"修改文章失败\"}"
// @Router /api/admin/articles/update/{id} [put]
func UpdateArticle(c *gin.Context) {
	//从路径参数中获取文章ID
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章ID不合法",
		})
		return
	}

	id := uint32(id64)

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误",
		})
		return

	}

	//校验置顶信息
	if req.IsTop != 0 && req.IsTop != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "置顶信息错误",
		})
		return
	}

	//校验tag是否存在
	ok, err := models.CheckTagsExist(req.TagIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "校验标签失败",
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "存在无效标签",
		})
		return
	}

	//校验分类是否存在
	categoryOK, err := models.CheckCategoryExist(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "校验分类失败",
		})
		return
	}
	if !categoryOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "分类不存在",
		})
		return
	}

	article := models.Article{
		Title:      req.Title,
		CategoryID: req.CategoryID,
		Summary:    req.Summary,
		Content:    req.Content,
		IsTop:      req.IsTop,
		CoverURL:   req.CoverURL,
	}

	if err := models.UpdateArticle(uint32(id), &article, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "修改文章失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id": article.ID,
		},
		"msg": "修改文章成功",
	})

}

// UpdateArticleStatus godoc
// @Summary 更新文章状态
// @Description 管理员发布或下架文章。
// @Tags admin/article
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "文章 ID"
// @Param body body UpdateArticleStatusRequest true "文章状态请求体。status 示例：1 表示发布，2 表示下架"
// @Success 200 {object} object "示例：{\"code\":200,\"msg\":\"文章发布成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"文章状态不合法,只允许发布或下架\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"更新文章状态失败\"}"
// @Router /api/admin/articles/updateStatus/{id} [put]
func UpdateArticleStatus(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章ID不合法",
		})
		return
	}

	id := uint32(id64)

	var req UpdateArticleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误",
		})
		return
	}

	status := *req.Status

	if status != 0 && status != 1 && status != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章状态不合法,只允许草稿、发布或下架",
		})
		return
	}

	if err := models.UpdateArticleStatus(id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新文章状态失败",
		})
		return
	}

	msg := "文章状态更新成功"
	if status == 0 {
		msg = "文章已经存入草稿箱，游客无法查看"
	} else if status == 1 {
		msg = "文章发布成功"
	} else {
		msg = "文章下架成功，游客无法查看"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  msg,
	})

}
