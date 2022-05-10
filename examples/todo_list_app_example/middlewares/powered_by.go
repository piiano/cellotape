package middlewares

import (
	r "github.com/piiano/restcontroller/router"
)

var PoweredByMiddleware = r.NewHandler(poweredByHandler)

func poweredByHandler(c r.Context, _ r.Request[r.Nil, r.Nil, r.Nil]) (r.Response[any], error) {
	c.Writer.Header().Add("X-Powered-By", "Piiano OpenAPI Router")
	_, err := c.Next()
	return r.Response[any]{}, err
}
