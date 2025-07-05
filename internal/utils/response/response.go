package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error string `json:"error"`
}

const (
	StatusOK = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {

	// set header for response
	w.Header().Set("Content-Type", "application/json")

	// set http response status code with provided header
	w.WriteHeader(status)

	// Return Response, encode method resturns error
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error: err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMessages []string
	
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field %s is required filed", err.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s is invalid filed", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error: strings.Join(errMessages, ", "),
	}
}
