package response

import (
	"encoding/json"
	"net/http"
)

type ErrorCode struct {
	Code    int
	Message string
}

var (
	ErrSystem = ErrorCode{
		Code:    500000,
		Message: "System error",
	}

	ErrUnauthorized = ErrorCode{
		Code:    310101,
		Message: "Unauthorized",
	}

	ErrNotFound = ErrorCode{
		Code:    404001,
		Message: "Not found",
	}
)

func Error(
	w http.ResponseWriter,
	err ErrorCode,
) error {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	resp := Response{
		Code:    err.Code,
		Message: err.Message,
		Data:    nil,
	}

	return json.NewEncoder(w).Encode(resp)
}
