package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/harshvijaythakkar/golang-students-api/internal/storage"
	"github.com/harshvijaythakkar/golang-students-api/internal/types"
	"github.com/harshvijaythakkar/golang-students-api/internal/utils/response"
)


func New(storage storage.Storage) http.HandlerFunc {
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

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

