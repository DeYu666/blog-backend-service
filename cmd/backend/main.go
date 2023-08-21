package main

import (
	"context"
	"log"
	"os"

	"github.com/DeYu666/blog-backend-service/config"
	"github.com/DeYu666/blog-backend-service/delivery"
	"github.com/DeYu666/blog-backend-service/delivery/middleware"
	"github.com/DeYu666/blog-backend-service/lib/client"
	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"github.com/DeYu666/blog-backend-service/service"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	// App      App      `mapstructure:"app" json:"app" yaml:"app"`
	Log      blog.Options      `mapstructure:"log" json:"log" yaml:"log"`
	Database client.SQLConfig  `mapstructure:"database" json:"database" yaml:"database"`
	Jwt      service.JwtConfig `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}

func main() {
	ctx := context.Background()

	configPath := "../../config.yaml"

	if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
		configPath = configEnv
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	err = yaml.Unmarshal(data, &config.Conf)
	if err != nil {
		log.Fatalf("Error unmarshalling yaml: %v", err)
	}

	if configDBHost := os.Getenv("DB_HOST"); configDBHost != "" {
		config.Conf.Database.Host = configDBHost
	}

	if configLogRootDir := os.Getenv("LOG_ROOT_DIR"); configLogRootDir != "" {
		config.Conf.Log.RootDir = configLogRootDir
		config.Conf.Database.LogRootDir = configLogRootDir + "/mysql"
	}

	// 初始化日志
	log := blog.InitializeLog(&config.Conf.Log)
	ctx = blog.InjectLogger(ctx, log)

	// 初始化连接数据库
	client.Mysql.Setup(ctx, config.Conf.Database)
	service.JwtService.Init(ctx, config.Conf.Jwt.Secret)

	var middlewareInjectLog = func() gin.HandlerFunc {
		return func(ginCtx *gin.Context) {
			blog.InjectLogger(ginCtx, log)
		}
	}

	r := gin.New()
	r.Use(middlewareInjectLog())
	r.Use(middleware.Cors())

	delivery.NewBlogHandler(r)
	delivery.NewBookHandler(r)
	delivery.NewDiaryHandler(r)
	delivery.NewLoveInfoHandler(r)
	delivery.NewProfileHandler(r)
	delivery.NewMemoHandler(r)
	delivery.NewAuthUserHandler(r)

	r.Run(":8080")
}
