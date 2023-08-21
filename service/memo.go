package service

import (
	"context"
	"time"

	"github.com/DeYu666/blog-backend-service/lib/client"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/repository"
	"gorm.io/gorm"
)

type MemoService interface {
	GetMemo(ctx context.Context) ([]model.Memo, error)
	GetMemoByStatusId(ctx context.Context, statusId int) ([]model.ShowMemo, error)
	CreateMemo(ctx context.Context, memo model.Memo) error
	UpdateMemo(ctx context.Context, memo model.Memo) error
	DeleteMemo(ctx context.Context, memoId uint) error
}

type memoService struct {
	memoRepo repository.Memo
}

func NewMemoService() MemoService {
	return &memoService{
		memoRepo: repository.NewMemo(),
	}
}

func (m *memoService) GetMemo(ctx context.Context) ([]model.Memo, error) {

	var memos []model.Memo
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindMemoArg{
			NoLimit: true,
		}
		memos, err = m.memoRepo.FindMemos(ctx, tx, cond)
		return err
	}, nil)

	return memos, err
}

func (m *memoService) GetMemoByStatusId(ctx context.Context, statusId int) ([]model.ShowMemo, error) {

	var memos []model.Memo
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindMemoArg{
			NoLimit:     true,
			StatusLikes: statusId,
		}
		memos, err = m.memoRepo.FindMemos(ctx, tx, cond)
		return err
	}, nil)

	if len(memos) < 1 {
		return nil, err
	}

	showMemos := make([]model.ShowMemo, 1)
	showMemos[0].CreateTime = memos[0].CreatedTime

	for _, memo := range memos {
		if !isSameDay(showMemos[len(showMemos)-1].CreateTime, memo.CreatedTime) {
			showMemos = append(showMemos, model.ShowMemo{CreateTime: memo.CreatedTime})
		}

		showMemos[len(showMemos)-1].Content = append(showMemos[len(showMemos)-1].Content, memo)
	}

	return showMemos, err
}

func (m *memoService) CreateMemo(ctx context.Context, memo model.Memo) error {

	var err error

	memo.CreatedTime = time.Now()
	memo.ModifiedTime = time.Now()
	memo.Status = 1

	if memo.ID.ID == 0 {
		memo.ID.ID = generateUnionId()
	}

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = m.memoRepo.CreateMemo(ctx, tx, memo)
		return err
	}, nil)

	return err
}

func (m *memoService) UpdateMemo(ctx context.Context, memo model.Memo) error {

	var err error

	memo.ModifiedTime = time.Now()

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = m.memoRepo.UpdateMemo(ctx, tx, memo)
		return err
	}, nil)

	return err
}

func (m *memoService) DeleteMemo(ctx context.Context, memoId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = m.memoRepo.DeleteMemo(ctx, tx, memoId)
		return err
	}, nil)

	return err
}

////////////////////////
//
//  private function
//
////////////////////////

func isSameDay(a time.Time, b time.Time) bool {
	if a.Year() == b.Year() {
		if a.Month().String() == b.Month().String() {
			if a.Day() == b.Day() {
				return true
			}
		}
	}

	return false
}
