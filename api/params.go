package api

import "github.com/gin-gonic/gin"

const (
	KeyParams = "params"

	LocaleEN      = "en"
	DefaultLocale = LocaleEN
)

type Params struct {
	Subject string
	Roles   []string
	Locale  string
	IPAddr  string
}

func ExtractParams(c *gin.Context) {
	var params Params

	params.Locale = c.GetHeader(HeaderAcceptLanguage)
	if params.Locale == "" {
		params.Locale = DefaultLocale
	}

	params.IPAddr = c.ClientIP()

	c.Set(KeyParams, params)
}

func GetParams(c *gin.Context) Params {
	v, exists := c.Get(KeyParams)
	if !exists {
		return Params{}
	}
	params, ok := v.(Params)
	if !ok {
		return Params{}
	}
	return params
}
