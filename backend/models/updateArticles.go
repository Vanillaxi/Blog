package models

import (
	"MyBlog/global"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Article struct {
	ID       uint32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `gorm:"not null" json:"title"`
	Content  string `gorm:"type:longtext;not null" json:"content"`
	Summary  string `json:"summary"`
	CoverURL string `gorm:"not null" json:"cover_url"`

	CategoryID uint32 `gorm:"not null" json:"category_id"`

	CommentCount uint32 `gorm:"not null" json:"comment_count"`

	Status    int8 `gorm:"not null;default:1" json:"status"`
	IsTop     int8 `gorm:"not null;default:0" json:"is_top"`
	IsDeleted int8 `gorm:"not null;default:0" json:"is_deleted"`

	PublishedTime *time.Time `gorm:"column:published_time" json:"published_time"`
	CreateTime    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"update_time"`
}

func CreateArticle(article *Article, tagIDs []uint32) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		//1.创建文章
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		//2. 没有标签就直接返回
		if len(tagIDs) == 0 {
			return nil
		}

		//3. 插入article_tag关系
		articleTags := make([]ArticleTag, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			articleTags = append(articleTags, ArticleTag{
				ArticleID: article.ID,
				TagID:     tagID,
			})
		}

		if err := tx.Create(&articleTags).Error; err != nil {
			return err
		}

		return nil
	})
}

func UpdateArticle(id uint32, article *Article, tagIDs []uint32) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"title":       article.Title,
			"summary":     article.Summary,
			"cover_url":   article.CoverURL,
			"category_id": article.CategoryID,
			"content":     article.Content,
			"is_top":      article.IsTop,
		}

		if err := tx.Model(&Article{}).
			Where("id=? and is_deleted=?", id, 0).
			Updates(updates).Error; err != nil {
			return err
		}

		//删除旧标签关系
		if err := tx.Where("article_id=?", id).
			Delete(&ArticleTag{}).Error; err != nil {
			return err
		}

		//插入新的标签关系
		if len(tagIDs) > 0 {
			articleTags := make([]ArticleTag, 0, len(tagIDs))
			for _, tagID := range tagIDs {
				articleTags = append(articleTags, ArticleTag{
					ArticleID: id,
					TagID:     tagID,
				})
			}

			if err := tx.Create(&articleTags).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func UpdateArticleStatus(id uint32, newStatus int8) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var article Article

		//1. 查询文章，并加行锁，防止并发重复发布/下架导致count错乱
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id=? and is_deleted=?", id, 0).
			First(&article).Error; err != nil {
			return err
		}

		oldStatus := article.Status

		//2.如果状态没变化，直接返回
		if oldStatus == newStatus {
			return nil
		}

		//3.查询文章绑定标签
		var tagIDs []uint32
		if err := tx.Model(&ArticleTag{}).
			Where("article_id=?", id).
			Pluck("tag_id", &tagIDs).Error; err != nil {
			return err
		}

		switch {
		//草稿->发布
		case oldStatus == 0 && newStatus == 1:
			now := time.Now()

			//更新文章状态
			if err := tx.Model(&Article{}).
				Where("id=?", id).
				Updates(map[string]interface{}{
					"status":         1,
					"published_time": gorm.Expr("IFNULL(published_time,?)", now),
				}).Error; err != nil {
				return err
			}

			if err := increaseArticleCount(tx, article.CategoryID, tagIDs); err != nil {
				return err
			}

		//发布->下架
		case oldStatus == 1 && newStatus == 2:
			if err := tx.Model(&Article{}).
				Where("id=? and is_deleted=?", id, 0).
				Update("status", 2).Error; err != nil {
				return err
			}

			if err := decreaseArticleCount(tx, article.CategoryID, tagIDs); err != nil {
				return err
			}

		//下架->发布
		case oldStatus == 2 && newStatus == 1:
			if err := tx.Model(&Article{}).
				Where("id=? and is_deleted=?", id, 0).
				Update("status", 1).Error; err != nil {
				return err
			}

			if err := increaseArticleCount(tx, article.CategoryID, tagIDs); err != nil {
				return err
			}

		default:
			return errors.New("非法文章状态流转")
		}

		return nil

	})
}

func DeleteArticle(id uint32) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var article Article

		//1.查文章，并加行锁，防止重复删除导致article_count重复扣
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id=? and is_deleted=?", id, 0).
			First(&article).Error; err != nil {
			return err
		}

		//2.查文章绑定的标签
		var tagIDs []uint32
		if err := tx.Model(&ArticleTag{}).
			Where("article_id=?", id).
			Pluck("tag_id", &tagIDs).Error; err != nil {
			return err
		}

		//3.如果已发布
		if article.Status == 1 {
			if err := decreaseArticleCount(tx, article.CategoryID, tagIDs); err != nil {
				return err
			}
		}

		//4.逻辑删除
		if err := tx.Model(&Article{}).
			Where("id=? and is_deleted=?", id, 0).
			Update("is_deleted", 1).Error; err != nil {
			return err
		}

		return nil
	})
}

// 辅助函数
func increaseArticleCount(tx *gorm.DB, categoryID uint32, tagIDs []uint32) error {
	//分类文章数+1
	if err := tx.Model(&Category{}).
		Where("id=?", categoryID).
		Update("article_count", gorm.Expr("article_count+1")).Error; err != nil {
		return err
	}

	//标签文章数+1
	if len(tagIDs) > 0 {
		if err := tx.Model(&Tag{}).
			Where("id IN ?", tagIDs).
			Update("article_count", gorm.Expr("article_count+1")).Error; err != nil {
			return err
		}
	}
	return nil
}

func decreaseArticleCount(tx *gorm.DB, categoryID uint32, tagIDs []uint32) error {
	//分类文章数-1
	if err := tx.Model(&Category{}).
		Where("id=?", categoryID).
		Update("article_count", gorm.Expr("GREATEST(article_count-1,0)")).Error; err != nil {
		return err
	}

	//标签文章数-1
	if len(tagIDs) > 0 {
		if err := tx.Model(&Tag{}).
			Where("id IN ?", tagIDs).
			Update("article_count", gorm.Expr("GREATEST(article_count-1,0)")).Error; err != nil {
			return err
		}
	}
	return nil
}
