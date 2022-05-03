package middlewares

import (
	"github.com/piiano/restcontroller/router"
	"log"
)

var LoggerMiddleware = router.NewHandler(func(c router.HandlerContext) (router.Response[any], error) {
	resp, err := router.NextResponse[any](c)
	if err != nil {
		log.Printf("[ERROR] %s %s - error ocured: %s\n", c.Request.Method, c.Request.URL.Path, err.Error())
		return resp, err
	}
	log.Printf("[INFO] %s %s - status %d (%d bytes)\n", c.Request.Method, c.Request.URL.Path, resp.Status, len(resp.Bytes))
	return resp, nil
})
