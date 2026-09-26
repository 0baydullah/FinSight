package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code       int         `json:"code"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Error      interface{} `json:"error"`
	Pagination interface{} `json:"pagination"`
}

type Error struct {
	Code    string      `json:"code"`
	Details interface{} `json:"details"`
}

func JSON(
	w http.ResponseWriter,
	statusCode int,
	message string,
	data interface{},
) {
	writeJSON(w, statusCode, Response{
		Code:       statusCode,
		Success:    true,
		Message:    message,
		Data:       data,
		Error:      nil,
		Pagination: nil,
	})
}

func ErrorJSON(
	w http.ResponseWriter,
	statusCode int,
	message string,
	errorCode string,
	details interface{},
) {
	writeJSON(w, statusCode, Response{
		Code:    statusCode,
		Success: false,
		Message: message,
		Data:    nil,
		Error: Error{
			Code:    errorCode,
			Details: details,
		},
		Pagination: nil,
	})
}

func writeJSON(
	w http.ResponseWriter,
	statusCode int,
	response Response,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(response)
}
