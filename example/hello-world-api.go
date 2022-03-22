package example

import (
	"fmt"
	"github.com/piiano/restcontroller/restcontroller"
	"time"
)

type Body struct {
	// Name of the person to greet.
	Name string
	// Multiline description
	// works too.
	DayOfBirth time.Time `json:"day_of_birth"`
}
type Response struct {
	Greeting string `json:"greeting" schema:"{Hello World!}"`
}

func Greet(name string) (string, error) {
	return fmt.Sprintf("Hello %s!", name), nil
}

type GreetPathParams struct {
	ID string `json:"id"`
}

var GreetController restcontroller.ControllerFn[Body, GreetPathParams, any, Response] = func(params restcontroller.Params[Body, GreetPathParams, any]) (Response, error) {
	greeting, err := Greet(params.Body.Name)
	return Response{Greeting: greeting}, err
}
