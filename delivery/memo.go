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

type MemoHandler struct {
	MemoService service.MemoService
}

func NewMemoHandler(router *gin.Engine) {
	handler := &MemoHandler{
		MemoService: service.NewMemoService(),
	}

	authRouter := router.Group("/inner").Use(middleware.JWTAuth())

	router.GET("/memo/status/:statusId", handler.GetMemoByStatusId)
	authRouter.POST("/memo", handler.CreateMemo)
	authRouter.PUT("/memo", handler.UpdateMemo)
	authRouter.DELETE("/memo/:id", handler.DeleteMemo)
}

func (m *MemoHandler) GetMemoByStatusId(ctx *gin.Context) {

	log := blog.Extract(ctx)

	param := ctx.Param("statusId")
	statusId, err := strconv.Atoi(param)
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	memo, err := m.MemoService.GetMemoByStatusId(ctx, statusId)
	if err != nil {
		log.Error("GetMemoByStatusId, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, memo)
}

func (m *MemoHandler) CreateMemo(ctx *gin.Context) {
	log := blog.Extract(ctx)

	memo, err := utils.BufferToStruct(ctx.Request.Body, &model.Memo{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = m.MemoService.CreateMemo(ctx, *memo)
	if err != nil {
		log.Error("CreateMemo, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (m *MemoHandler) UpdateMemo(ctx *gin.Context) {
	log := blog.Extract(ctx)

	memo, err := utils.BufferToStruct(ctx.Request.Body, &model.Memo{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	log.Info("UpdateMemo, memo: " + memo.Content)

	err = m.MemoService.UpdateMemo(ctx, *memo)
	if err != nil {
		log.Error("UpdateMemo, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (m *MemoHandler) DeleteMemo(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = m.MemoService.DeleteMemo(ctx, uint(id))
	if err != nil {
		log.Error("DeleteMemo, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}
