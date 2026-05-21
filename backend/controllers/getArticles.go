package controllers

import (
	"MyBlog/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetArticlesTimeline(c *gin.Context) {
	pageNum, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		pageNum = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		pageSize = 10
	}

	articles, total, err := models.GetArticleTimeline(pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取文章时间轴失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  articles,
			"total": total,
		},
		"msg": "获取文章时间轴成功",
	})

}

func GetArticlesByCategoryID(c *gin.Context) {
	categoryID64, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "分类ID不合法",
		})
		return
	}
	categoryID := uint32(categoryID64)

	pageNum, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		pageNum = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		pageSize = 10
	}

	//category>0时才校验分类是否存在
	if categoryID > 0 {
		ok, err := models.CheckCategoryExist(categoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "校验分类失败",
			})
			return
		}
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "分类不存在或已停用",
			})
			return
		}
	}

	articles, total, err := models.GetArticlesByCategoryID(categoryID, pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取文章列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  articles,
			"total": total,
		},
		"msg": "获取文章列表成功",
	})
}

func GetArticlesByTagID(c *gin.Context) {
	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "标签ID不合法",
		})
		return
	}

	tagID := uint32(id64)

	pageNum, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		pageNum = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		pageSize = 10
	}

	ok, err := models.CheckTagsExist([]uint32{tagID})
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
			"msg":  "标签不存在或已停用",
		})
		return
	}

	articles, total, err := models.GetArticlesByTagID(tagID, pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取标签文章列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  articles,
			"total": total,
		},
		"msg": "获取标签文章列表成功",
	})

}

func GetAdminArticles(c *gin.Context) {
	// page 默认 1
	pageNum, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	// pageSize 默认 10
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// status 默认 -1，表示全部状态
	status, err := strconv.Atoi(c.DefaultQuery("status", "-1"))
	if err != nil {
		status = -1
	}

	// is_deleted 默认 0，表示默认查未删除文章
	isDeleted, err := strconv.Atoi(c.DefaultQuery("is_deleted", "0"))
	if err != nil {
		isDeleted = 0
	}

	// category_id 默认 0，表示全部分类
	categoryID64, err := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 32)
	if err != nil {
		categoryID64 = 0
	}
	categoryID := uint32(categoryID64)

	keyword := c.Query("keyword")

	//参数合法性校验
	if status != -1 && status != 0 && status != 1 && status != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章状态参数错误",
		})
		return
	}

	if isDeleted != -1 && isDeleted != 0 && isDeleted != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "删除状态参数错误",
		})
		return
	}

	articles, total, err := models.GetAdminArticles(pageNum, pageSize, status, isDeleted, categoryID, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取后台文章列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  articles,
			"total": total,
			"page":  pageNum,
			"size":  pageSize,
		},
		"msg": "获取后台文章列表成功",
	})

}
