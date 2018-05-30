package shttp

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

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

func TestEngine_Sub(t *testing.T) {
	s := New()
	a := s.Sub("/a")
	assert.Equal(t, "/a", a.prefix)
}

func TestJoinPath(t *testing.T) {
	s := New()
	a := s.Sub("/a")
	assert.Equal(t, "/a/hello", a.joinPath("/hello"))
}

func TestSubroute(t *testing.T) {
	s := New()

	{
		sub1 := s.Sub("/sub1")
		sub1.GET("/hello", func(c *Context) error {
			return c.String(http.StatusOK, "/sub1/hello")
		})
		sub1.ErrorHandler = ErrorHandler(func(c *Context, err error) {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error from sub1"))
		})

		req := httptest.NewRequest(http.MethodGet, "/sub1/hello", nil)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		b, err := ioutil.ReadAll(rec.Body)
		assert.NoError(t, err)
		assert.Equal(t, MimePlainText, rec.Header().Get("Content-Type"))
		assert.Equal(t, "/sub1/hello", string(b))
	}
	{
		sub2 := s.Sub("/sub2")
		sub2.GET("/hello", func(c *Context) error {
			return c.String(http.StatusOK, "/sub2/hello")
		})
		sub2.ErrorHandler = ErrorHandler(func(c *Context, err error) {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error from sub2"))
		})

		req := httptest.NewRequest(http.MethodGet, "/sub2/hello", nil)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		b, err := ioutil.ReadAll(rec.Body)
		assert.NoError(t, err)
		assert.Equal(t, MimePlainText, rec.Header().Get("Content-Type"))
		assert.Equal(t, "/sub2/hello", string(b))
	}

}

func TestSubroute_ErrorHandler(t *testing.T) {
	s := New()
	s.GET("/", func(c *Context) error {
		return errors.New("")
	})

	sub := s.Sub("/a")
	sub.ErrorHandler = ErrorHandler(func(c *Context, err error) {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error from sub"))
	})
	sub.GET("/hello", func(c *Context) error {
		return errors.New("")
	})

	{
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	{
		req := httptest.NewRequest(http.MethodGet, "/a/hello", nil)
		rec := httptest.NewRecorder()

		s.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestContext_NoContent(t *testing.T) {
	s := New()
	s.GET("/", func(c *Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	s.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	b, err := ioutil.ReadAll(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(b))
}
