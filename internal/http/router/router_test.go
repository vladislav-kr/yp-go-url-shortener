package router

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Handlers struct{}
func (h *Handlers) SaveHandler(_ http.ResponseWriter, _ *http.Request)     {}
func (h *Handlers) RedirectHandler(_ http.ResponseWriter, _ *http.Request) {}

func TestNewRouter(t *testing.T) {

	t.Run("the router is created successfully", func(t *testing.T) {
		r, err := NewRouter(&Handlers{})

		assert.NoError(t, err)
		assert.NotEmpty(t, r)

	})
}
