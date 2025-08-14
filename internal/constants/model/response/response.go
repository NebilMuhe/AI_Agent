package response

import (
	"ai_agent/internal/constants/errors"
	"encoding/json"
	"net/http"
)

func SendSuccessResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Response{
		Ok:   true,
		Data: data,
	}); err != nil {
		w.WriteHeader(errors.ErrorMap[errors.ErrUnexpected])
		json.NewEncoder(w).Encode(Response{
			Ok: false,
			Error: &ErrorResponse{
				StausCode: errors.ErrorMap[errors.ErrUnexpected],
				Message:   errors.ErrUnexpected.Error(),
			},
		})
		return
	}
}

func SendErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	statusCode, ok := errors.ErrorMap[err]
	if !ok {
		w.WriteHeader(errors.ErrorMap[errors.ErrUnexpected])
		json.NewEncoder(w).Encode(Response{
			Ok: false,
			Error: &ErrorResponse{
				StausCode: errors.ErrorMap[errors.ErrUnexpected],
				Message:   errors.ErrUnexpected.Error(),
			},
		})
		return
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Response{
		Ok: false,
		Error: &ErrorResponse{
			StausCode: errors.ErrorMap[err],
			Message:   err.Error(),
		},
	}); err != nil {
		w.WriteHeader(errors.ErrorMap[errors.ErrUnexpected])
		json.NewEncoder(w).Encode(Response{
			Ok: false,
			Error: &ErrorResponse{
				StausCode: errors.ErrorMap[errors.ErrUnexpected],
				Message:   errors.ErrUnexpected.Error(),
			},
		})
	}
}

