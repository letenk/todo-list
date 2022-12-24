package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jellydator/ttlcache/v2"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/service"
)

type todoHandler struct {
	service service.TodoService
}

func NewTodoHandler(service service.TodoService) *todoHandler {
	return &todoHandler{service}
}

func (h *todoHandler) GetAll(c *gin.Context) {
	todoID, _ := strconv.Atoi(c.Query("activity_group_id"))
	key := "todos"
	keyTodoQuerySearch := "todo-search"
	keyTodoID := "todo-id"

	if todoID != 0 {
		cache.Remove(key)

		todosIdCache, _ := cache.Get(keyTodoID)
		if todosIdCache != todoID {
			cache.Remove(keyTodoQuerySearch)
		}

		todos, err := cache.Get(keyTodoQuerySearch)
		if err == ttlcache.ErrNotFound {
			cache.Remove(key)
			cache.Remove(keyTodoQuerySearch)
			// Get all
			todos, err := h.service.GetAll(uint64(todoID))
			if err != nil {
				resp := gin.H{}
				jsonResponse := web.JSONResponse(
					"Internal Server Error",
					"Internal Server Error",
					resp,
				)
				c.JSON(http.StatusInternalServerError, jsonResponse)
				return
			}

			formatResponseJSON := web.FormatTodos(todos)

			if len(todos) != 0 {
				// Cache data
				go cache.SetWithTTL(key, formatResponseJSON, time.Hour)
				go cache.SetWithTTL(keyTodoQuerySearch, formatResponseJSON, time.Hour)
				go cache.SetWithTTL(keyTodoID, todoID, time.Hour)
			} else {
				go cache.Remove(key)
			}

			jsonResponse := web.JSONResponse(
				"Success",
				"Success",
				formatResponseJSON,
			)
			c.JSON(http.StatusOK, jsonResponse)
			return
		}
		// Get all data from cache
		jsonResponse := web.JSONResponse(
			"Success",
			"Success",
			todos,
		)
		c.JSON(http.StatusOK, jsonResponse)
		return
	}

	_, err := cache.Get(keyTodoQuerySearch)
	if err != ttlcache.ErrNotFound {
		cache.Remove(key)
		cache.Remove(keyTodoQuerySearch)
	}

	// Get data from cache
	todos, err := cache.Get(key)
	if err == ttlcache.ErrNotFound {
		// Get all
		todos, err := h.service.GetAll(uint64(todoID))
		if err != nil {
			resp := gin.H{}
			jsonResponse := web.JSONResponse(
				"Internal Server Error",
				"Internal Server Error",
				resp,
			)
			c.JSON(http.StatusInternalServerError, jsonResponse)
			return
		}

		formatResponseJSON := web.FormatTodos(todos)

		if len(todos) != 0 {
			// Cache data
			go cache.SetWithTTL(key, formatResponseJSON, time.Hour)
		} else {
			go cache.Remove(key)
		}

		jsonResponse := web.JSONResponse(
			"Success",
			"Success",
			formatResponseJSON,
		)
		c.JSON(http.StatusOK, jsonResponse)
		return
	}

	// Get all data from cache
	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		todos,
	)
	c.JSON(http.StatusOK, jsonResponse)
}

func (h *todoHandler) GetOne(c *gin.Context) {
	var todoID web.TodoURI
	err := c.ShouldBindUri(&todoID)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Bad Request",
			"Uri id cannot be null",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Get data from cache
	key := fmt.Sprintf("todo-id-%d", todoID.ID)
	todo, err := cache.Get(key)
	if err == ttlcache.ErrNotFound {

		// Get all
		todo, err := h.service.GetOne(todoID.ID)
		if err != nil {
			resp := gin.H{}
			jsonResponse := web.JSONResponse(
				"Internal Server Error",
				"Internal Server Error",
				resp,
			)
			c.JSON(http.StatusInternalServerError, jsonResponse)
			return
		}

		// If not found
		if todo.ID == 0 {
			resp := gin.H{}
			message := fmt.Sprintf("Todo with ID %d Not Found", todoID.ID)
			jsonResponse := web.JSONResponse(
				"Not Found",
				message,
				resp,
			)
			c.JSON(http.StatusNotFound, jsonResponse)
			return
		}

		formatResponseJSON := web.FormatTodo(todo)
		// Cache data
		go cache.SetWithTTL(key, formatResponseJSON, time.Hour)

		jsonResponse := web.JSONResponse(
			"Success",
			"Success",
			formatResponseJSON,
		)
		c.JSON(http.StatusOK, jsonResponse)
		return
	}

	// Get one data from cache
	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		todo,
	)
	c.JSON(http.StatusOK, jsonResponse)
}

func (h *todoHandler) Create(c *gin.Context) {
	var req web.TodoCreateRequest
	err := c.ShouldBindJSON(&req)
	if req.Title == "" {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Bad Request",
			"title cannot be null",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	if req.ActivityGroupID == 0 {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Bad Request",
			"activity_group_id cannot be null",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Bad Request",
			"title, activity_group_id cannot be null",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Create
	newTodo, err := h.service.Create(req)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Internal Server Error",
			"Internal Server Error",
			resp,
		)
		c.JSON(http.StatusInternalServerError, jsonResponse)
		return
	}

	formatResponseJSON := web.FormatCreatedTodo(newTodo)
	// Cache
	if newTodo.ID != 0 {
		key := fmt.Sprintf("todo-id-%d", newTodo.ID)
		go cache.SetWithTTL(key, formatResponseJSON, time.Hour)
		go cache.Remove("todos")
	}

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		formatResponseJSON,
	)
	c.JSON(http.StatusCreated, jsonResponse)
}

func (h *todoHandler) Update(c *gin.Context) {
	var todoURI web.TodoURI
	err := c.ShouldBindUri(&todoURI)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Bad Request",
			"Uri id cannot be null",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	var req web.TodoUpdateRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Bad Request",
			"title or is_active cannot be null",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Get one by id
	todo, err := h.service.GetOne(todoURI.ID)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Internal Server Error",
			err.Error(),
			resp,
		)
		c.JSON(http.StatusInternalServerError, jsonResponse)
		return
	}

	// If not found
	if todo.ID == 0 {
		resp := gin.H{}
		message := fmt.Sprintf("Todo with ID %d Not Found", todoURI.ID)
		jsonResponse := web.JSONResponse(
			"Not Found",
			message,
			resp,
		)
		c.JSON(http.StatusNotFound, jsonResponse)
		return
	}

	if req.Title != "" {
		req.IsActive = true
	}

	// Update
	updatedTodo, err := h.service.Update(todo.ID, req)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Internal Server Error",
			err.Error(),
			resp,
		)
		c.JSON(http.StatusInternalServerError, jsonResponse)
		return
	}

	formatResponseJSON := web.FormatTodo(updatedTodo)
	// Cache
	if updatedTodo.ID != 0 {
		key := fmt.Sprintf("todo-id-%d", updatedTodo.ID)
		go cache.Remove(key)
		go cache.SetWithTTL(key, formatResponseJSON, time.Hour)
		go cache.Remove("todos")
	}

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		formatResponseJSON,
	)
	c.JSON(http.StatusOK, jsonResponse)
}

func (h *todoHandler) Delete(c *gin.Context) {
	var todoURI web.ActivityIdURI
	err := c.ShouldBindUri(&todoURI)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Bad Request",
			"Uri id cannot be null",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Find by id
	todo, err := h.service.GetOne(todoURI.ID)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Internal Server Error",
			err.Error(),
			resp,
		)
		c.JSON(http.StatusInternalServerError, jsonResponse)
		return
	}

	// If not found
	if todo.ID == 0 {
		resp := gin.H{}
		message := fmt.Sprintf("Todo with ID %d Not Found", todoURI.ID)
		jsonResponse := web.JSONResponse(
			"Not Found",
			message,
			resp,
		)
		c.JSON(http.StatusNotFound, jsonResponse)
		return
	}

	// Delete
	ok, err := h.service.Delete(todo.ID)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Internal Server Error",
			err.Error(),
			resp,
		)
		c.JSON(http.StatusInternalServerError, jsonResponse)
		return
	}

	// Cache and remove
	if ok {
		key := fmt.Sprintf("todo-id-%d", todo.ID)
		go cache.Remove(key)
		go cache.Remove("activities")
	}

	resp := gin.H{}
	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		resp,
	)
	c.JSON(http.StatusOK, jsonResponse)
}
