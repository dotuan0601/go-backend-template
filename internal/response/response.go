package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(
	w http.ResponseWriter,
	data interface{},
) error {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	resp := Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}

	return json.NewEncoder(w).Encode(resp)
}
