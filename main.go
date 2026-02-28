package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("secret_key")

func main() {
	r := chi.NewRouter()

	r.Post("/api/v1/auth/login", loginHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Username != "admin" || req.Password != "admin" {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString(secret)

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
