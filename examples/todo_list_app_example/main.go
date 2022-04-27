package main

import (
	_ "embed"
	models "github.com/piiano/restcontroller/examples/todo_list_app_example/rest"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	"github.com/piiano/restcontroller/router"
	"net/http"
)

//go:embed openapi.yaml
var specData []byte

func main() {
	tasksService := services.NewTasksService()

	spec, err := router.NewSpecFromData(specData)
	if err != nil {
		panic(err)
	}
	handler, err := router.NewOpenAPI(spec).
		WithGroup(models.TasksOperationsGroup(tasksService)).
		AsHandler()
	if err != nil {
		panic(err)
	}
	if err = http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
