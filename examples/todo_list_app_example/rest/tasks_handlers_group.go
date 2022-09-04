package rest

import (
	_ "embed"

	"github.com/piiano/cellotape/examples/todo_list_app_example/services"
	r "github.com/piiano/cellotape/router"
)

func TasksOperationsGroup(tasks services.TasksService) r.Group {
	return r.NewGroup().
		WithOperation("getTasksPage", getTasksPageOperation(tasks)).
		WithOperation("createNewTask", createNewTaskOperation(tasks)).
		WithOperation("getTaskByID", getTaskByIDOperation(tasks)).
		WithOperation("deleteTaskByID", deleteTaskByIDOperation(tasks)).
		WithOperation("updateTaskByID", updateTaskByIDOperation(tasks))

}

type idPathParam struct {
	// https://github.com/gin-gonic/gin/issues/2423
	ID string `uri:"id"`
}
