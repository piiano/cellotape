package models

import (
	_ "embed"
	"errors"
	"github.com/google/uuid"
	m "github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	r "github.com/piiano/restcontroller/router"
)

func TasksOperationsGroup(tasks services.TasksService) r.Group {
	return r.NewGroup().
		WithOperation("getTasksPage",
			r.OperationFunc(func(request r.Request[r.Nil, r.Nil, PaginationQueryParams]) (m.IdentifiableTasksPage, error) {
				return tasks.GetTasksPage(request.QueryParams.Page, request.QueryParams.PageSize), nil
			})).
		WithOperation("createNewTask",
			r.OperationFunc(func(request r.Request[m.Task, r.Nil, r.Nil]) (m.Identifiable, error) {
				return m.Identifiable{ID: tasks.CreateTask(request.Body)}, nil
			})).
		WithOperation("getTaskByID",
			r.OperationFunc(func(request r.Request[r.Nil, IDPathParam, r.Nil]) (m.Task, error) {
				if task, found := tasks.GetTaskByID(uuid.MustParse(request.PathParams.ID)); found {
					return task, nil
				}
				return m.Task{}, errors.New("404 not found") // TODO: return HTTP error
			})).
		WithOperation("deleteTaskByID",
			r.OperationFunc(func(request r.Request[r.Nil, IDPathParam, r.Nil]) (r.Nil, error) {
				if deleted := tasks.DeleteTaskByID(uuid.MustParse(request.PathParams.ID)); deleted {
					return nil, nil
				}
				return nil, errors.New("410 gone") // TODO: return HTTP error
			})).
		WithOperation("updateTaskByID",
			r.OperationFunc(func(request r.Request[m.Task, IDPathParam, r.Nil]) (r.Nil, error) {
				if updated := tasks.UpdateTaskByID(uuid.MustParse(request.PathParams.ID), request.Body); updated {
					return nil, nil
				}
				return nil, errors.New("404 not found") // TODO: return HTTP error
			}))
}

type PaginationQueryParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type IDPathParam struct {
	// https://github.com/gin-gonic/gin/issues/2423
	ID string `uri:"id"`
}
