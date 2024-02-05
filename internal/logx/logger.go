package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = &zLogger{}

type Logger interface {
	Print(...interface{})
	Println(...interface{})
	Printf(string, ...interface{})
	Printw(string, ...interface{})
	Printx(string, ...zapcore.Field)

	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})
	Debugw(string, ...interface{})
	Debugx(string, ...zapcore.Field)

	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})
	Infow(string, ...interface{})
	Infox(string, ...zapcore.Field)

	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})
	Warnw(string, ...interface{})
	Warnx(string, ...zapcore.Field)

	Warning(...interface{})
	Warningln(...interface{})
	Warningf(string, ...interface{})
	Warningw(string, ...interface{})
	Warningx(string, ...zapcore.Field)

	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})
	Errorw(string, ...interface{})
	Errorx(string, ...zapcore.Field)

	Panic(...interface{})
	Panicln(...interface{})
	Panicf(string, ...interface{})
	Panicw(string, ...interface{})
	Panicx(string, ...zapcore.Field)

	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})
	Fatalw(string, ...interface{})
	Fatalx(string, ...zapcore.Field)

	DPanic(...interface{})
	DPanicln(...interface{})
	DPanicf(string, ...interface{})
	DPanicw(string, ...interface{})
	DPanicx(string, ...zapcore.Field)

	GetSubLogger() Logger
	GetSubLoggerWithKeyValue(map[string]string) Logger
	GetSubLoggerWithOption(...zap.Option) Logger
}

type zLogger struct {
	logger  *zap.Logger
	slogger *zap.SugaredLogger
}

func newzLogger(logger *zap.Logger) *zLogger {
	aLogger := &zLogger{
		logger:  logger,
		slogger: logger.Sugar(),
	}
	return aLogger
}

// Print 输出"Info"级别日志信息；
func (z *zLogger) Print(args ...interface{}) {
	z.slogger.Info(args...)
}

// Println 输出"Infoln"级别日志信息；
func (z *zLogger) Println(args ...interface{}) {
	z.slogger.Infoln(args...)
}

// Printf 输出格式化的"Info"级别日志信息；
func (z *zLogger) Printf(template string, args ...interface{}) {
	z.slogger.Infof(template, args...)
}

// Printw 输出定制化的"Info"级别日志信息；
func (z *zLogger) Printw(msg string, keysAndValues ...interface{}) {
	z.slogger.Infow(msg, keysAndValues...)
}

// Printx 以zapfield方式，极速输出定制化的"Info"级别日志信息；
func (z *zLogger) Printx(msg string, fields ...zapcore.Field) {
	z.logger.Info(msg, fields...)
}

// Debug 输出"Debug"级别日志信息；
func (z *zLogger) Debug(args ...interface{}) {
	z.slogger.Debug(args...)
}

// Debugln 输出"Debugln"级别日志信息；
func (z *zLogger) Debugln(args ...interface{}) {
	z.slogger.Debugln(args...)
}

// Debugf 输出格式化的"Debug"级别日志信息；
func (z *zLogger) Debugf(template string, args ...interface{}) {
	z.slogger.Debugf(template, args...)
}

// Debugw 输出定制化的"Debug"级别日志信息；
func (z *zLogger) Debugw(msg string, keysAndValues ...interface{}) {
	z.slogger.Debugw(msg, keysAndValues...)
}

// Debugx 以zapfield方式，极速输出定制化的"Debug"级别日志信息；
func (z *zLogger) Debugx(msg string, fields ...zapcore.Field) {
	z.logger.Debug(msg, fields...)
}

// Info 输出"Info"级别日志信息；
func (z *zLogger) Info(args ...interface{}) {
	z.slogger.Info(args...)
}

// Infoln 输出"Infoln"级别日志信息；
func (z *zLogger) Infoln(args ...interface{}) {
	z.slogger.Infoln(args...)
}

// Infof 输出格式化的"Info"级别日志信息；
func (z *zLogger) Infof(template string, args ...interface{}) {
	z.slogger.Infof(template, args...)
}

// Infow 输出定制化的"Info"级别日志信息；
func (z *zLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.slogger.Infow(msg, keysAndValues...)
}

// Infox 以zapfield方式，极速输出定制化的"Info"级别日志信息；
func (z *zLogger) Infox(msg string, fields ...zapcore.Field) {
	z.logger.Info(msg, fields...)
}

// Warn 输出"Warn"级别日志信息；
func (z *zLogger) Warn(args ...interface{}) {
	z.slogger.Warn(args...)
}

// Warnln 输出"Warnln"级别日志信息；
func (z *zLogger) Warnln(args ...interface{}) {
	z.slogger.Warnln(args...)
}

// Warnf 输出格式化的"Warn"级别日志信息；
func (z *zLogger) Warnf(template string, args ...interface{}) {
	z.slogger.Warnf(template, args...)
}

// Warnw 输出定制化的"Warn"级别日志信息；
func (z *zLogger) Warnw(msg string, keysAndValues ...interface{}) {
	z.slogger.Warnw(msg, keysAndValues...)
}

// Warnx 以zapfield方式，极速输出定制化的"Warn"级别日志信息；
func (z *zLogger) Warnx(msg string, fields ...zapcore.Field) {
	z.logger.Warn(msg, fields...)
}

// Warning 输出"Warn"级别日志信息；
func (z *zLogger) Warning(args ...interface{}) {
	z.slogger.Warn(args...)
}

// Warningln 输出"Warnln"级别日志信息；
func (z *zLogger) Warningln(args ...interface{}) {
	z.slogger.Warnln(args...)
}

// Warningf 输出格式化的"Warn"级别日志信息；
func (z *zLogger) Warningf(template string, args ...interface{}) {
	z.slogger.Warnf(template, args...)
}

// Warningw 输出定制化的"Warn"级别日志信息；
func (z *zLogger) Warningw(msg string, keysAndValues ...interface{}) {
	z.slogger.Warnw(msg, keysAndValues...)
}

// Warningx 以zapfield方式，极速输出定制化的"Warn"级别日志信息；
func (z *zLogger) Warningx(msg string, fields ...zapcore.Field) {
	z.logger.Warn(msg, fields...)
}

// Error 输出"Error"级别日志信息；
func (z *zLogger) Error(args ...interface{}) {
	z.slogger.Error(args...)
}

// Errorln 输出"Errorln"级别日志信息；
func (z *zLogger) Errorln(args ...interface{}) {
	z.slogger.Errorln(args...)
}

// Errorf 输出格式化的"Error"级别日志信息；
func (z *zLogger) Errorf(template string, args ...interface{}) {
	z.slogger.Errorf(template, args...)
}

// Errorw 输出定制化的"Error"级别日志信息；
func (z *zLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.slogger.Errorw(msg, keysAndValues...)
}

// Errorx 以zapfield方式，极速输出定制化的"Error"级别日志信息；
func (z *zLogger) Errorx(msg string, fields ...zapcore.Field) {
	z.logger.Error(msg, fields...)
}

// DPanic 输出"DPanic"级别日志信息,但不引发程序Panic
func (z *zLogger) DPanic(args ...interface{}) {
	z.slogger.DPanic(args...)
}

// DPanicln 输出"DPanicln"级别日志信息,但不引发程序Panic
func (z *zLogger) DPanicln(args ...interface{}) {
	z.slogger.DPanicln(args...)
}

// DPanicf 输出格式化的"DPanic"级别日志信息，但不引发程序Panic
func (z *zLogger) DPanicf(template string, args ...interface{}) {
	z.slogger.DPanicf(template, args...)
}

// DPanicw 输出定制化的"Panic"级别日志信息,但不引发程序Panic()
func (z *zLogger) DPanicw(msg string, keysAndValues ...interface{}) {
	z.slogger.DPanicw(msg, keysAndValues...)
}

// DPanicx 以zapfield方式，极速输出定制化的"Panic"级别日志信息,但不引发程序Panic()
func (z *zLogger) DPanicx(msg string, fields ...zapcore.Field) {
	z.logger.DPanic(msg, fields...)
}

// Panic 输出"Panic"级别日志信息，并引发程序Panic；
func (z *zLogger) Panic(args ...interface{}) {
	z.slogger.Panic(args...)
}

// Panicln 输出"Panicln"级别日志信息，并引发程序Panic；
func (z *zLogger) Panicln(args ...interface{}) {
	z.slogger.Panicln(args...)
}

// Panicf 输出格式化的"Panic"级别日志信息，并引发程序Panic；
func (z *zLogger) Panicf(template string, args ...interface{}) {
	z.slogger.Panicf(template, args...)
}

// Panicw 输出定制化的"Panic"级别日志信息，并引发程序Panic；
func (z *zLogger) Panicw(msg string, keysAndValues ...interface{}) {
	z.slogger.Panicw(msg, keysAndValues...)
}

// Panicx 以zapfield方式，极速输出定制化的"Panic"级别日志信息，并引发程序Panic；
func (z *zLogger) Panicx(msg string, fields ...zapcore.Field) {
	z.logger.Panic(msg, fields...)
}

// Fatal 输出"Fatal"级别日志信息，并使程序退出（os.Exit(1)；
func (z *zLogger) Fatal(args ...interface{}) {
	z.slogger.Fatal(args...)
}

// Fatalln 输出"Fatalln"级别日志信息，并使程序退出（os.Exit(1)；
func (z *zLogger) Fatalln(args ...interface{}) {
	z.slogger.Fatalln(args...)
}

// Fatalf 输出格式化的"Fatal"级别日志信息，并使程序退出（os.Exit(1)；
func (z *zLogger) Fatalf(template string, args ...interface{}) {
	z.slogger.Fatalf(template, args...)
}

// Fatalw 输出定制化的"Fatal"级别日志信息，并使程序退出（os.Exit(1)；
func (z *zLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	z.slogger.Fatalw(msg, keysAndValues...)
}

// Fatalx 以zapfield方式，极速输出定制化的"Fatal"级别日志信息，并使程序退出（os.Exit(1)；
func (z *zLogger) Fatalx(msg string, fields ...zapcore.Field) {
	z.logger.Fatal(msg, fields...)
}

// Sync 将zapLogger缓冲内容刷写到输出端
func (z *zLogger) Sync() error {
	return z.logger.Sync()
}

// GetSubLogger 获取一个子Logger
func (z *zLogger) GetSubLogger() Logger {
	_zLogger := z.logger.WithOptions()
	tmp := newzLogger(_zLogger)
	return tmp
}

// GetSubLoggerWithKeyValue 获取一个子logger，并在子logger中，定制固定的输出内容
func (z *zLogger) GetSubLoggerWithKeyValue(keyAndValues map[string]string) Logger {
	fields := make([]zap.Field, 0)
	for key, value := range keyAndValues {
		fields = append(fields, zap.String(key, value))
	}
	_zLogger := z.logger.With(fields...)
	tmp := newzLogger(_zLogger)
	return tmp
}

func (z *zLogger) GetSubLoggerWithOption(opts ...zap.Option) Logger {
	_zLogger := z.logger.WithOptions(opts...)
	tmp := newzLogger(_zLogger)
	return tmp
}
