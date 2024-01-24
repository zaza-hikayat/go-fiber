package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/zaza-hikayat/go-fiber/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _logger *zap.Logger

func Logger() *zap.Logger {
	if _logger == nil {
		conf, _ := configs.LoadConfig()
		_logger = InitLogger(conf)
	}
	return _logger
}

func InitLogger(conf *configs.Config) *zap.Logger {
	var zapCores []zapcore.Core
	if conf.Server.IsProduction {
		timeFormatted := time.Now().Format("2006-01-02")
		file := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/app_%s.log", timeFormatted),
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     7, // days
		})

		level := zap.NewAtomicLevelAt(zap.WarnLevel)
		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		fileEncoder := zapcore.NewJSONEncoder(productionCfg)
		zapCores = append(zapCores, zapcore.NewCore(fileEncoder, file, level))
	}

	if conf.Server.IsLogging {
		level := zap.NewAtomicLevelAt(zap.InfoLevel)
		stdout := zapcore.AddSync(os.Stdout)
		developmentCfg := zap.NewDevelopmentEncoderConfig()
		consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
		zapCores = append(zapCores, zapcore.NewCore(consoleEncoder, stdout, level))
	}

	core := zapcore.NewTee(
		zapCores...,
	)
	_logger = zap.New(core)
	return _logger
}
