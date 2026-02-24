package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hperssn/projecter/internal/db"
	"github.com/hperssn/projecter/internal/handlers"
)

func main() {
	database := db.Connect("postgresql://dev:dev@localhost:5432/projecter?sslmode=disable")
	defer database.Close()

	h := handlers.New(database)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Route("/api/routes", func(r chi.Router) {
		r.Post("/", h.CreateRoute)
		r.Get("/{id}", h.GetRoute)
		r.Post("/{id}/holds", h.SaveHolds)
	})

	log.Println("listening on :8080")
	http.ListenAndServe(":8080", r)
}
