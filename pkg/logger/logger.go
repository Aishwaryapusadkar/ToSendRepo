package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logObject *zap.Logger

func Log(data ...context.Context) *zap.Logger {
	if data != nil {
		ctx := data[0]
		return logObject.With(zap.Any("messageType", ctx))
	} else {
		return logObject
	}
}

// func LoggerInit(logFilePath string, level zapcore.Level) {
// 	var (
// 		err error
// 	)
// 	fmt.Println("LOGGER INIT started")
// 	// logFile, _ := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

// 	cfg := zap.NewProductionConfig()
// 	cfg.Level = zap.NewAtomicLevelAt(level)
// 	cfg.EncoderConfig.FunctionKey = "f"
// 	//cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
// 	cfg.EncoderConfig.EncodeTime = syslogTimeEncoder
// 	cfg.EncoderConfig.ConsoleSeparator = " "
// 	cfg.EncoderConfig.EncodeCaller = MyCaller
// 	cfg.Encoding = "console"
// 	// if logFilePath != "" {
// 	// 	var paths []string
// 	// 	paths = append(paths, logFilePath)
// 	// 	cfg.OutputPaths = paths
// 	// }
// 	//TODO: Error handling
// 	prodConf := zap.NewProductionEncoderConfig()
// 	prodConf.EncodeTime = syslogTimeEncoder
// 	prodConf.EncodeLevel = CustomLevelEncoder
// 	logObject, err = cfg.Build()
// 	if err != nil {
// 		fmt.Println("failed to create custom production logger , Exiting system", err)
// 		os.Exit(0)
// 	} else if logObject == nil {
// 		logObject, err = zap.NewProduction()
// 		logObject.WithOptions(zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
// 		if err != nil {
// 			fmt.Println("failed to create production logger , Exiting system", err)
// 			os.Exit(0)
// 		}
// 		fmt.Println("Failed to create custom production logger, creating production logger")
// 	} else {
// 		logObject.WithOptions(zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
// 		fmt.Println("custom production logger created")
// 	}

// 	logObject.Info("Logger init successfully")
// }

func LoggerInit(logFilePath string, level zapcore.Level) {
	var err error
	fmt.Println("LOGGER INIT started")
	
	// Open the log file for writing
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("failed to open log file, Exiting system", err)
		os.Exit(0)
	}

	// Encoder configuration for both console and file output
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = syslogTimeEncoder
	encoderConfig.EncodeLevel = CustomLevelEncoder
	encoderConfig.EncodeCaller = MyCaller

	// Create two encoders: one for console and one for file output
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create a zapcore.WriteSyncer for the file
	fileWriter := zapcore.AddSync(logFile)

	// Create a zapcore.Core for the file logger with the provided log level
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, zap.NewAtomicLevelAt(level))

	// Optionally: create a console logger core for outputting to the console
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(level))

	// Combine both cores (console and file) if you want to log to both
	core := zapcore.NewTee(fileCore, consoleCore)

	// Build the logger from the combined core
	logObject = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.FatalLevel))

	if logObject == nil {
		fmt.Println("Failed to create logger, Exiting system")
		os.Exit(0)
	}

	logObject.Info("Logger initialized successfully")
}


// time logger
func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Jan 2 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func MyCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(filepath.Base(caller.FullPath()))
}
