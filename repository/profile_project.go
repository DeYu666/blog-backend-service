package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/DeYu666/blog-backend-service/model"
	"gorm.io/gorm"
)

type Project interface {
	CountProjects(ctx context.Context, tx *gorm.DB, condition FindProjectArg) (int64, error)
	FindProjects(ctx context.Context, tx *gorm.DB, condition FindProjectArg) ([]model.ProjectCv, error)
	GetProject(ctx context.Context, tx *gorm.DB, projectId uint) (model.ProjectCv, error)
	CreateProject(ctx context.Context, tx *gorm.DB, project model.ProjectCv) error
	UpdateProject(ctx context.Context, tx *gorm.DB, project model.ProjectCv) error
	DeleteProject(ctx context.Context, tx *gorm.DB, projectId uint) error

	FindProjectPwd(ctx context.Context, tx *gorm.DB) ([]model.ProjectCvPs, error)
	GetProjectPwd(ctx context.Context, tx *gorm.DB, projectPsId uint) (model.ProjectCvPs, error)
	CreatePwd(ctx context.Context, tx *gorm.DB, projectPs model.ProjectCvPs) error
	UpdateProjectPwd(ctx context.Context, tx *gorm.DB, projectPs model.ProjectCvPs) error
	DeleteProjectPwd(ctx context.Context, tx *gorm.DB, projectPsId uint) error
}

type project struct{}

func NewProject() Project {
	return &project{}
}

type FindProjectArg struct {
	IDs              []uint
	ProjectNames     []string
	ProjectNameLikes []string
	Offset           int32
	Limit            int32
	NoLimit          bool
}

func (p *project) CountProjects(ctx context.Context, tx *gorm.DB, condition FindProjectArg) (int64, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.IDs != nil {
		query.WriteString(fmt.Sprintf(" AND `id` IN (%s)", RepeatWithSep("?", len(condition.IDs), ",")))
		for _, id := range condition.IDs {
			args = append(args, id)
		}
	}

	if condition.ProjectNames != nil {
		query.WriteString(fmt.Sprintf(" AND `project_name` IN (%s)", RepeatWithSep("?", len(condition.ProjectNames), ",")))
		for _, projectName := range condition.ProjectNames {
			args = append(args, projectName)
		}
	}

	if condition.ProjectNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `project_name` IN (%s)", RepeatWithSep("?", len(condition.ProjectNameLikes), ",")))
		for _, projectName := range condition.ProjectNameLikes {
			args = append(args, projectName)
		}
	}

	var count int64
	if err := tx.Model(&model.ProjectCv{}).Where(query.String(), args...).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (p *project) FindProjects(ctx context.Context, tx *gorm.DB, condition FindProjectArg) ([]model.ProjectCv, error) {
	query := strings.Builder{}
	args := []interface{}{}

	query.WriteString("1 = 1")
	if condition.ProjectNames != nil {
		query.WriteString(fmt.Sprintf(" AND `project_name` IN (%s)", RepeatWithSep("?", len(condition.ProjectNames), ",")))
		for _, projectName := range condition.ProjectNames {
			args = append(args, projectName)
		}
	}
	if condition.ProjectNameLikes != nil {
		query.WriteString(fmt.Sprintf(" AND `project_name` IN (%s)", RepeatWithSep("?", len(condition.ProjectNameLikes), ",")))
		for _, projectName := range condition.ProjectNameLikes {
			args = append(args, projectName)
		}
	}

	if !condition.NoLimit {
		query.WriteString(" ORDER BY `id` DESC")
		query.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, condition.Limit, condition.Offset)
	}

	var projects []model.ProjectCv
	if err := tx.Where(query.String(), args...).Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *project) GetProject(ctx context.Context, tx *gorm.DB, projectId uint) (model.ProjectCv, error) {
	var project model.ProjectCv
	if err := tx.Model(&model.ProjectCv{}).Where("`id` = ?", projectId).First(&project).Error; err != nil {
		return model.ProjectCv{}, err
	}

	return project, nil
}

func (p *project) CreateProject(ctx context.Context, tx *gorm.DB, project model.ProjectCv) error {
	if err := tx.Model(&model.ProjectCv{}).Create(&project).Error; err != nil {
		return err
	}

	return nil
}

func (p *project) UpdateProject(ctx context.Context, tx *gorm.DB, project model.ProjectCv) error {
	if err := tx.Model(&model.ProjectCv{}).Where("`id` = ?", project.ID.ID).Updates(&project).Error; err != nil {
		return err
	}

	return nil
}

func (p *project) DeleteProject(ctx context.Context, tx *gorm.DB, projectId uint) error {
	if err := tx.Model(&model.ProjectCv{}).Where("`id` = ?", projectId).Delete(&model.ProjectCv{}).Error; err != nil {
		return err
	}

	return nil
}

func (p *project) FindProjectPwd(ctx context.Context, tx *gorm.DB) ([]model.ProjectCvPs, error) {
	var projectPwds []model.ProjectCvPs
	if err := tx.Model(&model.ProjectCvPs{}).Find(&projectPwds).Error; err != nil {
		return nil, err
	}

	return projectPwds, nil
}

func (p *project) GetProjectPwd(ctx context.Context, tx *gorm.DB, projectPsId uint) (model.ProjectCvPs, error) {
	var projectPs model.ProjectCvPs
	if err := tx.Model(&model.ProjectCvPs{}).Where("`id` = ?", projectPsId).First(&projectPs).Error; err != nil {
		return model.ProjectCvPs{}, err
	}

	return projectPs, nil
}

func (p *project) CreatePwd(ctx context.Context, tx *gorm.DB, projectPs model.ProjectCvPs) error {
	if err := tx.Model(&model.ProjectCvPs{}).Create(&projectPs).Error; err != nil {
		return err
	}

	return nil
}

func (p *project) UpdateProjectPwd(ctx context.Context, tx *gorm.DB, projectPs model.ProjectCvPs) error {
	if err := tx.Model(&model.ProjectCvPs{}).Where("`id` = ?", projectPs.ID.ID).Updates(&projectPs).Error; err != nil {
		return err
	}

	return nil
}

func (p *project) DeleteProjectPwd(ctx context.Context, tx *gorm.DB, projectPsId uint) error {
	if err := tx.Model(&model.ProjectCvPs{}).Where("`id` = ?", projectPsId).Delete(&model.ProjectCvPs{}).Error; err != nil {
		return err
	}

	return nil
}
