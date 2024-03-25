package cryptoutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

)

func BenchmarkGenerateRandomString10(b *testing.B) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 5",
			size: 1,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 15",
			size: 0,
		},
		{
			name: "size = 20",
			size: -1,
		},
	}
	for _, tt := range tests {
		tt := tt
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GenerateRandomString(tt.size)
			}
		})
	}

}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name         string
		size         int
		expectedSize int
		isError      bool
	}{
		{
			name:         "size = 1",
			size:         1,
			expectedSize: 1,
		},
		{
			name:         "size = 10",
			size:         10,
			expectedSize: 10,
		},
		{
			name:         "size = 0",
			size:         0,
			expectedSize: 0,
			isError:      true,
		},
		{
			name:         "size = -1",
			size:         -1,
			expectedSize: -1,
			isError:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id, err := GenerateRandomString(tt.size)

			if tt.isError {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tt.expectedSize, len(id))

		})
	}
}
