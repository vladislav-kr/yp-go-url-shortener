package mapkeeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper(t *testing.T) {

	stor := New()

	tests := []struct {
		name        string
		url         string
		expectedURL string
		isError     bool
	}{
		{
			name:        "ulr найден",
			url:         "https://ya.ru/",
			expectedURL: "https://ya.ru/",
		},
		{
			name:    "url не найден",
			isError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				id  string
				err error
			)

			if !tt.isError {
				id, err = stor.PostURL(tt.url)
				require.NoError(t, err)
				assert.NotEmpty(t, id)
			}

			url, err := stor.GetURL(id)

			if tt.isError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.url, url)

		})
	}

}
