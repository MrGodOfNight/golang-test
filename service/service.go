package service

import (
	"context"
	"errors"
	"testing-go/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Deposit(ctx context.Context, userID int, amount float64) error {
	tx, err := s.repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.UpdateUserBalance(ctx, userID, amount); err != nil {
		return err
	}

	if err := s.repo.CreateTransaction(ctx, userID, amount, "deposit"); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UserService) Transfer(ctx context.Context, senderID, receiverID int, amount float64) error {
	tx, err := s.repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	senderBalance, err := s.repo.GetUserBalance(ctx, senderID)
	if err != nil {
		return err
	}

	if senderBalance < amount {
		return errors.New("insufficient funds")
	}

	if err := s.repo.UpdateUserBalance(ctx, senderID, -amount); err != nil {
		return err
	}

	if err := s.repo.UpdateUserBalance(ctx, receiverID, amount); err != nil {
		return err
	}

	if err := s.repo.CreateTransaction(ctx, senderID, -amount, "transfer_out"); err != nil {
		return err
	}

	if err := s.repo.CreateTransaction(ctx, receiverID, amount, "transfer_in"); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UserService) GetLastTransactions(ctx context.Context, userID int, limit int) ([]map[string]interface{}, error) {
	return s.repo.GetLastTransactions(ctx, userID, limit)
}
