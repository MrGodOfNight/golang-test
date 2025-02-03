package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"testing-go/handler"
	"testing-go/repository"
	"testing-go/service"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// Download .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)

	r := gin.Default()
	r.POST("/users/:id/deposit", handler.Deposit)
	r.POST("/users/transfer", handler.Transfer)
	r.GET("/users/:id/transactions", handler.GetTransactions)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	log.Fatal(r.Run(":" + port))
}
