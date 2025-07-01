package service

import (
	"errors"
	"fmt"
	"wallet/internal/model"
	"wallet/internal/repository"
)

var (
	ErrWalletExists      = errors.New("кошелёк уже существует")
	ErrInvalidAmount     = errors.New("неверная сумма")
	ErrInsufficientFunds = errors.New("недостаточно средств")
)

type WalletService struct {
	repo repository.WalletDB
}

func NewWalletService(repo repository.WalletDB) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetBalance(walletID model.UUID) (float64, error) {
	balance, err := s.repo.GetBalance(walletID)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (s *WalletService) CreateWallet(walletID model.UUID) error {

	_, err := s.GetBalance(walletID)
	if err == nil {
		return ErrWalletExists
	}

	return s.repo.CreateWallet(walletID)
}

func (s *WalletService) Withdraw(walletID model.UUID, amount float64) error {

	if amount <= 0 {
		return ErrInvalidAmount
	}

	balance, err := s.GetBalance(walletID)
	if err != nil {
		return err
	}

	if balance < amount {
		return ErrInsufficientFunds
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Withdraw(walletID, amount); err != nil {
		return fmt.Errorf("failed to withdraw: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *WalletService) Deposit(walletID model.UUID, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Deposit(walletID, amount); err != nil {
		return fmt.Errorf("failed to deposit: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
