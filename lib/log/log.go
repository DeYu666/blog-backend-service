package blog

import (
	"os"
	"time"

	"github.com/DeYu666/blog-backend-service/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	Level    string `mapstructure:"level" json:"level" yaml:"level"`
	Format   string `mapstructure:"format" json:"format" yaml:"format"`
	ShowLine bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`

	WriteToFile bool   `mapstructure:"write_to_file" json:"write_to_file" yaml:"write_to_file"`
	RootDir     string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	Filename    string `mapstructure:"filename" json:"filename" yaml:"filename"`
	MaxBackups  int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxSize     int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"` // MB
	MaxAge      int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`    // day
	Compress    bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

type LoggerKey struct{}

var (
	level   zapcore.Level // zap 日志等级
	options []zap.Option  // zap 配置项
)

func InitializeLog(config *Options) *zap.Logger {

	if config.WriteToFile {
		// 创建根目录
		createRootDir(config)
	}

	// 设置日志等级
	setLogLevel(config)

	if config.ShowLine {
		options = append(options, zap.AddCaller())
	}

	// 初始化 zap
	logger := zap.New(getZapCore(config), options...)

	return logger
}

func createRootDir(config *Options) {
	if ok, _ := utils.PathExists(config.RootDir); !ok {
		_ = os.Mkdir(config.RootDir, os.ModePerm)
	}
}

func setLogLevel(config *Options) {
	switch config.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
}

// 扩展 Zap
func getZapCore(config *Options) zapcore.Core {
	var encoder zapcore.Encoder

	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}

	// 设置编码器
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, getLogWriter(config), level)
}

// 使用 lumberjack 作为日志写入器
func getLogWriter(config *Options) zapcore.WriteSyncer {
	if config.WriteToFile {
		file := &lumberjack.Logger{
			Filename:   config.RootDir + "/" + config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}

		return zapcore.AddSync(file)
	}

	return zapcore.AddSync(os.Stdout)
}
