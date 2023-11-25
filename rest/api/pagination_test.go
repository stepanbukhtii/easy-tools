package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPagination(t *testing.T) {
	tests := []struct {
		name         string
		pagination   Pagination
		countRecords int64
		expectLimit  int64
		expectOffset int64
		expectPages  int64
	}{
		{
			name: "page 1 per page 10",
			pagination: Pagination{
				Page:    1,
				PerPage: 10,
			},
			countRecords: 99,
			expectLimit:  10,
			expectOffset: 0,
			expectPages:  10,
		},
		{
			name: "page 5 per page 10",
			pagination: Pagination{
				Page:    5,
				PerPage: 10,
			},
			countRecords: 151,
			expectLimit:  10,
			expectOffset: 40,
			expectPages:  16,
		},
		{
			name:         "page 10 default per page",
			pagination:   Pagination{Page: 3},
			countRecords: 75,
			expectLimit:  10,
			expectOffset: 20,
			expectPages:  8,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			limit, offset := test.pagination.LimitOffset()
			require.Equal(t, test.expectLimit, limit)
			require.Equal(t, test.expectOffset, offset)
			require.Equal(t, test.expectPages, test.pagination.Pages(test.countRecords))
		})
	}
}
