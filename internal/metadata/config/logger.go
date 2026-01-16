package config

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupLogger() {
	f, _ := os.Create("logs/metadata.log")
	gin.DefaultWriter = io.MultiWriter(f)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}
