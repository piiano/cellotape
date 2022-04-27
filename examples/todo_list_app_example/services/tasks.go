package services

import (
	"github.com/google/uuid"
	"github.com/piiano/restcontroller/examples/todo_list_app_example/models"
)

// This is the service layer.
// The signatures of the functions here is pure to represent their purpose and is decoupled from any REST characteristics.
// These services are exposed as REST APIs only by their usage in a REST controller.
// They can be potentially be exposed in additional ways with other protocols, CLI, Golang SDK, etc.

type TasksService interface {
	GetTasksPage(page int, pageSize int) models.Page[models.IdentifiableTask]
	CreateTask(task models.Task) uuid.UUID
	DeleteTaskByID(id uuid.UUID) bool
	GetTaskByID(id uuid.UUID) (models.Task, bool)
}

type tasks map[uuid.UUID]models.Task

func NewTasksService() TasksService {
	return &tasks{}
}

func (t *tasks) GetTasksPage(page int, pageSize int) models.Page[models.IdentifiableTask] {
	tasksWithId := make([]models.IdentifiableTask, 0, page)
	for id, task := range *t {
		tasksWithId = append(tasksWithId, models.IdentifiableTask{
			Identifiable: models.Identifiable{ID: id},
			Task:         task,
		})
	}
	last := len(tasksWithId) - 1
	from := page * pageSize
	to := (page + 1) * pageSize
	isLast := last <= to
	if isLast {
		to = last
	}
	pageSlice := tasksWithId[from:to]
	return models.Page[models.IdentifiableTask]{
		Results:  pageSlice,
		Page:     page,
		PageSize: pageSize,
		IsLast:   isLast,
	}
}

func (t *tasks) CreateTask(task models.Task) uuid.UUID {
	id := uuid.New()
	(*t)[id] = task
	return id
}

func (t *tasks) DeleteTaskByID(id uuid.UUID) bool {
	_, ok := (*t)[id]
	if ok {
		delete(*t, id)
	}
	return ok
}

func (t *tasks) GetTaskByID(id uuid.UUID) (models.Task, bool) {
	task, ok := (*t)[id]
	return task, ok
}
