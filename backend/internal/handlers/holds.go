package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type holdRequest struct {
	XPct     float64 `json:"X_pct"`
	YPct     float64 `json:"Y_pct"`
	Note     string  `json:"note"`
	SeqOrder int     `json:"seq_order"`
}

type holdResponse struct {
	ID       string  `json:"id"`
	XPct     float64 `json:"x_pct"`
	YPct     float64 `json:"y_pct"`
	Note     string  `json:"note"`
	SeqOrder int     `json:"seq_order"`
}

func (h *Handler) SaveHolds(w http.ResponseWriter, r *http.Request) {
	routeID := chi.URLParam(r, "id")

	var holds []holdRequest
	if err := json.NewDecoder(r.Body).Decode(&holds); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var saved []holdResponse
	for _, hold := range holds {
		var id string
		err := h.db.QueryRowContext(r.Context(),
			`INSERT INTO holds (route_id, x_pct, y_pct, note, seq_order)
			 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			routeID, hold.XPct, hold.YPct, hold.Note, hold.SeqOrder,
		).Scan(&id)
		if err != nil {
			http.Error(w, "failed to save hold", http.StatusInternalServerError)
			return
		}
		saved = append(saved, holdResponse{
			ID:       id,
			XPct:     hold.XPct,
			YPct:     hold.YPct,
			Note:     hold.Note,
			SeqOrder: hold.SeqOrder,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saved)
}
