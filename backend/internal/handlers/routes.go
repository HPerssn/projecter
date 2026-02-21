package handlers

import (
	"encoding/json"
	"net/http"
)

type createRouteRequest struct {
	Name  string `json:"name"`
	Grade string `json:"grade"`
}

type createRouteResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Grade string `json:"grade"`
}

func (h *Handler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var req createRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var id string
	err := h.db.QueryRowContext(r.Context(),
		`INSERT INTO routes (name, grade) VALUES ($1, $2) RETURNING id`,
		req.Name, req.Grade,
	).Scan(&id)
	if err != nil {
		http.Error(w, "failed to create route", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createRouteResponse{
		ID:    id,
		Name:  req.Name,
		Grade: req.Grade,
	})
}

func (h *Handler) GetRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"ok": true})
}
