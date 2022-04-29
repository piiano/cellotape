package models

import (
	"fmt"
	"github.com/google/uuid"
	m "github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

type getTaskByIDResponses struct {
	OK         m.Task    `status:"200"`
	BadRequest httpError `status:"400"`
	NotFound   httpError `status:"404"`
}

func getTaskByIDOperation(tasks services.TasksService) r.OperationHandler {
	return r.OperationFunc(func(
		request r.Request[r.Nil, idPathParam, r.Nil],
		send r.Send[getTaskByIDResponses],
	) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			send(400, getTaskByIDResponses{
				BadRequest: httpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			})
			return
		}
		if task, found := tasks.GetTaskByID(id); found {
			send(200, getTaskByIDResponses{OK: task})
			return
		}
		send(404, getTaskByIDResponses{
			NotFound: httpError{
				Error:  "not found",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		})
	})
}
