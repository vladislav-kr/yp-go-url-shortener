package mapkeeper

import (
	"os"
	"path/filepath"
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

func TestKeeperSaveReadFile(t *testing.T) {

	absDir, err := filepath.Abs("")
	require.NoError(t, err)
	dir, err := os.MkdirTemp(absDir, "testfile*")

	t.Cleanup(func() {
		t.Helper()
		err = os.RemoveAll(dir)
		require.NoError(t, err)
	})

	storage := New()
	saveData := map[string]string{
		"fdfsfwewq2":  "https://ya.ru/",
		"fdfdd455654": "https://practicum.yandex.ru/",
	}
	path := dir + "/testfileKeeper.json"
	storage.storage = saveData

	err = storage.SaveToFile(path)
	require.NoError(t, err)

	storage.storage = map[string]string{}

	err = storage.LoadFromFile(path)
	require.NoError(t, err)

	assert.EqualValues(t, storage.storage, saveData)

}
