package rest

import (
	"fmt"
	"github.com/google/uuid"
	m "github.com/piiano/cellotape/examples/todo_list_app_example/models"
	"github.com/piiano/cellotape/examples/todo_list_app_example/services"
	r "github.com/piiano/cellotape/router"
	"net/http"
)

func updateTaskByIDOperation(tasks services.TasksService) r.Handler {
	return r.NewHandler(func(_ r.Context, request r.Request[m.Task, idPathParam, r.Nil]) (r.Response[updateTaskByIDResponses], error) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			return r.SendJSON(updateTaskByIDResponses{
				BadRequest: m.HttpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			}).Status(http.StatusBadRequest), nil
		}
		if updated := tasks.UpdateTaskByID(id, request.Body); updated {
			return r.SendJSON(updateTaskByIDResponses{}).Status(http.StatusNoContent), nil
		}
		return r.SendJSON(updateTaskByIDResponses{
			NotFound: m.HttpError{
				Error:  "not found",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		}).Status(http.StatusNotFound), nil
	})
}

type updateTaskByIDResponses struct {
	NoContent  r.Nil       `status:"204"`
	BadRequest m.HttpError `status:"400"`
	NotFound   m.HttpError `status:"404"`
}
