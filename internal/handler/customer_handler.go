package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/carlosmartinez/challenge-bi/internal/service"
)

type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	if err := parseJSON(r, &req); err != nil {
		writeError(w, errBadRequest("invalid request body"))
		return
	}
	if req.FullName == "" || req.Email == "" {
		writeError(w, errBadRequest("full_name and email are required"))
		return
	}

	customer, err := h.svc.Create(r.Context(), req.FullName, req.Email)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, errBadRequest("invalid customer id"))
		return
	}

	customer, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, customer)
}
