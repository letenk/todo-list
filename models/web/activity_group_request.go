package web

import (
	"time"

	"github.com/letenk/todo-list/models/domain"
)

type ActivityGroupIdURI struct {
	ID int64 `uri:"id" binding:"required"`
}
type ActivityGroupRequest struct {
	Title string `json:"title" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type ActivityGroupUpdateRequest struct {
	Title string `json:"title" binding:"required"`
}

type ActivityGroupCreateResponse struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ActivityGroupGetOneResponse struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Format for handle single response activity group
func FormatActivityGroup(activityGroup domain.ActivityGroup) ActivityGroupCreateResponse {
	formatter := ActivityGroupCreateResponse{
		ID:        activityGroup.ID,
		Title:     activityGroup.Title,
		Email:     activityGroup.Email,
		CreatedAt: activityGroup.CreatedAt,
		UpdatedAt: activityGroup.UpdatedAt,
	}
	return formatter
}

// Format for handle get One response activity group
func FormatActivityGroupGetOne(activityGroup domain.ActivityGroup) ActivityGroupGetOneResponse {
	formatter := ActivityGroupGetOneResponse{
		ID:        activityGroup.ID,
		Title:     activityGroup.Title,
		Email:     activityGroup.Email,
		CreatedAt: activityGroup.CreatedAt,
		UpdatedAt: activityGroup.UpdatedAt,
		DeletedAt: activityGroup.DeletedAt,
	}
	return formatter
}

// Format for handle multiples response activity group
func FormatActivitiesGroup(activityGroup []domain.ActivityGroup) []ActivityGroupCreateResponse {
	if len(activityGroup) == 0 {
		return []ActivityGroupCreateResponse{}
	}

	var formatters []ActivityGroupCreateResponse

	for _, data := range activityGroup {
		formatter := FormatActivityGroup(data)
		formatters = append(formatters, formatter)
	}

	return formatters
}
