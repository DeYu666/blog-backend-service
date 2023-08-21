package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type Skill interface {
	CountSkills(ctx context.Context, tx *gorm.DB, condition FindSkillsArg) (int64, error)
	FindSkills(ctx context.Context, tx *gorm.DB, condition FindSkillsArg) ([]model.SkillCv, error)
	GetSkill(ctx context.Context, tx *gorm.DB, skillId uint) (model.SkillCv, error)
	CreateSkill(ctx context.Context, tx *gorm.DB, experience model.SkillCv) error
	UpdateSkill(ctx context.Context, tx *gorm.DB, experience model.SkillCv) error
	DeleteSkill(ctx context.Context, tx *gorm.DB, skillId uint) error
}

type skill struct{}

func NewSkill() Skill {
	return &skill{}
}

type FindSkillsArg struct {
	IDs            []uint
	SkillNameLikes []string
	SkillInfoLikes []string
	Offset         int
	Limit          int
	NoLimit        bool
}

func (b *skill) CountSkills(ctx context.Context, tx *gorm.DB, condition FindSkillsArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.SkillNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `skill_name` Like (%s)", RepeatWithSep("?", len(condition.SkillNameLikes), ",")))
		for _, skillName := range condition.SkillNameLikes {
			args = append(args, skillName)
		}
	}

	if condition.SkillInfoLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `skill_info` Like (%s)", RepeatWithSep("?", len(condition.SkillInfoLikes), ",")))
		for _, skillInfo := range condition.SkillInfoLikes {
			args = append(args, skillInfo)
		}
	}

	var count int64
	if err := tx.Model(&model.SkillCv{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (b *skill) FindSkills(ctx context.Context, tx *gorm.DB, condition FindSkillsArg) ([]model.SkillCv, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.SkillNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `skill_name` Like (%s)", RepeatWithSep("?", len(condition.SkillNameLikes), ",")))
		for _, skillName := range condition.SkillNameLikes {
			args = append(args, skillName)
		}
	}

	if condition.SkillInfoLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `skill_info` Like (%s)", RepeatWithSep("?", len(condition.SkillInfoLikes), ",")))
		for _, skillInfo := range condition.SkillInfoLikes {
			args = append(args, skillInfo)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `id` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var skills []model.SkillCv
	if err := tx.Model(&model.SkillCv{}).Where(query.String(), args...).Find(&skills).Error; err != nil {
		return nil, err
	}

	return skills, nil
}

func (b *skill) GetSkill(ctx context.Context, tx *gorm.DB, skillId uint) (model.SkillCv, error) {
	var skill model.SkillCv
	if err := tx.Model(&model.SkillCv{}).Where("`id` = ?", skillId).First(&skill).Error; err != nil {
		return model.SkillCv{}, err
	}

	return skill, nil
}

func (b *skill) CreateSkill(ctx context.Context, tx *gorm.DB, skill model.SkillCv) error {
	if err := tx.Model(&model.SkillCv{}).Create(&skill).Error; err != nil {
		return err
	}

	return nil
}

func (b *skill) UpdateSkill(ctx context.Context, tx *gorm.DB, skill model.SkillCv) error {
	if err := tx.Model(&model.SkillCv{}).Where("`id` = ?", skill.ID.ID).Updates(&skill).Error; err != nil {
		return err
	}

	return nil
}

func (b *skill) DeleteSkill(ctx context.Context, tx *gorm.DB, skillId uint) error {
	if err := tx.Model(&model.SkillCv{}).Where("`id` = ?", skillId).Delete(&model.SkillCv{}).Error; err != nil {
		return err
	}

	return nil
}
