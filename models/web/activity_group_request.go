package web

import (
	"time"

	"github.com/letenk/todo-list/models/domain"
)

type ActivityIdURI struct {
	ID uint64 `uri:"id" binding:"required"`
}
type ActivityRequest struct {
	Title string `json:"title" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type ActivityUpdateRequest struct {
	Title string `json:"title" binding:"required"`
}

type ActivityCreateResponse struct {
	ID        uint64     `json:"id"`
	Title     string     `json:"title"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type ActivityGetOneResponse struct {
	ID        uint64     `json:"id"`
	Title     string     `json:"title"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Format for handle single response activity group
func FormatActivity(Activity domain.Activity) ActivityCreateResponse {
	formatter := ActivityCreateResponse{
		ID:        Activity.ID,
		Title:     Activity.Title,
		Email:     Activity.Email,
		CreatedAt: Activity.CreatedAt,
		UpdatedAt: Activity.UpdatedAt,
	}
	return formatter
}

// Format for handle get One response activity group
func FormatActivityGetOne(Activity domain.Activity) ActivityGetOneResponse {
	formatter := ActivityGetOneResponse{
		ID:        Activity.ID,
		Title:     Activity.Title,
		Email:     Activity.Email,
		CreatedAt: Activity.CreatedAt,
		UpdatedAt: Activity.UpdatedAt,
		DeletedAt: Activity.DeletedAt,
	}
	return formatter
}

// Format for handle multiples response activity group
func FormatActivitiesGroup(Activity []domain.Activity) []ActivityCreateResponse {
	if len(Activity) == 0 {
		return []ActivityCreateResponse{}
	}

	var formatters []ActivityCreateResponse

	for _, data := range Activity {
		formatter := FormatActivity(data)
		formatters = append(formatters, formatter)
	}

	return formatters
}
