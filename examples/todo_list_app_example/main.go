package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/piiano/cellotape/examples/todo_list_app_example/middlewares"
	"github.com/piiano/cellotape/examples/todo_list_app_example/rest"
	"github.com/piiano/cellotape/examples/todo_list_app_example/services"
	"github.com/piiano/cellotape/router"
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
		Use(middlewares.LoggerMiddleware, middlewares.AuthMiddleware, middlewares.PoweredByMiddleware).
		WithGroup(rest.TasksOperationsGroup(tasksService)).
		AsHandler()

	if err != nil {
		return err
	}
	port := 8080
	fmt.Printf("Starting HTTP server on port %d\n", port)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", port), handler); err != nil {
		return err
	}
	return nil
}
