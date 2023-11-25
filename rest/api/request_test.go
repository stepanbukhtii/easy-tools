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
	ID        string `uri:"id" json:",omitempty"`
	QueryText string `form:"query_text" json:",omitempty"`
	JSONText  string `json:"json_text"`
}

func TestQuery(t *testing.T) {

	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

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
			url:    "/?query_text=query_text",
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
			assert.NoError(t, Parse(c, &obj))

			assert.Equal(t, test.expectedObject.ID, obj.ID)
			assert.Equal(t, test.expectedObject.QueryText, obj.QueryText)
			assert.Equal(t, test.expectedObject.JSONText, obj.JSONText)
		})
	}
}
