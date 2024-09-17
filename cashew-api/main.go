package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/brianykl/cashew/cashew-api/db"
	"github.com/brianykl/cashew/cashew-api/handlers"
	"github.com/brianykl/cashew/cashew-api/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	dbHost := "localhost"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	sslMode := "disable"

	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbName, sslMode)
	transactionManager, err := db.NewTransactionManager(connString)
	handlers.TransactionManager = transactionManager

	tokenManager, err := db.NewTokenManager("localhost:6379")
	handlers.TokenManager = tokenManager

	http.HandleFunc("/link", handlers.LinkHandler)
	http.Handle("/protected/exchange", middleware.EnsureValidToken()(http.HandlerFunc(handlers.ExchangeHandler)))
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
