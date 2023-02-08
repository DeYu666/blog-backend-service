package frontstage

import (
	"strconv"

	"github.com/DeYu666/blog-backend-service/app_new/common/response"
	"github.com/DeYu666/blog-backend-service/app_new/dao/book"
	"github.com/gin-gonic/gin"
)

func GetBookContent(c *gin.Context) {
	param := c.Param("bookId")

	bookId, err := strconv.Atoi(param)
	if err != nil {
		response.ValidateFail(c, err.Error())
	}

	bookContents, err := book.GetBookContent(book.BookContentByBookId(bookId))
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	response.Success(c, bookContents)
}

func GetBooksListByCateId(c *gin.Context) {
	cateId, err := strconv.Atoi(c.Query("categoryId"))
	if err != nil {
		cateId = -1
	}
	pageNum, err := strconv.Atoi(c.Query("pageNum"))
	if err != nil {
		pageNum = 0
	}

	books, err := book.GetBooksLists(book.BookCategoryID(cateId), book.BookPaging(pageNum, 12), book.BookOrderByDesc())
	if err != nil {
		response.SqlFail(c, err.Error())
		return
	}

	tempBooks, _ := book.GetBooksLists(book.BookCategoryID(cateId))
	count := len(tempBooks)

	response.Success(c, gin.H{"booksContent": books, "pagesNum": count})
}
