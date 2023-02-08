package bootstrap

import (
	"fmt"

	"github.com/DeYu666/blog-backend-service/app_new/services/middleware"
	"github.com/DeYu666/blog-backend-service/global"
	"github.com/gin-gonic/gin"
)

// Option 不知道做什么的
type Option func(*gin.Engine)

var routersOptions []Option

// RouterInclude 注册app的路由配置
func RouterInclude(opts ...Option) {
	routersOptions = append(routersOptions, opts...)
}

// RunServer 初始化路由,并启动服务
func RunServer() {

	if global.App.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 跨域处理
	r.Use(middleware.Cors())

	for _, opt := range routersOptions {
		opt(r)
	}

	addr := "" + ":" + global.App.Config.App.Port

	if err := r.Run(addr); err != nil {
		global.App.Log.Error(fmt.Errorf("startup frontServices failed, err:%v\n", err).Error())
	}
}
