package cv

import (
	"testing"

	"github.com/DeYu666/blog-backend-service/app_new/models"
)

func TestSkill(t *testing.T) {
	skills, err := GetProjects()

	if len(skills) < 1 {
		t.Errorf("select book lists data error, it is %v, error is %v", skills, err)
		return
	}

	testName := "test_golang"

	skill := &models.SkillCv{SkillName: testName}

	err = AddSkill(skill)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
		return
	}

	info, err := GetSkills(SkillNameCv(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	skill = info[0]
	skill.SkillName = "test_modify_golang"

	err = ModifySkill(skill)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetSkills(SkillCvId(id))
	if err != nil || len(info) != 1 || info[0].SkillName != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteSkill(SkillCvId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetSkills(SkillCvId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}
