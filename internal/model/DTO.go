package model

type (
	CreateWalletRequest struct {
		WalletID UUID `json:"wallet_id"`
	}

	OperationRequest struct {
		WalletID      string    `json:"valletId"`
		Amount        float64   `json:"amount"`
		OperationType Operation `json:"operationType"`
	}

	BalanceResponse struct {
		WalletID UUID    `json:"wallet_id"`
		Balance  float64 `json:"balance"`
	}
)
