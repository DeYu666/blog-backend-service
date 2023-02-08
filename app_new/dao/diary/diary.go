package diary

import (
	"errors"
	"time"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

type Option func(db *gorm.DB)

func GetDiaries(options ...func(option *gorm.DB)) ([]*models.Diary, error) {
	return dbctl.GetDBData(&models.Diary{}, options...)
}
func DiaryId(id uint) Option {
	return func(db *gorm.DB) {
		db.Where("`id` = ?", id)
	}
}

func DiaryOrderByDesc() Option {
	return func(db *gorm.DB) {
		db.Order("created_time desc")
	}
}

func DiaryContent(content string) Option {
	return func(db *gorm.DB) {
		db.Where("`content` = ?", content)
	}
}

func AddDiary(diary *models.Diary) error {

	if diary.CreatedTime.IsZero() || diary.ModifiedTime.IsZero() {
		timeObj := time.Now()
		if diary.CreatedTime.IsZero() {
			diary.CreatedTime = timeObj
		}
		if diary.ModifiedTime.IsZero() {
			diary.ModifiedTime = timeObj
		}
	}

	return dbctl.AddDBData(diary)
}

func DeleteDiary(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.Diary{}, options...)
}

func ModifyDiary(diary *models.Diary) error {
	if diary.ID.ID == 0 {
		return errors.New("experienceCv id not exist, please check id")
	}

	data, err := GetDiaries(DiaryId(diary.ID.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist diary")
	}

	diary.ModifiedTime = time.Now()
	db := global.App.DB
	result := db.Save(&diary)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// 对日记中设置的密码进行 curd

func GetDiaryPwd(options ...func(option *gorm.DB)) ([]*models.DiaryPs, error) {
	return dbctl.GetDBData(&models.DiaryPs{}, options...)
}

func DiaryPwdId(id uint) Option {
	return func(db *gorm.DB) {
		db.Where("`id` = ?", id)
	}
}

func DiaryPassWord(pwd string) Option {
	return func(db *gorm.DB) {
		db.Where("`password` = ?", pwd)
	}
}

func AddDiaryPwd(DiaryPwd *models.DiaryPs) error {
	return dbctl.AddDBData(DiaryPwd)
}

func DeleteDiaryPwd(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.DiaryPs{}, options...)
}

func ModifyDiaryPwd(DiaryPwd *models.DiaryPs) error {
	if DiaryPwd.ID.ID == 0 {
		return errors.New("DiaryCv id not exist, please check id")
	}

	data, err := GetDiaryPwd(DiaryPwdId(DiaryPwd.ID.ID))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist Diary")
	}

	db := global.App.DB
	result := db.Save(&DiaryPwd)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
