package service

import (
	"context"

	"github.com/DeYu666/blog-backend-service/lib/client"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/repository"
	"gorm.io/gorm"
)

type DiaryService interface {
	GetDiary(ctx context.Context, password string) ([]model.Diary, error)
	CreateDiary(ctx context.Context, diary model.Diary) error
	UpdateDiary(ctx context.Context, diary model.Diary) error
	DeleteDiary(ctx context.Context, diaryId uint) error

	GetDiaryPs(ctx context.Context) ([]model.DiaryPs, error)
	CreateDiaryPs(ctx context.Context, diaryPs model.DiaryPs) error
	UpdateDiaryPs(ctx context.Context, diaryPs model.DiaryPs) error
	DeleteDiaryPs(ctx context.Context, diaryPsId uint) error
}

type diaryService struct {
	diaryRepo repository.Diary
}

func NewDiaryService() DiaryService {
	return &diaryService{
		diaryRepo: repository.NewDiary(),
	}
}

func (d *diaryService) GetDiary(ctx context.Context, password string) ([]model.Diary, error) {
	var diaries []model.Diary
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {

		isSuper := false

		diaryPs, err := d.diaryRepo.FindDiaryPs(ctx, tx)
		if err != nil {
			return err
		}

		for _, ps := range diaryPs {
			if ps.Password == password {
				isSuper = true
			}
		}

		cond := repository.FindDiaryArg{
			NoLimit: true,
		}
		diaries, err = d.diaryRepo.FindDiarys(ctx, tx, cond)
		if err != nil {
			return err
		}

		if !isSuper {
			for i := range diaries {
				if !diaries[i].IsOpen {
					lenContent := len(diaries[i].Content)
					diaries[i].Content = strDouble("*", lenContent)
				}
			}
		}

		return err
	}, nil)

	return diaries, err
}

func (d *diaryService) CreateDiary(ctx context.Context, diary model.Diary) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = d.diaryRepo.CreateDiary(ctx, tx, diary)
		return err
	}, nil)

	return err
}

func (d *diaryService) UpdateDiary(ctx context.Context, diary model.Diary) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = d.diaryRepo.UpdateDiary(ctx, tx, diary)
		return err
	}, nil)

	return err
}

func (d *diaryService) DeleteDiary(ctx context.Context, diaryId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = d.diaryRepo.DeleteDiary(ctx, tx, diaryId)
		return err
	}, nil)

	return err
}

func (d *diaryService) GetDiaryPs(ctx context.Context) ([]model.DiaryPs, error) {
	var diaryPs []model.DiaryPs
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		diaryPs, err = d.diaryRepo.FindDiaryPs(ctx, tx)
		return err
	}, nil)

	return diaryPs, err
}

func (d *diaryService) CreateDiaryPs(ctx context.Context, diaryPs model.DiaryPs) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = d.diaryRepo.CreateDiaryPs(ctx, tx, diaryPs)
		return err
	}, nil)

	return err
}

func (d *diaryService) UpdateDiaryPs(ctx context.Context, diaryPs model.DiaryPs) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = d.diaryRepo.UpdateDiaryPs(ctx, tx, diaryPs)
		return err
	}, nil)

	return err
}

func (d *diaryService) DeleteDiaryPs(ctx context.Context, diaryPsId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = d.diaryRepo.DeleteDiaryPs(ctx, tx, diaryPsId)
		return err
	}, nil)

	return err
}

// strDouble 将字符串翻倍，参数：字符串、翻倍的倍数
func strDouble(str string, n int) string {
	if n == 0 {
		return str
	}
	res := str
	for i := 0; i < n; i++ {
		res += str
	}
	return res
}
