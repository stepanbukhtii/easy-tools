package api

import (
	"github.com/gin-gonic/gin"
)

func Parse(c *gin.Context, obj any) error {
	if err := c.ShouldBindUri(obj); err != nil {
		return err
	}

	if err := c.ShouldBindQuery(obj); err != nil {
		return err
	}

	if c.Request.Body != nil {
		if err := c.ShouldBindJSON(obj); err != nil {
			return err
		}

	}

	return nil
}
