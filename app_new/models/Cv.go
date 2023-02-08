package models

import "time"

// ExperienceCv 工作经历
type ExperienceCv struct {
	ID
	WorkYear       time.Time `gorm:"type:time" json:"work_year"`
	EnterpriseName string    `gorm:"type:text" json:"enterprise_name"`
	WorkName       string    `json:"work_name"`
	WorkInfo       string    `json:"work_info"`
}

func (ExperienceCv) GetTableName() string {
	return "experience_cvs"
}

// SkillCv 技能简历
type SkillCv struct {
	ID
	SkillName    string `json:"skill_name"`
	SkillMastery int    `json:"skill_mastery"`
	SkillIntro   string `json:"skill_intro"`
}

func (SkillCv) GetTableName() string {
	return "skill_cvs"
}

// ProjectCv 项目经历
type ProjectCv struct {
	ID
	ProjectName    string    `json:"project_name"`
	ProjectContent string    `gorm:"type:text" json:"project_content"`
	ProjectImgSrc  string    `json:"project_img_src"`
	Abstract       string    `json:"abstract"`
	IsOpen         bool      `json:"is_open"`
	PublishTime    time.Time `gorm:"type:time" json:"publish_time"`
}

func (ProjectCv) GetTableName() string {
	return "project_cvs"
}

type ProjectCvPs struct {
	ID
	Password string `json:"password"`
}

func (ProjectCvPs) GetTableName() string {
	return "project_cv_ps"
}
