package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"wallet/internal/model"
	"wallet/internal/service"
)

type WalletHandler struct {
	service *service.WalletService
}

func NewWalletHandler(service *service.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	wallet := model.NewWallet()
	var req model.CreateWalletRequest
	req.WalletID = wallet.WalletID

	if err := h.service.CreateWallet(req.WalletID); err != nil {
		HandleError(w, err)
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]any{
		"walletID": req.WalletID,
	})
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	walletID := chi.URLParam(r, "WALLET_UUID")

	if !isValidUUID(walletID) {
		HandleError(w, fmt.Errorf("invalid wallet ID format"))
		return
	}

	balance, err := h.service.GetBalance(model.UUID(walletID))
	if err != nil {
		HandleError(w, err)
		return
	}

	response := model.BalanceResponse{
		WalletID: model.UUID(walletID),
		Balance:  balance,
	}

	sendJSONResponse(w, http.StatusOK, response)
}

func (h *WalletHandler) HandleOperation(w http.ResponseWriter, r *http.Request) {
	var req model.OperationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !isValidUUID(req.WalletID) {
		HandleError(w, fmt.Errorf("invalid wallet ID format"))
		return
	}

	walletID := model.UUID(req.WalletID)

	var err error
	switch req.OperationType {
	case model.DEPOSIT:
		err = h.service.Deposit(walletID, req.Amount)
	case model.WITHDRAW:
		err = h.service.Withdraw(walletID, req.Amount)
	default:
		http.Error(w, "Invalid operation type", http.StatusBadRequest)
		return
	}

	if err != nil {
		HandleError(w, err)
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]interface{}{"status": "success"})
}
