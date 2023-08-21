package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type GeneralCate interface {
	CountGeneralCate(ctx context.Context, tx *gorm.DB, condition FindGeneralCateArg) (int64, error)
	FindGeneralCate(ctx context.Context, tx *gorm.DB, condition FindGeneralCateArg) ([]model.BlogGeneralCategories, error)
	GetGeneralCate(ctx context.Context, tx *gorm.DB, generalCateId uint) (model.BlogGeneralCategories, error)
	CreateGeneralCate(ctx context.Context, tx *gorm.DB, cate model.BlogGeneralCategories) error
	UpdateGeneralCate(ctx context.Context, tx *gorm.DB, cate model.BlogGeneralCategories) error
	DeleteGeneralCate(ctx context.Context, tx *gorm.DB, generalCateId uint) error
}

type FindGeneralCateArg struct {
	GeneralCateIds []uint
	Names          []string
	Offset         int32
	Limit          int32
	NoLimit        bool
}

type generalCate struct{}

func NewGeneralCate() GeneralCate {
	return &generalCate{}
}

func (c *generalCate) CountGeneralCate(ctx context.Context, tx *gorm.DB, condition FindGeneralCateArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(" 1 = 1")
	if condition.GeneralCateIds != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.GeneralCateIds), ",")))
		for _, id := range condition.GeneralCateIds {
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
	if err := tx.Model(&model.BlogGeneralCategories{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (c *generalCate) FindGeneralCate(ctx context.Context, tx *gorm.DB, condition FindGeneralCateArg) ([]model.BlogGeneralCategories, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.GeneralCateIds != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.GeneralCateIds), ",")))
		for _, id := range condition.GeneralCateIds {
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

	var cates []model.BlogGeneralCategories
	if err := tx.Model(&model.BlogGeneralCategories{}).Where(query.String(), args...).Find(&cates).Error; err != nil {
		return nil, err
	}
	return cates, nil
}

func (c *generalCate) GetGeneralCate(ctx context.Context, tx *gorm.DB, generalCateId uint) (model.BlogGeneralCategories, error) {
	var cate model.BlogGeneralCategories
	if err := tx.Where("`id` = ?", generalCateId).First(&cate).Error; err != nil {
		return cate, err
	}
	return cate, nil
}

func (c *generalCate) CreateGeneralCate(ctx context.Context, tx *gorm.DB, cate model.BlogGeneralCategories) error {
	if err := tx.Create(&cate).Error; err != nil {
		return err
	}
	return nil
}

func (c *generalCate) UpdateGeneralCate(ctx context.Context, tx *gorm.DB, cate model.BlogGeneralCategories) error {
	if err := tx.Model(&model.BlogGeneralCategories{}).Where("`id` = ?", cate.ID.ID).Updates(cate).Error; err != nil {
		return err
	}
	return nil
}

func (c *generalCate) DeleteGeneralCate(ctx context.Context, tx *gorm.DB, generalCateId uint) error {
	if err := tx.Where("`id` = ?", generalCateId).Delete(&model.BlogGeneralCategories{}).Error; err != nil {
		return err
	}
	return nil
}
