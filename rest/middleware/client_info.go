package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stepanbukhtii/easy-tools/rest/client"
)

func ClientInfo(c *gin.Context) {
	var clientInfo client.Info

	clientInfo.Locale = c.GetHeader(api.HeaderAcceptLanguage)
	if clientInfo.Locale == "" {
		clientInfo.Locale = client.DefaultLocale
	}

	clientInfo.IPAddress = c.ClientIP()
	clientInfo.UserAgent = c.Request.UserAgent()
	clientInfo.DeviceID = c.GetHeader(api.HeaderDeviceID)

	c.Request = c.Request.WithContext(econtext.SetClientInfo(c.Request.Context(), clientInfo))
}
