package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jellydator/ttlcache/v2"
	"github.com/letenk/todo-list/models/domain"
	"github.com/letenk/todo-list/models/web"
	"github.com/letenk/todo-list/service"
)

// Create new instance ttlcache
var cache ttlcache.SimpleCache = ttlcache.NewCache()

type ActivityHandler struct {
	service service.ActivityService
}

func NewActivityHandler(service service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service}
}

func (h *ActivityHandler) GetAll(c *gin.Context) {
	// Get data from cache
	key := "activities"
	activities, err := cache.Get(key)
	if err == ttlcache.ErrNotFound {
		// Get all from database
		activities, err := h.service.GetAll()
		if err != nil {
			jsonResponse := web.JSONResponse(
				"Internal Server Error",
				"Internal Server Error",
				domain.Activity{},
			)
			c.JSON(http.StatusInternalServerError, jsonResponse)
			return
		}

		formatResponseJSON := web.FormatActivitiesGroup(activities)
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

	// Get all data from cache
	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		activities,
	)
	c.JSON(http.StatusOK, jsonResponse)

}

func (h *ActivityHandler) GetOne(c *gin.Context) {
	var activityID web.ActivityIdURI
	err := c.ShouldBindUri(&activityID)
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
	key := fmt.Sprintf("activity-id-%d", activityID.ID)
	activity, err := cache.Get(key)
	if err == ttlcache.ErrNotFound {
		// Find by id from database
		Activity, err := h.service.GetOne(activityID.ID)
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
			message := fmt.Sprintf("Activity with ID %d Not Found", activityID.ID)
			jsonResponse := web.JSONResponse(
				"Not Found",
				message,
				resp,
			)
			c.JSON(http.StatusNotFound, jsonResponse)
			return
		}

		formatResponseJSON := web.FormatActivityGetOne(Activity)
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
		activity,
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

	formatResponseJSON := web.FormatActivity(newActivity)
	// Cache
	if newActivity.ID != 0 {
		key := fmt.Sprintf("activity-id-%d", newActivity.ID)
		go cache.SetWithTTL(key, formatResponseJSON, time.Hour)
		go cache.Remove("activities")
	}

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		formatResponseJSON,
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
	formatResponseJSON := web.FormatActivityGetOne(updatedActivity)
	// Cache and remove
	if updatedActivity.ID != 0 {
		key := fmt.Sprintf("activity-id-%d", updatedActivity.ID)
		go cache.Remove(key)
		go cache.SetWithTTL(key, formatResponseJSON, time.Hour)
		go cache.Remove("activities")
	}

	jsonResponse := web.JSONResponse(
		"Success",
		"Success",
		formatResponseJSON,
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
	activity, err := h.service.GetOne(id.ID)
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
	if activity.ID == 0 {
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
	ok, err := h.service.Delete(activity.ID)
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
		key := fmt.Sprintf("activity-id-%d", activity.ID)
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
