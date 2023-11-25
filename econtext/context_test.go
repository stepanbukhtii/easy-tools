package econtext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	type testKey struct{}
	ctx := WithValue[testKey](context.Background(), "text value")
	assert.Equal(t, "text value", Value[testKey, string](ctx))

	valueDefault := ValueDefault[testKey, string](context.Background(), "default value")
	assert.Equal(t, "default value", valueDefault)
}

func TestSetLogField(t *testing.T) {
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
