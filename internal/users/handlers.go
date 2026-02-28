package users

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
	r.Get("/", h.ListUsersCRUD)
	r.Get("/{id}", h.GetUserCRUD)
	r.Post("/", h.CreateUserCRUD)
	r.Put("/{id}", h.UpdateUserCRUD)
	r.Delete("/{id}", h.DeleteUserCRUD)
}

// -------------------- GET --------------------

func (h *Handler) ListUsersCRUD(w http.ResponseWriter, r *http.Request) {
	users, err := GetAllUsers(h.storage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// -------------------- GET ONE --------------------

func (h *Handler) GetUserCRUD(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := GetUserByID(h.storage, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// -------------------- POST --------------------

func (h *Handler) CreateUserCRUD(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := CreateUser(h.storage, &u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// -------------------- PUT --------------------

func (h *Handler) UpdateUserCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	u.ID = id

	if err := UpdateUser(h.storage, &u); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// -------------------- DELETE --------------------

func (h *Handler) DeleteUserCRUD(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := DeleteUser(h.storage, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
