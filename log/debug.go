// +build !prod

package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logLevel levelType = DebugLevel
)

func defaultLogger() *zap.Logger {
	errSync := zapcore.AddSync(os.Stderr)
	outSync := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(defaultConsoleEncoder(), errSync, zap.LevelEnablerFunc(zapErrEnable)),
		zapcore.NewCore(defaultConsoleEncoder(), outSync, zap.LevelEnablerFunc(zapOutEnable)),
	)

	z := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel), zap.Development())
	return z
}

func DebugEnable() bool {
	return logEnable(DebugLevel)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}
