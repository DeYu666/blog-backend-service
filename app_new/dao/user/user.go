package user

import (
	"errors"

	"github.com/DeYu666/blog-backend-service/app_new/dao/common/dbctl"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"gorm.io/gorm"
)

type Option func(db *gorm.DB)

func GetUsers(options ...func(option *gorm.DB)) ([]*models.AuthUser, error) {
	return dbctl.GetDBData(&models.AuthUser{}, options...)
}

func UserId(id int) Option {
	return func(db *gorm.DB) {
		db.Where("`id` = ?", id)
	}
}

func Username(name string) Option {
	return func(db *gorm.DB) {
		db.Where("`username` = ?", name)
	}
}

func AddUser(user *models.AuthUser) error {
	return dbctl.AddDBData(user)
}

func DeleteUser(options ...func(option *gorm.DB)) error {
	return dbctl.DeleteDBData(&models.AuthUser{}, options...)
}

func ModifyUser(user *models.AuthUser) error {
	if user.ID == 0 {
		return errors.New("bookContent id not exist, please check id")
	}

	data, err := GetUsers(UserId(user.ID))

	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("database is not exist book content")
	}

	db := global.App.DB
	result := db.Save(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func IsExist(user *models.AuthUser) (bool, error) {

	resultUser, err := GetUsers(Username(user.Username))
	if err != nil {
		global.App.Log.Error(err.Error())
		return false, err
	}
	if len(resultUser) == 0 {
		return false, nil
	}

	for _, u := range resultUser {
		if u.Password == user.Password {
			return true, nil
		}
	}

	return false, nil
}
