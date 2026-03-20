package zlog

import (
	"fmt"
	"kama_chat_server/config"
	"os"
	"path"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	cfg := config.GetConfig().LogConfig
	logPath := cfg.LogPath
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 644)
	if err != nil {
		panic(err)
	}
	fileWriteSyncer := zapcore.AddSync(file)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel),
	)
	logger = zap.New(core)
}

func getCallerInfoForLog() (callerFields []zap.Field) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName)

	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}

func Error(pattern string, value ...any) {
	callerFields := getCallerInfoForLog()
	logger.Error(fmt.Sprintf(pattern, value...), callerFields...)
}

func Warn(pattern string, value ...any) {
	callerFields := getCallerInfoForLog()
	logger.Warn(fmt.Sprintf(pattern, value...), callerFields...)
}

func Fatal(pattern string, value ...any) {
	callerFields := getCallerInfoForLog()
	logger.Fatal(fmt.Sprintf(pattern, value...), callerFields...)
}

func Debug(pattern string, value ...any) {
	callerFields := getCallerInfoForLog()
	logger.Debug(fmt.Sprintf(pattern, value...), callerFields...)
}

func Info(pattern string, value ...any) {
	callerFields := getCallerInfoForLog()
	logger.Info(fmt.Sprintf(pattern, value...), callerFields...)
}
