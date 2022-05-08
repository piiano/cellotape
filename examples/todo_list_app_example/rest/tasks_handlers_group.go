package rest

import (
	_ "embed"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

func TasksOperationsGroup(tasks services.TasksService) r.Group {
	return r.NewGroup().
		WithOperation("getTasksPage", getTasksPageOperation(tasks)).
		WithOperation("createNewTask", createNewTaskOperation(tasks)).
		WithGroup(r.NewGroup().
			Use(r.NewHandler(func(c r.Context, req r.Request[r.Nil, idPathParam, r.Nil]) (r.Response[any], error) {
				_, err := c.NextFunc(c)
				return r.Response[any]{}, err
			})).
			WithOperation("getTaskByID", getTaskByIDOperation(tasks)).
			WithOperation("deleteTaskByID", deleteTaskByIDOperation(tasks)).
			WithOperation("updateTaskByID", updateTaskByIDOperation(tasks)),
		)

}

type idPathParam struct {
	// https://github.com/gin-gonic/gin/issues/2423
	ID string `uri:"id"`
}
