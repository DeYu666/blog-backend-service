package backstage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/common/tool"
	"github.com/DeYu666/blog-backend-service/app_new/dao/cv"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"github.com/gin-gonic/gin"
)

// AutoCvMigrate 自动更新表结构
func AutoCvMigrate(c *gin.Context) {
	result := global.App.DB.AutoMigrate(&models.ExperienceCv{}, &models.SkillCv{}, &models.ProjectCv{})
	//result := global.App.DB.AutoMigrate(&models.ExperienceCv{}, &models.SkillCv{}, &models.ProjectCv{}, &models.ProjectCvPs{})
	if result != nil {
		response.SqlFail(c, result.Error())
	} else {
		response.Success(c, "迁移数据成功")

	}
}

func GetCvExperience(c *gin.Context) {
	experienceCv, err := cv.GetExperiences()

	if err != nil {
		response.Fail(c, 10000, err.Error())
	}

	response.Success(c, experienceCv)
}

// AddCvExperience 添加工作经历
func AddCvExperience(c *gin.Context) {

	experienceCv, err := tool.BufferToStruct(c.Request.Body, &models.ExperienceCv{})

	fmt.Println(experienceCv)

	if err != nil {
		global.App.Log.Error(err.Error())
		response.ValidateFail(c, err.Error())
		return
	}

	result := cv.AddExperience(experienceCv)

	if result != nil {
		response.SqlFail(c, result.Error())
	} else {
		response.Success(c, "成功")
	}
}

// ModifyCvExperience 修改工作经历
func ModifyCvExperience(c *gin.Context) {

	experienceCv, err := tool.BufferToStruct(c.Request.Body, &models.ExperienceCv{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	result := cv.ModifyExperience(experienceCv)

	if result != nil {
		response.SqlFail(c, result.Error())
	} else {
		response.Success(c, "成功")
	}
}

// DeleteCvExperience 删除工作经历
func DeleteCvExperience(c *gin.Context) {

	experienceCv, err := tool.BufferToStruct(c.Request.Body, &models.ExperienceCv{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	result := cv.DeleteExperience(cv.ExperienceCvId(experienceCv.ID.ID))
	if result != nil {
		response.SqlFail(c, result.Error())
	} else {
		response.Success(c, "成功")
	}
}

// ----- 技能简历 ----

func GetCvSkill(c *gin.Context) {

	skillCv, err := cv.GetSkills()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, skillCv)
}

// AddCvSkill 添加技能简历
func AddCvSkill(c *gin.Context) {

	skillCv, err := tool.BufferToStruct(c.Request.Body, &models.SkillCv{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	result := cv.AddSkill(skillCv)
	if result != nil {
		response.SqlFail(c, result.Error())
		return
	}

	response.Success(c, "成功")

}

// ModifyCvSkill 修改技能简历
func ModifyCvSkill(c *gin.Context) {
	skillCv, err := tool.BufferToStruct(c.Request.Body, &models.SkillCv{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	result := cv.ModifySkill(skillCv)
	if result != nil {
		response.SqlFail(c, result.Error())
		return
	}

	response.Success(c, "成功")
}

// DeleteCvSkill 删除技能简历
func DeleteCvSkill(c *gin.Context) {
	skillCv, err := tool.BufferToStruct(c.Request.Body, &models.SkillCv{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	result := cv.DeleteSkill(cv.SkillCvId(skillCv.ID.ID))
	if result != nil {
		response.SqlFail(c, result.Error())
		return
	}

	response.Success(c, "成功")
}

// ----- 项目经历 ----
func GetCvProject(c *gin.Context) {

	projectCv, err := cv.GetProjects()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, projectCv)
}

func GetCvProjectByID(c *gin.Context) {

	param := c.Param("cvProjectId")

	id, err := strconv.Atoi(param)
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	projectCv, err := cv.GetProjects(cv.ProjectCvId(uint(id)))
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	if len(projectCv) == 0 {
		response.Success(c, nil)
		return
	}

	response.Success(c, projectCv[0])
}

// AddCvProject 添加项目经历
func AddCvProject(c *gin.Context) {

	projectCv, err := tool.BufferToStruct(c.Request.Body, &models.ProjectCv{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	projectCv.PublishTime = time.Now()

	result := cv.AddProject(projectCv)
	if result != nil {
		response.SqlFail(c, result.Error())
		return
	}

	response.Success(c, "成功")
}

// ModifyCvProject 修改项目经历
func ModifyCvProject(c *gin.Context) {
	projectCv, err := tool.BufferToStruct(c.Request.Body, &models.ProjectCv{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	result := cv.ModifyProject(projectCv)
	if result != nil {
		response.SqlFail(c, result.Error())
		return
	}

	response.Success(c, "成功")
}

// DeleteCvProject 删除项目经历
func DeleteCvProject(c *gin.Context) {
	projectCv, err := tool.BufferToStruct(c.Request.Body, &models.ProjectCv{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	result := cv.DeleteProject(cv.ProjectCvId(projectCv.ID.ID))
	if result != nil {
		response.SqlFail(c, result.Error())
		return
	}

	response.Success(c, "成功")
}

// 密码部分

func GetCvProjectPs(c *gin.Context) {
	projectCvPs, err := cv.GetProjectPwd()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, projectCvPs)
}

// ModifyCvProjectPs 修改项目观看密码
func ModifyCvProjectPs(c *gin.Context) {
	projectCvPs, err := tool.BufferToStruct(c.Request.Body, &models.ProjectCvPs{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	result := cv.ModifyProjectPwd(projectCvPs)
	if result != nil {
		response.SqlFail(c, result.Error())
		return
	}

	response.Success(c, "成功")

}
