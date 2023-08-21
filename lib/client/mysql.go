package client

import (
	"context"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	blog "github.com/DeYu666/blog-backend-service/lib/log"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLConfig struct {
	Driver              string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host                string `mapstructure:"host" json:"host" yaml:"host"`
	Port                int    `mapstructure:"port" json:"port" yaml:"port"`
	Database            string `mapstructure:"database" json:"database" yaml:"database"`
	UserName            string `mapstructure:"username" json:"username" yaml:"username"`
	Password            string `mapstructure:"password" json:"password" yaml:"password"`
	Charset             string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	LogMode             string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
	LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
	LogRootDir          string `mapstructure:"log_root_dir" json:"log_root_dir" yaml:"log_root_dir"`
	LogMaxBackups       int    `mapstructure:"log_max_backups" json:"log_max_backups" yaml:"log_max_backups"`
	LogMaxSize          int    `mapstructure:"log_max_size" json:"log_max_size" yaml:"log_max_size"`
	LogMaxAge           int    `mapstructure:"log_max_age" json:"log_max_age" yaml:"log_max_age"`
	LogCompress         bool   `mapstructure:"log_compress" json:"log_compress" yaml:"log_compress"`
}

var Mysql mysqlModel

type mysqlModel struct {
	db *gorm.DB
}

func (c *mysqlModel) DB() *gorm.DB {
	return c.db
}

// SetupDB 仅为测试用，正常使用时请使用Setup
func (c *mysqlModel) SetupDB(db *gorm.DB) {
	c.db = db
}

func (c *mysqlModel) Setup(ctx context.Context, config SQLConfig) {

	logger := blog.Extract(ctx)

	dsn := config.UserName + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" +
		config.Database + "?charset=" + config.Charset + "&parseTime=True&loc=Local"
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	// try connecting to mysql repeatedly
	for {
		if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: false,                 // 禁用自动创建外键约束
			Logger:                                   getGormLogger(config), // 使用自定义 Logger
		}); err != nil {

			// fmt.Println("mysql connect failed, err:", err)
			logger.Error("mysql connect failed, err:", zap.Any("err", err))

			time.Sleep(30 * time.Second)
		} else {
			sqlDB, _ := db.DB()
			sqlDB.SetMaxIdleConns(config.MaxIdleConns)
			sqlDB.SetMaxOpenConns(config.MaxOpenConns)
			c.db = db
			return
		}
	}
}

// 自定义 gorm Writer
func getGormLogWriter(config SQLConfig) logger.Writer {
	var writer io.Writer

	// 是否启用日志文件
	if config.EnableFileLogWriter {
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   config.LogRootDir + "/" + config.LogFilename,
			MaxSize:    config.LogMaxSize,
			MaxBackups: config.LogMaxBackups,
			MaxAge:     config.LogMaxAge,
			Compress:   config.LogCompress,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}

func getGormLogger(config SQLConfig) logger.Interface {
	var logMode logger.LogLevel

	switch config.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(config), logger.Config{
		SlowThreshold:             200 * time.Millisecond,      // 慢 SQL 阈值
		LogLevel:                  logMode,                     // 日志级别
		IgnoreRecordNotFoundError: false,                       // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  !config.EnableFileLogWriter, // 禁用彩色打印
	})
}
