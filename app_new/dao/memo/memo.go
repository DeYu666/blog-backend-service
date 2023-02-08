package memo

import (
	"errors"
	"time"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

type Option func(db *gorm.DB)

func GetMemos(options ...func(option *gorm.DB)) ([]*models.Memo, error) {
	return dbctl.GetDBData(&models.Memo{}, options...)
}

func MemoId(id uint) Option {
	return func(db *gorm.DB) {
		db.Where("`id` = ?", id)
	}
}

func MemoContent(content string) Option {
	return func(db *gorm.DB) {
		db.Where("`content` = ?", content)
	}
}

// MemoStatus 中参数 0 表示未完成； 1 表示已完成；2 表示全部
func MemoStatus(statusId int) Option {

	if statusId == 2 {
		return func(db *gorm.DB) {

		}
	}

	return func(db *gorm.DB) {
		db.Where("`status` = ?", statusId)
	}
}

func MemoOrderByDesc() Option {
	return func(db *gorm.DB) {
		db.Order("created_time desc")
	}
}

func AddMemo(memoContent *models.Memo) error {

	memoContent.Status = 0

	if memoContent.CreatedTime.IsZero() || memoContent.ModifiedTime.IsZero() {
		timeObj := time.Now()
		if memoContent.CreatedTime.IsZero() {
			memoContent.CreatedTime = timeObj
		}
		if memoContent.ModifiedTime.IsZero() {
			memoContent.ModifiedTime = timeObj
		}
	}

	return dbctl.AddDBData(memoContent)
}

func DeleteMemo(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.Memo{}, options...)
}

func ModifyMemo(memo *models.Memo) error {
	if memo.ID.ID == 0 {
		return errors.New("experienceCv id not exist, please check id")
	}

	data, err := GetMemos(MemoId(memo.ID.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist Memo")
	}

	memo.ModifiedTime = time.Now()
	db := global.App.DB
	result := db.Save(&memo)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
