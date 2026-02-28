package faculties

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Get("/", ListFacultiesCRUD)
	r.Get("/{id}", GetFacultyCRUD)
	r.Post("/", CreateFacultyCRUD)
	r.Put("/{id}", UpdateFacultyCRUD)
	r.Delete("/{id}", DeleteFacultyCRUD)
}

// -------------------- GET --------------------

func ListFacultiesCRUD(w http.ResponseWriter, r *http.Request) {
	faculties, err := GetAllFaculties()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(faculties)
}

// -------------------- GET ONE --------------------

func GetFacultyCRUD(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	faculty, err := GetFacultyByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(faculty)
}

// -------------------- POST --------------------

func CreateFacultyCRUD(w http.ResponseWriter, r *http.Request) {
	var f Faculty

	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if f.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if err := CreateFaculty(&f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f)
}

// -------------------- PUT --------------------

func UpdateFacultyCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var f Faculty
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	f.ID = id

	if err := UpdateFaculty(&f); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(f)
}

// -------------------- DELETE --------------------

func DeleteFacultyCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := DeleteFaculty(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
