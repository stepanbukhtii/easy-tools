package easycontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	ctx := context.Background()
	ctx = WithValue[traceIDKey](ctx, "text value")
	value := Value[traceIDKey, string](ctx)
	assert.Equal(t, "text value", value)

	traceID := GetDefault[localeKey, string](ctx, "default value")
	assert.Equal(t, "default value", traceID)
}

func TestTraceID(t *testing.T) {
	ctx := context.Background()
	ctx = SetTraceID(ctx, "trace id")
	traceID := TraceID(ctx)
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
		ctx = SetLogField(ctx, test.key, test.value)
		expected[test.key] = test.value
		fields := LogFields(ctx)
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
		fields := LogFields(ctx)
		assert.Equal(t, expected, fields)
	}
}
