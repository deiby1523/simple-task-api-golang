package store

import (
	"database/sql"
	"simple-task-api-golang/internal/model"
)

// Store defines the contract for task persistence.
// Any implementation must provide basic CRUD operations for tasks.
type Store interface {
	GetAll() ([]*model.Task, error)
	GetByID(id int) (*model.Task, error)
	Create(task *model.Task) (*model.Task, error)
	Update(id int, task *model.Task) (*model.Task, error)
	Delete(id int) error
}

// store is the concrete implementation of Store.
// It holds a reference to the database connection.
type store struct {
	db *sql.DB
}

// New creates a new Store instance using the provided database connection.
// It returns the Store interface to keep the implementation decoupled.
func New(db *sql.DB) Store {
	return &store{db: db}
}

// GetAll retrieves all tasks from the database.
// It returns a slice of Task pointers or an error if the query fails.
func (s *store) GetAll() ([]*model.Task, error) {
	var q = `SELECT id, title, description, completed FROM tasks`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	// Ensure rows are closed after processing
	defer rows.Close()

	var tasks []*model.Task

	// Iterate over the result set and scan each row into a Task struct
	for rows.Next() {
		t := model.Task{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil
}

// GetByID retrieves a single task by its ID.
// If the task does not exist or an error occurs, an error is returned.
func (s *store) GetByID(id int) (*model.Task, error) {
	var q = `SELECT id, title, description, completed FROM tasks WHERE id = ?`

	t := model.Task{}

	err := s.db.QueryRow(q, id).Scan(&t.ID, &t.Title, &t.Description, &t.Completed)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Create inserts a new task into the database.
// After insertion, the generated ID is assigned back to the task.
func (s *store) Create(task *model.Task) (*model.Task, error) {
	var q = `INSERT INTO tasks (title, description, completed) VALUES (?, ?, ?)`

	resp, err := s.db.Exec(q, task.Title, task.Description, task.Completed)
	if err != nil {
		return nil, err
	}

	// Retrieve the auto-generated ID from the database
	id, err := resp.LastInsertId()
	if err != nil {
		return nil, err
	}

	task.ID = int(id)
	return task, nil
}

// Update modifies an existing task identified by its ID.
// It returns the updated task or an error if the operation fails.
func (s *store) Update(id int, task *model.Task) (*model.Task, error) {
	var q = `UPDATE tasks SET title = ?, description = ?, completed = ? WHERE id = ?`

	_, err := s.db.Exec(q, task.Title, task.Description, task.Completed, id)
	if err != nil {
		return nil, err
	}

	task.ID = id
	return task, nil
}

// Delete removes a task from the database by its ID.
// It returns an error if the delete operation fails.
func (s *store) Delete(id int) error {
	var q = `DELETE FROM tasks WHERE id = ?`

	_, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}
