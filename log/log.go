package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"strings"
	"sync"
)

type LogConfig struct {
	FilePath    string `json:"file_path" mapstructure:"file_path" yaml:"file_path"`
	MaxSize     int    `json:"max_size" mapstructure:"max_size" yaml:"max_size"`
	MaxBackups  int    `json:"max_backups" mapstructure:"max_backups" yaml:"max_backups"`
	MaxAge      int    `json:"max_age" mapstructure:"max_age" yaml:"max_age"`
	Compress    bool   `json:"compress" mapstructure:"compress" yaml:"compress"`
	FileName    string `json:"file_name" mapstructure:"file_name" yaml:"file_name"`
	EncodeLevel string `json:"encode_level" mapstructure:"encode_level" yaml:"encode_level"`
	Mode        string `json:"mode" mapstructure:"mode" yaml:"mode"`
	Level       string `json:"level" mapstructure:"level" yaml:"level"`
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		FilePath:    "./logs",
		FileName:    "gind.log",
		Mode:        "dev",
		Level:       "debug",
		EncodeLevel: "CapitalColorLevelEncoder",
		MaxSize:     50,
		MaxBackups:  15,
		MaxAge:      15,
		Compress:    true,
	}
}

var logger *zap.Logger
var once sync.Once

var (
	DEBUG func(msg string, fields ...zap.Field)
	INFO  func(msg string, fields ...zap.Field)
	WARN  func(msg string, fields ...zap.Field)
	ERR   func(msg string, fields ...zap.Field)
	FATAL func(msg string, fields ...zap.Field)
)

func (C *LogConfig) Init() *zap.Logger {
	once.Do(func() {
		logger = C.loadCore()
		zap.ReplaceGlobals(logger)
	})
	DEBUG = logger.Debug
	INFO = logger.Info
	WARN = logger.Warn
	ERR = logger.Error
	FATAL = logger.Fatal
	return logger
}

/*
ioc 模块会自动执行 struct中的  Bean前缀的方法
*/
func (C *LogConfig) BeanInitLogger() *zap.Logger {
	once.Do(func() {
		logger = C.loadCore()
		zap.ReplaceGlobals(logger)
	})
	DEBUG = logger.Debug
	INFO = logger.Info
	WARN = logger.Warn
	ERR = logger.Error
	FATAL = logger.Fatal
	return logger
}

func (C *LogConfig) loadCore() *zap.Logger {
	switch C.Mode {
	case "dev":
		return zap.New(C.devCore())
	case "prod":
		return zap.New(C.prodCore())
	default:
		return zap.New(C.devCore())
	}
}

func (C *LogConfig) getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    C.getEncodeLevel(),                                 // 根据配置获取 编码器
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func (C *LogConfig) getEncodeLevel() zapcore.LevelEncoder {
	switch C.EncodeLevel {
	case "LowercaseLevelEncoder":
		// 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder":
		// 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder":
		// 大写编码器
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder":
		// 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func (C *LogConfig) prodCore() zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(
			C.getEncoderConfig(),
		),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(
				&lumberjack.Logger{
					Filename:   path.Join(C.FilePath, C.FileName),
					MaxSize:    C.MaxSize,
					MaxBackups: C.MaxBackups,
					MaxAge:     C.MaxAge,
					Compress:   C.Compress,
					LocalTime:  true,
				},
			),
		),
		zap.LevelEnablerFunc(
			func(level zapcore.Level) bool {
				return level >= C.getLevel()
			},
		),
	)
}

func (C *LogConfig) devCore() zapcore.Core {
	stdEncoderConfig := C.getEncoderConfig()
	stdEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(stdEncoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.LevelEnablerFunc(
			func(level zapcore.Level) bool {
				return level >= zapcore.DebugLevel
			},
		),
	)
}

func (C *LogConfig) getLevel() zapcore.Level {

	switch strings.ToLower(C.Level) {
	case "debug", "DEBUG":
		return zapcore.DebugLevel
	case "info", "INFO", "": // make the zero value useful
		return zapcore.InfoLevel
	case "warn", "WARN":
		return zapcore.WarnLevel
	case "error", "ERROR":
		return zapcore.ErrorLevel
	case "dpanic", "DPANIC":
		return zapcore.DPanicLevel
	case "panic", "PANIC":
		return zapcore.PanicLevel
	case "fatal", "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
