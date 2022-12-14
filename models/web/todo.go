package web

import (
	"time"

	"github.com/letenk/todo-list/models/domain"
)

type TodoURI struct {
	ID uint64 `uri:"id" binding:"required"`
}

type TodoCreateRequest struct {
	ActivityGroupID uint64 `json:"activity_group_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
}

type TodoUpdateRequest struct {
	Title    string `json:"title,omitempty"`
	IsActive bool   `json:"is_active"`
}
type TodoResponse struct {
	ID         uint64     `json:"id"`
	Title      string     `json:"title"`
	ActivityID uint64     `json:"activity_group_id"`
	IsActive   string     `json:"is_active"`
	Priority   string     `json:"priority"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletetAt  *time.Time `json:"deleted_at"`
}

type TodoCreatedResponse struct {
	ID         uint64     `json:"id"`
	Title      string     `json:"title"`
	ActivityID uint64     `json:"activity_group_id"`
	IsActive   bool       `json:"is_active"`
	Priority   string     `json:"priority"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletetAt  *time.Time `json:"deleted_at"`
}

// Format for handle single response todo
func FormatTodo(todo domain.Todo) TodoResponse {
	var isActive string

	// If isActive is false
	if !todo.IsActive {
		isActive = "0"
	} else {
		isActive = "1"
	}

	formatter := TodoResponse{
		ID:         todo.ID,
		Title:      todo.Title,
		ActivityID: todo.ActivityGroupID,
		IsActive:   isActive,
		Priority:   todo.Priority,
		CreatedAt:  todo.CreatedAt,
		UpdatedAt:  todo.UpdatedAt,
		DeletetAt:  todo.DeletedAt,
	}
	return formatter
}

// Format for handle single response todo
func FormatCreatedTodo(todo domain.Todo) TodoCreatedResponse {
	formatter := TodoCreatedResponse{
		ID:         todo.ID,
		Title:      todo.Title,
		ActivityID: todo.ActivityGroupID,
		IsActive:   todo.IsActive,
		Priority:   todo.Priority,
		CreatedAt:  todo.CreatedAt,
		UpdatedAt:  todo.UpdatedAt,
		DeletetAt:  todo.DeletedAt,
	}
	return formatter
}

// Format for handle single response todo
func FormatTodoResponse(todo domain.Todo) TodoResponse {
	var isActive string
	// If isActive is false
	if !todo.IsActive {
		isActive = "0"
	} else {
		isActive = "1"
	}

	formatter := TodoResponse{
		ID:         todo.ID,
		Title:      todo.Title,
		ActivityID: todo.ActivityGroupID,
		IsActive:   isActive,
		Priority:   todo.Priority,
		CreatedAt:  todo.CreatedAt,
		UpdatedAt:  todo.UpdatedAt,
		DeletetAt:  todo.DeletedAt,
	}
	return formatter
}

// Format for handle multiples response todo
func FormatTodos(todo []domain.Todo) []TodoResponse {
	if len(todo) == 0 {
		return []TodoResponse{}
	}

	var formatters []TodoResponse

	for _, data := range todo {
		formatter := FormatTodo(data)
		formatters = append(formatters, formatter)
	}

	return formatters
}
