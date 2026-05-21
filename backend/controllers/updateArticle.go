package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateArticleRequest struct {
	Title      string   `json:"title" binding:"required"`
	CategoryID uint32   `json:"category_id" binding:"required"`
	Summary    string   `json:"summary" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	IsTop      int8     `json:"is_top" binding:"required"`
	CoverURL   string   `json:"cover_url" binding:"required"`
	TagIDs     []uint32 `json:"tag_ids"`
}

type UpdateArticleStatusRequest struct {
	Status *int8 `json:"status" binding:"required"` //要用指针，不然传0可能会认为没传，导致格式错误
}

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

	if status != 1 && status != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章状态不合法,只允许发布或下架",
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
