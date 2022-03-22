package rest

import (
	m "github.com/piiano/cellotape/examples/todo_list_app_example/models"
	"github.com/piiano/cellotape/examples/todo_list_app_example/services"
	r "github.com/piiano/cellotape/router"
	"net/http"
)

func getTasksPageOperation(tasks services.TasksService) r.Handler {
	return r.NewHandler(func(_ r.Context, request r.Request[r.Nil, r.Nil, paginationQueryParams]) (r.Response[getTasksPageResponses], error) {
		tasksPage := tasks.GetTasksPage(request.QueryParams.Page, request.QueryParams.PageSize)
		return r.Send(200, getTasksPageResponses{OK: tasksPage}, http.Header{"Cache-Control": {"max-age=10"}})
	})
}

type paginationQueryParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}
type getTasksPageResponses struct {
	OK m.IdentifiableTasksPage `status:"200"`
}
