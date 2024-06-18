package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type Empty struct{}

type Request[T any] struct {
	Data   T
	Params Params
}

func (r *Request[T]) Parse(c *gin.Context) error {
	if err := ParseRequestURI(c, &r.Data); err != nil {
		return err
	}

	switch c.Request.Method {
	case http.MethodGet, http.MethodDelete:
		return c.ShouldBindQuery(&r.Data)
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return c.ShouldBindJSON(&r.Data)
	}

	r.Params = GetParams(c)

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

func ParseEmptyRequest(c *gin.Context) Request[Empty] {
	return Request[Empty]{
		Params: GetParams(c),
	}
}
