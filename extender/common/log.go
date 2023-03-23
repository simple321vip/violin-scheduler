package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type LogConfig struct {
	DebugFileName string `json:"debugFileName"`
	InfoFileName  string `json:"infoFileName"`
	WarnFileName  string `json:"warnFileName"`
	MaxSize       int    `json:"maxSize"`
	MaxAge        int    `json:"maxAge"`
	MaxBackups    int    `json:"maxBackups"`
}

var LG *zap.Logger

func InitConfig(cfg *LogConfig) {

	infoLogger := &lumberjack.Logger{
		Filename:   cfg.InfoFileName,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  true,
		Compress:   true,
	}

	debugLogger := &lumberjack.Logger{
		Filename:   cfg.DebugFileName,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  true,
		Compress:   true,
	}

	warnLogger := &lumberjack.Logger{
		Filename:   cfg.WarnFileName,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  true,
		Compress:   true,
	}

	infoWriteSyncer := zapcore.AddSync(infoLogger)
	debugWriteSyncer := zapcore.AddSync(debugLogger)
	warnWriteSyncer := zapcore.AddSync(warnLogger)

	encoder := getEncoder()
	infoCore := zapcore.NewCore(encoder, infoWriteSyncer, zapcore.InfoLevel)
	debugCore := zapcore.NewCore(encoder, debugWriteSyncer, zapcore.DebugLevel)
	warnCore := zapcore.NewCore(encoder, warnWriteSyncer, zapcore.WarnLevel)

	// 标准输出
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	std := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)

	// 集成四种log
	core := zapcore.NewTee(debugCore, infoCore, warnCore, std)

	LG = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(LG)
}

func getEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encoderConfig)

}
