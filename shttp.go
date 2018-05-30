package shttp

import (
	"net/http"

	"fmt"

	"encoding/json"

	"html/template"

	"github.com/julienschmidt/httprouter"
)

type (
	HandlerFunc  func(c *Context) error
	ErrorHandler func(c *Context, err error)
	Middleware   func(HandlerFunc) HandlerFunc
)

type Engine struct {
	router       *httprouter.Router
	ErrorHandler ErrorHandler
}

func defaultErrorHandler(c *Context, err error) {
	c.Response.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(c.Response, err.Error())
}

func New() *Engine {
	return &Engine{
		router:       httprouter.New(),
		ErrorHandler: defaultErrorHandler,
	}
}

func (e Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.ServeHTTP(w, r)
}

func (e Engine) wrapHandler(handler HandlerFunc) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, pathParams httprouter.Params) {
		ctx := &Context{
			Request:    r,
			Response:   w,
			PathParams: pathParams,
		}

		if err := handler(ctx); err != nil {
			e.ErrorHandler(ctx, err)
		}
	})
}

func (e Engine) GET(pattern string, handler HandlerFunc) {
	e.router.GET(pattern, e.wrapHandler(handler))
}

func (e Engine) POST(pattern string, handler HandlerFunc) {
	e.router.POST(pattern, e.wrapHandler(handler))
}

func (e Engine) PUT(pattern string, handler HandlerFunc) {
	e.router.PUT(pattern, e.wrapHandler(handler))
}

func (e Engine) PATCH(pattern string, handler HandlerFunc) {
	e.router.PATCH(pattern, e.wrapHandler(handler))
}

func (e Engine) DELETE(pattern string, handler HandlerFunc) {
	e.router.DELETE(pattern, e.wrapHandler(handler))
}

func (e Engine) OPTIONS(pattern string, handler HandlerFunc) {
	e.router.OPTIONS(pattern, e.wrapHandler(handler))
}

type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	PathParams httprouter.Params
}

const (
	MimePlainText = "plain/text"
	MimeJson      = "application/json; charset=utf8"
	MimeHTML      = "text/html"
)

func (c Context) String(status int, v string) error {
	c.Response.WriteHeader(status)
	c.Response.Header().Set("Content-Type", MimePlainText)
	_, err := c.Response.Write([]byte(v))
	return err
}

func (c Context) JSON(status int, v interface{}) error {
	c.Response.Header().Set("Content-Type", MimeJson)
	c.Response.WriteHeader(status)

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = c.Response.Write(b)
	return err
}

func (c Context) HTML(status int, tpl template.Template, data interface{}) error {
	c.Response.Header().Set("Content-Type", MimeHTML)
	c.Response.WriteHeader(status)
	return tpl.Execute(c.Response, data)
}
