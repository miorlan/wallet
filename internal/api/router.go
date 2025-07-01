package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"net/http"
)

func (h *WalletHandler) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.AllowContentType("application/json"),
		RateLimit,
	)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/wallet/create", h.CreateWallet)
		r.Get("/wallets/{WALLET_UUID}", h.GetBalance)

		r.Post("/wallet", h.HandleOperation)
	})

	return r
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
