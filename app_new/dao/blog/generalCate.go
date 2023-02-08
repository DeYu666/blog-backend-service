package blog

import (
	"errors"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetGeneralCategories(options ...func(option *gorm.DB)) ([]*models.BlogGeneralCategories, error) {
	return dbctl.GetDBData(&models.BlogGeneralCategories{}, options...)
}

func GeneralCateName(name string) Option {
	return func(db *gorm.DB) {
		db.Where("`name` = ? ", name)
	}
}

func GeneralCateId(id uint) Option {
	return setIdByUint(id)
}

func AddGeneralCate(cate *models.BlogGeneralCategories) error {
	return dbctl.AddDBData(cate)
}

func ModifyGeneralCate(cate *models.BlogGeneralCategories) error {

	if cate.ID.ID == 0 {
		return errors.New("generalCateId not exist, please check id")
	}

	data, err := GetGeneralCategories(GeneralCateId(cate.ID.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist general category")
	}

	db := global.App.DB
	result := db.Save(&cate)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteGeneralCate(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.BlogGeneralCategories{}, options...)
}
