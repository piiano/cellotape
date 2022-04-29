package models

import (
	_ "embed"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

type httpError struct {
	Error  string `json:"error"`
	Reason string `json:"reason,omitempty"`
}

type idPathParam struct {
	// https://github.com/gin-gonic/gin/issues/2423
	ID string `uri:"id"`
}

func TasksOperationsGroup(tasks services.TasksService) r.Group {
	return r.NewGroup().
		WithOperation("getTasksPage", getTasksPageOperation(tasks)).
		WithOperation("createNewTask", createNewTaskOperation(tasks)).
		WithOperation("getTaskByID", getTaskByIDOperation(tasks)).
		WithOperation("deleteTaskByID", deleteTaskByIDOperation(tasks)).
		WithOperation("updateTaskByID", updateTaskByIDOperation(tasks))
}
