package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/harshvijaythakkar/golang-students-api/internal/config"
	"github.com/harshvijaythakkar/golang-students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)


type Sqlite struct {
	Db *sql.DB
}

// initialise db and return sqlite object, error
func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

// CreateStudent creates a new record in database and returns lastId, error
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	// prepare sql query with placeholders
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?);")
	if err != nil {
		return 0, err
	}

	// close statement
	defer stmt.Close()

	// execute statement
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	
	// get last insurted ID
	lastid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastid, nil
}

// Get Student by Id
func (s *Sqlite) GetStudentByID(id int64) (types.Student, error) {
	
	// prepare statement
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1;")
	if err != nil {
		return types.Student{}, err
	}

	// close statement
	defer stmt.Close()

	var studentInfo types.Student

	// execute statement
	err = stmt.QueryRow(id).Scan(&studentInfo.Id, &studentInfo.Name, &studentInfo.Email, &studentInfo.Age)
	if err != nil {

		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no Student found with id: %s", fmt.Sprint(id))
		}

		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return studentInfo, nil
}

// Get all students
func (s *Sqlite) GetStudents() ([]types.Student, error) {

	// prepare statement
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students;")
	if err != nil {
		return nil, err
	}

	// close statement
	defer stmt.Close()

	// execute statement
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	// close
	defer rows.Close()

	var studentsInfo []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		studentsInfo = append(studentsInfo, student)
	}

	return  studentsInfo, nil
}

// Delete student by ID
func (s *Sqlite) DeleteStudent(id int64) (error) {

	// prepare statement
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?;")
	if err != nil {
		return err
	}

	// close statement
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	slog.Info("Records Deleted", slog.String("count", fmt.Sprint(rowsAffected)))

	return nil
}
