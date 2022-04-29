package main

import (
	_ "embed"
	"github.com/piiano/restcontroller/examples/hello_world_example/api"
	"github.com/piiano/restcontroller/router"
	"net/http"
)

//go:embed api/openapi.yaml
var spec []byte

func main() {
	spec, err := router.NewSpecFromData(spec)
	if err != nil {
		panic(err)
	}
	handler, err := router.NewOpenAPI(spec).
		WithContentType(router.JsonContentType{}).
		WithOperation("greet", api.GreetOperationHandler).
		AsHandler()
	if err != nil {
		panic(err)
	}
	if err = http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
