package models

import "github.com/google/uuid"

type Task struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
type Identifiable struct {
	ID uuid.UUID `json:"id"`
}
type IdentifiableTask struct {
	Identifiable
	Task
}

type page[T any] struct {
	Results  []T  `json:"results"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
	IsLast   bool `json:"isLast"`
}

type IdentifiableTasksPage page[IdentifiableTask]
