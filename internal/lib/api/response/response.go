package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"` // ok or error
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "ok"
	StatusError = "error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errsMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errsMsg = append(errsMsg, fmt.Sprintf("field %s is required", err.Field()))
		case "url":
			errsMsg = append(errsMsg, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errsMsg = append(errsMsg, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errsMsg, ","),
	}
}
