package middlewares

import (
	"log"
	"time"

	r "github.com/piiano/cellotape/router"
)

var LoggerMiddleware = r.NewHandler(loggerHandler)

func loggerHandler(c *r.Context, request r.Request[r.Nil, r.Nil, r.Nil]) (r.Response[any], error) {
	start := time.Now()
	response, err := c.Next()
	duration := time.Since(start)
	if err != nil {
		log.Printf("[ERROR] error occurred: %s. - %s - [%s] %s\n", err.Error(), duration, c.Request.Method, c.Request.URL.Path)
		return r.Response[any]{}, nil
	}
	log.Printf("[INFO] (status %d | %d bytes | %s) - [%s] %s\n", response.Status, len(response.Body), duration, c.Request.Method, c.Request.URL.Path)
	return r.Response[any]{}, nil
}
