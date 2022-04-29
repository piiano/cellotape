package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

type deleteTaskByIDResponses struct {
	NoContent  r.Nil     `status:"204"`
	BadRequest httpError `status:"400"`
	Gone       httpError `status:"410"`
}

func deleteTaskByIDOperation(tasks services.TasksService) r.OperationHandler {
	return r.OperationFunc(func(
		request r.Request[r.Nil, idPathParam, r.Nil],
		send r.Send[deleteTaskByIDResponses],
	) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			send(400, deleteTaskByIDResponses{
				BadRequest: httpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			})
			return
		}
		if deleted := tasks.DeleteTaskByID(id); deleted {
			send(204, deleteTaskByIDResponses{})
			return
		}
		send(410, deleteTaskByIDResponses{
			Gone: httpError{
				Error:  "gone",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		})
	})
}
