package main

import (
	"github.com/piiano/restcontroller/examples/hello_world_example"
	"github.com/piiano/restcontroller/router"
	"net/http"
)

func main() {
	spec, err := router.NewSpecFromData(hello_world_example.Spec)
	if err != nil {
		panic(err)
	}
	handler, err := router.NewOpenAPI().
		WithSpec(spec).
		WithOperation("greet", hello_world_example.GreetOperationHandler).
		AsHandler()
	if err != nil {
		panic(err)
	}
	if err = http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
