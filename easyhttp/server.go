package easyhttp

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/stepanbukhtii/easy-tools/api"
	"github.com/stepanbukhtii/easy-tools/config"
)

var DefaultCorsConfig = cors.Config{
	AllowAllOrigins: false,
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
	},
	AllowHeaders: []string{
		api.HeaderAuthorization,
		api.HeaderOrigin,
		api.HeaderContentLength,
		api.HeaderContentType,
		api.HeaderContentLanguage,
		api.HeaderAcceptLanguage,
		api.HeaderTraceID,
	},
	ExposeHeaders: []string{
		api.HeaderContentLength,
	},
	AllowCredentials: true,
	AllowWildcard:    true,
	MaxAge:           12 * time.Hour,
}

func NewRouter(c config.API) *gin.Engine {
	r := gin.New()

	corsConfig := DefaultCorsConfig
	if len(c.CORSOrigins) > 0 {
		corsConfig.AllowOrigins = c.CORSOrigins
	}

	r.Use(
		Recovery,
		MiddlewareLogger,
		TraceID,
		Locale,
		IPAddress,
		cors.New(corsConfig),
	)

	r.NoRoute(func(c *gin.Context) {
		api.RespondNotFound(c, nil)
	})

	r.GET("/health", SkipLogger, func(c *gin.Context) {
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
