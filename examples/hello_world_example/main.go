package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/piiano/cellotape/examples/hello_world_example/api"
	r "github.com/piiano/cellotape/router"
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
		WithContentType(r.JSONContentType{}).
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
