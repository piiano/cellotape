package middlewares

import (
	"github.com/piiano/restcontroller/router"
	"log"
	"time"
)

var LoggerMiddleware = router.NewHandler(func(c router.HandlerContext) (router.Response[any], error) {
	start := time.Now()
	resp, err := router.NextResponse[any](c)
	duration := time.Since(start)
	if err != nil {
		log.Printf("[ERROR] %s %s - error ocured: %s. - %s\n", c.Request.Method, c.Request.URL.Path, err.Error(), duration)
		return resp, err
	}
	log.Printf("[INFO] %s %s - status %d (%d bytes) - %s\n", c.Request.Method, c.Request.URL.Path, resp.Status, len(resp.Bytes), duration)
	return resp, nil
})
