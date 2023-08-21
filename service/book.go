package service

import (
	"context"
	"time"

	"github.com/DeYu666/blog-backend-service/lib/client"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/repository"
	"gorm.io/gorm"
)

type BookService interface {
	GetBooksList(ctx context.Context, status string, offset, limit int32) (model.ShelfArrWithCount, error)
	CreateBooksList(ctx context.Context, booksList model.BooksList) error
	UpdateBooksList(ctx context.Context, booksList model.BooksList) error
	DeleteBooksList(ctx context.Context, booksListId uint) error

	GetBookContent(ctx context.Context) ([]model.BookContent, error)
	GetBookContentByBookId(ctx context.Context, bookId uint) ([]model.BookContent, error)
	CreateBookContent(ctx context.Context, booksContent model.BookContent) error
	UpdateBookContent(ctx context.Context, booksContent model.BookContent) error
	DeleteBookContent(ctx context.Context, booksContentId uint) error
}

type bookService struct {
	bookShelfRepo   repository.BookShelf
	bookContentRepo repository.BookContent
}

func NewBookService() BookService {
	return &bookService{
		bookShelfRepo:   repository.NewBookShelf(),
		bookContentRepo: repository.NewBookContent(),
	}
}

func (b *bookService) GetBooksList(ctx context.Context, status string, offset, limit int32) (model.ShelfArrWithCount, error) {

	var booksList []model.BooksList
	var count int64
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindBookShelfArg{}
		cond.NoLimit = true

		if status != "" {
			cond.BookStatuses = []string{status}
		}

		count, err = b.bookShelfRepo.CountBookShelfs(ctx, tx, cond)
		if err != nil {
			return err
		}

		if count == 0 {
			return nil
		}

		if limit != 0 {
			cond.Limit = limit
			cond.Offset = offset
			cond.NoLimit = false
		}

		booksList, err = b.bookShelfRepo.FindBookShelfs(ctx, tx, cond)
		return err
	}, nil)

	booksListWithCount := model.ShelfArrWithCount{
		ShelfArr: booksList,
		Count:    count,
	}

	return booksListWithCount, err
}

func (b *bookService) CreateBooksList(ctx context.Context, booksList model.BooksList) error {

	var err error

	booksList.CreatedTime = time.Now()
	booksList.ModifiedTime = time.Now()

	booksList.BookStatus = "未读"

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.bookShelfRepo.CreateBookShelf(ctx, tx, booksList)
		return err
	}, nil)

	return err
}

func (b *bookService) UpdateBooksList(ctx context.Context, booksList model.BooksList) error {

	var err error

	booksList.ModifiedTime = time.Now()

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.bookShelfRepo.UpdateBookShelf(ctx, tx, booksList)
		return err
	}, nil)

	return err
}

func (b *bookService) DeleteBooksList(ctx context.Context, booksListId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.bookShelfRepo.DeleteBookShelf(ctx, tx, booksListId)
		return err
	}, nil)

	return err
}

func (b *bookService) GetBookContent(ctx context.Context) ([]model.BookContent, error) {

	var bookContent []model.BookContent
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindBookContentArg{
			NoLimit: true,
		}
		bookContent, err = b.bookContentRepo.FindBookContents(ctx, tx, cond)
		return err
	}, nil)

	return bookContent, err
}

func (b *bookService) GetBookContentByBookId(ctx context.Context, bookId uint) ([]model.BookContent, error) {

	var bookContents []model.BookContent
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		bookContents, err = b.bookContentRepo.GetBookContentByShelfId(ctx, tx, bookId)
		return err
	}, nil)

	return bookContents, err
}

func (b *bookService) CreateBookContent(ctx context.Context, bookContent model.BookContent) error {

	var err error

	bookContent.CreatedTime = time.Now()
	bookContent.ModifiedTime = time.Now()

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.bookContentRepo.CreateBookContent(ctx, tx, bookContent)
		return err
	}, nil)

	return err
}

func (b *bookService) UpdateBookContent(ctx context.Context, bookContent model.BookContent) error {

	var err error

	bookContent.ModifiedTime = time.Now()

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.bookContentRepo.UpdateBookContent(ctx, tx, bookContent)
		return err
	}, nil)

	return err
}

func (b *bookService) DeleteBookContent(ctx context.Context, bookContentId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = b.bookContentRepo.DeleteBookContent(ctx, tx, bookContentId)
		return err
	}, nil)

	return err
}
