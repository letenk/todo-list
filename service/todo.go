package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/repository"
)

type TodoService interface {
	Create(req web.TodoCreateRequest) (domain.Todo, error)
	GetAll(activityGroupID uint64) ([]domain.Todo, error)
	GetOne(id uint64) (domain.Todo, error)
	Update(id uint64, req web.TodoUpdateRequest) (domain.Todo, error)
	Delete(id uint64) (bool, error)
}

type todoService struct {
	repository repository.TodoRepository
}

func NewServiceTodo(repository repository.TodoRepository) *todoService {
	return &todoService{repository}
}

func (s *todoService) Create(req web.TodoCreateRequest) (domain.Todo, error) {
	todo := domain.Todo{
		ActivityGroupID: req.ActivityGroupId,
		Title:           req.Title,
	}

	newTodo, err := s.repository.Save(todo)
	if err != nil {
		return newTodo, err
	}

	return newTodo, err
}

func (s *todoService) GetAll(activityGroupID uint64) ([]domain.Todo, error) {
	if activityGroupID != 0 {
		// Find by activity group id
		todos, err := s.repository.FindByActivityGroupID(activityGroupID)
		if err != nil {
			return todos, err
		}
		return todos, nil
	}

	// Find all
	todos, err := s.repository.FindAll()

	if err != nil {
		return todos, err
	}

	return todos, nil
}

func (s *todoService) GetOne(id uint64) (domain.Todo, error) {
	// Find all
	todo, err := s.repository.FindOne(id)

	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (s *todoService) Update(id uint64, req web.TodoUpdateRequest) (domain.Todo, error) {
	// Find all
	todo, err := s.repository.FindOne(id)
	// If activity group not found
	if todo.ID == 0 {
		message := fmt.Sprintf("Todo with ID %d Not Found", id)
		return todo, errors.New(message)
	}

	if err != nil {
		return todo, err
	}

	// Change field title
	if req.Title != "" {
		todo.Title = req.Title
	}

	// Change field is active, if value req.IsActive is false
	if req.IsActive {
		todo.IsActive = true
	} else {
		todo.IsActive = req.IsActive

	}

	todo.UpdatedAt = time.Now()

	// Update
	updatedTodo, err := s.repository.Update(todo)
	if err != nil {
		return updatedTodo, err
	}

	return updatedTodo, nil
}

func (s *todoService) Delete(id uint64) (bool, error) {
	// Find one
	todo, err := s.repository.FindOne(id)
	// If activity group not found
	if todo.ID == 0 {
		message := fmt.Sprintf("Todo with ID %d Not Found", id)
		return false, errors.New(message)
	}

	if err != nil {
		return false, err
	}

	ok, err := s.repository.Delete(todo)
	if err != nil {
		return false, err
	}

	return ok, nil
}
