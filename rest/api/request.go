package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var Validator = binding.Validator

func ParseRequest(c *gin.Context, obj any) error {
	if c.Request.Body != nil {
		if err := c.ShouldBindJSON(obj); err != nil {
			return err
		}
	}

	if err := c.ShouldBindUri(obj); err != nil {
		return err
	}

	if err := c.ShouldBindQuery(obj); err != nil {
		return err
	}

	return Validator.ValidateStruct(obj)
}
