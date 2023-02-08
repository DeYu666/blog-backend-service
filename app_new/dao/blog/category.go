package blog

import (
	"errors"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetCategories(options ...func(option *gorm.DB)) ([]*models.BlogCategories, error) {
	category, err := dbctl.GetDBData(&models.BlogCategories{}, options...)
	if err != nil {
		return nil, err
	}

	for index, cate := range category {
		var generalCate models.BlogGeneralCategories
		err = global.App.DB.Model(&cate).Association("General").Find(&generalCate)
		if err != nil {
			return nil, err
		}

		category[index].General.ID = generalCate.ID
		category[index].General.Name = generalCate.Name
	}

	return category, nil
}

func CateId(id uint) Option {
	return setIdByUint(id)
}

func CategoryByGeneralID(generalId uint) Option {
	return func(db *gorm.DB) {
		db.Where("`general_id` = ?", generalId)
	}
}

func CateName(name string) Option {
	return func(db *gorm.DB) {
		db.Where("`name` = ?", name)
	}
}

func AddCategory(cate *models.BlogCategories) error {
	return dbctl.AddDBData(cate)
}

func ModifyCategory(cate *models.BlogCategories) error {

	if cate.ID.ID == 0 {
		return errors.New("categoryId not exist, please check id")
	}

	data, err := GetCategories(CateId(cate.ID.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist category")
	}

	db := global.App.DB

	result := db.Save(&cate)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteCategory(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.BlogCategories{}, options...)
}
