package models

import (
	"fmt"
	"github.com/google/uuid"
	m "github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

type updateTaskByIDResponses struct {
	NoContent  r.Nil     `status:"204"`
	BadRequest httpError `status:"400"`
	NotFound   httpError `status:"404"`
}

func updateTaskByIDOperation(tasks services.TasksService) r.OperationHandler {
	return r.OperationFunc(func(request r.Request[m.Task, idPathParam, r.Nil], send r.Send[updateTaskByIDResponses]) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			send(400, updateTaskByIDResponses{
				BadRequest: httpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			})
			return
		}
		if updated := tasks.UpdateTaskByID(id, request.Body); updated {
			send(204, updateTaskByIDResponses{})
			return
		}
		send(404, updateTaskByIDResponses{
			NotFound: httpError{
				Error:  "not found",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		})
	})
}
