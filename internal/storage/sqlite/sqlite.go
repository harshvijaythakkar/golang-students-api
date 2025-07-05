package sqlite

import (
	"database/sql"

	"github.com/harshvijaythakkar/golang-students-api/internal/config"
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
