package rest

import (
	"fmt"
	"github.com/google/uuid"
	m "github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

func updateTaskByIDOperation(tasks services.TasksService) r.Handler {
	return r.NewOperationHandler(func(request r.Request[m.Task, idPathParam, r.Nil]) (r.Response[updateTaskByIDResponses], error) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			return r.Send(400, updateTaskByIDResponses{
				BadRequest: m.HttpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			})
		}
		if updated := tasks.UpdateTaskByID(id, request.Body); updated {
			return r.Send(204, updateTaskByIDResponses{})
		}
		return r.Send(404, updateTaskByIDResponses{
			NotFound: m.HttpError{
				Error:  "not found",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		})
	})
}

type updateTaskByIDResponses struct {
	NoContent  r.Nil       `status:"204"`
	BadRequest m.HttpError `status:"400"`
	NotFound   m.HttpError `status:"404"`
}
