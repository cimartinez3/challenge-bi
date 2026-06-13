package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/carlosmartinez/challenge-bi/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionHandler struct {
	txSvc       *service.TransactionService
	transferSvc *service.TransferService
}

func NewTransactionHandler(txSvc *service.TransactionService, transferSvc *service.TransferService) *TransactionHandler {
	return &TransactionHandler{txSvc: txSvc, transferSvc: transferSvc}
}

func (h *TransactionHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, errBadRequest("invalid account id"))
		return
	}

	var req struct {
		Amount    string `json:"amount"`
		Reference string `json:"reference"`
	}
	if err := parseJSON(r, &req); err != nil {
		writeError(w, errBadRequest("invalid request body"))
		return
	}
	if req.Reference == "" {
		writeError(w, errBadRequest("reference is required"))
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		writeError(w, errBadRequest("invalid amount format"))
		return
	}

	tx, err := h.txSvc.Deposit(r.Context(), accountID, amount, req.Reference)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, tx)
}

func (h *TransactionHandler) Withdrawal(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, errBadRequest("invalid account id"))
		return
	}

	var req struct {
		Amount    string `json:"amount"`
		Reference string `json:"reference"`
	}
	if err := parseJSON(r, &req); err != nil {
		writeError(w, errBadRequest("invalid request body"))
		return
	}
	if req.Reference == "" {
		writeError(w, errBadRequest("reference is required"))
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		writeError(w, errBadRequest("invalid amount format"))
		return
	}

	tx, err := h.txSvc.Withdraw(r.Context(), accountID, amount, req.Reference)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, tx)
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromAccountID string `json:"from_account_id"`
		ToAccountID   string `json:"to_account_id"`
		Amount        string `json:"amount"`
		Reference     string `json:"reference"`
	}
	if err := parseJSON(r, &req); err != nil {
		writeError(w, errBadRequest("invalid request body"))
		return
	}
	if req.Reference == "" {
		writeError(w, errBadRequest("reference is required"))
		return
	}

	fromID, err := uuid.Parse(req.FromAccountID)
	if err != nil {
		writeError(w, errBadRequest("invalid from_account_id"))
		return
	}
	toID, err := uuid.Parse(req.ToAccountID)
	if err != nil {
		writeError(w, errBadRequest("invalid to_account_id"))
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		writeError(w, errBadRequest("invalid amount format"))
		return
	}

	result, err := h.transferSvc.Transfer(r.Context(), fromID, toID, amount, req.Reference)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *TransactionHandler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	accountID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, errBadRequest("invalid account id"))
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	var before *time.Time
	if b := r.URL.Query().Get("before"); b != "" {
		t, err := time.Parse(time.RFC3339, b)
		if err != nil {
			writeError(w, errBadRequest("invalid before format, use RFC3339"))
			return
		}
		before = &t
	}

	txs, err := h.txSvc.ListByAccount(r.Context(), accountID, limit, before)
	if err != nil {
		writeError(w, err)
		return
	}

	type response struct {
		Data    any  `json:"data"`
		HasMore bool `json:"has_more"`
	}
	writeJSON(w, http.StatusOK, response{
		Data:    txs,
		HasMore: len(txs) == limit,
	})
}

func (h *TransactionHandler) GetByReference(w http.ResponseWriter, r *http.Request) {
	ref := chi.URLParam(r, "reference")
	if ref == "" {
		writeError(w, errBadRequest("reference is required"))
		return
	}

	tx, err := h.txSvc.GetByReference(r.Context(), ref)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tx)
}
