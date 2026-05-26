package models

import "MyBlog/global"

type DashboardStats struct {
	ArticleCount        int64     `json:"article_count"`
	PublishedCount      int64     `json:"published_count"`
	DraftCount          int64     `json:"draft_count"`
	OfflineCount        int64     `json:"offline_count"`
	DeletedArticleCount int64     `json:"deleted_article_count"`
	CategoryCount       int64     `json:"category_count"`
	TagCount            int64     `json:"tag_count"`
	CommentCount        int64     `json:"comment_count"`
	GuestbookCount      int64     `json:"guestbook_count"`
	FriendlinkCount     int64     `json:"friendlink_count"`
	RecentArticles      []Article `json:"recent_articles"`
}

func CountArticles(status int8, isDeleted int8) (int64, error) {
	var total int64
	query := global.DB.Model(&Article{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if isDeleted >= 0 {
		query = query.Where("is_deleted = ?", isDeleted)
	}
	err := query.Count(&total).Error
	return total, err
}

func CountCategories() (int64, error) {
	var total int64
	err := global.DB.Model(&Category{}).Count(&total).Error
	return total, err
}

func CountTags() (int64, error) {
	var total int64
	err := global.DB.Model(&Tag{}).Count(&total).Error
	return total, err
}

func CountCommentsByTargetType(targetType int8) (int64, error) {
	var total int64
	err := global.DB.Model(&Comment{}).
		Where("target_type = ?", targetType).
		Count(&total).Error
	return total, err
}

func CountFriendLinks() (int64, error) {
	var total int64
	err := global.DB.Model(&FriendLink{}).Count(&total).Error
	return total, err
}

func GetRecentArticles(limit int) ([]Article, error) {
	var articles []Article
	if limit <= 0 {
		limit = 5
	}
	err := global.DB.Model(&Article{}).
		Select("id", "title", "summary", "cover_url", "category_id", "comment_count", "status", "is_top", "is_deleted", "published_time", "create_time", "update_time").
		Order("update_time DESC, id DESC").
		Limit(limit).
		Find(&articles).Error
	return articles, err
}
