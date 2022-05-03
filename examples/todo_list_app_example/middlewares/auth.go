package middlewares

import (
	"fmt"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/router"
)

const token = "secret"

var authHeader = fmt.Sprintf("Bearer %s", token)

var AuthMiddleware = router.NewHandler(func(c router.HandlerContext) (router.Response[authMiddlewareResponses], error) {
	if c.Request.Header.Get("Authorization") != authHeader {
		return router.Send(401, authMiddlewareResponses{Unauthorized: models.HttpError{
			Error:  "Unauthorized",
			Reason: "Authentication failed for provided token",
		}})
	}
	return router.NextResponse[authMiddlewareResponses](c)
})

type authMiddlewareResponses struct {
	Unauthorized models.HttpError `status:"401"`
}
