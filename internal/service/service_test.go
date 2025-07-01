package service_test

import (
	"errors"
	"testing"
	"wallet/internal/model"
	"wallet/internal/repository"
	"wallet/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletDB struct {
	mock.Mock
}

func (m *MockWalletDB) CreateWallet(walletID model.UUID) error {
	args := m.Called(walletID)
	return args.Error(0)
}

func (m *MockWalletDB) GetBalance(walletID model.UUID) (float64, error) {
	args := m.Called(walletID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockWalletDB) BeginTx() (repository.WalletTx, error) {
	args := m.Called()
	return args.Get(0).(repository.WalletTx), args.Error(1)
}

type MockWalletTx struct {
	mock.Mock
}

func (m *MockWalletTx) Deposit(walletID model.UUID, amount float64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockWalletTx) Withdraw(walletID model.UUID, amount float64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockWalletTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWalletTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func TestWalletService_GetBalance(t *testing.T) {
	tests := []struct {
		name        string
		walletID    model.UUID
		mockBalance float64
		mockError   error
		wantBalance float64
		wantError   error
	}{
		{
			name:        "success",
			walletID:    "test-wallet",
			mockBalance: 100.0,
			mockError:   nil,
			wantBalance: 100.0,
			wantError:   nil,
		},
		{
			name:        "wallet not found",
			walletID:    "non-existent",
			mockBalance: 0,
			mockError:   repository.ErrWalletNotFound,
			wantBalance: 0,
			wantError:   repository.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockWalletDB)
			mockRepo.On("GetBalance", tt.walletID).Return(tt.mockBalance, tt.mockError)

			service := service.NewWalletService(mockRepo)
			balance, err := service.GetBalance(tt.walletID)

			assert.Equal(t, tt.wantBalance, balance)
			assert.ErrorIs(t, err, tt.wantError)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestWalletService_CreateWallet(t *testing.T) {
	tests := []struct {
		name          string
		walletID      model.UUID
		getBalanceErr error
		createErr     error
		wantError     error
	}{
		{
			name:          "success",
			walletID:      "new-wallet",
			getBalanceErr: repository.ErrWalletNotFound,
			createErr:     nil,
			wantError:     nil,
		},
		{
			name:          "wallet already exists",
			walletID:      "existing-wallet",
			getBalanceErr: nil,
			createErr:     nil,
			wantError:     service.ErrWalletExists,
		},
		{
			name:          "create error",
			walletID:      "error-wallet",
			getBalanceErr: repository.ErrWalletNotFound,
			createErr:     errors.New("db error"),
			wantError:     errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockWalletDB)
			mockRepo.On("GetBalance", tt.walletID).Return(0.0, tt.getBalanceErr)
			if tt.getBalanceErr == repository.ErrWalletNotFound {
				mockRepo.On("CreateWallet", tt.walletID).Return(tt.createErr)
			}

			service := service.NewWalletService(mockRepo)
			err := service.CreateWallet(tt.walletID)

			if tt.wantError != nil {
				assert.ErrorContains(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestWalletService_Deposit(t *testing.T) {
	tests := []struct {
		name           string
		walletID       model.UUID
		amount         float64
		txError        error
		depositError   error
		commitError    error
		wantError      error
		expectRollback bool
	}{
		{
			name:         "success",
			walletID:     "wallet-1",
			amount:       50.0,
			txError:      nil,
			depositError: nil,
			commitError:  nil,
			wantError:    nil,
		},
		{
			name:      "invalid amount",
			walletID:  "wallet-1",
			amount:    0,
			wantError: service.ErrInvalidAmount,
		},
		{
			name:      "begin tx error",
			walletID:  "wallet-1",
			amount:    50.0,
			txError:   errors.New("tx error"),
			wantError: errors.New("failed to begin transaction"),
		},
		{
			name:           "deposit error",
			walletID:       "wallet-1",
			amount:         50.0,
			txError:        nil,
			depositError:   errors.New("deposit error"),
			wantError:      errors.New("failed to deposit"),
			expectRollback: true,
		},
		{
			name:           "commit error",
			walletID:       "wallet-1",
			amount:         50.0,
			txError:        nil,
			depositError:   nil,
			commitError:    errors.New("commit error"),
			wantError:      errors.New("failed to commit transaction"),
			expectRollback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockWalletDB)
			mockTx := new(MockWalletTx)

			if tt.amount > 0 {
				mockRepo.On("BeginTx").Return(mockTx, tt.txError)

				if tt.txError == nil {
					mockTx.On("Deposit", tt.walletID, tt.amount).Return(tt.depositError)

					if tt.depositError == nil {
						mockTx.On("Commit").Return(tt.commitError)
					}

					if tt.expectRollback {
						mockTx.On("Rollback").Return(nil)
					}
				}
			}

			service := service.NewWalletService(mockRepo)
			err := service.Deposit(tt.walletID, tt.amount)

			if tt.wantError != nil {
				assert.ErrorContains(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
			mockTx.AssertExpectations(t)
		})
	}
}

func TestWalletService_Withdraw(t *testing.T) {
	tests := []struct {
		name           string
		walletID       model.UUID
		amount         float64
		balance        float64
		getBalanceErr  error
		txError        error
		withdrawError  error
		commitError    error
		wantError      error
		expectRollback bool
	}{
		{
			name:          "success",
			walletID:      "wallet-1",
			amount:        30.0,
			balance:       100.0,
			getBalanceErr: nil,
			txError:       nil,
			withdrawError: nil,
			commitError:   nil,
			wantError:     nil,
		},
		{
			name:      "invalid amount",
			walletID:  "wallet-1",
			amount:    0,
			wantError: service.ErrInvalidAmount,
		},
		{
			name:          "wallet not found",
			walletID:      "non-existent",
			amount:        30.0,
			getBalanceErr: repository.ErrWalletNotFound,
			wantError:     repository.ErrWalletNotFound,
		},
		{
			name:          "insufficient funds",
			walletID:      "wallet-1",
			amount:        30.0,
			balance:       20.0,
			getBalanceErr: nil,
			wantError:     service.ErrInsufficientFunds,
		},
		{
			name:          "begin tx error",
			walletID:      "wallet-1",
			amount:        30.0,
			balance:       100.0,
			getBalanceErr: nil,
			txError:       errors.New("tx error"),
			wantError:     errors.New("failed to begin transaction"),
		},
		{
			name:           "withdraw error",
			walletID:       "wallet-1",
			amount:         30.0,
			balance:        100.0,
			getBalanceErr:  nil,
			txError:        nil,
			withdrawError:  errors.New("withdraw error"),
			wantError:      errors.New("failed to withdraw"),
			expectRollback: true,
		},
		{
			name:           "commit error",
			walletID:       "wallet-1",
			amount:         30.0,
			balance:        100.0,
			getBalanceErr:  nil,
			txError:        nil,
			withdrawError:  nil,
			commitError:    errors.New("commit error"),
			wantError:      errors.New("failed to commit transaction"),
			expectRollback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockWalletDB)
			mockTx := new(MockWalletTx)

			if tt.amount > 0 {
				mockRepo.On("GetBalance", tt.walletID).Return(tt.balance, tt.getBalanceErr)

				if tt.getBalanceErr == nil && tt.balance >= tt.amount {
					mockRepo.On("BeginTx").Return(mockTx, tt.txError)

					if tt.txError == nil {
						mockTx.On("Withdraw", tt.walletID, tt.amount).Return(tt.withdrawError)

						if tt.withdrawError == nil {
							mockTx.On("Commit").Return(tt.commitError)
						}

						if tt.expectRollback {
							mockTx.On("Rollback").Return(nil)
						}
					}
				}
			}

			service := service.NewWalletService(mockRepo)
			err := service.Withdraw(tt.walletID, tt.amount)

			if tt.wantError != nil {
				assert.ErrorContains(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
			mockTx.AssertExpectations(t)
		})
	}
}
