package blog

import (
	"testing"

	"github.com/DeYu666/blog-backend-service/app_new/models"
)

// 在 generalCate_test.go 中定义了 TestMain 方法，所以，在这里不需要重复定义。
// TestMain 方法，作用就是 初始化 连接数据库 的配置。

func TestCategory(t *testing.T) {
	generalCate, _ := GetGeneralCategories()

	if len(generalCate) < 1 {
		t.Error("select generalCate data error")
		return
	}

	testName := "test_golang"
	testGeneralCateId := generalCate[0].ID.ID

	cate := &models.BlogCategories{Name: testName, GeneralID: int(testGeneralCateId)}

	err := AddCategory(cate)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
	}

	info, err := GetCategories(CateName(testName))
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	cate = info[0]
	cate.Name = "test_modify_golang"

	err = ModifyCategory(cate)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetCategories(CateId(id))
	if err != nil || len(info) != 1 || info[0].Name != "test_modify_golang" {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteCategory(CateId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetCategories(CateId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}

}
