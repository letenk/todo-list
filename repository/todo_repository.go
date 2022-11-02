package repository

import (
	"github.com/letenk/todo-list/models/domain"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Save(todo domain.Todo) (domain.Todo, error)
	FindAll() ([]domain.Todo, error)
	FindOne(id uint64) (domain.Todo, error)
	Update(todo domain.Todo) (domain.Todo, error)
	Delete(todo domain.Todo) (bool, error)
}

type todoRepository struct {
	db *gorm.DB
}

func NewRepositoryTodo(db *gorm.DB) *todoRepository {
	return &todoRepository{db}
}

func (r *todoRepository) Save(todo domain.Todo) (domain.Todo, error) {
	err := r.db.Create(&todo).Error
	if err != nil {
		return todo, err
	}

	return todo, nil
}
