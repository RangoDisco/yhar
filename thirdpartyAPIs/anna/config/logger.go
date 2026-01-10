package config

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger() {
	file := &lumberjack.Logger{
		Filename:   "anna.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	gin.DefaultWriter = multiWriter
	gin.DefaultErrorWriter = multiWriter
}
