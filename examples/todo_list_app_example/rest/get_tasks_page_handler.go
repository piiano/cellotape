package models

import (
	m "github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

func getTasksPageOperation(tasks services.TasksService) r.OperationHandler {
	return r.OperationFunc(func(request r.Request[r.Nil, r.Nil, paginationQueryParams]) (int, getTasksPageResponses) {
		tasksPage := tasks.GetTasksPage(request.QueryParams.Page, request.QueryParams.PageSize)
		return 200, getTasksPageResponses{OK: tasksPage}
	})
}

type paginationQueryParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}
type getTasksPageResponses struct {
	OK m.IdentifiableTasksPage `status:"200"`
}
