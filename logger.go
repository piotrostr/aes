package main

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	flag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile("log.json", flag, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	encoder := zapcore.NewJSONEncoder(config.EncoderConfig)

	fWriter := zapcore.AddSync(f)
	fileCore := zapcore.NewCore(encoder, fWriter, config.Level)

	stdoutWriter := zapcore.AddSync(os.Stdout)
	sdoutCore := zapcore.NewCore(encoder, stdoutWriter, config.Level)

	core := zapcore.NewTee(fileCore, sdoutCore)

	prod := zap.New(core)
	defer prod.Sync()

	logger := prod.Sugar()
	return logger
}
