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

type LoveInfoHandler struct {
	LoveInfoService service.LoveInfoService
}

func NewLoveInfoHandler(router *gin.Engine) {
	handler := &LoveInfoHandler{
		LoveInfoService: service.NewLoveInfoService(),
	}

	authRouter := router.Group("/inner").Use(middleware.JWTAuth())

	router.GET("/love", handler.GetLoveInfo)
	authRouter.POST("/love", handler.CreateLoveInfo)
	authRouter.PUT("/love", handler.UpdateLoveInfo)
	authRouter.DELETE("/love/:id", handler.DeleteLoveInfo)
}

func (l *LoveInfoHandler) GetLoveInfo(ctx *gin.Context) {
	log := blog.Extract(ctx)

	loveInfo, err := l.LoveInfoService.GetLoveInfo(ctx)
	if err != nil {
		log.Error("GetLoveInfo, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, loveInfo)
}

func (l *LoveInfoHandler) CreateLoveInfo(ctx *gin.Context) {
	log := blog.Extract(ctx)

	loveInfo, err := utils.BufferToStruct(ctx.Request.Body, &model.LoveInfo{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = l.LoveInfoService.CreateLoveInfo(ctx, *loveInfo)
	if err != nil {
		log.Error("CreateLoveInfo, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (l *LoveInfoHandler) UpdateLoveInfo(ctx *gin.Context) {
	log := blog.Extract(ctx)

	loveInfo, err := utils.BufferToStruct(ctx.Request.Body, &model.LoveInfo{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = l.LoveInfoService.UpdateLoveInfo(ctx, *loveInfo)
	if err != nil {
		log.Error("UpdateLoveInfo, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (l *LoveInfoHandler) DeleteLoveInfo(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = l.LoveInfoService.DeleteLoveInfo(ctx, uint(id))
	if err != nil {
		log.Error("DeleteLoveInfo, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}
