package repository

import (
	"context"
	"fmt"
	"strings"

	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type Memo interface {
	CountMemos(ctx context.Context, tx *gorm.DB, condition FindMemoArg) (int64, error)
	FindMemos(ctx context.Context, tx *gorm.DB, condition FindMemoArg) ([]model.Memo, error)
	GetMemo(ctx context.Context, tx *gorm.DB, id uint) (model.Memo, error)
	CreateMemo(ctx context.Context, tx *gorm.DB, memo model.Memo) error
	UpdateMemo(ctx context.Context, tx *gorm.DB, memo model.Memo) error
	DeleteMemo(ctx context.Context, tx *gorm.DB, id uint) error
}

type memo struct{}

func NewMemo() Memo {
	return &memo{}
}

type FindMemoArg struct {
	IDs          []uint
	ContentLikes []string
	StatusLikes  int
	Offset       int
	Limit        int
	NoLimit      bool
}

func (b *memo) CountMemos(ctx context.Context, tx *gorm.DB, condition FindMemoArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(" 1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.ContentLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `content` Like (%s)", RepeatWithSep("?", len(condition.ContentLikes), ",")))
		for _, name := range condition.ContentLikes {
			args = append(args, name)
		}
	}

	if condition.StatusLikes != 0 {
		query.WriteString(" AND `status` = ?")
		args = append(args, condition.StatusLikes)
	}

	var count int64
	if err := tx.Model(&model.Memo{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (b *memo) FindMemos(ctx context.Context, tx *gorm.DB, condition FindMemoArg) ([]model.Memo, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.ContentLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `content` Like (%s)", RepeatWithSep("?", len(condition.ContentLikes), ",")))
		for _, name := range condition.ContentLikes {
			args = append(args, name)
		}
	}
	if condition.StatusLikes != 0 {
		query.WriteString(" AND `status` = ?")
		args = append(args, condition.StatusLikes)
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `id` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var memos []model.Memo
	if err := tx.Model(&model.Memo{}).Where(query.String(), args...).Find(&memos).Error; err != nil {
		return nil, err
	}

	return memos, nil
}

func (b *memo) GetMemo(ctx context.Context, tx *gorm.DB, id uint) (model.Memo, error) {
	var memo model.Memo
	if err := tx.Model(&model.Memo{}).Where("`id` = ?", id).First(&memo).Error; err != nil {
		return memo, err
	}

	return memo, nil
}

func (b *memo) CreateMemo(ctx context.Context, tx *gorm.DB, memo model.Memo) error {

	log := blog.Extract(ctx)

	log.Sugar().Infof("memo: %+v", memo)

	if err := tx.Create(memo).Error; err != nil {

		log.Sugar().Errorf("err: %+v", err)

		return err
	}

	return nil
}

func (b *memo) UpdateMemo(ctx context.Context, tx *gorm.DB, memo model.Memo) error {
	if err := tx.Model(&model.Memo{}).Where("`id` = ?", memo.ID.ID).Updates(memo).Error; err != nil {
		return err
	}

	return nil
}

func (b *memo) DeleteMemo(ctx context.Context, tx *gorm.DB, id uint) error {
	if err := tx.Model(&model.Memo{}).Where("`id` = ?", id).Delete(&model.Memo{}).Error; err != nil {
		return err
	}

	return nil
}
