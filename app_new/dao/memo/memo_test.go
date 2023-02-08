package memo

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

func TestMemo(t *testing.T) {

	testName := "test_golang"

	Memo := &models.Memo{Content: testName}

	err := AddMemo(Memo)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
		return
	}

	memos, err := GetMemos()
	if len(memos) < 1 {
		t.Errorf("select meme lists data error, it is %v, error is %v", memos, err)
		return
	}

	err = DeleteMemo(MemoId(memos[len(memos)-1].ID.ID))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err := GetMemos(MemoId(memos[len(memos)-1].ID.ID))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}
