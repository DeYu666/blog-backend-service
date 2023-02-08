package book

import (
	"fmt"
	"os"
	"testing"

	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/bootstrap"
	"github.com/DeYu666/blog-backend-service/global"
)

func TestMain(m *testing.M) {
	fmt.Println("先运行主测试函数，运行配置文件，连接数据库...")

	os.Setenv("VIPER_CONFIG", "../../../config.yaml")
	bootstrap.InitializeConfig()

	global.App.DB = bootstrap.InitializeDB()

	fmt.Println("开始测试...")
	m.Run()

	closeDbConnect()
}

// 释放数据库连接
func closeDbConnect() {
	if global.App.DB != nil {
		db, _ := global.App.DB.DB()
		err := db.Close()
		if err != nil {
			return
		}
	}
}

func TestBooksLists(t *testing.T) {
	books, err := GetBooksLists()

	if len(books) < 1 {
		t.Errorf("select book lists data error, it is %v, error is %v", books, err)
		return
	}

	testName := "test_golang"

	tag := &models.BooksList{BookName: testName, BookStatus: "", Abstract: ""}

	err = AddBooksLists(tag)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
	}

	info, err := GetBooksLists(BookName(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	tag = info[0]
	tag.BookName = "test_modify_golang"

	err = ModifyBooksLists(tag)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetBooksLists(BookId(id))
	if err != nil || len(info) != 1 || info[0].BookName != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteBooksLists(BookId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetBooksLists(BookId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}
