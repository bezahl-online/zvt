package command

import (
	"log"

	"github.com/bezahl-online/zvt/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func initLogger() {
	logfilePath := util.ENVFilePath("ZVT_LOGFILEPATH", "zvt.log.json")
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{logfilePath},
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
	if l, err := cfg.Build(); err != nil {
		log.Fatal(err)
	} else {
		Logger = l
	}
}
