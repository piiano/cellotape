package api

import (
	"fmt"
	"github.com/piiano/restcontroller/router"
	"time"
)

var GreetOperationHandler = router.OperationFunc(greetHandler)

func greetHandler(request router.Request[body, pathParams, queryParams]) (response, error) {
	if request.PathParams.Version != "v1" {
		return response{}, fmt.Errorf("unsupported version %q", request.PathParams.Version)
	}
	greeting, err := Greet(request.Body.Name, request.Body.DayOfBirth, request.QueryParams.GreetTemplate)
	return response{Greeting: greeting}, err
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
type response struct {
	Greeting string `json:"greeting" schema:"{Hello World!}"`
}
