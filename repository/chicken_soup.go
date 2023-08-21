package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type ChickenSoup interface {
	CountChickenSoups(ctx context.Context, tx *gorm.DB, condition FindChickenSoupArg) (int64, error)
	FindChickenSoups(ctx context.Context, tx *gorm.DB, condition FindChickenSoupArg) ([]model.ChickenSoup, error)
	GetChickenSoups(ctx context.Context, tx *gorm.DB, checkSoupId uint) (model.ChickenSoup, error)
	CreateChickenSoup(ctx context.Context, tx *gorm.DB, checkSoup model.ChickenSoup) error
	UpdateChickenSoup(ctx context.Context, tx *gorm.DB, checkSoup model.ChickenSoup) error
	DeleteChickenSoup(ctx context.Context, tx *gorm.DB, checkSoupId uint) error
}

type FindChickenSoupArg struct {
	IDs      []uint
	Sentence []string
	Offset   int32
	Limit    int32
	NoLimit  bool
}

type chickenSoup struct{}

func NewChickenSoup() ChickenSoup {
	return &chickenSoup{}
}

func (c *chickenSoup) CountChickenSoups(ctx context.Context, tx *gorm.DB, condition FindChickenSoupArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(" 1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.Sentence != nil {
		query.WriteString(fmt.Sprintf(" AND `sentence` IN (%s)", RepeatWithSep("?", len(condition.Sentence), ",")))
		for _, s := range condition.Sentence {
			args = append(args, s)
		}
	}

	var count int64
	if err := tx.Model(&model.ChickenSoup{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (c *chickenSoup) FindChickenSoups(ctx context.Context, tx *gorm.DB, condition FindChickenSoupArg) ([]model.ChickenSoup, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.Sentence != nil {
		query.WriteString(fmt.Sprintf(" AND `sentence` IN (%s)", RepeatWithSep("?", len(condition.Sentence), ",")))
		for _, s := range condition.Sentence {
			args = append(args, s)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `id` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var chickenSoups []model.ChickenSoup
	if err := tx.Model(&model.ChickenSoup{}).Where(query.String(), args...).Find(&chickenSoups).Error; err != nil {
		return nil, err
	}
	return chickenSoups, nil
}

func (c *chickenSoup) GetChickenSoups(ctx context.Context, tx *gorm.DB, checkSoupId uint) (model.ChickenSoup, error) {
	var chickenSoup model.ChickenSoup
	if err := tx.Model(&model.ChickenSoup{}).Where("`id` = ?", checkSoupId).First(&chickenSoup).Error; err != nil {
		return chickenSoup, err
	}
	return chickenSoup, nil
}

func (c *chickenSoup) CreateChickenSoup(ctx context.Context, tx *gorm.DB, checkSoup model.ChickenSoup) error {
	if err := tx.Model(&model.ChickenSoup{}).Create(&checkSoup).Error; err != nil {
		return err
	}
	return nil
}

func (c *chickenSoup) UpdateChickenSoup(ctx context.Context, tx *gorm.DB, checkSoup model.ChickenSoup) error {
	if err := tx.Model(&model.ChickenSoup{}).Where("`id` = ?", checkSoup.ID.ID).Updates(&checkSoup).Error; err != nil {
		return err
	}
	return nil
}

func (c *chickenSoup) DeleteChickenSoup(ctx context.Context, tx *gorm.DB, checkSoupId uint) error {
	if err := tx.Model(&model.ChickenSoup{}).Where("`id` = ?", checkSoupId).Delete(&model.ChickenSoup{}).Error; err != nil {
		return err
	}
	return nil
}
