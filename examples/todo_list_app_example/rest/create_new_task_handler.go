package models

import (
	m "github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

type createNewTaskResponses struct {
	OK m.Identifiable `status:"200"`
}

func createNewTaskOperation(tasks services.TasksService) r.OperationHandler {
	return r.OperationFunc(func(
		request r.Request[m.Task, r.Nil, r.Nil],
		send r.Send[createNewTaskResponses],
	) {
		id := tasks.CreateTask(request.Body)
		send(200, createNewTaskResponses{OK: m.Identifiable{ID: id}})
	})
}
