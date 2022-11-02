package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/letenk/todo-list/config"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/router"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var mutex sync.Mutex
var conn *gorm.DB
var route *gin.Engine

func TestMain(m *testing.M) {
	// Set env
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "root")
	os.Setenv("MYSQL_HOST", "127.0.0.1")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DBNAME", "todo4")

	// Open connection
	db := config.SetupDB()
	conn = db

	// setup router
	route = router.SetupRouter(db)
	m.Run()
}

func createRandomActivityGroupHandler(t *testing.T) web.ActivityGroupCreateResponse {
	data := web.ActivityGroupRequest{
		Title: jabufaker.RandomString(20),
		Email: jabufaker.RandomEmail(),
	}

	dataBody := fmt.Sprintf(`{"title": "%s", "email": "%s"}`, data.Title, data.Email)
	requestBody := strings.NewReader(dataBody)

	request := httptest.NewRequest(http.MethodPost, "http://localhost:3030/activity-groups", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 201, response.StatusCode)
	assert.Equal(t, "Success", responseBody["status"])
	assert.Equal(t, "Success", responseBody["message"])

	assert.NotEmpty(t, responseBody["data"])

	var contextData = responseBody["data"].(map[string]interface{})
	assert.NotEmpty(t, contextData["id"])
	assert.NotEmpty(t, contextData["created_at"])
	assert.NotEmpty(t, contextData["updated_at"])
	assert.Equal(t, data.Title, contextData["title"])
	assert.Equal(t, data.Email, contextData["email"])

	newActivityGroup := web.ActivityGroupCreateResponse{
		ID:    int64(contextData["id"].(float64)),
		Title: contextData["title"].(string),
		Email: contextData["email"].(string),
	}

	return newActivityGroup
}

func TestActivityGroupCreateHandler(t *testing.T) {
	t.Parallel()
	t.Run("Handler Create activity Group success", func(t *testing.T) {
		createRandomActivityGroupHandler(t)
	})

	t.Run("Handler Create activity Group failed title blank", func(t *testing.T) {
		dataBody := fmt.Sprintf(`{"title": "%s", "email": "%s"}`, "", jabufaker.RandomEmail())
		requestBody := strings.NewReader(dataBody)

		request := httptest.NewRequest(http.MethodPost, "http://localhost:3030/activity-groups", requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 400, response.StatusCode)
		assert.Equal(t, "Bad Request", responseBody["status"])
		assert.Equal(t, "title cannot be null", responseBody["message"])
		assert.Empty(t, responseBody["data"])
	})
}

func TestGetAllActivityGroupHandler(t *testing.T) {
	t.Parallel()
	// Create some random data
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			createRandomActivityGroupHandler(t)
			mutex.Unlock()
		}()
	}

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3030/activity-groups", nil)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "Success", responseBody["status"])
	assert.Equal(t, "Success", responseBody["message"])

	assert.NotEmpty(t, responseBody["data"])

	var contextData = responseBody["data"].([]interface{})

	assert.NotEqual(t, 0, len(contextData))
	// Data is not null

	for _, data := range contextData {
		list := data.(map[string]interface{})
		assert.NotEmpty(t, list["id"])
		assert.NotEmpty(t, list["title"])
		assert.NotEmpty(t, list["email"])
		assert.NotEmpty(t, list["created_at"])
		assert.NotEmpty(t, list["updated_at"])
	}
}

func TestGetOneActivityGroup(t *testing.T) {
	t.Parallel()
	newActivityGroup := createRandomActivityGroupHandler(t)

	t.Run("Success get one", func(t *testing.T) {
		id := fmt.Sprintf("%d", newActivityGroup.ID)
		request := httptest.NewRequest(http.MethodGet, "http://localhost:3030/activity-groups/"+id, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "Success", responseBody["status"])
		assert.Equal(t, "Success", responseBody["message"])

		assert.NotEmpty(t, responseBody["data"])

		var contextData = responseBody["data"].(map[string]interface{})
		assert.Equal(t, newActivityGroup.ID, int64(contextData["id"].(float64)))
		assert.Equal(t, newActivityGroup.Title, contextData["title"])
		assert.Equal(t, newActivityGroup.Email, contextData["email"])

		assert.NotEmpty(t, contextData["created_at"])
		assert.NotEmpty(t, contextData["updated_at"])
		assert.Nil(t, contextData["deteled_at"])
	})

	t.Run("Id not found", func(t *testing.T) {
		wrongId := "999999"
		request := httptest.NewRequest(http.MethodGet, "http://localhost:3030/activity-groups/"+wrongId, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 404, response.StatusCode)
		assert.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Activity with ID %s Not Found", wrongId)
		assert.Equal(t, message, responseBody["message"])
		assert.Empty(t, responseBody["data"])
	})
}

func TestUpdateActivityGroup(t *testing.T) {
	t.Parallel()
	newActivityGroup := createRandomActivityGroupHandler(t)

	t.Run("Success get one", func(t *testing.T) {
		data := web.ActivityGroupUpdateRequest{
			Title: jabufaker.RandomString(20),
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		id := fmt.Sprintf("%d", newActivityGroup.ID)
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:3030/activity-groups/"+id, requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "Success", responseBody["status"])
		assert.Equal(t, "Success", responseBody["message"])
		assert.NotEmpty(t, responseBody["data"])

		var contextData = responseBody["data"].(map[string]interface{})
		assert.Equal(t, newActivityGroup.ID, int64(contextData["id"].(float64)))
		assert.Equal(t, newActivityGroup.Email, contextData["email"])

		assert.NotEqual(t, newActivityGroup.UpdatedAt.String(), contextData["updated_at"])
		assert.NotEqual(t, newActivityGroup.Title, contextData["title"])

		assert.NotEmpty(t, contextData["created_at"])

	})

	t.Run("Body blank", func(t *testing.T) {
		dataBody := fmt.Sprintf(`{"title": "%s"}`, "")
		requestBody := strings.NewReader(dataBody)

		id := fmt.Sprintf("%d", newActivityGroup.ID)
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:3030/activity-groups/"+id, requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 400, response.StatusCode)
		assert.Equal(t, "Bad Request", responseBody["status"])
		assert.Equal(t, "title cannot be null", responseBody["message"])
		assert.Empty(t, responseBody["data"])
	})

	t.Run("Id not found", func(t *testing.T) {
		data := web.ActivityGroupUpdateRequest{
			Title: jabufaker.RandomString(20),
		}

		dataBody := fmt.Sprintf(`{"title": "%s"}`, data.Title)
		requestBody := strings.NewReader(dataBody)

		wrongId := "999999"
		request := httptest.NewRequest(http.MethodPatch, "http://localhost:3030/activity-groups/"+wrongId, requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 404, response.StatusCode)
		assert.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Activity with ID %s Not Found", wrongId)
		assert.Equal(t, message, responseBody["message"])
		assert.Empty(t, responseBody["data"])
	})
}

func TestDeleteActivityGroup(t *testing.T) {
	t.Parallel()
	newActivityGroup := createRandomActivityGroupHandler(t)

	t.Run("Deleted success", func(t *testing.T) {
		id := fmt.Sprintf("%d", newActivityGroup.ID)
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:3030/activity-groups/"+id, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 200, response.StatusCode)
		assert.Equal(t, "Success", responseBody["status"])
		assert.Equal(t, "Success", responseBody["message"])
		assert.Empty(t, responseBody["data"])
	})

	t.Run("Id not found", func(t *testing.T) {
		wrongId := "999999"
		request := httptest.NewRequest(http.MethodDelete, "http://localhost:3030/activity-groups/"+wrongId, nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		route.ServeHTTP(recorder, request)

		response := recorder.Result()

		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, 404, response.StatusCode)
		assert.Equal(t, "Not Found", responseBody["status"])
		message := fmt.Sprintf("Activity with ID %s Not Found", wrongId)
		assert.Equal(t, message, responseBody["message"])
		assert.Empty(t, responseBody["data"])
	})
}
