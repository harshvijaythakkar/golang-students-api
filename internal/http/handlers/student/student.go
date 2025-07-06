package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

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

// Get student by Id
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting Student by Id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}
		student, err := storage.GetStudentByID(intId)

		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, student)

	}
}

// Get all students
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Geting all students")
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}

// Delete student by ID
func Delete(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Deleting student", slog.String("id", fmt.Sprint(id)))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		_, err = storage.GetStudentByID(intId)
		if err != nil {
			slog.Error("Student not found", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		err = storage.DeleteStudent(intId)
		if err != nil {
			slog.Error("Error in deleting student", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}
		response.WriteJson(w, http.StatusNoContent, map[string]bool{"result": true})
	}
}
