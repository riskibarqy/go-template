package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	validator "gopkg.in/go-playground/validator.v9"
)

// ErrorResponse represents the default error response
type ErrorResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Fields  []*FieldError `json:"fields"`
}

// FieldError represents error message for each field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Responder represents the http responder interface
type Responder struct {
}

func makeFieldError(field string, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

// JSON writes json http response
func (r *Responder) JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// HTML writes html http response
func (r *Responder) HTML(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	fmt.Fprint(w, data)
}

// Error writes error http response
func (r *Responder) Error(w http.ResponseWriter, status int, err error) {
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
		errorCode = "UnprocessableEntity"
	case http.StatusInternalServerError:
		errorCode = "InternalServerError"
	}

	errorFields := []*FieldError{}
	switch err := err.(type) {
	case validator.ValidationErrors:
		for _, e := range err {
			errorFields = append(errorFields,
				makeFieldError(e.Field(), e.ActualTag()))
		}
	}

	if status == http.StatusInternalServerError {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    errorCode,
			Message: "Server error",
			Fields:  errorFields,
		})

		errMessage := fmt.Sprintf("%+v\n%s", err, string(debug.Stack()))
		log.Println(errMessage)
		return
	}

	if len(errorFields) > 0 {
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    errorCode,
			Message: "validation error",
			Fields:  errorFields,
		})
		return
	}

	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    errorCode,
		Message: err.Error(),
		Fields:  errorFields,
	})
}

// NewResponder creates a new http responder
func NewResponder() *Responder {
	return &Responder{}
}
