package main

import (
	_ "embed"
	"fmt"
	"github.com/piiano/restcontroller/examples/hello_world_example/api"
	"github.com/piiano/restcontroller/router"
	"net/http"
	"os"
)

//go:embed openapi.yaml
var specData []byte

func main() {
	spec, err := router.NewSpecFromData(specData)
	if err != nil {
		if err != nil {
			fmt.Println("failed loading the spec")
			fmt.Println(err)
			os.Exit(2)
		}
	}
	handler, err := router.NewOpenAPI(spec).
		WithContentType(router.JsonContentType{}).
		WithOperation("greet", api.GreetOperationHandler).
		AsHandler()
	if err != nil {
		fmt.Println("failed creating an handler from the spec")
		fmt.Println(err)
		os.Exit(2)
	}
	if err = http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
