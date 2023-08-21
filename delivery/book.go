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

type BookHandler struct {
	BookService service.BookService
}

func NewBookHandler(router *gin.Engine) {
	handler := &BookHandler{
		BookService: service.NewBookService(),
	}

	authRouter := router.Group("/inner").Use(middleware.JWTAuth())

	router.GET("/books_list", handler.GetBooksList)
	authRouter.POST("/books_list", handler.CreateBooksList)
	authRouter.PUT("/books_list", handler.UpdateBooksList)
	authRouter.DELETE("/books_list/:id", handler.DeleteBooksList)

	router.GET("/book_content/:bookId", handler.GetBookContentByBookId)
	authRouter.POST("/book_content", handler.CreateBookContent)
	authRouter.PUT("/book_content", handler.UpdateBookContent)
	authRouter.DELETE("/book_content/:id", handler.DeleteBookContent)
}

func (b *BookHandler) GetBooksList(ctx *gin.Context) {

	log := blog.Extract(ctx)

	cateId, err := strconv.Atoi(ctx.Query("categoryId"))
	if err != nil {
		cateId = -1
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		offset = 1
	}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 0
	}

	bookStatus := ""
	if cateId == 0 {
		bookStatus = "在看"
	} else if cateId == 1 {
		bookStatus = "已看"
	} else if cateId == 2 {
		bookStatus = "未读"
	}

	booksList, err := b.BookService.GetBooksList(ctx, bookStatus, int32(offset), int32(limit))
	if err != nil {
		log.Error("GetBooksList, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, booksList)
}

func (b *BookHandler) CreateBooksList(ctx *gin.Context) {
	log := blog.Extract(ctx)

	book, err := utils.BufferToStruct(ctx.Request.Body, &model.BooksList{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = b.BookService.CreateBooksList(ctx, *book)
	if err != nil {
		log.Error("CreateBooksList, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (b *BookHandler) UpdateBooksList(ctx *gin.Context) {
	log := blog.Extract(ctx)

	book, err := utils.BufferToStruct(ctx.Request.Body, &model.BooksList{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = b.BookService.UpdateBooksList(ctx, *book)
	if err != nil {
		log.Error("UpdateBooksList, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (b *BookHandler) DeleteBooksList(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = b.BookService.DeleteBooksList(ctx, uint(id))
	if err != nil {
		log.Error("DeleteBooksList, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (b *BookHandler) GetBookContentByBookId(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("bookId")
	bookId, err := strconv.Atoi(param)
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	bookContents, err := b.BookService.GetBookContentByBookId(ctx, uint(bookId))
	if err != nil {
		log.Error("GetBookContentByBookId, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, bookContents)
}

func (b *BookHandler) CreateBookContent(ctx *gin.Context) {
	log := blog.Extract(ctx)

	bookContent, err := utils.BufferToStruct(ctx.Request.Body, &model.BookContent{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = b.BookService.CreateBookContent(ctx, *bookContent)
	if err != nil {
		log.Error("CreateBookContent, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (b *BookHandler) UpdateBookContent(ctx *gin.Context) {
	log := blog.Extract(ctx)

	bookContent, err := utils.BufferToStruct(ctx.Request.Body, &model.BookContent{})
	if err != nil {
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = b.BookService.UpdateBookContent(ctx, *bookContent)
	if err != nil {
		log.Error("UpdateBookContent, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}

func (b *BookHandler) DeleteBookContent(ctx *gin.Context) {
	log := blog.Extract(ctx)

	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		log.Error("DeleteCate, err " + err.Error())
		response.ValidateFail(ctx, err.Error())
		return
	}

	err = b.BookService.DeleteBookContent(ctx, uint(id))
	if err != nil {
		log.Error("DeleteBookContent, err " + err.Error())
		response.FailByError(ctx, response.Errors.ServeError)
		return
	}

	response.Success(ctx, "success")
}
