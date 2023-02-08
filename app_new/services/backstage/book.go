package backstage

import (
	"strconv"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/common/tool"
	"github.com/DeYu666/blog-backend-service/app_new/dao/book"
	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/global"
	"github.com/gin-gonic/gin"
)

// AutoBookMigrate 自动更新表结构
func AutoBookMigrate(c *gin.Context) {
	result := global.App.DB.AutoMigrate(&models.BookContent{}, &models.BooksList{})
	c.JSON(200, result)
}

func GetBooksList(c *gin.Context) {

	booksList, err := book.GetBooksLists()
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, booksList)
}

func GetBookContent(c *gin.Context) {
	param := c.Param("bookId")

	bookId, err := strconv.Atoi(param)
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	bookContents, err := book.GetBookContent(book.BookContentByBookId(bookId), book.BookContentOrderByDesc())
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, bookContents)
}

// AddBooksList 添加书籍列表
func AddBooksList(c *gin.Context) {

	booksList, err := tool.BufferToStruct(c.Request.Body, &models.BooksList{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = book.AddBooksLists(booksList)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}

// AddBookContent 添加书籍内容
func AddBookContent(c *gin.Context) {

	bookContent, err := tool.BufferToStruct(c.Request.Body, &models.BookContent{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = book.AddBookContent(bookContent)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")

}

// ModifyBooksList 修改书籍列表
func ModifyBooksList(c *gin.Context) {
	booksList, err := tool.BufferToStruct(c.Request.Body, &models.BooksList{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = book.ModifyBooksLists(booksList)
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, "成功")
}

// DeleteBooksList 删除书籍列表
func DeleteBooksList(c *gin.Context) {
	booksList, err := tool.BufferToStruct(c.Request.Body, &models.BooksList{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = book.DeleteBooksLists(book.BookId(booksList.ID.ID))
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}

// ModifyBookContent 修改书籍内容
func ModifyBookContent(c *gin.Context) {

	bookContent, err := tool.BufferToStruct(c.Request.Body, &models.BookContent{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = book.ModifyBookContent(bookContent)
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}

// DeleteBookContent 删除书籍内容
func DeleteBookContent(c *gin.Context) {

	bookContent, err := tool.BufferToStruct(c.Request.Body, &models.BookContent{})
	if err != nil {
		response.ValidateFail(c, err.Error())
		return
	}

	err = book.DeleteBookContent(book.BookContentId(bookContent.ID.ID))
	if err != nil {
		response.SqlFail(c, err.Error())
	}

	response.Success(c, "成功")
}
