package faculties

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
	r.Get("/", h.ListFacultiesCRUD)
	r.Get("/{id}", h.GetFacultyCRUD)
	r.Post("/", h.CreateFacultyCRUD)
	r.Put("/{id}", h.UpdateFacultyCRUD)
	r.Delete("/{id}", h.DeleteFacultyCRUD)
}

// -------------------- GET --------------------

func (h *Handler) ListFacultiesCRUD(w http.ResponseWriter, r *http.Request) {
	faculties, err := GetAllFaculties(h.storage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(faculties)
}

// -------------------- GET ONE --------------------

func (h *Handler) GetFacultyCRUD(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	faculty, err := GetFacultyByID(h.storage, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(faculty)
}

// -------------------- POST --------------------

func (h *Handler) CreateFacultyCRUD(w http.ResponseWriter, r *http.Request) {
	var f Faculty

	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if f.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if err := CreateFaculty(h.storage, &f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

// -------------------- PUT --------------------

func (h *Handler) UpdateFacultyCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var f Faculty
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	f.ID = id

	if err := UpdateFaculty(h.storage, &f); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(f)
}

// -------------------- DELETE --------------------

func (h *Handler) DeleteFacultyCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := DeleteFaculty(h.storage, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
