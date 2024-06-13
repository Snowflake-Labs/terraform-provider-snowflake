package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnumToTerraformStringList(t *testing.T) {
	type customString string
	tests := []struct {
		name   string
		values []customString
		want   string
	}{
		{
			name:   "one element",
			values: []customString{"a"},
			want:   `["a"]`,
		}, {
			name:   "more elements",
			values: []customString{"a", "b"},
			want:   `["a" "b"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, EnumToTerraformStringList(tt.values))
		})
	}
}
