package middlewares

import (
	r "github.com/piiano/cellotape/router"
	"github.com/piiano/cellotape/router/utils"
)

var PoweredByMiddleware = r.NewHandler(poweredByHandler)

func poweredByHandler(c *r.Context, _ r.Request[utils.Nil, utils.Nil, utils.Nil]) (r.Response[any], error) {
	c.Writer.Header().Add("X-Powered-By", "Piiano OpenAPI Router")
	_, err := c.Next()
	return r.Response[any]{}, err
}
