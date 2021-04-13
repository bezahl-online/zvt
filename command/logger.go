package command

import (
	"log"

	"github.com/bezahl-online/zvt/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogger() *zap.Logger {
	logfilePath := util.ENVFilePath("ZVT_LOGFILEPATH", "zvt.log.json")
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{logfilePath}, // "stdout"},
		// ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	zapLogger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	return zapLogger
}
