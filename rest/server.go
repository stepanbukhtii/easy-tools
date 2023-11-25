package rest

import (
	"net/http"
	"slices"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/config"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stepanbukhtii/easy-tools/rest/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var DefaultOTELSkipPath = []string{"/health", "/swagger-ui", "/swagger-config", "/swagger"}

func NewRouter(apiConfig config.API, serviceConfig config.Service) *gin.Engine {
	r := gin.New()

	gin.DisableBindValidation()

	corsConfig := DefaultCorsConfig
	if len(apiConfig.CORSOrigins) > 0 {
		corsConfig.AllowOrigins = apiConfig.CORSOrigins
	}

	otelFilter := func(ctx *gin.Context) bool { return !slices.Contains(DefaultOTELSkipPath, ctx.Request.URL.Path) }

	r.Use(
		middleware.Recovery,
		otelgin.Middleware(serviceConfig.Name, otelgin.WithGinFilter(otelFilter)),
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
