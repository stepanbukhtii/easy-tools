package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func ParseRequest(c *gin.Context, obj any) error {
	if err := ParseRequestURI(c, obj); err != nil {
		return err
	}

	switch c.Request.Method {
	case http.MethodGet, http.MethodDelete:
		return c.ShouldBindQuery(obj)
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return c.ShouldBindJSON(obj)
	}

	return nil
}

func ParseRequestURI(c *gin.Context, obj any) error {
	if len(c.Params) == 0 {
		return nil
	}

	m := make(map[string][]string, len(c.Params))
	for _, v := range c.Params {
		m[v.Key] = []string{v.Value}
	}
	return binding.MapFormWithTag(obj, m, "uri")
}
