package students

import (
	"encoding/json"
	"net/http"
	"strconv"
	"student-api/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	storage *postgres.Storage
}

func NewHandler(storage *postgres.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.ListStudentsCRUD)
	r.Post("/", h.CreateStudentCRUD)
	r.Get("/{id}", h.GetStudentCRUD)
	r.Put("/{id}", h.UpdateStudentCRUD)
	r.Post("/{id}/archive", h.ArchiveStudentCRUD)
	r.Delete("/{id}", h.DeleteStudentCRUD)
}

// GET /api/v1/students
func (h *Handler) ListStudentsCRUD(w http.ResponseWriter, r *http.Request) {
	rows, err := h.storage.DB().Query(`
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
func (h *Handler) GetStudentCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	s, err := GetStudentByID(h.storage, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(s)
}

// POST /api/v1/students
func (h *Handler) CreateStudentCRUD(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ok, _ := IsTicketUnique(h.storage, s.StudentTicket)
	if !ok {
		http.Error(w, "student_ticket must be unique", http.StatusBadRequest)
		return
	}

	if err := CreateStudent(h.storage, &s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

// PUT /api/v1/students/{id}
func (h *Handler) UpdateStudentCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	s.ID = id

	if err := UpdateStudent(h.storage, &s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(s)
}

// POST /api/v1/students/{id}/archive
func (h *Handler) ArchiveStudentCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := ArchiveStudent(h.storage, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DELETE /api/v1/students/{id}
func (h *Handler) DeleteStudentCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := DeleteStudent(h.storage, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
