package models

import (
	"fmt"
	"github.com/google/uuid"
	m "github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

func updateTaskByIDOperation(tasks services.TasksService) r.OperationHandler {
	return r.OperationFunc(func(request r.Request[m.Task, idPathParam, r.Nil]) (int, updateTaskByIDResponses) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			return 400, updateTaskByIDResponses{
				BadRequest: httpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			}
		}
		if updated := tasks.UpdateTaskByID(id, request.Body); updated {
			return 204, updateTaskByIDResponses{}
		}
		return 404, updateTaskByIDResponses{
			NotFound: httpError{
				Error:  "not found",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		}
	})
}

type updateTaskByIDResponses struct {
	NoContent  r.Nil     `status:"204"`
	BadRequest httpError `status:"400"`
	NotFound   httpError `status:"404"`
}
