package services

import (
	"github.com/google/uuid"
	"github.com/piiano/cellotape/examples/todo_list_app_example/models"
)

// This is the service layer.
// The signatures of the functions here is pure to represent their purpose and is decoupled from any REST characteristics.
// These services are exposed as REST APIs only by their usage in a REST controller.
// They can be potentially be exposed in additional ways with other protocols, CLI, Golang SDK, etc.

type TasksService interface {
	GetTasksPage(page int, pageSize int) models.IdentifiableTasksPage
	CreateTask(models.Task) uuid.UUID
	GetTaskByID(uuid.UUID) (models.Task, bool)
	UpdateTaskByID(uuid.UUID, models.Task) bool
	DeleteTaskByID(uuid.UUID) bool
}

type tasks map[uuid.UUID]models.Task

func NewTasksService() TasksService {
	return &tasks{}
}

func (t *tasks) GetTasksPage(page int, pageSize int) models.IdentifiableTasksPage {
	tasksWithIds := make([]models.IdentifiableTask, 0, page)
	for id, task := range *t {
		tasksWithIds = append(tasksWithIds, models.IdentifiableTask{
			Identifiable: models.Identifiable{ID: id},
			Task:         task,
		})
	}
	from := page * pageSize
	to := (page + 1) * pageSize
	isLast := to >= len(tasksWithIds)
	if isLast {
		to = len(tasksWithIds)
	}
	pageSlice := tasksWithIds[from:to]
	return models.IdentifiableTasksPage{
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

func (t *tasks) UpdateTaskByID(id uuid.UUID, task models.Task) bool {
	_, ok := (*t)[id]
	if ok {
		(*t)[id] = task
	}
	return ok
}
