package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"wallet/internal/model"
)

var ErrWalletNotFound = errors.New("Кошелёк не найден")

type DBWallet struct {
	db *sql.DB
}

func NewDBWallet(db *sql.DB) *DBWallet {
	return &DBWallet{db: db}
}

func (w *DBWallet) CreateWallet(walletID model.UUID) error {
	_, err := w.db.Exec(`
	INSERT INTO wallets(wallet_id, balance)
	VALUES($1, 0.0)`, walletID)

	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (w *DBWallet) GetBalance(walletID model.UUID) (float64, error) {
	var balance float64

	err := w.db.QueryRow(`
        SELECT balance 
        FROM wallets 
        WHERE wallet_id = $1`,
		walletID,
	).Scan(&balance)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrWalletNotFound
	}

	return balance, err
}
