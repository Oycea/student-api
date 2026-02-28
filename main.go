package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"student-api/internal/config"
	"student-api/internal/storage/postgres"
	"student-api/internal/users"
	"time"

	"student-api/internal/students"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("secret_key")

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

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()

	// Подключение БД
	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Fatal("Failed to init storage")
		os.Exit(1)
	}

	_ = storage

	r := chi.NewRouter()

	r.Post("/api/v1/auth/login", loginHandler)
	r.Post("/api/v1/auth/logout", logoutHandler)

	r.Route("/api/v1/students", func(r chi.Router) {
		r.Use(jwtMiddleware)
		students.RegisterRoutes(r) // регистрируем CRUD-эндпойнты
	})

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(jwtMiddleware)
		users.RegisterRoutes(r)
	})

	log.Println("Server started on :8080!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
