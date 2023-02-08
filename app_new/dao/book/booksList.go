package book

import (
	"errors"
	"time"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetBooksLists(options ...func(option *gorm.DB)) ([]*models.BooksList, error) {
	return dbctl.GetDBData(&models.BooksList{}, options...)
}

func BookId(id uint) Option {
	return setIdByUint(id)
}

func BookName(name string) Option {
	return func(db *gorm.DB) {
		db.Where("`book_name` = ?", name)
	}
}

/*
BookPaging 分页逻辑
pageNum 指获取第 pageNum 页的数据
pageSize 指每页中有多少条数据
*/
func BookPaging(pageNum int, pageSize int) Option {
	if pageNum == 0 {
		pageNum = 1
	}
	offset := (pageNum - 1) * pageSize

	return func(db *gorm.DB) {
		db.Limit(pageSize).Offset(offset)
	}
}

func BookOrderByDesc() Option {
	return func(db *gorm.DB) {
		db.Order("created_time desc")
	}
}

// BookCategoryID 中 cateId 指 book 的分类，当 -1 时，获取全部分类
func BookCategoryID(cateId int) Option {

	if cateId == -1 {
		return func(db *gorm.DB) {

		}
	}

	bookStatus := ""
	if cateId == 0 {
		bookStatus = "在看"
	} else if cateId == 1 {
		bookStatus = "已看"
	} else if cateId == 2 {
		bookStatus = "未读"
	}

	return func(db *gorm.DB) {
		db.Where("`book_status` = ? ", bookStatus)
	}
}

func AddBooksLists(book *models.BooksList) error {

	if book.CreatedTime.IsZero() || book.ModifiedTime.IsZero() {
		timeObj := time.Now()
		if book.CreatedTime.IsZero() {
			book.CreatedTime = timeObj
		}
		if book.ModifiedTime.IsZero() {
			book.ModifiedTime = timeObj
		}
	}

	book.BookStatus = "未读"

	return dbctl.AddDBData(book)
}

func DeleteBooksLists(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.BooksList{}, options...)
}

func ModifyBooksLists(book *models.BooksList) error {

	if book.ID.ID == 0 {
		return errors.New("book id not exist, please check id")
	}

	data, err := GetBooksLists(BookId(book.ID.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist category")
	}

	db := global.App.DB

	book.ModifiedTime = time.Now()

	result := db.Save(&book)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
