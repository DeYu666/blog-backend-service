package global

import (
	"github.com/DeYu666/blog-backend-service/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/**
定义 Application 结构体，用来存放一些项目启动时的变量，便于调用.
*/

type Application struct {
	ConfigViper *viper.Viper
	Config      config.Configuration
	Log         *zap.Logger
	DB          *gorm.DB
}

var App = new(Application)
