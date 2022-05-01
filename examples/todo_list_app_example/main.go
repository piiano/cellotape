package main

import (
	_ "embed"
	"fmt"
	models "github.com/piiano/restcontroller/examples/todo_list_app_example/rest"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	"github.com/piiano/restcontroller/router"
	"log"
	"net/http"
)

//go:embed openapi.yaml
var specData []byte

func main() {
	if err := mainHandler(); err != nil {
		log.Fatal(err)
	}
}

func mainHandler() error {
	spec, err := router.NewSpecFromData(specData)
	if err != nil {
		return err
	}
	tasksService := services.NewTasksService()
	handler, err := router.NewOpenAPIRouter(spec).
		WithContentType(router.JsonContentType{}).
		WithGroup(models.TasksOperationsGroup(tasksService)).
		AsHandler()
	if err != nil {
		return err
	}
	port := 8080
	fmt.Printf("Starting HTTP server on port %d", port)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		return err
	}
	return nil
}
