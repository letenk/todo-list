package web

type ActivityGroupIdURI struct {
	ID int `uri:"id" binding:"required"`
}
type ActivityGroupRequest struct {
	Title string `json:"title" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type ActivityGroupUpdateRequest struct {
	Title string `json:"title" binding:"required"`
}
