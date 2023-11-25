package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/rest/swagger/swaggertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterSwagger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	assert.NoError(t, RegisterSwagger(r, swaggertest.FilesFS))

	swaggerFile, err := swaggertest.FilesFS.ReadFile("swagger.yaml")
	assert.NoError(t, err)

	tests := []struct {
		name         string
		url          string
		expectedBody string
	}{
		{
			name: "swagger ui",
			url:  "/swagger-ui/",
			expectedBody: `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui.css" />
    <link rel="stylesheet" type="text/css" href="index.css" />
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
  </head>

  <body>
    <div id="swagger-ui"></div>
    <script src="./swagger-ui-bundle.js" charset="UTF-8"> </script>
    <script src="./swagger-ui-standalone-preset.js" charset="UTF-8"> </script>
    <script src="./swagger-initializer.js" charset="UTF-8"> </script>
  </body>
</html>
`,
		}, {
			name:         "swagger config",
			url:          "/swagger-config",
			expectedBody: `{"urls":[{"name":"Swagger Petstore","url":"/swagger/swagger.yaml"}]}`,
		}, {
			name:         "swagger",
			url:          "/swagger/swagger.yaml",
			expectedBody: string(swaggerFile),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, test.url, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, test.expectedBody, strings.ReplaceAll(w.Body.String(), "\r", ""))
		})
	}
}
