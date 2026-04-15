package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(level string) (*zap.Logger, error) {
	// Niveau
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// Encoder (JSON)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Output fichier
	var logFile = os.Getenv("LOG_FILE")
	fmt.Println("Log dans le fichier suivant :", logFile)
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100, // MB
		MaxBackups: 5,
		MaxAge:     7, // jours
		Compress:   true,
	})

	writer := zapcore.AddSync(fileWriter)

	consoleWriter := zapcore.AddSync(os.Stdout)

	// combine les writers
	multiWriter := zapcore.NewMultiWriteSyncer(writer, consoleWriter)

	core := zapcore.NewCore(encoder, multiWriter, zapLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel))

	return logger, nil
}
