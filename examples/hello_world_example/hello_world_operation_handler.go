package hello_world_example

import (
	_ "embed"
	"github.com/piiano/restcontroller/router"
	"time"
)

//go:embed hello_world_openapi.yaml
var Spec []byte

// We export the REST controller to be loaded by the REST engine.

var GreetOperationHandler = router.NewOperationHandler(greetHandler)

// All types of the REST controller are defined in package scope.
// This should always be the case to clarify they are strictly used for defining the REST interface and for nothing else.
// This way it is much easier to maintain healthy decoupling between the REST controllers to the service layer.

func greetHandler(request router.Request[body, pathParams, queryParams]) (response, error) {
	if request.PathParameters.Version == "v1" {
		greeting, err := GreetV1(request.Body.Name)
		return response{Greeting: greeting}, err
	}
	greeting, err := GreetV2(request.Body.Name, request.Body.DayOfBirth, request.QueryParameters.GreetTemplate)
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
