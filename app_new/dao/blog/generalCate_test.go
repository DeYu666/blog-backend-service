package blog

import (
	"fmt"
	"os"
	"testing"

	"github.com/DeYu666/blog-backend-service/global"

	"github.com/DeYu666/blog-backend-service/app_new/models"
	"github.com/DeYu666/blog-backend-service/bootstrap"
)

func TestGeneralCate(t *testing.T) {
	testName := "test_golang"
	cate := &models.BlogGeneralCategories{}
	cate.Name = testName

	err := AddGeneralCate(cate)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
		return
	}

	info, err := GetGeneralCategories(GeneralCateName(testName))

	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	cate = info[0]
	cate.Name = "test_modify_golang"

	err = ModifyGeneralCate(cate)

	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetGeneralCategories(GeneralCateId(id))

	if err != nil || len(info) != 1 || info[0].Name != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteGeneralCate(GeneralCateId(id))

	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetGeneralCategories(GeneralCateId(id))

	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}

func TestMain(m *testing.M) {
	os.Setenv("VIPER_CONFIG", "../../../config.yaml")
	bootstrap.InitializeConfig()

	global.App.DB = bootstrap.InitializeDB()

	fmt.Println("write setup code here...")

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
