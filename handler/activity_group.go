package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/service"
)

type ActivityHandler struct {
	service service.ActivityService
}

func NewActivityHandler(service service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service}
}

func (h *ActivityHandler) GetAll(c *gin.Context) {
	// Get all
	Activity, err := h.service.GetAll()
	if err != nil {
		jsonResponse := web.JSONResponse(
			"Internal Server Error",
			"Internal Server Error",
			domain.Activity{},
		)
		c.JSON(http.StatusInternalServerError, jsonResponse)
		return
	}

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		web.FormatActivitiesGroup(Activity),
	)
	c.JSON(http.StatusOK, jsonResponse)
}

func (h *ActivityHandler) GetOne(c *gin.Context) {
	var id web.ActivityIdURI
	err := c.ShouldBindUri(&id)
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
	Activity, err := h.service.GetOne(id.ID)
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
	if Activity.ID == 0 {
		resp := gin.H{}
		message := fmt.Sprintf("Activity with ID %d Not Found", id.ID)
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
		web.FormatActivityGetOne(Activity),
	)
	c.JSON(http.StatusOK, jsonResponse)

}

func (h *ActivityHandler) Create(c *gin.Context) {
	var req web.ActivityRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Bad Request",
			"title cannot be null",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Create
	newActivity, err := h.service.Create(req)
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
		web.FormatActivity(newActivity),
	)
	c.JSON(http.StatusCreated, jsonResponse)
}

func (h *ActivityHandler) Update(c *gin.Context) {
	var id web.ActivityIdURI
	err := c.ShouldBindUri(&id)
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

	var req web.ActivityUpdateRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		resp := gin.H{}
		jsonResponse := web.JSONResponse(
			"Bad Request",
			"title cannot be null",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Find by id
	Activity, err := h.service.GetOne(id.ID)
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
	if Activity.ID == 0 {
		resp := gin.H{}
		message := fmt.Sprintf("Activity with ID %d Not Found", id.ID)
		jsonResponse := web.JSONResponse(
			"Not Found",
			message,
			resp,
		)
		c.JSON(http.StatusNotFound, jsonResponse)
		return
	}

	// Update
	updatedActivity, err := h.service.Update(Activity.ID, req)
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
		web.FormatActivityGetOne(updatedActivity),
	)
	c.JSON(http.StatusOK, jsonResponse)

}

func (h *ActivityHandler) Delete(c *gin.Context) {
	var id web.ActivityIdURI
	err := c.ShouldBindUri(&id)
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
	Activity, err := h.service.GetOne(id.ID)
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
	if Activity.ID == 0 {
		resp := gin.H{}
		message := fmt.Sprintf("Activity with ID %d Not Found", id.ID)
		jsonResponse := web.JSONResponse(
			"Not Found",
			message,
			resp,
		)
		c.JSON(http.StatusNotFound, jsonResponse)
		return
	}

	// Delete
	_, err = h.service.Delete(Activity.ID)
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
