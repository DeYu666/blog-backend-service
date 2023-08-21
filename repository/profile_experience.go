package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type Experience interface {
	CountExperiences(ctx context.Context, tx *gorm.DB, condition FindExperiencesArg) (int64, error)
	FindExperiences(ctx context.Context, tx *gorm.DB, condition FindExperiencesArg) ([]model.ExperienceCv, error)
	GetExperience(ctx context.Context, tx *gorm.DB, experienceId uint) (model.ExperienceCv, error)
	CreateExperience(ctx context.Context, tx *gorm.DB, experience model.ExperienceCv) error
	UpdateExperience(ctx context.Context, tx *gorm.DB, experience model.ExperienceCv) error
	DeleteExperience(ctx context.Context, tx *gorm.DB, experienceId uint) error
}

type experience struct{}

func NewExperience() Experience {
	return &experience{}
}

type FindExperiencesArg struct {
	IDs                 []uint
	ExperienceNameLikes []string
	WorkNameLikes       []string
	WorkInfoLikes       []string
	Offset              int
	Limit               int
	NoLimit             bool
}

func (b *experience) CountExperiences(ctx context.Context, tx *gorm.DB, condition FindExperiencesArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.ExperienceNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `experience_name` Like (%s)", RepeatWithSep("?", len(condition.ExperienceNameLikes), ",")))
		for _, experienceName := range condition.ExperienceNameLikes {
			args = append(args, experienceName)
		}
	}

	if condition.WorkNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `work_name` Like (%s)", RepeatWithSep("?", len(condition.WorkNameLikes), ",")))
		for _, workName := range condition.WorkNameLikes {
			args = append(args, workName)
		}
	}

	if condition.WorkInfoLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `work_info` Like (%s)", RepeatWithSep("?", len(condition.WorkInfoLikes), ",")))
		for _, workInfo := range condition.WorkInfoLikes {
			args = append(args, workInfo)
		}
	}

	var count int64
	if err := tx.Model(&model.ExperienceCv{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (b *experience) FindExperiences(ctx context.Context, tx *gorm.DB, condition FindExperiencesArg) ([]model.ExperienceCv, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.ExperienceNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `experience_name` Like (%s)", RepeatWithSep("?", len(condition.ExperienceNameLikes), ",")))
		for _, experienceName := range condition.ExperienceNameLikes {
			args = append(args, experienceName)
		}
	}

	if condition.WorkNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `work_name` Like (%s)", RepeatWithSep("?", len(condition.WorkNameLikes), ",")))
		for _, workName := range condition.WorkNameLikes {
			args = append(args, workName)
		}
	}

	if condition.WorkInfoLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `work_info` Like (%s)", RepeatWithSep("?", len(condition.WorkInfoLikes), ",")))
		for _, workInfo := range condition.WorkInfoLikes {
			args = append(args, workInfo)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var experiences []model.ExperienceCv
	if err := tx.Model(&model.ExperienceCv{}).Where(query.String(), args...).Find(&experiences).Error; err != nil {
		return nil, err
	}

	return experiences, nil
}

func (b *experience) GetExperience(ctx context.Context, tx *gorm.DB, experienceId uint) (model.ExperienceCv, error) {
	var experience model.ExperienceCv
	if err := tx.Model(&model.ExperienceCv{}).Where("`id` = ?", experienceId).First(&experience).Error; err != nil {
		return model.ExperienceCv{}, err
	}

	return experience, nil
}

func (b *experience) CreateExperience(ctx context.Context, tx *gorm.DB, experience model.ExperienceCv) error {
	if err := tx.Model(&model.ExperienceCv{}).Create(&experience).Error; err != nil {
		return err
	}

	return nil
}

func (b *experience) UpdateExperience(ctx context.Context, tx *gorm.DB, experience model.ExperienceCv) error {
	if err := tx.Model(&model.ExperienceCv{}).Where("`id` = ?", experience.ID.ID).Updates(&experience).Error; err != nil {
		return err
	}

	return nil
}

func (b *experience) DeleteExperience(ctx context.Context, tx *gorm.DB, experienceId uint) error {
	if err := tx.Model(&model.ExperienceCv{}).Where("`id` = ?", experienceId).Delete(&model.ExperienceCv{}).Error; err != nil {
		return err
	}

	return nil
}
