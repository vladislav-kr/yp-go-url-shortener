package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {

	var envs = map[string]string{
		"SERVER_ADDRESS": ":9091",
		"BASE_URL":       "http://localhost:9091",
	}

	tests := []struct {
		name         string
		beforeFunc   func(t *testing.T)
		expectedAddr string
		expectedURL  string
	}{
		{
			name: "env variables exist",
			beforeFunc: func(t *testing.T) {
				for k, v := range envs {
					t.Setenv(k, v)
				}
			},
			expectedAddr: ":9091",
			expectedURL:  "http://localhost:9091",
		},
		{
			name:         "env variables don't exist",
			beforeFunc:   func(t *testing.T) {},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc(t)
			cfg, err := LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, tt.expectedAddr, cfg.HTTP.Host)
			assert.Equal(t, tt.expectedURL, cfg.URLShortener.RedirectHost)

		})
	}

}
