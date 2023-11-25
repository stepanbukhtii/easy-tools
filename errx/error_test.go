package errx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorIs(t *testing.T) {
	baseError := errors.New("test error")
	customErr := Wrap(baseError, "custom_error")
	customErr2 := Wrap(baseError, "custom_error")
	customErr3 := New("custom_error")
	tests := []struct {
		name      string
		err       error
		targetErr error
		expected  bool
	}{
		{
			name:      "custom error is target error",
			err:       customErr,
			targetErr: baseError,
			expected:  true,
		}, {
			name:      "two custom error",
			err:       customErr,
			targetErr: customErr2,
			expected:  true,
		},
		{
			name:      "different errors",
			err:       customErr,
			targetErr: customErr3,
			expected:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, errors.Is(test.err, test.targetErr))
		})
	}
}
