package web

type TodoCreateRequest struct {
	ActivityGroupId uint64 `json:"activity_group_id" binding:"required"`
	Title           string `json:"title" binding:"required"`
}

type TodoUpdateRequest struct {
	Title    string `json:"title,omitempty"`
	IsActive bool   `json:"is_active"`
}
