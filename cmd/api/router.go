package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosmartinez/challenge-bi/internal/handler"
	"github.com/carlosmartinez/challenge-bi/internal/repository"
	"github.com/carlosmartinez/challenge-bi/internal/service"
)

func newRouter(pool *pgxpool.Pool) http.Handler {
	customerRepo := repository.NewCustomerRepo(pool)
	accountRepo := repository.NewAccountRepo(pool)
	txRepo := repository.NewTransactionRepo(pool)

	customerSvc := service.NewCustomerService(customerRepo)
	accountSvc := service.NewAccountService(accountRepo, customerRepo)
	txSvc := service.NewTransactionService(accountRepo, txRepo)
	transferSvc := service.NewTransferService(accountRepo, txRepo)

	customerHandler := handler.NewCustomerHandler(customerSvc)
	accountHandler := handler.NewAccountHandler(accountSvc)
	txHandler := handler.NewTransactionHandler(txSvc, transferSvc)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(handler.LoggerMiddleware)
	r.Use(handler.APIKeyMiddleware)

	// Health endpoint is exempt from auth — required by Kubernetes liveness/readiness probes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Group(func(r chi.Router) {
		r.Post("/customers", customerHandler.Create)
		r.Get("/customers/{id}", customerHandler.Get)

		r.Post("/accounts", accountHandler.Create)
		r.Get("/accounts/{id}", accountHandler.Get)
		r.Patch("/accounts/{id}/status", accountHandler.UpdateStatus)
		r.Get("/accounts/{id}/transactions", txHandler.ListByAccount)

		r.Post("/accounts/{id}/deposit", txHandler.Deposit)
		r.Post("/accounts/{id}/withdrawal", txHandler.Withdrawal)
		r.Post("/transfers", txHandler.Transfer)
		r.Get("/transactions/{reference}", txHandler.GetByReference)
	})

	return r
}
