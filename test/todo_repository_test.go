package test

import (
	"sync"
	"testing"
	"time"

	"github.com/letenk/todo-list/helper"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/repository"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestFindAllTodoRepository(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomTodoRepository(t)
			mutex.Unlock()
		}()
	}

	todoRepository := repository.NewRepositoryTodo(ConnTest)

	// Find all
	todos, err := todoRepository.FindAll()
	helper.ErrLogPanic(err)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(todos))

	for _, data := range todos {
		require.NotEmpty(t, data.ID)
		require.NotEmpty(t, data.Title)
		require.NotEmpty(t, data.ActivityGroupID)
		require.NotEmpty(t, data.Priority)
		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)

		require.NotNil(t, data.IsActive)

		require.Empty(t, data.DeletedAt)
	}

}

func TestFindByActivityGroupTodoRepository(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	// todos for store `new data todos` from channel
	var todos []domain.Todo

	// channel for store data `new data todos` from process create random todo
	channel := make(chan domain.Todo)
	defer close(channel)

	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			newTodos := createRandomTodoRepository(t)
			channel <- newTodos
			mutex.Unlock()
		}()
		todos = append(todos, <-channel)
	}

	todoRepository := repository.NewRepositoryTodo(ConnTest)

	// Find by actiivity group
	todos, err := todoRepository.FindByActivityGroupID(todos[0].ActivityGroupID)
	helper.ErrLogPanic(err)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))

	for _, data := range todos {
		// Equal
		require.Equal(t, todos[0].ID, data.ID)
		require.Equal(t, todos[0].Title, data.Title)
		require.Equal(t, todos[0].ActivityGroupID, data.ActivityGroupID)
		require.Equal(t, todos[0].IsActive, data.IsActive)
		require.Equal(t, todos[0].Priority, data.Priority)

		require.NotEmpty(t, data.CreatedAt)
		require.NotEmpty(t, data.UpdatedAt)
		require.Empty(t, data.DeletedAt)
	}

}

func TestFindOneTodoRepository(t *testing.T) {
	t.Parallel()
	newTodo := createRandomTodoRepository(t)

	todoRepository := repository.NewRepositoryTodo(ConnTest)

	// Find One
	todo, err := todoRepository.FindOne(newTodo.ID)
	helper.ErrLogPanic(err)

	assert.NoError(t, err)

	// Equal
	require.Equal(t, todo.ID, newTodo.ID)
	require.Equal(t, todo.Title, newTodo.Title)
	require.Equal(t, todo.ActivityGroupID, newTodo.ActivityGroupID)
	require.Equal(t, todo.IsActive, newTodo.IsActive)
	require.Equal(t, todo.Priority, newTodo.Priority)

	require.NotEmpty(t, newTodo.CreatedAt)
	require.NotEmpty(t, newTodo.UpdatedAt)
	require.Empty(t, newTodo.DeletedAt)

}

func TestCreateTodoRepository(t *testing.T) {
	t.Parallel()
	createRandomTodoRepository(t)
}

func TestUpdateTodoRepository(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	// todos for store `new data todos` from channel
	var todos []domain.Todo

	// channel for store data `new data todos` from process create random todo
	channel := make(chan domain.Todo)
	defer close(channel)

	// Create some random data
	for i := 0; i < 2; i++ {
		go func() {
			mutex.Lock()
			newTodos := createRandomTodoRepository(t)
			channel <- newTodos
			mutex.Unlock()
		}()
		todos = append(todos, <-channel)
	}

	todoRepository := repository.NewRepositoryTodo(ConnTest)

	dataUpdate := domain.Todo{
		ID:              todos[0].ID,
		ActivityGroupID: todos[1].ActivityGroupID,
		Title:           jabufaker.RandomString(20),
		IsActive:        false,
		Priority:        helper.RandomPriority(),
		CreatedAt:       todos[0].CreatedAt,
		UpdatedAt:       time.Now(),
		DeletedAt:       nil,
	}

	// Update
	todo, err := todoRepository.Update(dataUpdate)
	helper.ErrLogPanic(err)

	assert.NoError(t, err)

	// Test
	require.Equal(t, todo.ID, todos[0].ID)

	require.NotEqual(t, todo.Title, todos[0].Title)
	require.NotEqual(t, todo.ActivityGroupID, todos[0].ActivityGroupID)
	require.NotEqual(t, todo.IsActive, todos[0].IsActive)
	require.NotEqual(t, todo.Priority, todos[0].Priority)

	require.NotEmpty(t, todo.CreatedAt)
	require.NotEmpty(t, todo.UpdatedAt)

	require.Empty(t, todos[0].DeletedAt)
}

func TestDeleteTodoRepository(t *testing.T) {
	t.Parallel()
	newTodo := createRandomTodoRepository(t)

	todoRepository := repository.NewRepositoryTodo(ConnTest)

	// Update
	ok, err := todoRepository.Delete(newTodo)
	helper.ErrLogPanic(err)

	assert.NoError(t, err)
	assert.True(t, ok)

	todo, err := todoRepository.FindOne(newTodo.ID)
	helper.ErrLogPanic(err)
	nullId := uint64(0)
	assert.Equal(t, nullId, todo.ID)
}
