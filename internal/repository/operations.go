package repository

import (
	"database/sql"
	"wallet/internal/model"
)

type postgresTX struct {
	tx *sql.Tx
}

func (r *DBWallet) BeginTx() (WalletTx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	return &postgresTX{tx: tx}, nil
}

func (w *postgresTX) Deposit(walletID model.UUID, amount float64) error {

	_, err := w.tx.Exec(`
        UPDATE wallets 
        SET balance = balance + $1 
        WHERE wallet_id = $2`,
		amount, walletID,
	)
	if err != nil {
		return err
	}

	_, err = w.tx.Exec(`
        INSERT INTO transactions (wallet_id, operation_type, amount, created_at)
        VALUES ($1, 'DEPOSIT', $2, NOW())`,
		walletID, amount,
	)
	if err != nil {
		return err
	}

	return nil
}

func (w *postgresTX) Withdraw(walletID model.UUID, amount float64) error {

	_, err := w.tx.Exec(`
        UPDATE wallets 
        SET balance = balance - $1 
        WHERE wallet_id = $2`,
		amount, walletID,
	)
	if err != nil {
		return err
	}

	_, err = w.tx.Exec(`
        INSERT INTO transactions (wallet_id, operation_type, amount, created_at)
        VALUES ($1, 'WITHDRAW', $2, NOW())`,
		walletID, amount,
	)
	if err != nil {
		return err
	}

	return nil
}

func (w *postgresTX) Commit() error {
	return w.tx.Commit()
}

func (w *postgresTX) Rollback() error {
	return w.tx.Rollback()
}
