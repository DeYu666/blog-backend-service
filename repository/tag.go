package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type Tag interface {
	CountTags(ctx context.Context, tx *gorm.DB, condition FindTagArg) (int64, error)
	FindTags(ctx context.Context, tx *gorm.DB, condition FindTagArg) ([]model.BlogTag, error)
	GetTag(ctx context.Context, tx *gorm.DB, tagId uint) (model.BlogTag, error)
	CreateTag(ctx context.Context, tx *gorm.DB, tag model.BlogTag) error
	UpdateTag(ctx context.Context, tx *gorm.DB, tag model.BlogTag) error
	DeleteTag(ctx context.Context, tx *gorm.DB, tagId uint) error
}

type tag struct{}

func NewTag() Tag {
	return &tag{}
}

type FindTagArg struct {
	IDs     []uint
	Names   []string
	Offset  int32
	Limit   int32
	NoLimit bool
}

func (t *tag) CountTags(ctx context.Context, tx *gorm.DB, condition FindTagArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(" 1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.Names != nil {
		query.WriteString(fmt.Sprintf(" AND `name` IN (%s)", RepeatWithSep("?", len(condition.Names), ",")))
		for _, name := range condition.Names {
			args = append(args, name)
		}
	}

	var count int64

	err := tx.Model(&model.BlogTag{}).Where(query.String(), args...).Count(&count).Error
	return count, err
}

func (t *tag) FindTags(ctx context.Context, tx *gorm.DB, condition FindTagArg) ([]model.BlogTag, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.Names != nil {
		query.WriteString(fmt.Sprintf(" AND `name` IN (%s)", RepeatWithSep("?", len(condition.Names), ",")))
		for _, name := range condition.Names {
			args = append(args, name)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `id` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var tags []model.BlogTag
	err := tx.Model(&model.BlogTag{}).Where(query.String(), args...).Find(&tags).Error
	return tags, err
}

func (t *tag) GetTag(ctx context.Context, tx *gorm.DB, tagId uint) (model.BlogTag, error) {
	var tag model.BlogTag
	err := tx.Model(&model.BlogTag{}).Where("`id` = ?", tagId).First(&tag).Error
	return tag, err
}

func (t *tag) CreateTag(ctx context.Context, tx *gorm.DB, tag model.BlogTag) error {
	return tx.Model(&model.BlogTag{}).Create(&tag).Error
}

func (t *tag) UpdateTag(ctx context.Context, tx *gorm.DB, tag model.BlogTag) error {
	return tx.Model(&model.BlogTag{}).Where("`id` = ?", tag.ID.ID).Updates(&tag).Error
}

func (t *tag) DeleteTag(ctx context.Context, tx *gorm.DB, tagId uint) error {
	return tx.Model(&model.BlogTag{}).Where("`id` = ?", tagId).Delete(&model.BlogTag{}).Error
}
