package controllers

import (
	"MyBlog/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateArticleRequest struct {
	Title      string   `json:"title" binding:"required"`
	CategoryID uint32   `json:"category_id" binding:"required"`
	Summary    string   `json:"summary" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	IsTop      int8     `json:"is_top" binding:"required"`
	CoverURL   string   `json:"cover_url" binding:"required"`
	TagIDs     []uint32 `json:"tag_ids"`
}

func CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误",
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

	//校验置顶信息
	if req.IsTop != 0 && req.IsTop != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "置顶信息错误",
		})
		return
	}

	article := models.Article{
		Title:        req.Title,
		CategoryID:   req.CategoryID,
		Summary:      req.Summary,
		Content:      req.Content,
		Status:       0,
		IsDeleted:    0,
		IsTop:        req.IsTop,
		CoverURL:     req.CoverURL,
		CommentCount: 0,
	}

	if err := models.CreateArticle(&article, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建文章失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id": article.ID,
		},
		"msg": "创建文章成功",
	})

}
