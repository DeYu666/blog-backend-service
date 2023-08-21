package service

import (
	"context"

	"github.com/DeYu666/blog-backend-service/lib/client"
	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/repository"
	"gorm.io/gorm"
)

type AuthUserService interface {
	Login(ctx context.Context, user model.AuthUser) (model.AuthUser, error)
	GetUsers(ctx context.Context) ([]model.AuthUser, error)
	Register(ctx context.Context, user model.AuthUser) error
	UpdateUser(ctx context.Context, user model.AuthUser) error
	DeleteUser(ctx context.Context, userId uint) error
}

type authUserService struct {
	authUserRepo repository.AuthUser
}

func NewAuthUserService() AuthUserService {
	return &authUserService{
		authUserRepo: repository.NewAuthUser(),
	}
}

func (a *authUserService) Login(ctx context.Context, user model.AuthUser) (model.AuthUser, error) {

	log := blog.Extract(ctx)

	log.Sugar().Info("user", user)

	var authUser []model.AuthUser
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindAuthUserArg{
			UserNames: []string{user.Username},
			Passwords: []string{user.Password},
			NoLimit:   true,
		}
		authUser, err = a.authUserRepo.FindAuthUser(ctx, tx, cond)
		return err
	}, nil)

	log.Sugar().Info("authUser", authUser)

	if len(authUser) == 0 {
		return model.AuthUser{}, err
	}

	return authUser[0], err
}

func (a *authUserService) GetUsers(ctx context.Context) ([]model.AuthUser, error) {

	var authUsers []model.AuthUser
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindAuthUserArg{
			NoLimit: true,
		}
		authUsers, err = a.authUserRepo.FindAuthUser(ctx, tx, cond)
		return err
	}, nil)

	return authUsers, err
}

func (a *authUserService) Register(ctx context.Context, user model.AuthUser) error {
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = a.authUserRepo.CreateAuthUser(ctx, tx, user)
		return err
	}, nil)

	return err
}

func (a *authUserService) DeleteUser(ctx context.Context, userId uint) error {
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {

		cond := repository.FindAuthUserArg{
			NoLimit:    true,
			ExcludeIds: []uint{userId},
		}

		if count, err := a.authUserRepo.CountAuthUser(ctx, tx, cond); err != nil || count == 0 {
			return err
		}

		err = a.authUserRepo.DeleteAuthUser(ctx, tx, userId)
		return err
	}, nil)

	return err
}

func (a *authUserService) UpdateUser(ctx context.Context, user model.AuthUser) error {
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = a.authUserRepo.UpdateAuthUser(ctx, tx, user)
		return err
	}, nil)

	return err
}
