package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

func New(o *Options) *gin.Engine {
	var options Options

	if o == nil {
		options = *o
	}

	corsConfig := DefaultCorsConfig
	if options.CorsConfig != nil {
		corsConfig = *options.CorsConfig
	}
	if len(options.CorsAllowOrigins) > 0 {
		corsConfig.AllowOrigins = options.CorsAllowOrigins
	}

	r := gin.New()

	r.Use(
		Recovery(),
		MiddlewareLogger(options.LoggerSkipPath),
		ExtractTraceID,
		ExtractParams,
		cors.New(corsConfig),
	)

	r.NoRoute(func(c *gin.Context) {
		RespondNotFound(c, ErrorCodeInvalidRoute, PathNotFound)
	})

	r.GET("/health", func(c *gin.Context) {
		RespondOK(c)
	})

	return r
}
