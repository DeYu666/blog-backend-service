package diary

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

func TestDiary(t *testing.T) {
	diaries, err := GetDiaries()

	if len(diaries) < 1 {
		t.Errorf("select book lists data error, it is %v, error is %v", diaries, err)
		return
	}

	testName := "test_golang"

	diary := &models.Diary{Content: testName}

	err = AddDiary(diary)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
		return
	}

	info, err := GetDiaries(DiaryContent(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	diary = info[0]
	diary.Content = "test_modify_golang"

	err = ModifyDiary(diary)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetDiaries(DiaryId(id))
	if err != nil || len(info) != 1 || info[0].Content != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteDiary(DiaryId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetDiaries(DiaryId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}

func TestDiaryPwd(t *testing.T) {
	DiaryPwd, err := GetDiaryPwd()

	if len(DiaryPwd) < 1 {
		t.Errorf("select Diary password data error, it is %v, error is %v", DiaryPwd, err)
		return
	}

	testName := "test_golang"

	DiaryPwd1 := &models.DiaryPs{Password: testName}

	err = AddDiaryPwd(DiaryPwd1)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
	}

	info, err := GetDiaryPwd(DiaryPassWord(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	DiaryPwd1 = info[0]
	DiaryPwd1.Password = "test_modify_golang"

	err = ModifyDiaryPwd(DiaryPwd1)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetDiaryPwd(DiaryPwdId(id))
	if err != nil || len(info) != 1 || info[0].Password != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteDiaryPwd(DiaryPwdId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetDiaryPwd(DiaryPwdId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}
