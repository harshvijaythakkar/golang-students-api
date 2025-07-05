package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/harshvijaythakkar/golang-students-api/internal/types"
	"github.com/harshvijaythakkar/golang-students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)


func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating student")

		// we need to serialise request body into student struct
		var student types.Student

		// Deocde request body using json package and store it in student struct
		err := json.NewDecoder(r.Body).Decode(&student)

		// check for empty body error
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return 
		}

		// check for other error
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			
			// type cast err into validator.ValidationErrors type
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))

			return 
		}


		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "Ok"})
	}
}

