package delivery

import (
	"strconv"

	"github.com/DeYu666/blog-backend-service/delivery/middleware"
	"github.com/DeYu666/blog-backend-service/delivery/response"
	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/service"
	"github.com/DeYu666/blog-backend-service/utils"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	ProfileService service.ProfileService
}

func NewProfileHandler(router *gin.Engine) {
	handler := &ProfileHandler{
		ProfileService: service.NewProfile(),
	}

	authRouter := router.Group("/inner").Use(middleware.JWTAuth())

	router.GET("/profile/experience", handler.GetExperience)
	authRouter.POST("/profile/experience", handler.CreateExperience)
	authRouter.PUT("/profile/experience", handler.UpdateExperience)
	authRouter.DELETE("/profile/experience/:id", handler.DeleteExperience)

	router.GET("/profile/skill", handler.GetSkill)
	authRouter.POST("/profile/skill", handler.CreateSkill)
	authRouter.PUT("/profile/skill", handler.UpdateSkill)
	authRouter.DELETE("/profile/skill/:id", handler.DeleteSkill)

	router.GET("/profile/project", handler.GetProject)
	router.GET("/profile/projectById/:id", handler.GetProjectById)
	authRouter.POST("/profile/project", handler.CreateProject)
	authRouter.PUT("/profile/project", handler.UpdateProject)
	authRouter.DELETE("/profile/project/:id", handler.DeleteProject)

	authRouter.GET("/profile/projectPs", handler.GetProjectPs)
	authRouter.POST("/profile/projectPs", handler.CreateProjectPs)
	authRouter.PUT("/profile/projectPs", handler.UpdateProjectPs)
	authRouter.DELETE("/profile/projectPs/:id", handler.DeleteProjectPs)
}

func (p *ProfileHandler) GetExperience(ctx *gin.Context) {
	log := blog.Extract(ctx)

	experiences, err := p.ProfileService.GetExperience(ctx)
	if err != nil {
		log.Error("GetExperience, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, experiences)
}

func (p *ProfileHandler) CreateExperience(ctx *gin.Context) {
	log := blog.Extract(ctx)

	experience, err := utils.BufferToStruct(ctx.Request.Body, &model.ExperienceCv{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.CreateExperience(ctx, *experience)
	if err != nil {
		log.Error("CreateExperience, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) UpdateExperience(ctx *gin.Context) {
	log := blog.Extract(ctx)

	experience, err := utils.BufferToStruct(ctx.Request.Body, &model.ExperienceCv{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.UpdateExperience(ctx, *experience)
	if err != nil {
		log.Error("UpdateExperience, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) DeleteExperience(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.DeleteExperience(ctx, uint(id))
	if err != nil {
		log.Error("DeleteExperience, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) GetSkill(ctx *gin.Context) {
	log := blog.Extract(ctx)

	skills, err := p.ProfileService.GetSkill(ctx)
	if err != nil {
		log.Error("GetSkill, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, skills)
}

func (p *ProfileHandler) CreateSkill(ctx *gin.Context) {
	log := blog.Extract(ctx)

	skill, err := utils.BufferToStruct(ctx.Request.Body, &model.SkillCv{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.CreateSkill(ctx, *skill)
	if err != nil {
		log.Error("CreateSkill, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) UpdateSkill(ctx *gin.Context) {
	log := blog.Extract(ctx)

	skill, err := utils.BufferToStruct(ctx.Request.Body, &model.SkillCv{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.UpdateSkill(ctx, *skill)
	if err != nil {
		log.Error("UpdateSkill, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) DeleteSkill(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.DeleteSkill(ctx, uint(id))
	if err != nil {
		log.Error("DeleteSkill, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) GetProject(ctx *gin.Context) {
	log := blog.Extract(ctx)

	projects, err := p.ProfileService.GetProject(ctx)
	if err != nil {
		log.Error("GetProject, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, projects)
}

func (p *ProfileHandler) GetProjectById(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	projectId, err := strconv.Atoi(param)
	if err != nil {
		log.Error("GetProjectById, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	log.Debug("GetProjectById, projectId " + param)

	projects, err := p.ProfileService.GetProjectById(ctx, uint(projectId))
	if err != nil {
		log.Error("GetProjectById, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, projects)
}

func (p *ProfileHandler) CreateProject(ctx *gin.Context) {
	log := blog.Extract(ctx)

	project, err := utils.BufferToStruct(ctx.Request.Body, &model.ProjectCv{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.CreateProject(ctx, *project)
	if err != nil {
		log.Error("CreateProject, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) UpdateProject(ctx *gin.Context) {
	log := blog.Extract(ctx)

	project, err := utils.BufferToStruct(ctx.Request.Body, &model.ProjectCv{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.UpdateProject(ctx, *project)
	if err != nil {
		log.Error("UpdateProject, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) DeleteProject(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}
	err = p.ProfileService.DeleteProject(ctx, uint(id))
	if err != nil {
		log.Error("DeleteProject, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) GetProjectPs(ctx *gin.Context) {
	log := blog.Extract(ctx)

	projects, err := p.ProfileService.GetProjectPs(ctx)
	if err != nil {
		log.Error("GetProjectPs, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, projects)
}

func (p *ProfileHandler) CreateProjectPs(ctx *gin.Context) {
	log := blog.Extract(ctx)

	project, err := utils.BufferToStruct(ctx.Request.Body, &model.ProjectCvPs{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.CreateProjectPs(ctx, *project)
	if err != nil {
		log.Error("CreateProjectPs, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) UpdateProjectPs(ctx *gin.Context) {
	log := blog.Extract(ctx)

	project, err := utils.BufferToStruct(ctx.Request.Body, &model.ProjectCvPs{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.UpdateProjectPs(ctx, *project)
	if err != nil {
		log.Error("UpdateProjectPs, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (p *ProfileHandler) DeleteProjectPs(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = p.ProfileService.DeleteProjectPs(ctx, uint(id))
	if err != nil {
		log.Error("DeleteProjectPs, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}
