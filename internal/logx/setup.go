package logx

import (
	"io"
	"os"
	"strings"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// levelController 日志输出基本控制器
	levelController = zap.NewAtomicLevelAt(zap.DebugLevel)

	// DefaultConfig 日志配置
	DefaultConfig = zapcore.EncoderConfig{
		CallerKey:     "line",
		LevelKey:      "level",
		MessageKey:    "message",
		TimeKey:       "time",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(firstUpper(level.String()))
		},
		EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("[" + caller.FullPath() + "]")
		},
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	}
)

// initDefaultLogger 在没有外部调用Setup进行日志库设置的情况下，进行默认的日志库配置；
// 以便开发单独的小应用的使用时候；
func initDefaultLogger() {
	SetupLogger("", zap.AddStacktrace(zapcore.ErrorLevel))
}

// CloseLogger 系统运行结束时，将日志落盘；
func CloseLogger() {
	_ = rootLogger.Sync()
}

func SetupLogger(logfile string, options ...zap.Option) {
	encoder := zapcore.NewConsoleEncoder(DefaultConfig)

	// 将日志输出到屏幕
	core := zapcore.NewCore(encoder, os.Stdout, levelController)
	// 将日志输出到滚动切割文件中
	if logfile != "" {
		lumberWriterSync := zapcore.AddSync(fileWriter(logfile))
		core = zapcore.NewCore(encoder, lumberWriterSync, levelController)
	}
	options = append(options, zap.AddCaller(), zap.AddCallerSkip(2))
	// 生产根logger，设置输出调度点(上跳2行）
	_zLogger := zap.New(core, options...)

	rootLogger = newzLogger(_zLogger)
}

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func SetLevel(l zapcore.Level) {
	levelController.SetLevel(l)
}

func fileWriter(path string) io.Writer {
	out := &lumberjack.Logger{
		Filename:   path,
		MaxBackups: 7,
		MaxSize:    50,
		MaxAge:     7,
		LocalTime:  true, // use local time zone
	}
	c := cron.New()
	_, _ = c.AddFunc("@daily", func() {
		_ = out.Rotate()
	})
	c.Start()
	return out
}
