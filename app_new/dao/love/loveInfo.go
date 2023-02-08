package love

import (
	"errors"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

type Option func(db *gorm.DB)

func GetLoveInfos(options ...func(option *gorm.DB)) ([]*models.LoveInfo, error) {
	return dbctl.GetDBData(&models.LoveInfo{}, options...)
}

func LoveInfoId(id uint) Option {
	return func(db *gorm.DB) {
		db.Where("`id` = ?", id)
	}
}

func AddLoveInfo(LoveInfo *models.LoveInfo) error {
	LoveInfo.ID.ID = 0
	return dbctl.AddDBData(LoveInfo)
}

func DeleteLoveInfo(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.LoveInfo{}, options...)
}

func ModifyLoveInfo(LoveInfo *models.LoveInfo) error {
	if LoveInfo.ID.ID == 0 {
		return errors.New("LoveInfo id not exist, please check id")
	}

	data, err := GetLoveInfos(LoveInfoId(LoveInfo.ID.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist LoveInfo")
	}

	db := global.App.DB
	result := db.Save(&LoveInfo)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
