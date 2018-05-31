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

type Router struct {
	prefix       string
	router       *httprouter.Router
	middlewares  []Middleware
	ErrorHandler ErrorHandler
}

func defaultErrorHandler(c *Context, err error) {
	c.Response.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(c.Response, err.Error())
}

func New() *Router {
	return &Router{
		router:       httprouter.New(),
		ErrorHandler: defaultErrorHandler,
	}
}

func (e Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.ServeHTTP(w, r)
}

func (e *Router) Use(middleware ...Middleware) {
	e.middlewares = append(e.middlewares, middleware...)
}

func (e Router) Sub(prefix string) *Router {
	ne := &Router{
		prefix:       prefix,
		router:       e.router,
		ErrorHandler: e.ErrorHandler,
	}
	copy(ne.middlewares, e.middlewares)

	return ne
}

func (e Router) wrapHandler(handler HandlerFunc) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, pathParams httprouter.Params) {
		ctx := &Context{
			Request:    r,
			Response:   w,
			PathParams: pathParams,
		}

		h := handler
		for i := len(e.middlewares) - 1; i >= 0; i-- {
			h = e.middlewares[i](h)
		}

		if err := h(ctx); err != nil {
			e.ErrorHandler(ctx, err)
		}
	})
}

func (e Router) joinPath(pattern string) string {
	return e.prefix + pattern
}

func (e Router) Static(pattern string, root string) {
	e.router.ServeFiles(pattern+"/*filepath", http.Dir(root))
}

func (e Router) GET(pattern string, handler HandlerFunc) {
	e.router.GET(e.joinPath(pattern), e.wrapHandler(handler))
}

func (e Router) POST(pattern string, handler HandlerFunc) {
	e.router.POST(e.joinPath(pattern), e.wrapHandler(handler))
}

func (e Router) PUT(pattern string, handler HandlerFunc) {
	e.router.PUT(e.joinPath(pattern), e.wrapHandler(handler))
}

func (e Router) PATCH(pattern string, handler HandlerFunc) {
	e.router.PATCH(e.joinPath(pattern), e.wrapHandler(handler))
}

func (e Router) DELETE(pattern string, handler HandlerFunc) {
	e.router.DELETE(e.joinPath(pattern), e.wrapHandler(handler))
}

func (e Router) OPTIONS(pattern string, handler HandlerFunc) {
	e.router.OPTIONS(e.joinPath(pattern), e.wrapHandler(handler))
}

type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	PathParams httprouter.Params
	Values     map[string]interface{}
}

func (c *Context) Set(key string, v interface{}) {
	if c.Values == nil {
		c.Values = make(map[string]interface{})
	}

	c.Values[key] = v
}

func (c Context) Get(key string) (interface{}, bool) {
	v, ok := c.Values[key]
	return v, ok
}

func (c Context) MustGet(key string) interface{} {
	v, ok := c.Get(key)
	if !ok {
		panic(fmt.Errorf("key %s is not exist", key))
	}
	return v
}

const (
	MimePlainText = "plain/text"
	MimeJson      = "application/json; charset=utf8"
	MimeHTML      = "text/html"
)

func (c Context) NoContent(status int) error {
	c.Response.WriteHeader(status)
	return nil
}

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
