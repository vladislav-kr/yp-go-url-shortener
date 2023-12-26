package fileutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFullPathFromRelative(t *testing.T) {

	absPath, err := filepath.Abs("")
	require.NoError(t, err)

	cases := []struct {
		name         string
		relativePath string
		fullPath     string
		isError      bool
	}{
		{
			name:         "one slash",
			relativePath: "/tmp/short-url-db.json",
			fullPath:     absPath + "\\tmp\\short-url-db.json",
		},
		{
			name:         "one backslash",
			relativePath: `\tmp\short-url-db.json`,
			fullPath:     absPath + "\\tmp\\short-url-db.json",
		},
		{
			name:         "repeat slash",
			relativePath: "//tmp//short-url-db.json",
			fullPath:     absPath + "\\tmp\\short-url-db.json",
			isError:      true,
		},
		{
			name:         "repeat backslash",
			relativePath: `\\tmp\\short-url-db.json`,
			fullPath:     absPath + "\\tmp\\short-url-db.json",
			isError:      true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path, err := CreateFullPathFromRelative(tc.relativePath)
			require.NoError(t, err)
			if tc.isError {
				assert.NotEqual(t, tc.fullPath, path)
				return
			}
			assert.Equal(t, tc.fullPath, path)
		})
	}
}
