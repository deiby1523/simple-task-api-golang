package service

import (
	"errors"
	"simple-task-api-golang/internal/model"
	"simple-task-api-golang/internal/store"
)

// Service represents the business logic layer of the application.
// It acts as an intermediary between handlers (HTTP/API layer)
// and the data persistence layer (store).
type Service struct {
	store store.Store
}

// New creates a new Service instance.
// It receives a Store implementation via dependency injection
// to keep the service layer decoupled from the persistence layer.
func New(s store.Store) *Service {
	return &Service{store: s}
}

// GetTasks retrieves all tasks from the store.
func (s *Service) GetTasks() ([]*model.Task, error) {
	return s.store.GetAll()
}

// GetTaskByID retrieves a single task by its ID.
func (s *Service) GetTaskByID(id int) (*model.Task, error) {
	return s.store.GetByID(id)
}

// CreateTask validates the task data and creates a new task.
// It ensures required business rules are met before persistence.
func (s *Service) CreateTask(task model.Task) (*model.Task, error) {
	// Basic validation: title is required
	if task.Title == "" {
		return nil, errors.New("task title is required")
	}

	return s.store.Create(&task)
}

// UpdateTask validates the task data and updates an existing task
// identified by its ID.
func (s *Service) UpdateTask(id int, task model.Task) (*model.Task, error) {
	// Basic validation: title is required
	if task.Title == "" {
		return nil, errors.New("task title is required")
	}

	return s.store.Update(id, &task)
}

// DeleteTask removes a task by its ID.
func (s *Service) DeleteTask(id int) error {
	return s.store.Delete(id)
}
