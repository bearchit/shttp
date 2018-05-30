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

func TestMiddleware(t *testing.T) {
	mw1 := func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			c.String(http.StatusOK, "mw1")
			return next(c)
		}
	}

	mw2 := func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			c.String(http.StatusOK, "mw2")
			return next(c)
		}
	}

	s := New()
	s.Use(mw1, mw2)
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
	assert.Equal(t, "mw1mw2Hello", string(b))
}
