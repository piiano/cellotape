package example

import (
	"github.com/piiano/restcontroller/restcontroller"
	"time"
)

// We export the REST controller to be loaded by the REST engine.

var GreetController = restcontroller.NewController(greetController)

// All types of the REST controller are defined in package scope.
// This should always be the case to clarify they are strictly used for defining the REST interface and for nothing else.
// This way it is much easier to maintain healthy decoupling between the REST controllers to the service layer.

func greetController(params restcontroller.Params[body, pathParams, queryParams]) (response, error) {
	if params.Path.Version == "v1" {
		greeting, err := GreetV1(params.Body.Name)
		return response{Greeting: greeting}, err
	}
	greeting, err := GreetV2(params.Body.Name, params.Body.DayOfBirth, params.Query.GreetTemplate)
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
