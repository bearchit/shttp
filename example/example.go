package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/bearchit/shttp"
)

func main() {
	r := shttp.New()
	r.ErrorHandler = func(c *shttp.Context, err error) {
		c.String(http.StatusInternalServerError, err.Error())
	}

	r.Use(func(next shttp.HandlerFunc) shttp.HandlerFunc {
		return func(c *shttp.Context) error {
			log.Println(c.Request.RequestURI)
			return next(c)
		}
	})

	r.GET("/", func(c *shttp.Context) error {
		return c.String(http.StatusOK, "Hello")
	})

	r.GET("/error", func(c *shttp.Context) error {
		return errors.New("error from root")
	})

	sub := r.Sub("/sub")
	sub.ErrorHandler = func(c *shttp.Context, err error) {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	sub.GET("/", func(c *shttp.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello from sub",
		})
	})

	sub.GET("/error", func(c *shttp.Context) error {
		return errors.New("error from sub")
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
