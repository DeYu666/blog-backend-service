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

type DiaryHandler struct {
	DiaryService service.DiaryService
}

func NewDiaryHandler(router *gin.Engine) {
	handler := &DiaryHandler{
		DiaryService: service.NewDiaryService(),
	}

	authRouter := router.Group("/inner").Use(middleware.JWTAuth())

	router.GET("/diary", handler.GetDiary)
	authRouter.POST("/diary", handler.CreateDiary)
	authRouter.PUT("/diary", handler.UpdateDiary)
	authRouter.DELETE("/diary/:id", handler.DeleteDiary)

	authRouter.GET("/diary_ps", handler.GetDiaryPs)
	authRouter.POST("/diary_ps", handler.CreateDiaryPs)
	authRouter.PUT("/diary_ps", handler.UpdateDiaryPs)
	authRouter.DELETE("/diary_ps/:id", handler.DeleteDiaryPs)
}

func (d *DiaryHandler) GetDiary(ctx *gin.Context) {
	log := blog.Extract(ctx)

	diaryPassword := ctx.Query("diaryPassword")

	tokenStr := ctx.Request.Header.Get("Authorization")
	token, err := service.JwtService.ValidateToken(tokenStr)
	if token == nil || err != nil {
		log.Info("GetDiary, token " + tokenStr)
	} else {
		pds, err := d.DiaryService.GetDiaryPs(ctx)
		if err != nil {
			log.Error("GetDiaryPs, err " + err.Error())
			response.FailByError(ctx, response.Errors.ServeError)
			return
		}
		if len(pds) > 0 {
			diaryPassword = pds[0].Password
		}
	}

	diaries, err := d.DiaryService.GetDiary(ctx, diaryPassword)
	if err != nil {
		log.Error("GetDiary, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, diaries)
}

func (d *DiaryHandler) CreateDiary(ctx *gin.Context) {
	log := blog.Extract(ctx)

	diary, err := utils.BufferToStruct(ctx.Request.Body, &model.Diary{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = d.DiaryService.CreateDiary(ctx, *diary)
	if err != nil {
		log.Error("CreateDiary, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (d *DiaryHandler) UpdateDiary(ctx *gin.Context) {
	log := blog.Extract(ctx)

	diary, err := utils.BufferToStruct(ctx.Request.Body, &model.Diary{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = d.DiaryService.UpdateDiary(ctx, *diary)
	if err != nil {
		log.Error("UpdateDiary, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (d *DiaryHandler) DeleteDiary(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = d.DiaryService.DeleteDiary(ctx, uint(id))
	if err != nil {
		log.Error("DeleteDiary, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (d *DiaryHandler) GetDiaryPs(ctx *gin.Context) {
	log := blog.Extract(ctx)

	diaries, err := d.DiaryService.GetDiaryPs(ctx)
	if err != nil {
		log.Error("GetDiaryPs, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, diaries)
}

func (d *DiaryHandler) CreateDiaryPs(ctx *gin.Context) {
	log := blog.Extract(ctx)

	diaryPs, err := utils.BufferToStruct(ctx.Request.Body, &model.DiaryPs{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = d.DiaryService.CreateDiaryPs(ctx, *diaryPs)
	if err != nil {
		log.Error("CreateDiaryPs, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (d *DiaryHandler) UpdateDiaryPs(ctx *gin.Context) {
	log := blog.Extract(ctx)

	diaryPs, err := utils.BufferToStruct(ctx.Request.Body, &model.DiaryPs{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = d.DiaryService.UpdateDiaryPs(ctx, *diaryPs)
	if err != nil {
		log.Error("UpdateDiaryPs, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (d *DiaryHandler) DeleteDiaryPs(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = d.DiaryService.DeleteDiaryPs(ctx, uint(id))
	if err != nil {
		log.Error("DeleteDiaryPs, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}
