package cv

import (
	"errors"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetSkills(options ...func(option *gorm.DB)) ([]*models.SkillCv, error) {
	return dbctl.GetDBData(&models.SkillCv{}, options...)
}

func SkillCvId(id uint) Option {
	return setIdByUint(id)
}

func SkillNameCv(name string) Option {
	return func(db *gorm.DB) {
		db.Where("`skill_name` = ?", name)
	}
}

func AddSkill(skillCv *models.SkillCv) error {
	return dbctl.AddDBData(skillCv)
}

func DeleteSkill(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.SkillCv{}, options...)
}

func ModifySkill(skill *models.SkillCv) error {
	if skill.ID.ID == 0 {
		return errors.New("skillCv id not exist, please check id")
	}

	data, err := GetSkills(SkillCvId(skill.ID.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist skillCv")
	}

	db := global.App.DB
	result := db.Save(&skill)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
