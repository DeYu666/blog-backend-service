package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type LoveInfo interface {
	CountLoveInfos(ctx context.Context, tx *gorm.DB, condition FindLoveInfoArg) (int64, error)
	FindLoveInfos(ctx context.Context, tx *gorm.DB, condition FindLoveInfoArg) ([]model.LoveInfo, error)
	GetLoveInfo(ctx context.Context, tx *gorm.DB, id uint) (model.LoveInfo, error)
	CreateLoveInfo(ctx context.Context, tx *gorm.DB, loveInfo model.LoveInfo) error
	UpdateLoveInfo(ctx context.Context, tx *gorm.DB, loveInfo model.LoveInfo) error
	DeleteLoveInfo(ctx context.Context, tx *gorm.DB, id uint) error
}

type loveInfo struct{}

func NewLoveInfo() LoveInfo {
	return &loveInfo{}
}

type FindLoveInfoArg struct {
	IDs            []uint
	LoveNames      []string
	ExtraInfoLikes []string
	Offset         int
	Limit          int
	NoLimit        bool
}

func (b *loveInfo) CountLoveInfos(ctx context.Context, tx *gorm.DB, condition FindLoveInfoArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(" 1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.LoveNames != nil {
		query.WriteString(fmt.Sprintf(" AND `love_name` IN (%s)", RepeatWithSep("?", len(condition.LoveNames), ",")))
		for _, name := range condition.LoveNames {
			args = append(args, name)
		}
	}
	if condition.ExtraInfoLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `extra_info` Like (%s)", RepeatWithSep("?", len(condition.ExtraInfoLikes), ",")))
		for _, name := range condition.ExtraInfoLikes {
			args = append(args, name)
		}
	}

	var count int64
	if err := tx.Model(&model.LoveInfo{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (b *loveInfo) FindLoveInfos(ctx context.Context, tx *gorm.DB, condition FindLoveInfoArg) ([]model.LoveInfo, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.LoveNames != nil {
		query.WriteString(fmt.Sprintf(" AND `love_name` IN (%s)", RepeatWithSep("?", len(condition.LoveNames), ",")))
		for _, name := range condition.LoveNames {
			args = append(args, name)
		}
	}
	if condition.ExtraInfoLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `extra_info` Like (%s)", RepeatWithSep("?", len(condition.ExtraInfoLikes), ",")))
		for _, name := range condition.ExtraInfoLikes {
			args = append(args, name)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `confession_time` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var loveInfos []model.LoveInfo
	if err := tx.Model(&model.LoveInfo{}).Where(query.String(), args...).Find(&loveInfos).Error; err != nil {
		return nil, err
	}

	return loveInfos, nil
}

func (b *loveInfo) GetLoveInfo(ctx context.Context, tx *gorm.DB, id uint) (model.LoveInfo, error) {
	var loveInfo model.LoveInfo
	if err := tx.Model(&model.LoveInfo{}).Where("`id` = ?", id).First(&loveInfo).Error; err != nil {
		return loveInfo, err
	}

	return loveInfo, nil
}

func (b *loveInfo) CreateLoveInfo(ctx context.Context, tx *gorm.DB, loveInfo model.LoveInfo) error {
	if err := tx.Model(&model.LoveInfo{}).Create(loveInfo).Error; err != nil {
		return err
	}

	return nil
}

func (b *loveInfo) UpdateLoveInfo(ctx context.Context, tx *gorm.DB, loveInfo model.LoveInfo) error {
	if err := tx.Model(&model.LoveInfo{}).Where("`id` = ?", loveInfo.ID.ID).Updates(loveInfo).Error; err != nil {
		return err
	}

	return nil
}

func (b *loveInfo) DeleteLoveInfo(ctx context.Context, tx *gorm.DB, id uint) error {
	if err := tx.Model(&model.LoveInfo{}).Where("`id` = ?", id).Delete(&model.LoveInfo{}).Error; err != nil {
		return err
	}

	return nil
}
