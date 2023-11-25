package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/rest/api"
	"github.com/stepanbukhtii/easy-tools/rest/client"
	"github.com/stretchr/testify/assert"
)

func TestClientInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	deviceID := uuid.NewString()
	ipAddress := "192.168.1.1"
	userAgent := "User-Agent"

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.RemoteAddr = fmt.Sprintf("%s:%s", ipAddress, "56000")
	c.Request.Header.Set(api.HeaderUserAgent, userAgent)
	c.Request.Header.Set(api.HeaderAcceptLanguage, client.LocaleEN)
	c.Request.Header.Set(api.HeaderDeviceID, deviceID)

	ClientInfo(c)

	expectedClientInfo := client.Info{
		Locale:    client.LocaleEN,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		DeviceID:  deviceID,
	}
	assert.Equal(t, expectedClientInfo, econtext.ClientInfo(c.Request.Context()))
}
