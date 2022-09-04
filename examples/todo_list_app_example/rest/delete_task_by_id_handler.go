package rest

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	m "github.com/piiano/cellotape/examples/todo_list_app_example/models"
	"github.com/piiano/cellotape/examples/todo_list_app_example/services"
	r "github.com/piiano/cellotape/router"
)

func deleteTaskByIDOperation(tasks services.TasksService) r.Handler {
	return r.NewHandler(func(_ r.Context, request r.Request[r.Nil, idPathParam, r.Nil]) (r.Response[deleteTaskByIDResponses], error) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			return r.SendJSON(deleteTaskByIDResponses{
				BadRequest: m.HttpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			}).Status(http.StatusBadRequest), nil
		}
		if deleted := tasks.DeleteTaskByID(id); deleted {
			return r.SendJSON(deleteTaskByIDResponses{}).Status(http.StatusNoContent), nil
		}
		return r.SendJSON(deleteTaskByIDResponses{
			Gone: m.HttpError{
				Error:  "gone",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		}).Status(http.StatusGone), nil
	})
}

type deleteTaskByIDResponses struct {
	NoContent  r.Nil       `status:"204"`
	BadRequest m.HttpError `status:"400"`
	Gone       m.HttpError `status:"410"`
}
