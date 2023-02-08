package book

import (
	"testing"

	"github.com/DeYu666/blog-backend-service/app_new/models"
)

func TestBookContent(t *testing.T) {
	data, err := GetBookContent()

	if len(data) < 1 {
		t.Errorf("select book lists data error, it is %v, error is %v", data, err)
		return
	}

	testName := "test_golang"

	books, _ := GetBooksLists()

	content := &models.BookContent{BookContent: testName, BookId: int(books[0].ID.ID)}

	err = AddBookContent(content)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
	}

	info, err := GetBookContent(BookContent(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	content = info[0]
	content.BookContent = "test_modify_golang"

	err = ModifyBookContent(content)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetBookContent(BookContentId(id))
	if err != nil || len(info) != 1 || info[0].BookContent != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteBookContent(BookContentId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetBookContent(BookContentId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}
