package rest

import (
	"fmt"
	"github.com/google/uuid"
	m "github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

func deleteTaskByIDOperation(tasks services.TasksService) r.Handler {
	return r.NewOperationHandler(func(request r.Request[r.Nil, idPathParam, r.Nil]) (r.Response[deleteTaskByIDResponses], error) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			return r.Send(400, deleteTaskByIDResponses{
				BadRequest: m.HttpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			})
		}
		if deleted := tasks.DeleteTaskByID(id); deleted {
			return r.Send(204, deleteTaskByIDResponses{})
		}
		return r.Send(410, deleteTaskByIDResponses{
			Gone: m.HttpError{
				Error:  "gone",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		})
	})
}

type deleteTaskByIDResponses struct {
	NoContent  r.Nil       `status:"204"`
	BadRequest m.HttpError `status:"400"`
	Gone       m.HttpError `status:"410"`
}
