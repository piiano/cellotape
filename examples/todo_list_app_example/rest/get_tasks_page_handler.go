package rest

import (
	"net/http"

	m "github.com/piiano/cellotape/examples/todo_list_app_example/models"
	"github.com/piiano/cellotape/examples/todo_list_app_example/services"
	r "github.com/piiano/cellotape/router"
	"github.com/piiano/cellotape/router/utils"
)

func getTasksPageOperation(tasks services.TasksService) r.Handler {
	return r.NewHandler(func(_ *r.Context, request r.Request[utils.Nil, utils.Nil, paginationQueryParams]) (r.Response[getTasksPageResponses], error) {
		tasksPage := tasks.GetTasksPage(request.QueryParams.Page, request.QueryParams.PageSize)
		return r.SendOKJSON(getTasksPageResponses{OK: tasksPage}, http.Header{"Cache-Control": {"max-age=10"}}), nil
	})
}

type paginationQueryParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}
type getTasksPageResponses struct {
	OK m.IdentifiableTasksPage `status:"200"`
}
