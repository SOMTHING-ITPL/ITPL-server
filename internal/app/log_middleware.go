package server

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *strings.Builder
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewLogger() (*zap.Logger, error) {
	infoFile, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	errorFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel), //stdout 콘솔 출력

		zapcore.NewCore(fileEncoder, zapcore.AddSync(infoFile), infoLevel),   //info
		zapcore.NewCore(fileEncoder, zapcore.AddSync(errorFile), errorLevel), //error
	)

	logger := zap.New(core, zap.AddCaller())
	return logger, nil
}

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		bw := &bodyWriter{body: &strings.Builder{}, ResponseWriter: c.Writer}
		c.Writer = bw

		c.Next()

		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		if status >= 400 {
			respBody := bw.body.String()

			logger.Error("HTTP error : ",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.String("response", respBody),
			)
		} else {
			logger.Info("HTTP request",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
			)
		}

	}
}
