package rest

import (
	m "github.com/piiano/cellotape/examples/todo_list_app_example/models"
	"github.com/piiano/cellotape/examples/todo_list_app_example/services"
	r "github.com/piiano/cellotape/router"
	"github.com/piiano/cellotape/router/utils"
)

func createNewTaskOperation(tasks services.TasksService) r.Handler {
	return r.NewHandler(func(c *r.Context, request r.Request[m.Task, utils.Nil, utils.Nil]) (r.Response[createNewTaskResponses], error) {
		id := tasks.CreateTask(request.Body)
		return r.SendOKJSON(createNewTaskResponses{OK: m.Identifiable{ID: id}}), nil
	})
}

type createNewTaskResponses struct {
	OK m.Identifiable `status:"200"`
}
