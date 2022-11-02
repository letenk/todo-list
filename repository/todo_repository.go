package repository

import (
	"github.com/letenk/todo-list/models/domain"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Save(todo domain.Todo) (domain.Todo, error)
	FindAll() ([]domain.Todo, error)
	FindByActivityGroupID(activityGroupID uint64) ([]domain.Todo, error)
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

func (r *todoRepository) FindAll() ([]domain.Todo, error) {
	var todos []domain.Todo

	err := r.db.Find(&todos).Error
	if err != nil {
		return todos, nil
	}

	return todos, nil
}

func (r *todoRepository) FindByActivityGroupID(activityGroupID uint64) ([]domain.Todo, error) {
	var todos []domain.Todo

	err := r.db.Where("activity_group_id = ?", activityGroupID).Find(&todos).Error
	if err != nil {
		return todos, nil
	}

	return todos, nil
}

func (r *todoRepository) FindOne(id uint64) (domain.Todo, error) {
	var todo domain.Todo

	err := r.db.Where("id = ?", id).Find(&todo).Error
	if err != nil {
		return todo, nil
	}

	return todo, nil
}

func (r *todoRepository) Save(todo domain.Todo) (domain.Todo, error) {
	err := r.db.Create(&todo).Error
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (r *todoRepository) Update(todo domain.Todo) (domain.Todo, error) {
	err := r.db.Save(&todo).Error
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (r *todoRepository) Delete(todo domain.Todo) (bool, error) {
	err := r.db.Delete(&todo).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
