package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	activityGroupID, _ := strconv.Atoi(c.Query("activity_group_id"))

	// Get all
	todos, err := h.service.GetAll(uint64(activityGroupID))
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

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		web.FormatTodos(todos),
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
	// Get all
	todos, err := h.service.GetOne(todoID.ID)
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
	if todos.ID == 0 {
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

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		web.FormatTodo(todos),
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

	if req.ActivityGroupId == 0 {
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

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		web.FormatTodo(newTodo),
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

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		web.FormatTodo(updatedTodo),
	)
	c.JSON(http.StatusOK, jsonResponse)
}

func (h *todoHandler) Delete(c *gin.Context) {
	var todoURI web.ActivityGroupIdURI
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
	_, err = h.service.Delete(todo.ID)
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

	resp := gin.H{}
	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		resp,
	)
	c.JSON(http.StatusOK, jsonResponse)

}
