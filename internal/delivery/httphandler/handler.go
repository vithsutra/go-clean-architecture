package httphandler

import (
	"encoding/json"
	"id-generator/internal/app"

	"net/http"
)

type Handler struct {
	service *app.IDService
}

func NewHandler(service *app.IDService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetID(w http.ResponseWriter, r *http.Request) {
	id, err := h.service.GetID()
	if err != nil {
		http.Error(w, "failed to generate ID", http.StatusInternalServerError)
		return
	}
	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
