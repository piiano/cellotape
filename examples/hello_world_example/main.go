package main

import (
	_ "embed"
	"github.com/piiano/restcontroller/examples/hello_world_example/api"
	r "github.com/piiano/restcontroller/router"
	"log"
	"net/http"
)

//go:embed openapi.yaml
var specData []byte

func main() {
	if err := handleMain(); err != nil {
		log.Fatal(err)
	}
}

func handleMain() error {
	spec, err := r.NewSpecFromData(specData)
	if err != nil {
		return err
	}
	handler, err := r.NewOpenAPIRouter(spec).
		WithOperation("greet", api.GreetOperationHandler).
		AsHandler()
	if err != nil {
		return err
	}
	if err = http.ListenAndServe(":8080", handler); err != nil {
		return err
	}
	return nil
}
