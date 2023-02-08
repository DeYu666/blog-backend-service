package blog

import (
	"errors"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetTag(options ...func(option *gorm.DB)) ([]*models.BlogTag, error) {
	return dbctl.GetDBData(&models.BlogTag{}, options...)
}

func TagId(id uint) Option {
	return func(db *gorm.DB) {
		db.Where("`id` = ?", id)
	}
}

func TagName(name string) Option {
	return func(db *gorm.DB) {
		db.Where("`name` = ?", name)
	}
}

func AddTag(tag *models.BlogTag) error {
	return dbctl.AddDBData(tag)
}

func DeleteTag(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.BlogTag{}, options...)
}

func ModifyTag(tag *models.BlogTag) error {
	if tag.ID.ID == 0 {
		return errors.New("categoryId not exist, please check id")
	}

	data, err := GetTag(TagId(tag.ID.ID))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist category")
	}

	db := global.App.DB

	result := db.Save(&tag)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
