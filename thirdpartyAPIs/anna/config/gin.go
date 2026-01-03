package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	return e
}
