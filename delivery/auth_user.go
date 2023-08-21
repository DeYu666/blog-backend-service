package delivery

import (
	"fmt"
	"strconv"

	"github.com/DeYu666/blog-backend-service/delivery/middleware"
	"github.com/DeYu666/blog-backend-service/delivery/response"
	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/model"
	"github.com/DeYu666/blog-backend-service/service"
	"github.com/DeYu666/blog-backend-service/utils"
	"github.com/gin-gonic/gin"
)

type AuthUserHandler struct {
	AuthUserService service.AuthUserService
}

func NewAuthUserHandler(router *gin.Engine) {
	handler := &AuthUserHandler{
		AuthUserService: service.NewAuthUserService(),
	}

	authRouter := router.Group("/inner").Use(middleware.JWTAuth())

	router.POST("/login", handler.Login)

	authRouter.GET("/users", handler.GetUsers)
	authRouter.POST("/register", handler.Register)
	authRouter.PUT("/user", handler.UpdateUser)
	authRouter.DELETE("/user/:id", handler.DeleteUser)
	// route.POST("/logout", handler.Logout)
	// route.POST("/refresh_token", handler.RefreshToken)
}

type ResponseLogin struct {
	Token string `json:"access_token"`
}

func (a *AuthUserHandler) Login(ctx *gin.Context) {
	log := blog.Extract(ctx)

	log.Sugar().Info("Login", fmt.Sprintf("%v", ctx.Request.Body))

	user, err := utils.BufferToStruct(ctx.Request.Body, &model.AuthUser{})
	if err != nil {
		log.Error("Login, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	authUser, err := a.AuthUserService.Login(ctx, *user)
	if err != nil {
		log.Error("Login, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	if authUser.Username == "" {
		response.Fail(ctx, 10003, "用户名或密码错误")
		return
	}
	tokenStr, err := service.JwtService.GenerateToken(authUser)
	if err != nil {
		response.Fail(ctx, 10004, err.Error())
		return
	}

	resp := ResponseLogin{
		Token: tokenStr,
	}

	response.Success(ctx, resp)
}

func (a *AuthUserHandler) GetUsers(ctx *gin.Context) {
	log := blog.Extract(ctx)

	users, err := a.AuthUserService.GetUsers(ctx)
	if err != nil {
		log.Error("GetUsers, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, users)
}

func (a *AuthUserHandler) Register(ctx *gin.Context) {
	log := blog.Extract(ctx)

	user, err := utils.BufferToStruct(ctx.Request.Body, &model.AuthUser{})
	if err != nil {
		log.Error("Register, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = a.AuthUserService.Register(ctx, *user)
	if err != nil {
		log.Error("Register, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (a *AuthUserHandler) UpdateUser(ctx *gin.Context) {
	log := blog.Extract(ctx)

	user, err := utils.BufferToStruct(ctx.Request.Body, &model.AuthUser{})
	if err != nil {
		log.Error("Register, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = a.AuthUserService.UpdateUser(ctx, *user)
	if err != nil {
		log.Error("Register, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (a *AuthUserHandler) DeleteUser(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = a.AuthUserService.DeleteUser(ctx, uint(id))
	if err != nil {
		log.Error("DeleteUser, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}
