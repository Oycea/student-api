package main

import (
	"log"
	"net/http"
	"student-api/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	storage.Init()

	log.Println("Server started on :8080!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
