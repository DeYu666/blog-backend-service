package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type Diary interface {
	CountDiarys(ctx context.Context, tx *gorm.DB, condition FindDiaryArg) (int64, error)
	FindDiarys(ctx context.Context, tx *gorm.DB, condition FindDiaryArg) ([]model.Diary, error)
	GetDiary(ctx context.Context, tx *gorm.DB, diaryId uint) (model.Diary, error)
	CreateDiary(ctx context.Context, tx *gorm.DB, diary model.Diary) error
	UpdateDiary(ctx context.Context, tx *gorm.DB, diary model.Diary) error
	DeleteDiary(ctx context.Context, tx *gorm.DB, diaryId uint) error

	FindDiaryPs(ctx context.Context, tx *gorm.DB) ([]model.DiaryPs, error)
	GetDiaryPs(ctx context.Context, tx *gorm.DB, diaryId uint) (model.DiaryPs, error)
	CreateDiaryPs(ctx context.Context, tx *gorm.DB, diaryPs model.DiaryPs) error
	UpdateDiaryPs(ctx context.Context, tx *gorm.DB, diaryPs model.DiaryPs) error
	DeleteDiaryPs(ctx context.Context, tx *gorm.DB, diaryPsId uint) error
}

type diary struct{}

func NewDiary() Diary {
	return &diary{}
}

type FindDiaryArg struct {
	IDs          []uint
	ContentLikes []string
	Offset       int
	Limit        int
	NoLimit      bool
}

func (b *diary) CountDiarys(ctx context.Context, tx *gorm.DB, condition FindDiaryArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString(" 1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.ContentLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `content` Like (%s)", RepeatWithSep("?", len(condition.ContentLikes), ",")))
		for _, content := range condition.ContentLikes {
			args = append(args, content)
		}
	}

	var count int64
	if err := tx.Model(&model.Diary{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (b *diary) FindDiarys(ctx context.Context, tx *gorm.DB, condition FindDiaryArg) ([]model.Diary, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.ContentLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `content` Like (%s)", RepeatWithSep("?", len(condition.ContentLikes), ",")))
		for _, content := range condition.ContentLikes {
			args = append(args, content)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `id` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var diarys []model.Diary
	if err := tx.Model(&model.Diary{}).Where(query.String(), args...).Find(&diarys).Error; err != nil {
		return nil, err
	}

	return diarys, nil
}

func (b *diary) GetDiary(ctx context.Context, tx *gorm.DB, diaryId uint) (model.Diary, error) {
	var diary model.Diary
	if err := tx.Model(&model.Diary{}).Where("`id` = ?", diaryId).First(&diary).Error; err != nil {
		return model.Diary{}, err
	}

	return diary, nil
}

func (b *diary) CreateDiary(ctx context.Context, tx *gorm.DB, diary model.Diary) error {
	if err := tx.Model(&model.Diary{}).Create(&diary).Error; err != nil {
		return err
	}

	return nil
}

func (b *diary) UpdateDiary(ctx context.Context, tx *gorm.DB, diary model.Diary) error {
	if err := tx.Model(&model.Diary{}).Where("`id` = ?", diary.ID.ID).Updates(&diary).Error; err != nil {
		return err
	}

	return nil
}

func (b *diary) DeleteDiary(ctx context.Context, tx *gorm.DB, diaryId uint) error {
	if err := tx.Model(&model.Diary{}).Where("`id` = ?", diaryId).Delete(&model.Diary{}).Error; err != nil {
		return err
	}

	return nil
}

func (b *diary) FindDiaryPs(ctx context.Context, tx *gorm.DB) ([]model.DiaryPs, error) {
	var diaryPs []model.DiaryPs
	if err := tx.Model(&model.DiaryPs{}).Find(&diaryPs).Error; err != nil {
		return nil, err
	}

	return diaryPs, nil
}

func (b *diary) GetDiaryPs(ctx context.Context, tx *gorm.DB, diaryId uint) (model.DiaryPs, error) {
	var diaryPs model.DiaryPs
	if err := tx.Model(&model.DiaryPs{}).Where("`diary_id` = ?", diaryId).Find(&diaryPs).Error; err != nil {
		return diaryPs, err
	}

	return diaryPs, nil
}

func (b *diary) CreateDiaryPs(ctx context.Context, tx *gorm.DB, diaryPs model.DiaryPs) error {
	if err := tx.Model(&model.DiaryPs{}).Create(&diaryPs).Error; err != nil {
		return err
	}

	return nil
}

func (b *diary) UpdateDiaryPs(ctx context.Context, tx *gorm.DB, diaryPs model.DiaryPs) error {
	if err := tx.Model(&model.DiaryPs{}).Where("`id` = ?", diaryPs.ID.ID).Updates(&diaryPs).Error; err != nil {
		return err
	}

	return nil
}

func (b *diary) DeleteDiaryPs(ctx context.Context, tx *gorm.DB, diaryPsId uint) error {
	if err := tx.Model(&model.DiaryPs{}).Where("`id` = ?", diaryPsId).Delete(&model.DiaryPs{}).Error; err != nil {
		return err
	}

	return nil
}
