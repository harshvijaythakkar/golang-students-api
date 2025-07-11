package storage

import "github.com/harshvijaythakkar/golang-students-api/internal/types"

type Storage interface{
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	DeleteStudent(id int64) (error)
	UpdateStudent(id int64, data map[string]interface{}) (error)
}

