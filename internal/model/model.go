package model

import "github.com/google/uuid"

type Wallet struct {
	WalletID UUID
	Balance  float64
}

type UUID string

type Operation string

const (
	DEPOSIT  Operation = "DEPOSIT"
	WITHDRAW Operation = "WITHDRAW"
)

func NewWallet() *Wallet {
	return &Wallet{
		WalletID: UUID(uuid.New().String()),
		Balance:  0.0,
	}
}
