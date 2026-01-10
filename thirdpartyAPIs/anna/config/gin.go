package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() http.Handler {
	SetupLogger()
	e := gin.New()

	return e
}
