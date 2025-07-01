package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sajimenezher_meli/meli-frescos-8/internal/services"
)

type SellerHandler struct {
	service services.SellerService
}

func NewSellerHandler(service services.SellerService) *SellerHandler {
	return &SellerHandler{service: service}
}

func (h *SellerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	sellers, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "error getting sellers", http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(sellers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
