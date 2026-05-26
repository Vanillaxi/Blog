package controllers

import (
	"MyBlog/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateArticleRequest struct {
	Title      string   `json:"title" binding:"required" example:"我的第一篇博客"`
	CategoryID uint32   `json:"category_id" binding:"required" example:"1"`
	Summary    string   `json:"summary" example:"这是一篇关于 Go 后端开发的文章摘要"`
	Content    string   `json:"content" binding:"required" example:"这里是文章正文内容，包含 Gin、GORM 和 MySQL 的实践记录。"`
	IsTop      int8     `json:"is_top" example:"0"`
	CoverURL   string   `json:"cover_url" example:"https://example.com/cover.jpg"`
	TagIDs     []uint32 `json:"tag_ids" example:"1,2"`
}

// CreateArticle godoc
// @Summary 创建文章
// @Description 管理员创建文章，文章初始状态由业务逻辑设置。
// @Tags admin/article
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param body body CreateArticleRequest true "创建文章请求体。需要先准备 category_id 和 tag_ids，示例会创建一篇草稿文章"
// @Success 200 {object} object "示例：{\"code\":200,\"data\":{\"id\":1},\"msg\":\"创建文章成功\"}"
// @Failure 400 {object} object "示例：{\"code\":400,\"msg\":\"参数格式错误\"}"
// @Failure 500 {object} object "示例：{\"code\":500,\"msg\":\"创建文章失败\"}"
// @Router /api/admin/articles/create [post]
func CreateArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateArticle bind failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数格式错误: " + err.Error(),
		})
		return
	}

	log.Printf("CreateArticle payload: title=%q category_id=%d tag_ids=%v is_top=%d cover_url=%q content_len=%d", req.Title, req.CategoryID, req.TagIDs, req.IsTop, req.CoverURL, len(req.Content))

	//校验tag是否存在
	ok, err := models.CheckTagsExist(req.TagIDs)
	if err != nil {
		log.Printf("CreateArticle check tags failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "校验标签失败: " + err.Error(),
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
		log.Printf("CreateArticle check category failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "校验分类失败: " + err.Error(),
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
		log.Printf("CreateArticle create failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建文章失败: " + err.Error(),
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
