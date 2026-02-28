package main

import (
	"log"
	"net/http"
	"os"

	"student-api/internal/auth"
	"student-api/internal/config"
	faculties "student-api/internal/departments"
	"student-api/internal/storage/postgres"
	"student-api/internal/students"
	"student-api/internal/users"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Загрузка конфига
	cfg := config.MustLoad()

	// Подключение БД
	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Fatal("Failed to init storage")
		os.Exit(1)
	}

	r := chi.NewRouter()

	r.Post("/api/v1/auth/login", auth.LoginHandler)
	r.Post("/api/v1/auth/logout", auth.LogoutHandler)

	studentHandler := students.NewHandler(storage)
	r.Route("/api/v1/students", func(r chi.Router) {
		r.Use(auth.JWTMiddleware)
		studentHandler.RegisterRoutes(r)
	})

	userHandler := users.NewHandler(storage)
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Use(auth.JWTMiddleware)
		userHandler.RegisterRoutes(r)
	})

	facultieHandler := faculties.NewHandler(storage)
	r.Route("/api/v1/faculties", func(r chi.Router) {
		r.Use(auth.JWTMiddleware)
		facultieHandler.RegisterRoutes(r)
	})

	log.Println("Server started on :8080!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
