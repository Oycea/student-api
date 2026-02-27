package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
