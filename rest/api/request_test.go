package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testBody struct {
	ID        string `uri:"id"`
	QueryText string `form:"query_text"`
	JSONText  string `json:"json_text"`
}

type testBodyRequired struct {
	ID        string `uri:"id" binding:"required"`
	QueryText string `form:"query_text" binding:"required"`
	JSONText  string `json:"json_text" binding:"required"`
}

func TestParse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c := gin.CreateTestContextOnly(httptest.NewRecorder(), gin.New())

	jsonData, err := json.Marshal(testBody{JSONText: "json_text"})
	assert.NoError(t, err)

	tests := []struct {
		name           string
		method         string
		url            string
		body           io.Reader
		params         []gin.Param
		expectedObject testBody
	}{
		{
			name:           "success get with params",
			method:         http.MethodGet,
			url:            "/123456789",
			params:         []gin.Param{{Key: "id", Value: "123456789"}},
			expectedObject: testBody{ID: "123456789"},
		}, {
			name:           "success get with query",
			method:         http.MethodGet,
			url:            "/?query_text=query_text",
			expectedObject: testBody{QueryText: "query_text"},
		}, {
			name:           "success post with body",
			method:         http.MethodPost,
			url:            "/",
			body:           bytes.NewBuffer(jsonData),
			expectedObject: testBody{JSONText: "json_text"},
		}, {
			name:   "success post with params, query and body",
			method: http.MethodPost,
			url:    "/123456789?query_text=query_text",
			params: []gin.Param{{Key: "id", Value: "123456789"}},
			body:   bytes.NewBuffer(jsonData),
			expectedObject: testBody{
				ID:        "123456789",
				QueryText: "query_text",
				JSONText:  "json_text",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.Request, _ = http.NewRequest(test.method, test.url, test.body)
			c.Params = test.params

			var obj testBody
			assert.NoError(t, ParseRequest(c, &obj))

			assert.Equal(t, test.expectedObject.ID, obj.ID)
			assert.Equal(t, test.expectedObject.QueryText, obj.QueryText)
			assert.Equal(t, test.expectedObject.JSONText, obj.JSONText)
		})
	}
}

func TestParseValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	gin.DisableBindValidation()

	c := gin.CreateTestContextOnly(httptest.NewRecorder(), gin.New())

	jsonData, err := json.Marshal(testBody{JSONText: "json_text"})
	assert.NoError(t, err)

	tests := []struct {
		name          string
		method        string
		url           string
		body          io.Reader
		params        []gin.Param
		expectedError string
	}{
		{
			name:   "success validation",
			method: http.MethodPost,
			url:    "/123456789?query_text=query_text",
			params: []gin.Param{{Key: "id", Value: "123456789"}},
			body:   bytes.NewBuffer(jsonData),
		}, {
			name:          "error validation params",
			method:        http.MethodGet,
			url:           "/?query_text=query_text",
			body:          bytes.NewBuffer(jsonData),
			expectedError: "Key: 'testBodyRequired.ID' Error:Field validation for 'ID' failed on the 'required' tag",
		}, {
			name:          "error validation query",
			method:        http.MethodGet,
			url:           "/123456789",
			params:        []gin.Param{{Key: "id", Value: "123456789"}},
			body:          bytes.NewBuffer(jsonData),
			expectedError: "Key: 'testBodyRequired.QueryText' Error:Field validation for 'QueryText' failed on the 'required' tag",
		}, {
			name:          "error validation body",
			method:        http.MethodPost,
			url:           "/123456789?query_text=query_text",
			params:        []gin.Param{{Key: "id", Value: "123456789"}},
			expectedError: "Key: 'testBodyRequired.JSONText' Error:Field validation for 'JSONText' failed on the 'required' tag",
		}, {
			name:   "error validation params, query and body",
			method: http.MethodPost,
			url:    "/",
			expectedError: `Key: 'testBodyRequired.ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'testBodyRequired.QueryText' Error:Field validation for 'QueryText' failed on the 'required' tag
Key: 'testBodyRequired.JSONText' Error:Field validation for 'JSONText' failed on the 'required' tag`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.Request, _ = http.NewRequest(test.method, test.url, test.body)
			c.Params = test.params

			var obj testBodyRequired
			if test.expectedError != "" {
				assert.EqualError(t, ParseRequest(c, &obj), test.expectedError)
				return
			}

			assert.NoError(t, ParseRequest(c, &obj))
		})
	}
}
