package api

import (
	"fmt"
	"github.com/piiano/restcontroller/router"
	"time"
)

var GreetOperationHandler = router.OperationFunc(greetHandler)

func greetHandler(request router.Request[body, pathParams, queryParams], send router.Send[responses]) {
	if request.PathParams.Version != "v1" && request.PathParams.Version != "1" && request.PathParams.Version != "1.0" {
		errMessage := fmt.Sprintf("unsupported version %q", request.PathParams.Version)
		send(400, responses{BadRequest: BadRequest{Message: errMessage}})
		return
	}
	greeting, err := Greet(request.Body.Name, request.Body.DayOfBirth, request.QueryParams.GreetTemplate)
	if err != nil {
		send(400, responses{BadRequest: BadRequest{Message: err.Error()}})
		return
	}
	send(200, responses{OK: OK{Greeting: greeting}})
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
	OK         `status:"200"`
	BadRequest `status:"400"`
}

type OK struct {
	Greeting string `json:"greeting" schema:"{Hello World!}"`
}
type BadRequest struct {
	Message string `json:"message"`
}
