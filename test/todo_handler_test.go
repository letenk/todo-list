package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/letenk/todo-list/models/web"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomTodoHandler(t *testing.T) web.TodoCreateResponse {
	newActivityGroup := createRandomActivityGroupHandler(t)
	data := web.TodoCreateResponse{
		Title:           jabufaker.RandomString(20),
		ActivityGroupID: newActivityGroup.ID,
	}

	dataBody := fmt.Sprintf(`{"title": "%s", "activity_group_id": %d}`, data.Title, data.ActivityGroupID)
	requestBody := strings.NewReader(dataBody)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:3030/todo-items", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	Route.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	require.Equal(t, 201, response.StatusCode)
	require.Equal(t, "Success", responseBody["status"])
	require.Equal(t, "Success", responseBody["message"])

	require.NotEmpty(t, responseBody["data"])

	var contextData = responseBody["data"].(map[string]interface{})
	fmt.Println(contextData["priority"])
	require.NotEmpty(t, contextData["id"])
	require.NotEmpty(t, contextData["created_at"])
	require.NotEmpty(t, contextData["updated_at"])

	require.Equal(t, "very-high", contextData["priority"])
	require.Equal(t, data.Title, contextData["title"])
	require.Equal(t, data.ActivityGroupID, uint64(contextData["activity_group_id"].(float64)))

	require.Equal(t, "1", contextData["is_active"].(string))

	newtodo := web.TodoCreateResponse{
		ID:              uint64(contextData["id"].(float64)),
		Title:           contextData["title"].(string),
		ActivityGroupID: uint64(contextData["activity_group_id"].(float64)),
		IsActive:        contextData["is_active"].(string),
		Priority:        contextData["priority"].(string),
	}

	return newtodo
}

func TestCreateTodoHandler(t *testing.T) {
	t.Parallel()

	t.Run("create new todo success", func(t *testing.T) {
		createRandomTodoHandler(t)
	})

	t.Run("create new todo without title", func(t *testing.T) {
		newActivityGroup := createRandomActivityGroupHandler(t)
		data := web.TodoCreateResponse{
			ActivityGroupID: newActivityGroup.ID,
		}

		dataBody := fmt.Sprintf(`{"activity_group_id": "%d"}`, data.ActivityGroupID)
		requestBody := strings.NewReader(dataBody)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3030/todo-items", requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 400, response.StatusCode)
		require.Equal(t, "Bad Request", responseBody["status"])
		require.Equal(t, "title cannot be null", responseBody["message"])

		require.Empty(t, responseBody["data"])
	})

	t.Run("create new todo withoud activity group id", func(t *testing.T) {
		data := web.TodoCreateResponse{
			Title: jabufaker.RandomString(20),
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3030/todo-items", requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 400, response.StatusCode)
		require.Equal(t, "Bad Request", responseBody["status"])
		require.Equal(t, "activity_group_id cannot be null", responseBody["message"])

		require.Empty(t, responseBody["data"])
	})
}

func TestGetAllTodoHandler(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	var newTodos []web.TodoCreateResponse

	// Create channel for store new todos created
	channel := make(chan web.TodoCreateResponse)
	defer close(channel)
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			newTodo := createRandomTodoHandler(t)
			channel <- newTodo
			mutex.Unlock()
		}()
		newTodos = append(newTodos, <-channel)
	}

	t.Run("Get all todo without query activity group id", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:3030/todo-items", nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])

		require.NotEmpty(t, responseBody["data"])

		contextBody := responseBody["data"].([]interface{})
		// Length todos must be greater than 0
		require.NotEqual(t, 0, len(contextBody))

		for _, data := range contextBody {
			list := data.(map[string]interface{})
			require.NotEmpty(t, list["id"])
			require.NotEmpty(t, list["title"])
			require.NotEmpty(t, list["is_active"])
			require.NotEmpty(t, list["priority"])
			require.NotEmpty(t, list["created_at"])
			require.NotEmpty(t, list["updated_at"])
			require.Nil(t, list["deleted_at"])
		}
	})

	t.Run("Get all todo with query activity group id", func(t *testing.T) {
		activityGroupID := fmt.Sprintf("%d", newTodos[0].ActivityGroupID)
		request := httptest.NewRequest(http.MethodGet, "http://localhost:3030/todo-items?activity_group_id="+activityGroupID, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])

		require.NotEmpty(t, responseBody["data"])

		contextBody := responseBody["data"].([]interface{})
		// Length todos must be 1
		require.Equal(t, 1, len(contextBody))

		for _, data := range contextBody {
			list := data.(map[string]interface{})
			require.Equal(t, newTodos[0].ID, uint64(list["id"].(float64)))
			require.Equal(t, newTodos[0].Title, list["title"])
			require.Equal(t, newTodos[0].IsActive, list["is_active"])
			require.Equal(t, newTodos[0].Priority, list["priority"])

			require.NotEmpty(t, list["created_at"])
			require.NotEmpty(t, list["updated_at"])

			require.Nil(t, list["deleted_at"])
		}
	})

}

func TestGetOneTodoHandler(t *testing.T) {
	t.Parallel()
	var mutex sync.Mutex
	var newTodos []web.TodoCreateResponse

	// Create channel for store new todos created
	channel := make(chan web.TodoCreateResponse)
	defer close(channel)
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			newTodo := createRandomTodoHandler(t)
			channel <- newTodo
			mutex.Unlock()
		}()
		newTodos = append(newTodos, <-channel)
	}

	t.Run("Get one todo success", func(t *testing.T) {
		todoId := fmt.Sprintf("%d", newTodos[0].ID)
		request := httptest.NewRequest(http.MethodGet, "http://localhost:3030/todo-items/"+todoId, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])

		var contextData = responseBody["data"].(map[string]interface{})
		require.NotEmpty(t, contextData)

		require.Equal(t, newTodos[0].ID, uint64(contextData["id"].(float64)))
		require.Equal(t, newTodos[0].Title, contextData["title"])
		require.Equal(t, newTodos[0].IsActive, contextData["is_active"])
		require.Equal(t, newTodos[0].Priority, contextData["priority"])

		require.NotEmpty(t, contextData["created_at"])
		require.NotEmpty(t, contextData["updated_at"])

		require.Nil(t, contextData["deleted_at"])
	})

	t.Run("Get one ID not found", func(t *testing.T) {
		wrongID := fmt.Sprintf("%d", 9999999)
		request := httptest.NewRequest(http.MethodGet, "http://localhost:3030/todo-items/"+wrongID, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 404, response.StatusCode)
		require.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Todo with ID %s Not Found", wrongID)
		require.Equal(t, message, responseBody["message"])

		require.Empty(t, responseBody["data"])
	})
}

func TestUpdateTodo(t *testing.T) {
	t.Parallel()

	t.Run("Success update todo with field title", func(t *testing.T) {
		newTodo := createRandomTodoHandler(t)
		data := web.TodoUpdateRequest{
			Title: jabufaker.RandomString(20),
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		id := fmt.Sprintf("%d", newTodo.ID)
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:3030/todo-items/"+id, requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "Success", responseBody["status"])
		assert.Equal(t, "Success", responseBody["message"])
		assert.NotEmpty(t, responseBody["data"])

		var contextData = responseBody["data"].(map[string]interface{})
		fmt.Println(contextData)
		assert.Equal(t, newTodo.ID, uint64(contextData["id"].(float64)))
		assert.Equal(t, newTodo.ActivityGroupID, uint64(contextData["activity_group_id"].(float64)))
		assert.Equal(t, newTodo.IsActive, contextData["is_active"])
		assert.Equal(t, newTodo.Priority, contextData["priority"])

		assert.NotEmpty(t, contextData["created_at"])

		assert.NotEqual(t, newTodo.UpdatedAt.String(), contextData["updated_at"])
		assert.NotEqual(t, newTodo.Title, contextData["title"])

		assert.Nil(t, newTodo.DeletetAt)
	})

	t.Run("Success update todo with field is_active", func(t *testing.T) {
		newTodo := createRandomTodoHandler(t)
		data := web.TodoUpdateRequest{
			IsActive: false,
		}

		dataBody := fmt.Sprintf(`{"is_active": %t}`, data.IsActive)
		requestBody := strings.NewReader(dataBody)
		fmt.Println(dataBody)
		id := fmt.Sprintf("%d", newTodo.ID)
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:3030/todo-items/"+id, requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "Success", responseBody["status"])
		assert.Equal(t, "Success", responseBody["message"])
		assert.NotEmpty(t, responseBody["data"])

		var contextData = responseBody["data"].(map[string]interface{})
		assert.Equal(t, newTodo.ID, uint64(contextData["id"].(float64)))
		assert.Equal(t, newTodo.ActivityGroupID, uint64(contextData["activity_group_id"].(float64)))
		assert.Equal(t, newTodo.Title, contextData["title"])
		assert.Equal(t, newTodo.Priority, contextData["priority"])

		assert.NotEmpty(t, contextData["created_at"])

		assert.NotEqual(t, newTodo.UpdatedAt.String(), contextData["updated_at"])
		assert.NotEqual(t, newTodo.IsActive, contextData["is_active"])

		assert.Nil(t, newTodo.DeletetAt)

	})

	t.Run("Id not found", func(t *testing.T) {
		data := web.ActivityGroupUpdateRequest{
			Title: jabufaker.RandomString(20),
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		wrongId := "999999"
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:3030/todo-items/"+wrongId, requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 404, response.StatusCode)
		assert.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Todo with ID %s Not Found", wrongId)
		assert.Equal(t, message, responseBody["message"])
		assert.Empty(t, responseBody["data"])
	})
}

func TestDeleteTodo(t *testing.T) {
	t.Parallel()
	newTodo := createRandomTodoHandler(t)

	t.Run("Deleted success", func(t *testing.T) {
		id := fmt.Sprintf("%d", newTodo.ID)
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:3030/todo-items/"+id, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 200, response.StatusCode)
		require.Equal(t, "Success", responseBody["status"])
		require.Equal(t, "Success", responseBody["message"])
		require.Empty(t, responseBody["data"])
	})

	t.Run("Id not found", func(t *testing.T) {
		wrongId := "999999"
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:3030/todo-items/"+wrongId, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		Route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		require.Equal(t, 404, response.StatusCode)
		require.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Todo with ID %s Not Found", wrongId)
		require.Equal(t, message, responseBody["message"])
		require.Empty(t, responseBody["data"])
	})
}
