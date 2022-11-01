package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/service"
)

type activityGroupHandler struct {
	service service.ActivityGroupService
}

func NewActivityGroupHandler(service service.ActivityGroupService) *activityGroupHandler {
	return &activityGroupHandler{service}
}

func (h *activityGroupHandler) GetAll(c *gin.Context) {
	// Get all
	activityGroup, err := h.service.GetAll()
	if err != nil {
		jsonResponse := web.JSONResponse(
			"Internal Server Error",
			"Internal Server Error",
			domain.ActivityGroup{},
		)
		c.JSON(http.StatusInternalServerError, jsonResponse)
		return
	}

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		web.FormatActivitiesGroup(activityGroup),
	)
	c.JSON(http.StatusOK, jsonResponse)
}

func (h *activityGroupHandler) GetOne(c *gin.Context) {
	var id web.ActivityGroupIdURI
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
	activityGroup, err := h.service.GetOne(int(id.ID))
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
	if activityGroup.ID == 0 {
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
		web.FormatActivityGroupGetOne(activityGroup),
	)
	c.JSON(http.StatusOK, jsonResponse)

}

func (h *activityGroupHandler) Create(c *gin.Context) {
	var req web.ActivityGroupRequest
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
	newActivityGroup, err := h.service.Create(req)
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
		web.FormatActivityGroup(newActivityGroup),
	)
	c.JSON(http.StatusCreated, jsonResponse)
}

func (h *activityGroupHandler) Update(c *gin.Context) {
	var id web.ActivityGroupIdURI
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

	var req web.ActivityGroupUpdateRequest
	err = c.ShouldBindJSON(&req)
	fmt.Println(err)
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
	activityGroup, err := h.service.GetOne(int(id.ID))
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
	if activityGroup.ID == 0 {
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
	updatedActivityGroup, err := h.service.Update(activityGroup.ID, req)
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
		web.FormatActivityGroupGetOne(updatedActivityGroup),
	)
	c.JSON(http.StatusOK, jsonResponse)

}

func (h *activityGroupHandler) Delete(c *gin.Context) {
	var id web.ActivityGroupIdURI
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
	activityGroup, err := h.service.GetOne(int(id.ID))
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
	if activityGroup.ID == 0 {
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
	_, err = h.service.Delete(activityGroup.ID)
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
