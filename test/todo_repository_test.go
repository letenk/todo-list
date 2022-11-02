package test

import (
	"testing"

	"github.com/letenk/todo-list/helper"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/repository"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
)

func createRandomTodoRepository(t *testing.T) domain.Todo {
	todoRepository := repository.NewRepositoryTodo(ConnTest)

	newActivityGroup := createRandomActivityGroupRepository(t)

	todo := domain.Todo{
		ActivityGroupID: newActivityGroup.ID,
		Title:           jabufaker.RandomString(20),
		IsActive:        true,
		Priority:        helper.RandomPriority(),
	}

	// Save to db
	newTodo, err := todoRepository.Save(todo)
	helper.ErrLogPanic(err)

	// Test pas
	assert.NoError(t, err)

	assert.NotEmpty(t, newTodo.ID)
	assert.NotEmpty(t, newActivityGroup.CreatedAt)
	assert.NotEmpty(t, newActivityGroup.UpdatedAt)
	assert.Empty(t, newActivityGroup.DeletedAt)

	assert.Equal(t, todo.ActivityGroupID, newTodo.ActivityGroupID)
	assert.Equal(t, todo.Title, newTodo.Title)
	assert.Equal(t, todo.IsActive, newTodo.IsActive)
	assert.Equal(t, todo.Priority, newTodo.Priority)

	return newTodo
}

func TestCreateTodoRepository(t *testing.T) {
	t.Parallel()
	createRandomTodoRepository(t)
}
