package api

import (
	"fmt"
	"github.com/piiano/restcontroller/router"
	"time"
)

var GreetOperationHandler = router.OperationFunc(greetHandler)

func greetHandler(request router.Request[body, pathParams, queryParams]) (int, responses) {
	if request.PathParams.Version != "v1" && request.PathParams.Version != "1" && request.PathParams.Version != "1.0" {
		errMessage := fmt.Sprintf("unsupported version %q", request.PathParams.Version)
		return 400, responses{BadRequest: badRequest{Message: errMessage}}
	}
	greeting, err := greet(request.Body.Name, request.Body.DayOfBirth, request.QueryParams.GreetTemplate)
	if err != nil {
		return 400, responses{BadRequest: badRequest{Message: err.Error()}}
	}
	return 200, responses{OK: ok{Greeting: greeting}}
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
