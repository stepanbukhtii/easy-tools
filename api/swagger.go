package api

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stepanbukhtii/easy-tools/api/swaggerui"
	"net/http"
	"strings"
)

type SwaggerFile struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func RegisterSwagger(r *gin.Engine, docsFs embed.FS, nameToURLFiles map[string]string) {
	var swaggerFiles []SwaggerFile
	for name, url := range nameToURLFiles {
		if strings.HasPrefix(url, "/") {
			url = fmt.Sprintf("/swagger%s", url)
		} else {
			url = fmt.Sprintf("/swagger/%s", url)
		}
		swaggerFiles = append(swaggerFiles, SwaggerFile{
			URL:  url,
			Name: name,
		})
	}

	r.StaticFS("/swagger-ui", http.FS(swaggerui.FilesFS))

	r.GET("/swagger-config", func(c *gin.Context) {
		c.JSON(200, map[string][]SwaggerFile{"urls": swaggerFiles})
	})

	r.StaticFS("/swagger", http.FS(docsFs))
}
