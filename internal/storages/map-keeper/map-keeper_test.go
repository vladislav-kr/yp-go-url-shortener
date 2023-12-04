package mapkeeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper(t *testing.T) {

	t.Run("saving and receiving", func(t *testing.T) {

		const testURL = "https://ya.ru/"

		s := New()

		id, err := s.PostURL(testURL)
		require.NoError(t, err)
		assert.NotEmpty(t, id)

		url, err := s.GetURL(id)
		require.NoError(t, err)
		assert.Equal(t, testURL, url)

	})

}
