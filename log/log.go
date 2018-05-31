package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type levelType int8

const (
	DebugLevel levelType = iota + 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

var (
	logger *zap.SugaredLogger
)

func init() {
	z := defaultLogger()
	logger = z.Sugar()
}

func defaultConsoleEncoder() zapcore.Encoder {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = TimeEncoder
	return zapcore.NewConsoleEncoder(cfg)
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// 错误级别及以上在stderr输出
func zapErrEnable(lvl zapcore.Level) bool {
	return lvl >= zapcore.ErrorLevel
}

func zapOutEnable(lvl zapcore.Level) bool {
	return lvl < zapcore.ErrorLevel
}

func SetLevel(level levelType) {
	logLevel = level
}

func InfoEnable() bool {
	return logEnable(InfoLevel)
}

func WarnEnable() bool {
	return logEnable(WarnLevel)
}

func ErrorEnable() bool {
	return logEnable(ErrorLevel)
}

func FatalEnable() bool {
	return logEnable(FatalLevel)
}

func PanicEnable() bool {
	return logEnable(PanicLevel)
}

func DPanicEnable() bool {
	return logEnable(DPanicLevel)
}

func logEnable(level levelType) bool {
	return level >= logLevel
}

func Info(args ...interface{}) {
	if !InfoEnable() {
		return
	}
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	if !WarnEnable() {
		return
	}
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	if !ErrorEnable() {
		return
	}
	logger.Error(args...)
}

func DPanic(args ...interface{}) {
	if !DPanicEnable() {
		return
	}
	logger.DPanic(args...)
}

func Panic(args ...interface{}) {
	if !PanicEnable() {
		return
	}
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	if !FatalEnable() {
		return
	}
	logger.Fatal(args...)
}

func Infof(template string, args ...interface{}) {
	if !InfoEnable() {
		return
	}
	logger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	if !WarnEnable() {
		return
	}
	logger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	if !ErrorEnable() {
		return
	}
	logger.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	if !DPanicEnable() {
		return
	}
	logger.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	if !PanicEnable() {
		return
	}
	logger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	if !FatalEnable() {
		return
	}
	logger.Fatalf(template, args...)
}

func Sync() {
	logger.Sync()
}
