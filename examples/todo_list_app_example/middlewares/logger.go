package middlewares

import (
	"log"
	"time"

	r "github.com/piiano/cellotape/router"
)

var LoggerMiddleware = r.RawHandler(loggerHandler)

func loggerHandler(c *r.Context) error {
	start := time.Now()
	response, err := c.Next()
	durations := struct {
		Handler time.Duration
		Read    time.Duration
		Write   time.Duration
	}{
		Handler: time.Since(start) - c.Durations.ReadDuration(),
		Read:    c.Durations.ReadDuration(),
		Write:   c.Durations.WriteDuration(),
	}

	if err != nil {
		log.Printf("[ERROR] error occurred: %s. - durations: %+v - [%s] %s\n", err.Error(), durations, c.Request.Method, c.Request.URL.Path)
		return err
	}
	log.Printf("[INFO] (status %d | %d bytes | durations: %+v) - [%s] %s\n", response.Status, len(response.Body), durations, c.Request.Method, c.Request.URL.Path)
	return nil
}
