package models

import (
	_ "embed"
	"errors"
	"github.com/google/uuid"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/models"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/services"
	"github.com/piiano/restcontroller/router"
)

func TasksOperationsGroup(taskService services.TasksService) router.Group {
	getAllTasksOperationHandler := router.NewOperationHandler(func(request router.Request[router.Nil, router.Nil, PaginationQueryParams]) (models.Page[models.IdentifiableTask], error) {
		page := taskService.GetTasksPage(request.QueryParameters.Page, request.QueryParameters.PageSize)
		return page, nil
	})
	createTaskOperationHandler := router.NewOperationHandler(func(request router.Request[models.Task, router.Nil, router.Nil]) (models.Identifiable, error) {
		id := taskService.CreateTask(request.Body)
		return models.Identifiable{ID: id}, nil
	})
	getTaskByIDOperationHandler := router.NewOperationHandler(func(request router.Request[router.Nil, IDPathParam, router.Nil]) (models.Task, error) {
		if task, found := taskService.GetTaskByID(request.PathParameters.ID); found {
			return task, nil
		}
		return models.Task{}, errors.New("not found") // TODO: return HTTP error
	})
	deleteTaskByIDOperationHandler := router.NewOperationHandler(func(request router.Request[router.Nil, IDPathParam, router.Nil]) (router.Nil, error) {
		if deleted := taskService.DeleteTaskByID(request.PathParameters.ID); deleted {
			return 0, nil
		}
		return 0, errors.New("gone") // TODO: return HTTP error
	})
	return router.NewGroup().
		WithOperation("getTasksPage", getAllTasksOperationHandler).
		WithOperation("createNewTask", createTaskOperationHandler).
		WithOperation("getTaskByID", getTaskByIDOperationHandler).
		WithOperation("deleteTaskByID", deleteTaskByIDOperationHandler)
}

type PaginationQueryParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type IDPathParam struct {
	ID uuid.UUID `uri:"id"`
}
