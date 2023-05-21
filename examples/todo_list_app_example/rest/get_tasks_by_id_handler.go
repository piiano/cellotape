package rest

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	m "github.com/piiano/cellotape/examples/todo_list_app_example/models"
	"github.com/piiano/cellotape/examples/todo_list_app_example/services"
	r "github.com/piiano/cellotape/router"
	"github.com/piiano/cellotape/router/utils"
)

func getTaskByIDOperation(tasks services.TasksService) r.Handler {
	return r.NewHandler(func(_ *r.Context, request r.Request[utils.Nil, idPathParam, utils.Nil]) (r.Response[getTaskByIDResponses], error) {
		id, err := uuid.Parse(request.PathParams.ID)
		if err != nil {
			return r.SendJSON(getTaskByIDResponses{
				BadRequest: m.HttpError{
					Error:  "bad request",
					Reason: err.Error(),
				},
			}).Status(http.StatusBadRequest), nil
		}
		if task, found := tasks.GetTaskByID(id); found {
			return r.SendOKJSON(getTaskByIDResponses{OK: task}), nil
		}
		return r.SendJSON(getTaskByIDResponses{
			NotFound: m.HttpError{
				Error:  "not found",
				Reason: fmt.Sprintf("task with id %q is not found", request.PathParams.ID),
			},
		}).Status(http.StatusNotFound), nil
	})
}

type getTaskByIDResponses struct {
	OK         m.Task      `status:"200"`
	BadRequest m.HttpError `status:"400"`
	NotFound   m.HttpError `status:"404"`
}
