package cv

import (
	"errors"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

func GetProjects(options ...func(option *gorm.DB)) ([]*models.ProjectCv, error) {
	return dbctl.GetDBData(&models.ProjectCv{}, options...)
}

func ProjectCvId(id uint) Option {
	return setIdByUint(id)
}

func ProjectNameCv(name string) Option {
	return func(db *gorm.DB) {
		db.Where("`project_name` = ?", name)
	}
}

func AddProject(projectCv *models.ProjectCv) error {
	return dbctl.AddDBData(projectCv)
}

func DeleteProject(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.ProjectCv{}, options...)
}

func ModifyProject(project *models.ProjectCv) error {
	if project.ID.ID == 0 {
		return errors.New("projectCv id not exist, please check id")
	}

	data, err := GetProjects(ProjectCvId(project.ID.ID))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist ProjectCv")
	}

	db := global.App.DB
	result := db.Save(&project)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// 对项目中设置的密码进行 curd

func GetProjectPwd(options ...func(option *gorm.DB)) ([]*models.ProjectCvPs, error) {
	return dbctl.GetDBData(&models.ProjectCvPs{}, options...)
}

func ProjectPwdId(id uint) Option {
	return setIdByUint(id)
}

func ProjectPassWord(pwd string) Option {
	return func(db *gorm.DB) {
		db.Where("`password` = ?", pwd)
	}
}

func AddProjectPwd(projectPwd *models.ProjectCvPs) error {
	return dbctl.AddDBData(projectPwd)
}

func DeleteProjectPwd(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.ProjectCvPs{}, options...)
}

func ModifyProjectPwd(projectPwd *models.ProjectCvPs) error {
	if projectPwd.ID.ID == 0 {
		return errors.New("projectCv id not exist, please check id")
	}

	data, err := GetProjectPwd(ProjectPwdId(projectPwd.ID.ID))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist ProjectCv")
	}

	db := global.App.DB
	result := db.Save(&projectPwd)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
