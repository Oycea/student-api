package students

import (
	"database/sql"
	"errors"
	"time"

	"student-api/internal/storage"
)

type Student struct {
	ID             int64           `json:"id"`
	StudentTicket  string          `json:"student_ticket"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	MiddleName     string          `json:"middle_name,omitempty"`
	DepartmentID   int64           `json:"department_id"`
	AdmissionYear  int             `json:"admission_year"`
	Degree         string          `json:"degree"`
	ThesisTitle    string          `json:"thesis_title,omitempty"`
	GraduationYear sql.NullInt32   `json:"graduation_year,omitempty"`
	Graduated      bool            `json:"graduated"`
	Grade          sql.NullFloat64 `json:"grade,omitempty"`
	Archived       bool            `json:"archived"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// Проверка уникальности student_ticket
func IsTicketUnique(ticket string) (bool, error) {
	var exists bool
	err := storage.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM students WHERE student_ticket=$1)", ticket).Scan(&exists)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

// Создание студента
func CreateStudent(s *Student) error {
	query := `
		INSERT INTO students 
		(student_ticket, first_name, last_name, middle_name, department_id, admission_year, degree, thesis_title, graduation_year, graduated, grade, archived, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW())
		RETURNING id, created_at, updated_at
	`
	return storage.DB.QueryRow(query,
		s.StudentTicket,
		s.FirstName,
		s.LastName,
		s.MiddleName,
		s.DepartmentID,
		s.AdmissionYear,
		s.Degree,
		s.ThesisTitle,
		s.GraduationYear,
		s.Graduated,
		s.Grade,
		s.Archived,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

// Получение студента по ID
func GetStudentByID(id int64) (*Student, error) {
	s := &Student{}
	query := `SELECT id, student_ticket, first_name, last_name, middle_name, department_id, admission_year, degree, thesis_title, graduation_year, graduated, grade, archived, created_at, updated_at
			  FROM students WHERE id=$1`
	err := storage.DB.QueryRow(query, id).Scan(
		&s.ID,
		&s.StudentTicket,
		&s.FirstName,
		&s.LastName,
		&s.MiddleName,
		&s.DepartmentID,
		&s.AdmissionYear,
		&s.Degree,
		&s.ThesisTitle,
		&s.GraduationYear,
		&s.Graduated,
		&s.Grade,
		&s.Archived,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("student not found")
	}
	return s, err
}

// Обновление студента
func UpdateStudent(s *Student) error {
	query := `
		UPDATE students SET first_name=$1, last_name=$2, middle_name=$3, department_id=$4, admission_year=$5, degree=$6, thesis_title=$7, graduation_year=$8, graduated=$9, grade=$10, updated_at=NOW()
		WHERE id=$11 AND archived=false
	`
	res, err := storage.DB.Exec(query,
		s.FirstName,
		s.LastName,
		s.MiddleName,
		s.DepartmentID,
		s.AdmissionYear,
		s.Degree,
		s.ThesisTitle,
		s.GraduationYear,
		s.Graduated,
		s.Grade,
		s.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("cannot update archived or non-existent student")
	}
	return nil
}

// Архивация студента
func ArchiveStudent(id int64) error {
	res, err := storage.DB.Exec("UPDATE students SET archived=true, updated_at=NOW() WHERE id=$1", id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("student not found")
	}
	return nil
}

// Удаление студента (soft delete)
func DeleteStudent(id int64) error {
	return ArchiveStudent(id)
}
