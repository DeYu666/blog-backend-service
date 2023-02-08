package frontstage

import (
	"strconv"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/dao/cv"
	"github.com/gin-gonic/gin"
)

func GetCvExperience(c *gin.Context) {
	experienceCv, err := cv.GetExperiences()

	if err != nil {
		response.Fail(c, 10000, err.Error())
	}

	response.Success(c, experienceCv)
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
