package web

type ResponseWithData struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSONResponse(status string, message string, data interface{}) ResponseWithData {
	jsonResponse := ResponseWithData{
		Status:  status,
		Message: message,
		Data:    data,
	}

	return jsonResponse
}
