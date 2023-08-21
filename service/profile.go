package service

import (
	"context"
	"time"

	"github.com/DeYu666/blog-backend-service/lib/client"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/repository"
	"gorm.io/gorm"
)

type ProfileService interface {
	GetExperience(ctx context.Context) ([]model.ExperienceCv, error)
	CreateExperience(ctx context.Context, experience model.ExperienceCv) error
	UpdateExperience(ctx context.Context, experience model.ExperienceCv) error
	DeleteExperience(ctx context.Context, experienceId uint) error

	GetSkill(ctx context.Context) ([]model.SkillCv, error)
	CreateSkill(ctx context.Context, skill model.SkillCv) error
	UpdateSkill(ctx context.Context, skill model.SkillCv) error
	DeleteSkill(ctx context.Context, skillId uint) error

	GetProject(ctx context.Context) ([]model.ProjectCv, error)
	GetProjectById(ctx context.Context, projectId uint) (model.ProjectCv, error)
	CreateProject(ctx context.Context, project model.ProjectCv) error
	UpdateProject(ctx context.Context, project model.ProjectCv) error
	DeleteProject(ctx context.Context, projectId uint) error

	GetProjectPs(ctx context.Context) ([]model.ProjectCvPs, error)
	CreateProjectPs(ctx context.Context, projectPs model.ProjectCvPs) error
	UpdateProjectPs(ctx context.Context, projectPs model.ProjectCvPs) error
	DeleteProjectPs(ctx context.Context, projectPsId uint) error
}

type profileService struct {
	experienceRepo repository.Experience
	skillRepo      repository.Skill
	projectRepo    repository.Project
}

func NewProfile() ProfileService {
	return &profileService{
		experienceRepo: repository.NewExperience(),
		skillRepo:      repository.NewSkill(),
		projectRepo:    repository.NewProject(),
	}
}

func (p *profileService) GetExperience(ctx context.Context) ([]model.ExperienceCv, error) {
	var experiences []model.ExperienceCv
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindExperiencesArg{
			NoLimit: true,
		}
		experiences, err = p.experienceRepo.FindExperiences(ctx, tx, cond)
		return err
	}, nil)

	return experiences, err
}

func (p *profileService) CreateExperience(ctx context.Context, experience model.ExperienceCv) error {
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.experienceRepo.CreateExperience(ctx, tx, experience)
		return err
	}, nil)

	return err
}

func (p *profileService) UpdateExperience(ctx context.Context, experience model.ExperienceCv) error {
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.experienceRepo.UpdateExperience(ctx, tx, experience)
		return err
	}, nil)

	return err
}

func (p *profileService) DeleteExperience(ctx context.Context, experienceId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.experienceRepo.DeleteExperience(ctx, tx, experienceId)
		return err
	}, nil)

	return err
}

func (p *profileService) GetSkill(ctx context.Context) ([]model.SkillCv, error) {
	var skills []model.SkillCv
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindSkillsArg{
			NoLimit: true,
		}
		skills, err = p.skillRepo.FindSkills(ctx, tx, cond)
		return err
	}, nil)

	return skills, err
}

func (p *profileService) CreateSkill(ctx context.Context, skill model.SkillCv) error {
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.skillRepo.CreateSkill(ctx, tx, skill)
		return err
	}, nil)

	return err
}

func (p *profileService) UpdateSkill(ctx context.Context, skill model.SkillCv) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.skillRepo.UpdateSkill(ctx, tx, skill)
		return err
	}, nil)

	return err
}

func (p *profileService) DeleteSkill(ctx context.Context, skillId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.skillRepo.DeleteSkill(ctx, tx, skillId)
		return err
	}, nil)

	return err
}

func (p *profileService) GetProject(ctx context.Context) ([]model.ProjectCv, error) {
	var projects []model.ProjectCv
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		cond := repository.FindProjectArg{
			NoLimit: true,
		}
		projects, err = p.projectRepo.FindProjects(ctx, tx, cond)
		return err
	}, nil)

	return projects, err
}

func (p *profileService) GetProjectById(ctx context.Context, projectId uint) (model.ProjectCv, error) {
	var project model.ProjectCv
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		project, err = p.projectRepo.GetProject(ctx, tx, projectId)
		return err
	}, nil)

	return project, err
}

func (p *profileService) CreateProject(ctx context.Context, project model.ProjectCv) error {

	var err error

	project.PublishTime = time.Now()

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.projectRepo.CreateProject(ctx, tx, project)
		return err
	}, nil)

	return err
}

func (p *profileService) UpdateProject(ctx context.Context, project model.ProjectCv) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.projectRepo.UpdateProject(ctx, tx, project)
		return err
	}, nil)

	return err

}

func (p *profileService) DeleteProject(ctx context.Context, projectId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.projectRepo.DeleteProject(ctx, tx, projectId)
		return err
	}, nil)

	return err
}

func (p *profileService) GetProjectPs(ctx context.Context) ([]model.ProjectCvPs, error) {
	var projects []model.ProjectCvPs
	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		projects, err = p.projectRepo.FindProjectPwd(ctx, tx)
		return err
	}, nil)

	return projects, err
}

func (p *profileService) CreateProjectPs(ctx context.Context, projectPs model.ProjectCvPs) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.projectRepo.CreatePwd(ctx, tx, projectPs)
		return err
	}, nil)

	return err
}

func (p *profileService) UpdateProjectPs(ctx context.Context, projectPs model.ProjectCvPs) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.projectRepo.UpdateProjectPwd(ctx, tx, projectPs)
		return err
	}, nil)

	return err

}

func (p *profileService) DeleteProjectPs(ctx context.Context, projectPsId uint) error {

	var err error

	err = client.Mysql.DB().Transaction(func(tx *gorm.DB) error {
		err = p.projectRepo.DeleteProjectPwd(ctx, tx, projectPsId)
		return err
	}, nil)

	return err
}
