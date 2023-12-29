package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/config"
	"net/http"
	"time"
)

const (
	HeaderAuthenticate    = "WWW-Authenticate"
	HeaderAuthorization   = "Authorization"
	HeaderOrigin          = "Origin"
	HeaderContentLength   = "Content-Length"
	HeaderContentType     = "Content-Type"
	HeaderContentLanguage = "Content-Language"
	HeaderAcceptLanguage  = "Accept-Language"
	HeaderTraceID         = "Trace-ID"
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
		HeaderAuthorization,
		HeaderOrigin,
		HeaderContentLength,
		HeaderContentType,
		HeaderAcceptLanguage,
		HeaderContentLanguage,
		HeaderTraceID,
	},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	AllowWildcard:    true,
	MaxAge:           12 * time.Hour,
}

type Options struct {
	LoggerSkipPath   []string
	CorsAllowOrigins []string
	CorsConfig       *cors.Config
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
		ExtractTraceID,
		ExtractParams,
		cors.New(corsConfig),
	)

	r.NoRoute(func(c *gin.Context) {
		RespondNotFound(c, ErrorCodeInvalidRoute, PathNotFound)
	})

	r.GET("/health", MiddlewareSkipLogger, func(c *gin.Context) {
		RespondOK(c)
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
