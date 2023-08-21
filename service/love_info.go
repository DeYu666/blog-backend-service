package service

import (
	"context"
	"time"

	"github.com/DeYu666/blog-backend-service/lib/client"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/repository"
	"gorm.io/gorm"
)

type LoveInfoService interface {
	GetLoveInfo(ctx context.Context) ([]model.LoveInfo, error)
	CreateLoveInfo(ctx context.Context, loveInfo model.LoveInfo) error
	UpdateLoveInfo(ctx context.Context, loveInfo model.LoveInfo) error
	DeleteLoveInfo(ctx context.Context, loveInfoId uint) error
}

type loveInfoService struct {
	loveInfoRepo repository.LoveInfo
}

func NewLoveInfoService() LoveInfoService {
	return &loveInfoService{
		loveInfoRepo: repository.NewLoveInfo(),
	}
}

func (l *loveInfoService) GetLoveInfo(ctx context.Context) ([]model.LoveInfo, error) {
	var loveInfos []model.LoveInfo
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindLoveInfoArg{
			NoLimit: true,
		}
		loveInfos, err = l.loveInfoRepo.FindLoveInfos(ctx, tx, cond)
		return err
	}, nil)

	return loveInfos, err
}

func (l *loveInfoService) CreateLoveInfo(ctx context.Context, loveInfo model.LoveInfo) error {

	var err error

	if loveInfo.ID.ID == 0 {
		loveInfo.ID.ID = generateUnionId()
	}

	if loveInfo.KnownTime.IsZero() {
		loveInfo.KnownTime = time.Now()
	}

	if loveInfo.ConfessionTime.IsZero() {
		loveInfo.ConfessionTime = time.Now()
	}

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = l.loveInfoRepo.CreateLoveInfo(ctx, tx, loveInfo)
		return err
	}, nil)

	return err
}

func (l *loveInfoService) UpdateLoveInfo(ctx context.Context, loveInfo model.LoveInfo) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = l.loveInfoRepo.UpdateLoveInfo(ctx, tx, loveInfo)
		return err
	}, nil)

	return err
}

func (l *loveInfoService) DeleteLoveInfo(ctx context.Context, loveInfoId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = l.loveInfoRepo.DeleteLoveInfo(ctx, tx, loveInfoId)

		return err
	}, nil)

	return err
}
