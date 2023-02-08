package cv

import (
	"testing"
	"time"

	"github.com/DeYu666/blog-backend-service/app_new/models"
)

func TestProject(t *testing.T) {
	projects, err := GetProjects()

	if len(projects) < 1 {
		t.Errorf("select book lists data error, it is %v, error is %v", projects, err)
		return
	}

	testName := "test_golang"

	project := &models.ProjectCv{PublishTime: time.Now(), ProjectName: testName, Abstract: ""}

	err = AddProject(project)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
	}

	info, err := GetProjects(ProjectNameCv(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	project = info[0]
	project.ProjectName = "test_modify_golang"

	err = ModifyProject(project)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetProjects(ProjectCvId(id))
	if err != nil || len(info) != 1 || info[0].ProjectName != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteProject(ProjectCvId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetProjects(ProjectCvId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}

func TestProjectPwd(t *testing.T) {
	projectPwd, err := GetProjectPwd()

	if len(projectPwd) < 1 {
		t.Errorf("select project password data error, it is %v, error is %v", projectPwd, err)
		return
	}

	testName := "test_golang"

	projectPwd1 := &models.ProjectCvPs{Password: testName}

	err = AddProjectPwd(projectPwd1)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
	}

	info, err := GetProjectPwd(ProjectPassWord(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	projectPwd1 = info[0]
	projectPwd1.Password = "test_modify_golang"

	err = ModifyProjectPwd(projectPwd1)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetProjectPwd(ProjectPwdId(id))
	if err != nil || len(info) != 1 || info[0].Password != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteProjectPwd(ProjectPwdId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetProjectPwd(ProjectPwdId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}
