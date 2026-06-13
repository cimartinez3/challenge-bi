package handler

import (
	"net/http"

	"github.com/carlosmartinez/challenge-bi/internal/domain"
	"github.com/carlosmartinez/challenge-bi/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID string `json:"customer_id"`
		Type       string `json:"type"`
	}
	if err := parseJSON(r, &req); err != nil {
		writeError(w, errBadRequest("invalid request body"))
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		writeError(w, errBadRequest("invalid customer_id"))
		return
	}

	accType := domain.AccountType(req.Type)
	if accType != domain.TypeSavings && accType != domain.TypeChecking {
		writeError(w, errBadRequest("type must be 'savings' or 'checking'"))
		return
	}

	account, err := h.svc.Create(r.Context(), customerID, accType)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, account)
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, errBadRequest("invalid account id"))
		return
	}

	account, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, account)
}

func (h *AccountHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, errBadRequest("invalid account id"))
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := parseJSON(r, &req); err != nil {
		writeError(w, errBadRequest("invalid request body"))
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), id, domain.AccountStatus(req.Status)); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": req.Status})
}
