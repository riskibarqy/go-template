package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/riskibarqy/99backend-challenge/internal/types"

	validator "gopkg.in/go-playground/validator.v9"

	"github.com/pkg/errors"
)

// FieldError represents error message for each field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents error message
type ErrorResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Fields  []*FieldError `json:"fields"`
}

// MakeFieldError create field error object
func MakeFieldError(field string, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

// Error writes error http response
func Error(w http.ResponseWriter, data string, status int, err types.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var errorCode string
	switch status {
	case http.StatusUnauthorized:
		errorCode = "Unauthorized"
	case http.StatusNotFound:
		errorCode = "NotFound"
	case http.StatusBadRequest:
		errorCode = "BadRequest"
	case http.StatusUnprocessableEntity:
		errorCode = "ValidationError"
	}

	errorFields := []*FieldError{}

	switch err.Error.(type) {
	case validator.ValidationErrors:
		data = "Bad Request"
		for _, err := range err.Error.(validator.ValidationErrors) {
			e := MakeFieldError(
				err.Field(),
				err.ActualTag())

			errorFields = append(errorFields, e)
		}
	}

	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    errorCode,
		Message: data,
		Fields:  errorFields,
	})

	if err.Error != nil {
		log.Printf("INFO: %v\n", err.Error.Error())
		log.Printf("DETAIL [%s - %s]: %s\n", err.Path, err.Type, err.Message)
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}

		var st errors.StackTrace
		if err, ok := err.Error.(stackTracer); ok {
			st = err.StackTrace()
			fmt.Printf("INFO: %+v\n", st[0])
		}
	}
}
