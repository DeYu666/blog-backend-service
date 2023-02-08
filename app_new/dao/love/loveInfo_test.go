package love

import (
	"fmt"
	"os"
	"testing"
	"time"

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

func TestLoveInfo(t *testing.T) {

	loveInfo := &models.LoveInfo{ConfessionTime: time.Now()}

	err := AddLoveInfo(loveInfo)
	if err != nil {
		t.Errorf("insert Data error, it`s %v", err)
		return
	}

	info, err := GetLoveInfos()
	if err != nil || len(info) < 1 {
		t.Errorf("select Data error, error is %v, info is %v", err, info)
		return
	}

	id := info[0].ID.ID

	loveInfo = info[0]
	loveInfo.ConfessionTime = time.Now()

	err = ModifyLoveInfo(loveInfo)
	if err != nil {
		t.Errorf("modify Data error, it`s %v", err)
		return
	}

	info, err = GetLoveInfos(LoveInfoId(id))
	if err != nil || len(info) != 1 {
		t.Errorf("modify Data error, error`s %v, info is %v", err, info)
	}

	err = DeleteLoveInfo(LoveInfoId(id))
	if err != nil {
		t.Errorf("delete Data error, it`s %v", err)
		return
	}

	info, err = GetLoveInfos(LoveInfoId(id))
	if err != nil || len(info) != 0 {
		t.Errorf("delete Data error, error`s %v, info is %v", err, info)
		return
	}
}
