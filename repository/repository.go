package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{Db: db}
}

func (r *UserRepository) GetUserBalance(ctx context.Context, userID int) (float64, error) {
	var balance float64
	query := "SELECT balance FROM users WHERE id = $1"
	err := r.Db.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get user balance: %w", err)
	}
	return balance, nil
}

func (r *UserRepository) UpdateUserBalance(ctx context.Context, userID int, amount float64) error {
	query := "UPDATE users SET balance = balance + $1 WHERE id = $2"
	_, err := r.Db.ExecContext(ctx, query, amount, userID)
	return err
}

func (r *UserRepository) CreateTransaction(ctx context.Context, userID int, amount float64, transactionType string) error {
	query := "INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, $3)"
	_, err := r.Db.ExecContext(ctx, query, userID, amount, transactionType)
	return err
}

func (r *UserRepository) GetLastTransactions(ctx context.Context, userID int, limit int) ([]map[string]interface{}, error) {
	query := `
        SELECT amount, type, created_at 
        FROM transactions 
        WHERE user_id = $1 
        ORDER BY created_at DESC 
        LIMIT $2
    `
	rows, err := r.Db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []map[string]interface{}
	for rows.Next() {
		var amount float64
		var transactionType string
		var createdAt string
		if err := rows.Scan(&amount, &transactionType, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, map[string]interface{}{
			"amount":     amount,
			"type":       transactionType,
			"created_at": createdAt,
		})
	}
	return transactions, nil
}
