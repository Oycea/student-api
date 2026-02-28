package students

import (
	"encoding/json"
	"net/http"
	"strconv"
	"student-api/internal/storage"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Get("/", ListStudentsCRUD)
	r.Post("/", CreateStudentCRUD)
	r.Get("/{id}", GetStudentCRUD)
	r.Put("/{id}", UpdateStudentCRUD)
	r.Post("/{id}/archive", ArchiveStudentCRUD)
	r.Delete("/{id}", DeleteStudentCRUD)
}

// GET /api/v1/students
func ListStudentsCRUD(w http.ResponseWriter, r *http.Request) {
	rows, err := storage.DB.Query(`
		SELECT id, student_ticket, first_name, last_name, middle_name, department_id, admission_year, degree, thesis_title, graduation_year, graduated, grade, archived, created_at, updated_at 
		FROM students 
		WHERE archived=false`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		var s Student
		rows.Scan(
			&s.ID, &s.StudentTicket, &s.FirstName, &s.LastName, &s.MiddleName,
			&s.DepartmentID, &s.AdmissionYear, &s.Degree, &s.ThesisTitle,
			&s.GraduationYear, &s.Graduated, &s.Grade, &s.Archived, &s.CreatedAt, &s.UpdatedAt,
		)
		students = append(students, s)
	}

	json.NewEncoder(w).Encode(students)
}

// GET /api/v1/students/{id}
func GetStudentCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	s, err := GetStudentByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(s)
}

// POST /api/v1/students
func CreateStudentCRUD(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ok, _ := IsTicketUnique(s.StudentTicket)
	if !ok {
		http.Error(w, "student_ticket must be unique", http.StatusBadRequest)
		return
	}

	if err := CreateStudent(&s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

// PUT /api/v1/students/{id}
func UpdateStudentCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	s.ID = id

	if err := UpdateStudent(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(s)
}

// POST /api/v1/students/{id}/archive
func ArchiveStudentCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := ArchiveStudent(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DELETE /api/v1/students/{id}
func DeleteStudentCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := DeleteStudent(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
