package repository

import (
	"context"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type AuthUser interface {
	CountAuthUser(ctx context.Context, tx *gorm.DB, cond FindAuthUserArg) (int64, error)
	FindAuthUser(ctx context.Context, tx *gorm.DB, cond FindAuthUserArg) ([]model.AuthUser, error)
	GetAuthUser(ctx context.Context, tx *gorm.DB, id uint) (model.AuthUser, error)
	CreateAuthUser(ctx context.Context, tx *gorm.DB, authUser model.AuthUser) error
	UpdateAuthUser(ctx context.Context, tx *gorm.DB, authUser model.AuthUser) error
	DeleteAuthUser(ctx context.Context, tx *gorm.DB, authUserId uint) error
}

type authUser struct{}

func NewAuthUser() AuthUser {
	return &authUser{}
}

type FindAuthUserArg struct {
	IDs        []uint
	ExcludeIds []uint
	UserNames  []string
	Passwords  []string
	Offset     int32
	Limit      int32
	NoLimit    bool
}

func (a *authUser) CountAuthUser(ctx context.Context, tx *gorm.DB, cond FindAuthUserArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if cond.IDs != nil {
		query.WriteString(" AND `id` IN (?)")
		args = append(args, cond.IDs)
	}

	if cond.ExcludeIds != nil {
		query.WriteString(" AND `id` NOT IN (?)")
		args = append(args, cond.ExcludeIds)
	}

	if cond.UserNames != nil {
		query.WriteString(" AND `username` IN (?)")
		args = append(args, cond.UserNames)
	}

	if cond.Passwords != nil {
		query.WriteString(" AND `password` IN (?)")
		args = append(args, cond.Passwords)
	}

	if !cond.NoLimit {
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, cond.Limit, cond.Offset)
	}

	var count int64
	if err := tx.Model(&model.AuthUser{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (a *authUser) FindAuthUser(ctx context.Context, tx *gorm.DB, cond FindAuthUserArg) ([]model.AuthUser, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if cond.IDs != nil {
		query.WriteString(" AND `id` IN (?)")
		args = append(args, cond.IDs)
	}

	if cond.ExcludeIds != nil {
		query.WriteString(" AND `id` NOT IN (?)")
		args = append(args, cond.ExcludeIds)
	}

	if cond.UserNames != nil {
		query.WriteString(" AND `username` IN (?)")
		args = append(args, cond.UserNames)
	}

	if cond.Passwords != nil {
		query.WriteString(" AND `password` IN (?)")
		args = append(args, cond.Passwords)
	}

	if !cond.NoLimit {
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, cond.Limit, cond.Offset)
	}

	var authUsers []model.AuthUser

	err := tx.Model(&model.AuthUser{}).Where(query.String(), args...).Find(&authUsers).Error
	return authUsers, err
}

func (a *authUser) GetAuthUser(ctx context.Context, tx *gorm.DB, id uint) (model.AuthUser, error) {
	var authUser model.AuthUser

	err := tx.Model(&model.AuthUser{}).Where("`id` = ?", id).First(&authUser).Error
	return authUser, err
}

func (a *authUser) CreateAuthUser(ctx context.Context, tx *gorm.DB, authUser model.AuthUser) error {
	err := tx.Model(&model.AuthUser{}).Create(&authUser).Error
	return err
}

func (a *authUser) UpdateAuthUser(ctx context.Context, tx *gorm.DB, authUser model.AuthUser) error {
	err := tx.Model(&model.AuthUser{}).Where("`id` = ?", authUser.ID).Updates(&authUser).Error
	return err
}

func (a *authUser) DeleteAuthUser(ctx context.Context, tx *gorm.DB, authUserId uint) error {
	err := tx.Model(&model.AuthUser{}).Where("`id` = ?", authUserId).Delete(&model.AuthUser{}).Error
	return err
}
