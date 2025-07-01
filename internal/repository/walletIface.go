package repository

import "wallet/internal/model"

type WalletDB interface {
	CreateWallet(walletID model.UUID) error
	GetBalance(walletID model.UUID) (float64, error)
	BeginTx() (WalletTx, error)
}

type WalletTx interface {
	Deposit(walletID model.UUID, amount float64) error
	Withdraw(walletID model.UUID, amount float64) error
	Commit() error
	Rollback() error
}
