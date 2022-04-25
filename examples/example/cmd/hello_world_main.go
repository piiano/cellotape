package main

import (
	"github.com/piiano/restcontroller/examples/example"
	"github.com/piiano/restcontroller/router"
	"net/http"
)

func main() {
	spec, err := router.NewSpecFromData(example.Spec)
	if err != nil {
		panic(err)
	}
	handler, err := router.NewOpenAPI().
		WithSpec(spec).
		WithOperation("greet", example.GreetOperationHandler).
		AsHandler()
	if err != nil {
		panic(err)
	}
	if err = http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
