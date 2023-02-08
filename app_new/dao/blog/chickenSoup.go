package blog

import (
	"errors"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetChickenSoups(options ...func(option *gorm.DB)) ([]*models.ChickenSoup, error) {
	return dbctl.GetDBData(&models.ChickenSoup{}, options...)
}

func ChickenSoupId(id uint) Option {
	return setIdByUint(id)
}

func ChickenSoupSentence(sentence string) Option {
	return func(db *gorm.DB) {
		db.Where("`sentence` = ?", sentence)
	}
}

func ChickenSoupOrderByDesc() Option {
	return func(db *gorm.DB) {
		db.Order("id desc")
	}
}

func AddChickenSoup(data *models.ChickenSoup) error {
	return dbctl.AddDBData(data)
}

func DeleteChickenSoup(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.ChickenSoup{}, options...)
}

func ModifyChickenSoup(data *models.ChickenSoup) error {
	if data.ID.ID == 0 {
		return errors.New("chickenSoup id not exist, please check id")
	}

	info, err := GetChickenSoups(ChickenSoupId(data.ID.ID))

	if err != nil {
		return err
	}
	if len(info) == 0 {
		return errors.New("database is not exist general category")
	}

	db := global.App.DB
	result := db.Save(&data)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
