package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

func deleteTaskByIDOperation(tasks services.TasksService) r.OperationHandler {
	return r.OperationFunc(func(request r.Request[r.Nil, idPathParam, r.Nil]) (int, deleteTaskByIDResponses) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			return 400, deleteTaskByIDResponses{
				BadRequest: httpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			}
		}
		if deleted := tasks.DeleteTaskByID(id); deleted {
			return 204, deleteTaskByIDResponses{}
		}
		return 410, deleteTaskByIDResponses{
			Gone: httpError{
				Error:  "gone",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		}
	})
}

type deleteTaskByIDResponses struct {
	NoContent  r.Nil     `status:"204"`
	BadRequest httpError `status:"400"`
	Gone       httpError `status:"410"`
}
