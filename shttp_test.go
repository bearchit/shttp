package shttp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	s := New()

	s.GET("/", func(c *Context) error {
		return c.String(http.StatusOK, "Hello")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	s.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	b, err := ioutil.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, MimePlainText, rec.Header().Get("Content-Type"))
	assert.Equal(t, "Hello", string(b))
}
