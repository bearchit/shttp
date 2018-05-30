## Usage

```go
package main

import (
    "log"
    "net/http"
        
    "github.com/bearchit/shttp"
)

func main() {
	s := shttp.New()
	s.GET("/", func(c *shttp.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello, World",
		})
	})
	
	log.Fatal(http.ListenAndServe(":8080", s))
}
```