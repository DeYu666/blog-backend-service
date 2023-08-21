package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type BookContent interface {
	CountBookContents(ctx context.Context, tx *gorm.DB, condition FindBookContentArg) (int64, error)
	FindBookContents(ctx context.Context, tx *gorm.DB, condition FindBookContentArg) ([]model.BookContent, error)
	GetBookContentByShelfId(ctx context.Context, tx *gorm.DB, shelfId uint) ([]model.BookContent, error)
	GetBookContent(ctx context.Context, tx *gorm.DB, bookContentId uint) (model.BookContent, error)
	CreateBookContent(ctx context.Context, tx *gorm.DB, bookContent model.BookContent) error
	UpdateBookContent(ctx context.Context, tx *gorm.DB, bookContent model.BookContent) error
	DeleteBookContent(ctx context.Context, tx *gorm.DB, bookContentId uint) error
}

type bookContent struct{}

func NewBookContent() BookContent {
	return &bookContent{}
}

type FindBookContentArg struct {
	IDs         []uint
	ConentLikes []string
	BookIDs     []uint
	Offset      int
	Limit       int
	NoLimit     bool
}

func (b *bookContent) CountBookContents(ctx context.Context, tx *gorm.DB, condition FindBookContentArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(" 1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.ConentLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `book_content` Like (%s)", RepeatWithSep("?", len(condition.ConentLikes), ",")))
		for _, name := range condition.ConentLikes {
			args = append(args, name)
		}
	}
	if condition.BookIDs != nil {
		query.WriteString(fmt.Sprintf(" AND `book_id` IN (%s)", RepeatWithSep("?", len(condition.BookIDs), ",")))
		for _, id := range condition.BookIDs {
			args = append(args, id)
		}
	}

	var count int64

	err := tx.Model(&model.BookContent{}).Where(query.String(), args...).Count(&count).Error
	return count, err
}

func (b *bookContent) FindBookContents(ctx context.Context, tx *gorm.DB, condition FindBookContentArg) ([]model.BookContent, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.ConentLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `book_content` Like (%s)", RepeatWithSep("?", len(condition.ConentLikes), ",")))
		for _, name := range condition.ConentLikes {
			args = append(args, name)
		}
	}
	if condition.BookIDs != nil {
		query.WriteString(fmt.Sprintf(" AND `book_id` IN (%s)", RepeatWithSep("?", len(condition.BookIDs), ",")))
		for _, id := range condition.BookIDs {
			args = append(args, id)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `id` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var bookContents []model.BookContent

	err := tx.Model(&model.BookContent{}).Where(query.String(), args...).Find(&bookContents).Error
	return bookContents, err
}

func (b *bookContent) GetBookContent(ctx context.Context, tx *gorm.DB, bookContentId uint) (model.BookContent, error) {
	var bookContent model.BookContent
	err := tx.Model(&model.BookContent{}).Where("`id` = ?", bookContentId).First(&bookContent).Error
	return bookContent, err
}

func (b *bookContent) GetBookContentByShelfId(ctx context.Context, tx *gorm.DB, shelfId uint) ([]model.BookContent, error) {
	var bookContents []model.BookContent
	err := tx.Model(&model.BookContent{}).Where("`book_id` = ?", shelfId).Find(&bookContents).Error
	return bookContents, err
}

func (b *bookContent) CreateBookContent(ctx context.Context, tx *gorm.DB, bookContent model.BookContent) error {
	err := tx.Model(&model.BookContent{}).Create(&bookContent).Error
	return err
}

func (b *bookContent) UpdateBookContent(ctx context.Context, tx *gorm.DB, bookContent model.BookContent) error {
	err := tx.Model(&model.BookContent{}).Where("`id` = ?", bookContent.ID.ID).Updates(&bookContent).Error
	return err
}

func (b *bookContent) DeleteBookContent(ctx context.Context, tx *gorm.DB, bookContentId uint) error {
	err := tx.Model(&model.BookContent{}).Where("`id` = ?", bookContentId).Delete(&model.BookContent{}).Error
	return err
}
