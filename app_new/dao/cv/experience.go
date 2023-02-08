package cv

import (
	"errors"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetExperiences(options ...func(option *gorm.DB)) ([]*models.ExperienceCv, error) {
	return dbctl.GetDBData(&models.ExperienceCv{}, options...)
}

func ExperienceCvId(id uint) Option {
	return setIdByUint(id)
}

func EnterpriseName(name string) Option {
	return func(db *gorm.DB) {
		db.Where("`enterprise_name` = ?", name)
	}
}

func AddExperience(experienceCv *models.ExperienceCv) error {
	experienceCv.ID.ID = 0
	return dbctl.AddDBData(experienceCv)
}

func DeleteExperience(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.ExperienceCv{}, options...)
}

func ModifyExperience(experience *models.ExperienceCv) error {
	if experience.ID.ID == 0 {
		return errors.New("experienceCv id not exist, please check id")
	}

	data, err := GetExperiences(ExperienceCvId(experience.ID.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist experienceCv")
	}

	db := global.App.DB
	result := db.Save(&experience)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
