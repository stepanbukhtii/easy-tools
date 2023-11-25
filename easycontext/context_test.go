package easycontext

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTraceID(t *testing.T) {
	ctx := context.Background()
	ctx = AddTraceID(ctx, "trace id")
	traceID, ok := GetTraceID(ctx)
	assert.True(t, ok)
	assert.Equal(t, "trace id", traceID)
}

func TestAddLogField(t *testing.T) {
	ctx := context.Background()
	expected := make(map[string]string)

	var tests = []struct {
		key   string
		value string
	}{
		{key: "key 1", value: "value 1"},
		{key: "key 2", value: "value 2"},
		{key: "key 3", value: "value 3"},
		{key: "key 1", value: "value 12"},
	}
	for _, test := range tests {
		ctx = AddLogField(ctx, test.key, test.value)
		expected[test.key] = test.value
		fields, ok := GetLogFields(ctx)
		assert.True(t, ok)
		assert.Equal(t, expected, fields)
	}
}

func TestCopyLogField(t *testing.T) {
	ctx := context.Background()
	expected := make(map[string]string)

	var tests = []struct {
		key   string
		value string
	}{
		{key: "key 1", value: "value 1"},
		{key: "key 2", value: "value 2"},
		{key: "key 3", value: "value 3"},
		{key: "key 1", value: "value 12"},
	}
	for _, test := range tests {
		ctx = CopyLogFields(ctx, test.key, test.value)
		expected[test.key] = test.value
		fields, ok := GetLogFields(ctx)
		assert.True(t, ok)
		assert.Equal(t, expected, fields)
	}
}
