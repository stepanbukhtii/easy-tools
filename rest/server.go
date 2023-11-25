package rest

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/config"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stepanbukhtii/easy-tools/rest/middleware"
)

func NewRouter(c config.API) *gin.Engine {
	r := gin.New()

	gin.DisableBindValidation()

	corsConfig := DefaultCorsConfig
	if len(c.CORSOrigins) > 0 {
		corsConfig.AllowOrigins = c.CORSOrigins
	}

	r.Use(
		middleware.Recovery,
		middleware.ClientInfo,
		middleware.Logger,
		cors.New(corsConfig),
	)

	r.NoRoute(func(c *gin.Context) { api.RespondNotFound(c, nil) })

	r.GET("/health", middleware.SkipLogger, func(c *gin.Context) {
		api.RespondOK(c)
	})

	return r
}

func NewServer(c config.API, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:           c.Address,
		Handler:        router,
		ReadTimeout:    c.Timeout,
		WriteTimeout:   c.Timeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}
