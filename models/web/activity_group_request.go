package web

type ActivityGroupRequest struct {
	Title string `json:"title" binding:"required"`
	Email string `json:"email" binding:"required"`
}
