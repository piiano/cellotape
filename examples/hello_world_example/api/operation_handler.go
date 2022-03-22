package api

import (
	"fmt"
	r "github.com/piiano/cellotape/router"
	"time"
)

var GreetOperationHandler = r.NewHandler(greetHandler)

func greetHandler(_ r.Context, request r.Request[body, pathParams, queryParams]) (r.Response[responses], error) {
	if request.PathParams.Version != "v1" && request.PathParams.Version != "1" && request.PathParams.Version != "1.0" {
		errMessage := fmt.Sprintf("unsupported version %q", request.PathParams.Version)
		return r.Send(400, responses{BadRequest: badRequest{Message: errMessage}})
	}
	greeting, err := greet(request.Body.Name, request.Body.DayOfBirth, request.QueryParams.GreetTemplate)
	if err != nil {
		return r.Send(400, responses{BadRequest: badRequest{Message: err.Error()}})
	}
	return r.Send(200, responses{OK: ok{Greeting: greeting}})
}

type body struct {
	Name       string    `json:"name"`
	DayOfBirth time.Time `json:"day_of_birth"`
}
type pathParams struct {
	Version string `uri:"version"`
}
type queryParams struct {
	GreetTemplate string `form:"greetTemplate"`
}
type responses struct {
	OK         ok         `status:"200"`
	BadRequest badRequest `status:"400"`
}

type ok struct {
	Greeting string `json:"greeting" schema:"{Hello World!}"`
}
type badRequest struct {
	Message string `json:"message"`
}
