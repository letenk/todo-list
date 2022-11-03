package test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/letenk/todo-list/helper"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/repository"
	"github.com/letenk/todo-list/service"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/require"
)

func createRandomTodoService(t *testing.T) domain.Todo {
	repository := repository.NewRepositoryTodo(ConnTest)
	service := service.NewServiceTodo(repository)

	newActivity := createRandomActivityRepository(t)

	data := web.TodoCreateRequest{
		ActivityGroupID: newActivity.ID,
		Title:           jabufaker.RandomString(20),
	}

	// Create
	newTodo, err := service.Create(data)
	helper.ErrLogPanic(err)

	// Test
	require.NoError(t, err)

	require.Equal(t, data.ActivityGroupID, newTodo.ActivityGroupID)
	require.Equal(t, data.Title, newTodo.Title)
	require.Equal(t, "very-high", newTodo.Priority)

	require.True(t, newTodo.IsActive)

	require.NotEmpty(t, newTodo.ID)
	require.NotEmpty(t, newTodo.CreatedAt)
	require.NotEmpty(t, newTodo.UpdatedAt)

	require.Empty(t, newTodo.DeletedAt)

	return newTodo
}

func TestCreateTodoService(t *testing.T) {
	t.Parallel()
	createRandomTodoService(t)
}

func TestGetAllTodoServices(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	var newTodos []domain.Todo

	// Create channel for store new todos created
	channel := make(chan domain.Todo)
	defer close(channel)
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			newTodo := createRandomTodoService(t)
			channel <- newTodo
			mutex.Unlock()
		}()
		newTodos = append(newTodos, <-channel)
	}

	repository := repository.NewRepositoryTodo(ConnTest)
	service := service.NewServiceTodo(repository)

	t.Run("Get all todos without query activity_group_id", func(t *testing.T) {
		// Get activity groups
		ActivityId := int64(0)
		todos, err := service.GetAll(uint64(ActivityId))
		helper.ErrLogPanic(err)

		// Length todos must be greater than 0
		require.NotEqual(t, 0, len((todos)))

		for _, data := range todos {
			require.NotEmpty(t, data.ID)
			require.NotEmpty(t, data.Title)
			require.NotEmpty(t, data.ActivityGroupID)

			require.NotNil(t, data.IsActive)

			require.NotEmpty(t, data.Priority)
			require.NotEmpty(t, data.CreatedAt)
			require.NotEmpty(t, data.UpdatedAt)

			require.Nil(t, data.DeletedAt)
		}
	})

	t.Run("Get all todos with query activity_group_id", func(t *testing.T) {
		// Get activity groups
		ActivityId := newTodos[0].ActivityGroupID
		todos, err := service.GetAll(uint64(ActivityId))
		helper.ErrLogPanic(err)

		// Length todos must be 1
		require.Equal(t, 1, len((todos)))

		for _, data := range todos {
			require.Equal(t, newTodos[0].ID, data.ID)
			require.Equal(t, newTodos[0].Title, data.Title)
			require.Equal(t, newTodos[0].ActivityGroupID, data.ActivityGroupID)
			require.Equal(t, newTodos[0].IsActive, data.IsActive)
			require.Equal(t, newTodos[0].Priority, data.Priority)

			require.Equal(t, todos[0].CreatedAt, data.CreatedAt)
			require.Equal(t, todos[0].UpdatedAt, data.UpdatedAt)

			require.Nil(t, data.DeletedAt)
		}
	})

}

func TestGetOneTodoServices(t *testing.T) {
	t.Parallel()
	newTodo := createRandomTodoService(t)

	repository := repository.NewRepositoryTodo(ConnTest)
	service := service.NewServiceTodo(repository)

	// Get activity groups
	todo, err := service.GetOne(newTodo.ID)
	helper.ErrLogPanic(err)

	require.Equal(t, newTodo.ID, todo.ID)
	require.Equal(t, newTodo.Title, todo.Title)
	require.Equal(t, newTodo.ActivityGroupID, todo.ActivityGroupID)
	require.Equal(t, newTodo.IsActive, todo.IsActive)
	require.Equal(t, newTodo.Priority, todo.Priority)

	require.NotEmpty(t, todo.CreatedAt)
	require.NotEmpty(t, todo.UpdatedAt)
	require.Empty(t, todo.DeletedAt)
}

func TestUpdateTodoService(t *testing.T) {
	t.Parallel()

	repository := repository.NewRepositoryTodo(ConnTest)
	service := service.NewServiceTodo(repository)

	t.Run("Update success", func(t *testing.T) {
		// Create random data
		newTodo := createRandomTodoService(t)
		dataUpdated := web.TodoUpdateRequest{
			Title:    jabufaker.RandomString(20),
			IsActive: false,
		}

		updatedTodo, err := service.Update(newTodo.ID, dataUpdated)
		helper.ErrLogPanic(err)

		require.Equal(t, newTodo.ID, updatedTodo.ID)
		require.Equal(t, newTodo.ActivityGroupID, updatedTodo.ActivityGroupID)
		require.Equal(t, newTodo.CreatedAt, updatedTodo.CreatedAt)

		require.NotEqual(t, newTodo.Title, updatedTodo.Title)
		require.NotEqual(t, newTodo.IsActive, updatedTodo.IsActive)

		require.NotEqual(t, newTodo.UpdatedAt, updatedTodo.UpdatedAt)

		require.Nil(t, updatedTodo.DeletedAt)

	})

	t.Run("Update success without field is_active", func(t *testing.T) {
		// Create random data
		newTodo := createRandomTodoService(t)
		dataUpdated := web.TodoUpdateRequest{
			Title:    jabufaker.RandomString(20),
			IsActive: true, // this sample and change type do it in handler, when checking field is false or true do in handler
		}

		updatedTodo, err := service.Update(newTodo.ID, dataUpdated)
		helper.ErrLogPanic(err)

		require.Equal(t, newTodo.ID, updatedTodo.ID)
		require.Equal(t, newTodo.ActivityGroupID, updatedTodo.ActivityGroupID)
		require.Equal(t, newTodo.CreatedAt, updatedTodo.CreatedAt)
		require.Equal(t, newTodo.IsActive, updatedTodo.IsActive)

		require.NotEqual(t, newTodo.Title, updatedTodo.Title)

		require.NotEqual(t, newTodo.UpdatedAt, updatedTodo.UpdatedAt)

		require.Nil(t, updatedTodo.DeletedAt)
	})

	t.Run("Update success without field title", func(t *testing.T) {
		// Create random data
		newTodo := createRandomTodoService(t)
		dataUpdated := web.TodoUpdateRequest{
			IsActive: false,
		}

		updatedTodo, err := service.Update(newTodo.ID, dataUpdated)
		helper.ErrLogPanic(err)

		require.Equal(t, newTodo.ID, updatedTodo.ID)
		require.Equal(t, newTodo.ActivityGroupID, updatedTodo.ActivityGroupID)
		require.Equal(t, newTodo.CreatedAt, updatedTodo.CreatedAt)
		require.Equal(t, newTodo.Title, updatedTodo.Title)

		require.NotEqual(t, newTodo.IsActive, updatedTodo.IsActive)

		require.NotEqual(t, newTodo.UpdatedAt, updatedTodo.UpdatedAt)

		require.Nil(t, updatedTodo.DeletedAt)
	})

	t.Run("Update failed todo not found", func(t *testing.T) {
		dataUpdated := web.TodoUpdateRequest{
			Title:    jabufaker.RandomString(20),
			IsActive: false,
		}

		_, err := service.Update(7329323, dataUpdated)
		require.Error(t, err)

		message := fmt.Sprintf("Todo with ID %d Not Found", 7329323)
		require.Equal(t, message, err.Error())

	})
}

func TestDeleteTodoService(t *testing.T) {
	t.Parallel()
	// Create random data
	newTodo := createRandomTodoService(t)

	repository := repository.NewRepositoryTodo(ConnTest)
	service := service.NewServiceTodo(repository)

	t.Run("Delete success", func(t *testing.T) {

		ok, err := service.Delete(newTodo.ID)
		helper.ErrLogPanic(err)

		require.True(t, ok)
	})

	t.Run("Delete failed todo not found", func(t *testing.T) {
		ok, err := service.Delete(7329323)
		require.Error(t, err)
		require.False(t, ok)

		message := fmt.Sprintf("Todo with ID %d Not Found", 7329323)
		require.Equal(t, message, err.Error())

	})
}
