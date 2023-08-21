package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type Category interface {
	CountCategory(ctx context.Context, tx *gorm.DB, condition FindCategoriesArg) (int64, error)
	FindCategory(ctx context.Context, tx *gorm.DB, condition FindCategoriesArg) ([]model.BlogCategories, error)
	GetCategories(ctx context.Context, tx *gorm.DB, cateId uint) (model.BlogCategories, error)
	CreateCategory(ctx context.Context, tx *gorm.DB, cate model.BlogCategories) error
	UpdateCategory(ctx context.Context, tx *gorm.DB, cate model.BlogCategories) error
	DeleteCategory(ctx context.Context, tx *gorm.DB, cateId uint) error
}

type FindCategoriesArg struct {
	IDs            []uint
	Names          []string
	GeneralCateIds []uint
	Offset         int32
	Limit          int32
	NoLimit        bool
}

type category struct{}

func NewCategory() Category {
	return &category{}
}

func (c *category) CountCategory(ctx context.Context, tx *gorm.DB, condition FindCategoriesArg) (int64, error) {

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

	if condition.GeneralCateIds != nil {
		query.WriteString(fmt.Sprintf(" AND `general_id` IN (%s)", RepeatWithSep("?", len(condition.GeneralCateIds), ",")))
		for _, generalCateId := range condition.GeneralCateIds {
			args = append(args, generalCateId)
		}
	}

	var count int64
	if err := tx.Model(&model.BlogCategories{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (c *category) FindCategory(ctx context.Context, tx *gorm.DB, condition FindCategoriesArg) ([]model.BlogCategories, error) {
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

	if condition.GeneralCateIds != nil {
		query.WriteString(fmt.Sprintf(" AND `general_id` IN (%s)", RepeatWithSep("?", len(condition.GeneralCateIds), ",")))
		for _, generalCateId := range condition.GeneralCateIds {
			args = append(args, generalCateId)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `id` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var categories []model.BlogCategories
	if err := tx.Where(query.String(), args...).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *category) GetCategories(ctx context.Context, tx *gorm.DB, cateId uint) (model.BlogCategories, error) {
	var category model.BlogCategories
	if err := tx.Where("id = ?", cateId).First(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (c *category) CreateCategory(ctx context.Context, tx *gorm.DB, cate model.BlogCategories) error {
	if err := tx.Create(&cate).Error; err != nil {
		return err
	}
	return nil
}

func (c *category) UpdateCategory(ctx context.Context, tx *gorm.DB, cate model.BlogCategories) error {
	if err := tx.Model(&cate).Where("id = ?", cate.ID.ID).Updates(&cate).Error; err != nil {
		return err
	}
	return nil
}

func (c *category) DeleteCategory(ctx context.Context, tx *gorm.DB, cateId uint) error {
	if err := tx.Where("id = ?", cateId).Delete(&model.BlogCategories{}).Error; err != nil {
		return err
	}
	return nil
}
