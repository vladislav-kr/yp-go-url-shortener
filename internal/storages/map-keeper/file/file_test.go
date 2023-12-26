package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
)

func TestReadWriteFile(t *testing.T) {
	absDir, err := filepath.Abs("")
	require.NoError(t, err)
	dir, err := os.MkdirTemp(absDir, "testfile*")
	require.NoError(t, err)
	f, err := os.CreateTemp(dir, "test*.json")
	require.NoError(t, err)
	name := f.Name()
	f.Close()

	t.Cleanup(func() {
		t.Helper()
		err = os.RemoveAll(dir)
		require.NoError(t, err)
	})

	p, err := NewProducer(name)
	require.NoError(t, err)
	saveData := models.FileURL{
		ShortURL:    "jsjkfkjdf",
		OriginalURL: "https://ya.ru/",
	}

	err = p.Write(&saveData)
	require.NoError(t, err)
	err = p.Close()
	require.NoError(t, err)

	fileData := models.FileURL{}
	c, err := NewConsumer(name)
	require.NoError(t, err)
	c.Decode(&fileData)
	err = c.Close()
	require.NoError(t, err)

	assert.EqualValues(t, saveData, fileData)

}
