package api

import (
	"embed"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	"github.com/stepanbukhtii/easy-tools/rest/swagger/swaggerui"
)

type swaggerContent struct {
	Info info `yaml:"info"`
}

type info struct {
	Title string `yaml:"title"`
}

type swaggerConfigUrl struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func RegisterSwagger(r *gin.Engine, docsFs embed.FS) error {
	files, err := docsFs.ReadDir(".")
	if err != nil {
		return err
	}

	var swaggerConfigURLs []swaggerConfigUrl
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		file, err := docsFs.Open(f.Name())
		if err != nil {
			return err
		}

		fileContent, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		var data swaggerContent
		if err := yaml.Unmarshal(fileContent, &data); err != nil {
			continue
		}

		swaggerConfigURLs = append(swaggerConfigURLs, swaggerConfigUrl{
			Name: data.Info.Title,
			URL:  fmt.Sprintf("/swagger/%s", f.Name()),
		})
	}

	r.StaticFS("/swagger-ui", http.FS(swaggerui.FilesFS))

	r.GET("/swagger-config", func(c *gin.Context) {
		c.JSON(200, map[string][]swaggerConfigUrl{"urls": swaggerConfigURLs})
	})

	r.StaticFS("/swagger", http.FS(docsFs))

	return nil
}
