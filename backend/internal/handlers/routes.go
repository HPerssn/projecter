package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

type routeResponse struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Grade string         `json:"grade"`
	Holds []holdResponse `json:"holds"`
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
	id := chi.URLParam(r, "id")

	var route routeResponse
	err := h.db.QueryRowContext(r.Context(),
		`SELECT id, name, grade FROM routes WHERE id = $1`, id,
	).Scan(&route.ID, &route.Name, &route.Grade)
	if err != nil {
		http.Error(w, "route not found", http.StatusNotFound)
		return
	}

	rows, err := h.db.QueryContext(r.Context(),
		`SELECT id, X_pct, Y_pct, COALESCE(note, ''), seq_order
		 FROM holds WHERE route_id = $1 ORDER BY seq_order`, id,
	)
	if err != nil {
		http.Error(w, "failed to fetch holds", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	route.Holds = []holdResponse{}
	for rows.Next() {
		var hold holdResponse
		if err := rows.Scan(&hold.ID, &hold.XPct, &hold.YPct, &hold.Note, &hold.SeqOrder); err != nil {
			http.Error(w, "failed to scan hold", http.StatusInternalServerError)
			return
		}
		route.Holds = append(route.Holds, hold)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(route)
}
