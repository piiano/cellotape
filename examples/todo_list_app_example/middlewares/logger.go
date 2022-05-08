package middlewares

import (
	r "github.com/piiano/restcontroller/router"
	"log"
	"time"
)

var LoggerMiddleware = r.NewHandler(loggerHandler)

func loggerHandler(c r.Context, request r.Request[r.Nil, r.Nil, r.Nil]) (r.Response[any], error) {
	start := time.Now()
	response, err := c.Next()
	duration := time.Since(start)
	if err != nil {
		log.Printf("[ERROR] error ocured: %s. - %s - [%s] %s\n", err.Error(), duration, request.Method, request.URL.Path)
		return r.Response[any]{}, nil
	}
	log.Printf("[INFO] (status %d | %d bytes | %s) - [%s] %s\n", response.Status, len(response.Body), duration, request.Method, request.URL.Path)
	return r.Response[any]{}, nil
}
