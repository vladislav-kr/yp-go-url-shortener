package router

import (
	"testing"

	"net/http"

	"github.com/stretchr/testify/assert"
)

type Handlers struct{}

func (h *Handlers) SaveHandler(_ http.ResponseWriter, r *http.Request)     {}
func (h *Handlers) RedirectHandler(_ http.ResponseWriter, _ *http.Request) {}

func TestNewRouter(t *testing.T) {

	r := NewRouter(&Handlers{})

	assert.NotEmpty(t, r)
}
