package main

import (
	_ "embed"
	"fmt"
	models "github.com/piiano/restcontroller/examples/todo_list_app_example/rest"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	"github.com/piiano/restcontroller/router"
	"net/http"
	"os"
)

//go:embed openapi.yaml
var specData []byte

func main() {
	spec, err := router.NewSpecFromData(specData)
	if err != nil {
		fmt.Println("failed loading the spec")
		fmt.Println(err)
		os.Exit(2)
	}
	tasksService := services.NewTasksService()
	handler, err := router.NewOpenAPI(spec).
		WithContentType(router.JsonContentType{}).
		WithGroup(models.TasksOperationsGroup(tasksService)).
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
