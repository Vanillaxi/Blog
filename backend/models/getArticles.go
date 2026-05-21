package models

import (
	"MyBlog/global"

	"gorm.io/gorm"
)

func (Article) TableName() string {
	return "article"
}

func GetArticleTimeline(pageNum int, pageSize int) ([]Article, int64, error) {
	var articles []Article
	var total int64

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	query := global.DB.Model(&Article{}).
		Where("status=? and is_deleted=? and published_time is not null", 1, 0)

	//查总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//按时间轴列表
	err := query.
		Select(
			"id",
			"title",
			"summary",
			"cover_url",
			"category_id",
			"comment_count",
			"is_top",
			"published_time",
			"create_time",
			"update_time").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Order("published_time desc,id desc").
		Find(&articles).Error

	if err != nil {
		return nil, 0, err
	}
	return articles, total, nil

}

func GetArticlesByCategoryID(categoryID uint32, pageNum int, pageSize int) ([]Article, int64, error) {
	var articles []Article
	var total int64

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	query := global.DB.Model(&Article{}).
		Where("status=? and is_deleted=?", 1, 0)

	//如果categoryID>0，则添加分类筛选条件
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	//先查总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//分页查询，不计算总数
	err := query.
		Select(
			"id",
			"title",
			"summary",
			"cover_url",
			"category_id",
			"comment_count",
			"is_top",
			"published_time",
			"create_time",
			"update_time").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Order("is_top desc,published_time desc,id desc").
		Find(&articles).Error

	if err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

func GetArticlesByTagID(tagID uint32, pageNum int, pageSize int) ([]Article, int64, error) {
	var articles []Article
	var total int64

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	baseQuery := global.DB.Table("article AS a").
		Joins("JOIN article_tag AS at ON a.id =at.article_id").
		Where("at.tag_id=? and a.status=? and a.deleted=?", tagID, 1, 0)

	//先查总数
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//分页查询
	err := baseQuery.
		Select(
			"a.id",
			"a.title",
			"a.summary",
			"a.cover_url",
			"a.category_id",
			"a.comment_count",
			"a.is_top",
			"a.published_time",
			"a.create_time",
			"a.update_time",
		).
		Order("a.is_top desc,a.published_time desc,a.id desc").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&articles).Error

	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil

}

// 只查已发布+未删除
func GetArticleDetail(id uint32) (Article, error) {
	var article Article
	err := global.DB.
		Where("id=? and status=? and is_deleted=?", id, 1, 0).
		First(&article).Error
	return article, err
}

// 后台查全部
func GetAdminArticles(pageNum int, pageSize int, status int, isDeleted int, categoryID uint32, keyword string) ([]Article, int64, error) {
	var articles []Article
	var total int64

	query := global.DB.Model(&Article{})

	//status=-1表示全部状态
	if status != -1 {
		query = query.Where("status=?", status)
	}

	//isDeleted=-1表示全部
	if isDeleted != -1 {
		query = query.Where("is_deleted=?", isDeleted)
	}

	// categoryID = 0 表示全部分类
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	//标题模糊查询
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	//总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再分页查询
	err := query.
		Select("id", "title", "summary", "cover_url", "category_id",
			"comment_count", "status", "is_top", "is_deleted",
			"published_time", "create_time", "update_time").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Order("update_time DESC, id DESC").
		Find(&articles).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return articles, total, nil
		}
		return nil, 0, err
	}

	return articles, total, nil

}
