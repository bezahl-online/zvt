package command

import (
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func initLogger() {
	cfg := zap.NewProductionConfig()
	var logfilePath string
	if v := os.Getenv("ZVT_LOGFILEPATH"); len(v) > 0 {
		logfilePath = v
	} else {
		log.Fatal("please set environment varibalbe 'ZVT_LOGFILEPATH'")
	}
	logfilePath = strings.TrimRight(logfilePath, "/")
	logfilePath += "/zvt.log.json"
	cfg.OutputPaths = []string{
		logfilePath,
	}
	if l, err := cfg.Build(); err != nil {
		log.Fatal(err)
	} else {
		Logger = l.Sugar()
	}
}
