package book

import (
	"errors"
	"time"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetBookContent(options ...func(option *gorm.DB)) ([]*models.BookContent, error) {
	return dbctl.GetDBData(&models.BookContent{}, options...)
}

func BookContentId(id uint) Option {
	return setIdByUint(id)
}

func BookContentByBookId(bookId int) Option {
	return func(db *gorm.DB) {
		db.Where("book_id = ?", bookId)
	}
}

func BookContentOrderByDesc() Option {
	return func(db *gorm.DB) {
		db.Order("created_time desc")
	}
}

func BookContent(content string) Option {
	return func(db *gorm.DB) {
		db.Where("`book_content` = ?", content)
	}
}

func DeleteBookContent(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.BookContent{}, options...)
}

func AddBookContent(content *models.BookContent) error {

	if content.CreatedTime.IsZero() || content.ModifiedTime.IsZero() {
		timeObj := time.Now()
		if content.CreatedTime.IsZero() {
			content.CreatedTime = timeObj
		}
		if content.ModifiedTime.IsZero() {
			content.ModifiedTime = timeObj
		}
	}

	err := dbctl.AddDBData(content)
	if err != nil {
		return err
	}

	// 更新书籍状态
	books, err := GetBooksLists(BookId(uint(content.BookId)))
	if err != nil {
		return err
	}
	if len(books) == 0 {
		return errors.New("未查到此书籍")
	}

	if books[0].BookStatus != "在看" {
		books[0].BookStatus = "在看"
		err = ModifyBooksLists(books[0])
		if err != nil {
			return err
		}
	}

	return nil
}

func ModifyBookContent(content *models.BookContent) error {
	if content.ID.ID == 0 {
		return errors.New("bookContent id not exist, please check id")
	}

	data, err := GetBookContent(BookContentId(content.ID.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist book content")
	}

	content.ModifiedTime = time.Now()

	db := global.App.DB
	result := db.Save(&content)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
