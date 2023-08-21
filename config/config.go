package config

import (
	"github.com/DeYu666/blog-backend-service/lib/client"
	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/service"
)

var Conf Configuration

type Configuration struct {
	Log      blog.Options                 `mapstructure:"log" json:"log" yaml:"log"`
	Database client.SQLConfig             `mapstructure:"database" json:"database" yaml:"database"`
	Jwt      service.JwtConfig            `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	QiNiuYun client.QiNiuYunConfiguration `mapstructure:"qiniuyun" json:"qiniuyun" yaml:"qiniuyun"`
}
