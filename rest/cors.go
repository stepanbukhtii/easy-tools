package rest

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/stepanbukhtii/easy-tools/rest/api"
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
		//api.HeaderTraceID,
	},
	ExposeHeaders: []string{
		api.HeaderContentLength,
	},
	AllowCredentials: true,
	AllowWildcard:    true,
	MaxAge:           12 * time.Hour,
}
