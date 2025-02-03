package service_test

import (
	"context"
	"testing"

	"testing-go/repository"
	"testing-go/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Инициализируем репозиторий и сервис
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)

	// Тестовые данные
	userID := 1
	amount := 100.0

	// Ожидаемые SQL-запросы
	mock.ExpectBegin() // Начало транзакции
	mock.ExpectExec("UPDATE users SET balance = balance").
		WithArgs(amount, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactions").
		WithArgs(userID, amount, "deposit").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit() // Фиксация транзакции

	// Выполняем метод Deposit
	err = service.Deposit(context.Background(), userID, amount)
	assert.NoError(t, err)

	// Проверяем, что все ожидания выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestTransfer(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Инициализируем репозиторий и сервис
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)

	// Тестовые данные
	senderID := 1
	receiverID := 2
	amount := 50.0
	senderBalance := 100.0

	// Ожидаемые SQL-запросы
	mock.ExpectBegin() // Добавляем ожидание вызова Begin
	mock.ExpectQuery("SELECT balance FROM users WHERE id = \\$1").
		WithArgs(senderID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(senderBalance))

	mock.ExpectExec("UPDATE users SET balance = balance").
		WithArgs(-amount, senderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE users SET balance = balance").
		WithArgs(amount, receiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactions").
		WithArgs(senderID, -amount, "transfer_out").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO transactions").
		WithArgs(receiverID, amount, "transfer_in").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit() // Фиксация транзакции

	// Выполняем метод Transfer
	err = service.Transfer(context.Background(), senderID, receiverID, amount)
	assert.NoError(t, err)

	// Проверяем, что все ожидания выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Инициализируем репозиторий и сервис
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)

	// Тестовые данные
	senderID := 1
	receiverID := 2
	amount := 150.0
	senderBalance := 100.0

	// Ожидаемые SQL-запросы
	mock.ExpectBegin() // Добавляем ожидание вызова Begin
	mock.ExpectQuery("SELECT balance FROM users WHERE id = \\$1").
		WithArgs(senderID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(senderBalance))

	// Выполняем метод Transfer
	err = service.Transfer(context.Background(), senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())

	// Проверяем, что все ожидания выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
func TestGetLastTransactions(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Инициализируем репозиторий и сервис
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)

	// Тестовые данные
	userID := 1
	limit := 2

	// Ожидаемый SQL-запрос
	rows := sqlmock.NewRows([]string{"amount", "type", "created_at"}).
		AddRow(100.0, "deposit", "2023-10-01T12:00:00Z").
		AddRow(-50.0, "transfer_out", "2023-10-02T12:00:00Z")

	mock.ExpectQuery("SELECT amount, type, created_at FROM transactions WHERE user_id = \\$1 ORDER BY created_at DESC LIMIT \\$2").
		WithArgs(userID, limit).
		WillReturnRows(rows)

	// Выполняем метод GetLastTransactions
	transactions, err := service.GetLastTransactions(context.Background(), userID, limit)
	assert.NoError(t, err)

	// Проверяем результат
	expected := []map[string]interface{}{
		{"amount": 100.0, "type": "deposit", "created_at": "2023-10-01T12:00:00Z"},
		{"amount": -50.0, "type": "transfer_out", "created_at": "2023-10-02T12:00:00Z"},
	}
	assert.Equal(t, expected, transactions)

	// Проверяем, что все ожидания выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
