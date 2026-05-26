package models

import (
	"MyBlog/global"
	"MyBlog/utils"
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Comment struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetType   int8      `gorm:"not null" json:"target_type"` //1文章 2留言板
	TargetID     uint64    `gorm:"default:0" json:"target_id"`
	RootID       uint64    `gorm:"not null;default:0" json:"root_id"`
	ParentID     uint64    `gorm:"not null;default:0" json:"parent_id"`
	Nickname     string    `gorm:"size:50;not null" json:"nickname"`
	Email        string    `gorm:"size:100" json:"email"`
	Website      string    `gorm:"size:255" json:"website"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	Content      string    `gorm:"size:1000;not null" json:"content"`
	IP           string    `gorm:"size:64" json:"-"`
	IPLocation   string    `gorm:"column:ip_location;size:100" json:"ip_location"`
	UserAgent    string    `gorm:"size:500" json:"user_agent"`
	IsAdminReply int8      `gorm:"not null;default:0" json:"is_admin_reply"`
	IsDeleted    int8      `gorm:"not null;default:0" json:"is_deleted"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`

	ReplyToNickname string    `gorm:"-" json:"reply_to_nickname,omitempty"`
	Browser         string    `gorm:"-" json:"browser,omitempty"`
	OS              string    `gorm:"-" json:"os,omitempty"`
	Children        []Comment `gorm:"foreignKey:ParentID" json:"children"` //用于楼中楼嵌套
}

type AdminCommentItem struct {
	ID           uint64    `json:"id"`
	TargetType   int8      `json:"target_type"`
	TargetID     uint64    `json:"target_id"`
	ArticleTitle string    `json:"article_title"`
	ParentID     uint64    `json:"parent_id"`
	Nickname     string    `json:"nickname"`
	Email        string    `json:"email"`
	Website      string    `json:"website"`
	Content      string    `json:"content"`
	IP           string    `json:"ip"`
	IPLocation   string    `json:"ip_location"`
	UserAgent    string    `json:"user_agent"`
	IsDeleted    int8      `json:"is_deleted"`
	CreateTime   time.Time `json:"create_time"`
}

type CommentValidationError struct {
	Message string
}

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func (err *CommentValidationError) Error() string {
	return err.Message
}

func IsCommentValidationError(err error) bool {
	var validationErr *CommentValidationError
	return errors.As(err, &validationErr)
}

func NormalizeVisitorEmail(email string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(email))
	if value == "" {
		return "", &CommentValidationError{Message: "请填写邮箱"}
	}
	if !emailPattern.MatchString(value) {
		return "", &CommentValidationError{Message: "邮箱格式不正确"}
	}
	return value, nil
}

func FindVisitorNicknameByEmail(email string) (string, bool, error) {
	value, err := NormalizeVisitorEmail(email)
	if err != nil {
		return "", false, err
	}

	var identity VisitorIdentity
	err = global.DB.Where("email = ?", value).First(&identity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, err
	}
	if err == nil {
		return identity.Nickname, true, nil
	}

	var legacy Comment
	err = global.DB.
		Where("LOWER(email) = ?", value).
		Where("email <> ''").
		Order("id ASC").
		First(&legacy).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, err
	}
	if err == nil {
		nickname := strings.TrimSpace(legacy.Nickname)
		return nickname, nickname != "", nil
	}

	return "", false, nil
}

// TableName
func (Comment) TableName() string {
	return "comment"
}

// 生成Gravatar头像URL
func GenerateGravatar(email string) string {
	return ""
}

// 获取树形评论
func GetComments(targetType int8, targetID uint64, pageNum, pageSize int) ([]Comment, int64, error) {
	var topComments []Comment
	var total int64

	//一级评论
	db := global.DB.Model(&Comment{}).
		Where("target_type = ? and target_id=? and parent_id=0 and is_deleted=0", targetType, targetID)

	db.Count(&total)

	err := db.Order("create_time ASC").
		Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&topComments).Error
	if err != nil {
		return nil, 0, err
	}

	topIDs := make([]uint64, 0, len(topComments))
	for i := range topComments {
		topIDs = append(topIDs, topComments[i].ID)
		enrichCommentPublicMeta(&topComments[i])
	}

	if len(topIDs) == 0 {
		return topComments, total, nil
	}

	var replies []Comment
	err = global.DB.
		Where("target_type = ? and target_id = ? and parent_id <> 0 and is_deleted = 0", targetType, targetID).
		Order("create_time ASC").
		Find(&replies).Error
	if err != nil {
		return nil, 0, err
	}

	topIndexByID := make(map[uint64]int, len(topComments))
	commentByID := make(map[uint64]*Comment, len(topComments)+len(replies))
	for i := range topComments {
		topIndexByID[topComments[i].ID] = i
		commentByID[topComments[i].ID] = &topComments[i]
	}
	for i := range replies {
		commentByID[replies[i].ID] = &replies[i]
	}

	for i := range replies {
		enrichCommentPublicMeta(&replies[i])
		if parent := commentByID[replies[i].ParentID]; parent != nil {
			replies[i].ReplyToNickname = parent.Nickname
		}

		rootID := replies[i].RootID
		if rootID == 0 {
			rootID = replies[i].ParentID
		}
		if _, ok := topIndexByID[rootID]; !ok {
			if parent := commentByID[replies[i].ParentID]; parent != nil && parent.RootID != 0 {
				rootID = parent.RootID
			}
		}
		if topIndex, ok := topIndexByID[rootID]; ok {
			topComments[topIndex].Children = append(topComments[topIndex].Children, replies[i])
		}
	}

	return topComments, total, nil
}

func enrichCommentPublicMeta(comment *Comment) {
	if comment.Avatar == "" {
		comment.Avatar = GenerateGravatar(comment.Email)
	}
	if comment.IPLocation == "" {
		comment.IPLocation = utils.ResolveIPLocation(comment.IP)
	}
	comment.Browser = utils.ParseBrowser(comment.UserAgent)
	comment.OS = utils.ParseOS(comment.UserAgent)
}

// 添加评论
func AddComment(comment *Comment) error {
	email, err := NormalizeVisitorEmail(comment.Email)
	if err != nil {
		return err
	}
	comment.Email = email
	comment.Nickname = strings.TrimSpace(comment.Nickname)
	comment.Content = strings.TrimSpace(comment.Content)

	if comment.Content == "" {
		return &CommentValidationError{Message: "请填写内容"}
	}

	if comment.Avatar == "" {
		comment.Avatar = GenerateGravatar(comment.Email)
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if comment.ParentID == 0 {
			comment.RootID = 0
		} else {
			var parent Comment
			if err := tx.
				Where("id = ? and target_type = ? and target_id = ? and is_deleted = 0", comment.ParentID, comment.TargetType, comment.TargetID).
				First(&parent).Error; err != nil {
				return &CommentValidationError{Message: "回复的评论不存在"}
			}
			comment.RootID = parent.RootID
			if comment.RootID == 0 {
				comment.RootID = parent.ID
			}
		}

		if err := ensureVisitorIdentity(tx, comment); err != nil {
			return err
		}
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		if comment.TargetType == 1 {
			return updateArticleCommentCount(tx, comment.TargetID, 1)
		}
		return nil
	})
}

func ensureVisitorIdentity(tx *gorm.DB, comment *Comment) error {
	var identity VisitorIdentity
	err := tx.Where("email = ?", comment.Email).First(&identity).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err == nil {
		comment.Nickname = strings.TrimSpace(identity.Nickname)
		if comment.Nickname == "" {
			return &CommentValidationError{Message: "请填写昵称"}
		}
		return tx.Model(&VisitorIdentity{}).
			Where("id = ?", identity.ID).
			Updates(map[string]interface{}{
				"last_ip":    comment.IP,
				"user_agent": comment.UserAgent,
				"avatar":     comment.Avatar,
			}).Error
	}

	var legacy Comment
	legacyErr := tx.
		Where("LOWER(email) = ?", comment.Email).
		Where("email <> ''").
		Order("id ASC").
		First(&legacy).Error
	if legacyErr != nil && !errors.Is(legacyErr, gorm.ErrRecordNotFound) {
		return legacyErr
	}

	boundNickname := comment.Nickname
	firstIP := comment.IP
	if legacyErr == nil {
		boundNickname = strings.TrimSpace(legacy.Nickname)
		firstIP = legacy.IP
	}
	if boundNickname == "" {
		return &CommentValidationError{Message: "请填写昵称"}
	}
	comment.Nickname = boundNickname

	identity = VisitorIdentity{
		Email:     comment.Email,
		Nickname:  boundNickname,
		Avatar:    comment.Avatar,
		FirstIP:   firstIP,
		LastIP:    comment.IP,
		UserAgent: comment.UserAgent,
	}
	return tx.Create(&identity).Error
}

// 管理员删除评论
func DeleteComment(commentID uint64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var comment Comment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? and is_deleted = 0", commentID).
			First(&comment).Error; err != nil {
			return err
		}

		deleteIDs := []uint64{comment.ID}
		deletedCount := int64(1)
		if comment.TargetType == 1 && comment.ParentID == 0 {
			var childIDs []uint64
			if err := tx.Model(&Comment{}).
				Where("target_type = ? and target_id = ? and root_id = ? and is_deleted = 0", comment.TargetType, comment.TargetID, comment.ID).
				Pluck("id", &childIDs).Error; err != nil {
				return err
			}
			deleteIDs = append(deleteIDs, childIDs...)
			deletedCount += int64(len(childIDs))
		}

		if err := tx.Model(&Comment{}).
			Where("id IN ?", deleteIDs).
			Update("is_deleted", 1).Error; err != nil {
			return err
		}
		if comment.TargetType == 1 {
			return updateArticleCommentCount(tx, comment.TargetID, -deletedCount)
		}
		return nil
	})
}

func RestoreComment(commentID uint64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var comment Comment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? and is_deleted = 1", commentID).
			First(&comment).Error; err != nil {
			return err
		}

		if err := tx.Model(&Comment{}).
			Where("id = ? and is_deleted = 1", commentID).
			Update("is_deleted", 0).Error; err != nil {
			return err
		}
		if comment.TargetType == 1 {
			return updateArticleCommentCount(tx, comment.TargetID, 1)
		}
		return nil
	})
}

func updateArticleCommentCount(tx *gorm.DB, articleID uint64, delta int64) error {
	if delta == 0 {
		return nil
	}

	expr := gorm.Expr("comment_count + ?", delta)
	if delta < 0 {
		expr = gorm.Expr("GREATEST(comment_count - ?, 0)", -delta)
	}

	result := tx.Model(&Article{}).
		Where("id = ?", articleID).
		Update("comment_count", expr)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &CommentValidationError{Message: "文章不存在"}
	}
	return nil
}

func GetAdminComments(targetType int8, targetID uint64, includeDeleted bool, pageNum, pageSize int) ([]AdminCommentItem, int64, error) {
	var comments []AdminCommentItem
	var total int64

	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	query := global.DB.Table("comment AS c").
		Select("c.id", "c.target_type", "c.target_id", "article.title AS article_title", "c.parent_id", "c.nickname", "c.email", "c.website", "c.content", "c.ip", "c.ip_location", "c.user_agent", "c.is_deleted", "c.create_time").
		Joins("LEFT JOIN article ON c.target_type = 1 AND c.target_id = article.id")

	if targetType > 0 {
		query = query.Where("c.target_type = ?", targetType)
	}
	if targetID > 0 {
		query = query.Where("c.target_id = ?", targetID)
	}
	if !includeDeleted {
		query = query.Where("c.is_deleted = 0")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("c.create_time DESC, c.id DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Scan(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	for i := range comments {
		if comments[i].IPLocation == "" {
			comments[i].IPLocation = utils.ResolveIPLocation(comments[i].IP)
		}
	}

	return comments, total, nil
}
