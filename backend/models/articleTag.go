package models

type ArticleTag struct {
	ArticleID uint32 `gorm:"column:article_id;primaryKey" json:"article_id"`
	TagID     uint32 `gorm:"column:tag_id;primaryKey" json:"tag_id"`
}

func (ArticleTag) TableName() string {
	return "article_tag"
}
