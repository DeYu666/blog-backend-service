package main

import (
	"github.com/DeYu666/blog-backend-service/app_new/services"
	"github.com/DeYu666/blog-backend-service/bootstrap"
	"github.com/DeYu666/blog-backend-service/global"
)

// set CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o blog_back_end_go main.go
// scp ./blog_back_end_go root@49.234.223.32:/root/
func main() {

	bootstrap.InitializeConfig()

	// 初始化日志
	global.App.Log = bootstrap.InitializeLog()
	global.App.Log.Info("log init success!")

	// 初始化数据库
	global.App.DB = bootstrap.InitializeDB()
	// 程序关闭前，释放数据库连接
	defer closeDbConnect()

	// 加载多个APP的路由配置
	// bootstrap.RouterInclude(backstage.Routers)
	// bootstrap.RouterInclude(frontstage.Routers)
	bootstrap.RouterInclude(services.RoutersBackStage)
	bootstrap.RouterInclude(services.RoutersFrontStage)
	bootstrap.RouterInclude(services.RoutersShuangPin)

	bootstrap.RunServer()
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
