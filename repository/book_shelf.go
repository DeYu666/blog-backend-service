package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type BookShelf interface {
	CountBookShelfs(ctx context.Context, tx *gorm.DB, condition FindBookShelfArg) (int64, error)
	FindBookShelfs(ctx context.Context, tx *gorm.DB, condition FindBookShelfArg) ([]model.BooksList, error)
	GetBookShelf(ctx context.Context, tx *gorm.DB, bookShelfId uint) (model.BooksList, error)
	CreateBookShelf(ctx context.Context, tx *gorm.DB, bookShelf model.BooksList) error
	UpdateBookShelf(ctx context.Context, tx *gorm.DB, bookShelf model.BooksList) error
	DeleteBookShelf(ctx context.Context, tx *gorm.DB, bookShelfId uint) error
}

type bookShelf struct{}

func NewBookShelf() BookShelf {
	return &bookShelf{}
}

type FindBookShelfArg struct {
	IDs           []uint
	BookNames     []string
	BookNameLikes []string
	BookStatuses  []string
	Offset        int32
	Limit         int32
	NoLimit       bool
}

func (b *bookShelf) CountBookShelfs(ctx context.Context, tx *gorm.DB, condition FindBookShelfArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(" 1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.BookNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `book_name` Like (%s)", RepeatWithSep("?", len(condition.BookNameLikes), ",")))
		for _, name := range condition.BookNameLikes {
			args = append(args, name)
		}
	}

	if condition.BookNames != nil {
		query.WriteString(fmt.Sprintf(" AND `book_name` IN (%s)", RepeatWithSep("?", len(condition.BookNames), ",")))
		for _, id := range condition.BookNames {
			args = append(args, id)
		}
	}

	if condition.BookStatuses != nil {
		query.WriteString(fmt.Sprintf(" AND `book_status` IN (%s)", RepeatWithSep("?", len(condition.BookStatuses), ",")))
		for _, id := range condition.BookStatuses {
			args = append(args, id)
		}
	}

	var count int64

	err := tx.Model(&model.BooksList{}).Where(query.String(), args...).Count(&count).Error
	return count, err
}

func (b *bookShelf) FindBookShelfs(ctx context.Context, tx *gorm.DB, condition FindBookShelfArg) ([]model.BooksList, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}
	if condition.BookNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `book_name` Like (%s)", RepeatWithSep("?", len(condition.BookNameLikes), ",")))
		for _, name := range condition.BookNameLikes {
			args = append(args, name)
		}
	}

	if condition.BookNames != nil {
		query.WriteString(fmt.Sprintf(" AND `book_name` IN (%s)", RepeatWithSep("?", len(condition.BookNames), ",")))
		for _, id := range condition.BookNames {
			args = append(args, id)
		}
	}

	if condition.BookStatuses != nil {
		query.WriteString(fmt.Sprintf(" AND `book_status` IN (%s)", RepeatWithSep("?", len(condition.BookStatuses), ",")))
		for _, id := range condition.BookStatuses {
			args = append(args, id)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `id` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var bookShelfs []model.BooksList

	err := tx.Model(&model.BooksList{}).Where(query.String(), args...).Find(&bookShelfs).Error
	return bookShelfs, err
}

func (b *bookShelf) GetBookShelf(ctx context.Context, tx *gorm.DB, bookShelfId uint) (model.BooksList, error) {
	var bookShelf model.BooksList
	err := tx.Model(&model.BooksList{}).Where("`id` = ?", bookShelfId).First(&bookShelf).Error
	return bookShelf, err
}

func (b *bookShelf) CreateBookShelf(ctx context.Context, tx *gorm.DB, bookShelf model.BooksList) error {
	err := tx.Model(&model.BooksList{}).Create(&bookShelf).Error
	return err
}

func (b *bookShelf) UpdateBookShelf(ctx context.Context, tx *gorm.DB, bookShelf model.BooksList) error {
	err := tx.Model(&model.BooksList{}).Where("`id` = ?", bookShelf.ID.ID).Updates(&bookShelf).Error
	return err
}

func (b *bookShelf) DeleteBookShelf(ctx context.Context, tx *gorm.DB, bookShelfId uint) error {
	err := tx.Model(&model.BooksList{}).Where("`id` = ?", bookShelfId).Delete(&model.BooksList{}).Error
	return err
}
