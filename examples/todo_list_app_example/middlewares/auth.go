package middlewares

import (
	"fmt"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	r "github.com/piiano/restcontroller/router"
)

const token = "secret"

var authHeader = fmt.Sprintf("Bearer %s", token)

var AuthMiddleware = r.NewHandler(func(c r.Context, req r.Request[r.Nil, r.Nil, r.Nil]) (r.Response[authResponses], error) {
	if req.Headers.Get("Authorization") != authHeader {
		return r.Send(401, authResponses{Unauthorized: models.HttpError{
			Error:  "Unauthorized",
			Reason: "Authentication failed for provided token",
		}})
	}
	_, err := c.NextFunc(c)
	return r.Response[authResponses]{}, err
})

type authResponses struct {
	Unauthorized models.HttpError `status:"401"`
}
