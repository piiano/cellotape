package rest

import (
	"fmt"
	"github.com/google/uuid"
	m "github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

func getTaskByIDOperation(tasks services.TasksService) r.Handler {
	return r.NewOperationHandler(func(request r.Request[r.Nil, idPathParam, r.Nil]) (r.Response[getTaskByIDResponses], error) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			return r.Send(400, getTaskByIDResponses{
				BadRequest: m.HttpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			})
		}
		if task, found := tasks.GetTaskByID(id); found {
			return r.Send(200, getTaskByIDResponses{OK: task})
		}
		return r.Send(404, getTaskByIDResponses{
			NotFound: m.HttpError{
				Error:  "not found",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		})
	})
}

type getTaskByIDResponses struct {
	OK         m.Task      `status:"200"`
	BadRequest m.HttpError `status:"400"`
	NotFound   m.HttpError `status:"404"`
}
