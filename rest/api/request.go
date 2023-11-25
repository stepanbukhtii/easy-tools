package api

import (
	"errors"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var Validator = binding.Validator

func ParseRequest(c *gin.Context, obj any) bool {
	if err := ParseRequestError(c, obj); err != nil {
		RespondBadRequest(c, err)
		return false
	}
	return true
}

func ParseRequestError(c *gin.Context, obj any) error {
	if c.Request.Body != nil {
		if err := c.ShouldBindJSON(obj); err != nil && !errors.Is(err, io.EOF) {
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
