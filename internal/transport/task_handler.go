package transport

import (
	"encoding/json"
	"net/http"
	"simple-task-api-golang/internal/model"
	"simple-task-api-golang/internal/service"
	"strconv"
	"strings"
)

// TaskHandler handles HTTP requests related to task resources.
// It acts as the transport layer between HTTP clients and the service layer.
type TaskHandler struct {
	service *service.Service
}

// New creates a new TaskHandler instance.
// It injects the service dependency to keep the handler decoupled
// from the business logic implementation.
func New(s *service.Service) *TaskHandler {
	return &TaskHandler{service: s}
}

// HandleTasks handles requests to the /tasks endpoint.
// Supported methods:
//   - GET: retrieves all tasks
//   - POST: creates a new task
func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		// Retrieve all tasks
		tasks, err := h.service.GetTasks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		// Decode request body into Task model
		var task model.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create task using the service layer
		created, err := h.service.CreateTask(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(created)

	default:
		// Method not allowed
		http.Error(w, "Method is not available", http.StatusMethodNotAllowed)
	}
}

// HandleTaskById handles requests to the /task/{id} endpoint.
// Supported methods:
//   - GET: retrieves a task by ID
//   - PUT: updates a task by ID
//   - DELETE: deletes a task by ID
func (h *TaskHandler) HandleTaskById(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from URL path
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Not found", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		// Retrieve task by ID
		task, err := h.service.GetTaskByID(id)
		if err != nil {
			http.Error(w, "Not found", http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)

	case http.MethodPut:
		// Decode request body into Task model
		var task model.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Update task using the service layer
		updated, err := h.service.UpdateTask(id, task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)

	case http.MethodDelete:
		// Delete task by ID
		if err := h.service.DeleteTask(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		// Method not allowed
		http.Error(w, "Method is not available", http.StatusMethodNotAllowed)
	}
}
