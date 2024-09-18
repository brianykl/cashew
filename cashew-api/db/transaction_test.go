package db

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

func TestDatabaseOperations(t *testing.T) {
	// Load environment variables
	err := godotenv.Load("../.env.local")
	if err != nil {
		t.Fatalf("Error loading .env.local file: %v", err)
	}

	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Log connection details (be careful with logging passwords in production)
	t.Logf("DB_HOST: %s", dbHost)
	t.Logf("DB_USER: %s", dbUser)
	t.Logf("DB_NAME: %s", dbName)
	t.Logf("DB_SSLMODE: %s", dbSSLMode)

	// Construct the connection string
	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbName, dbSSLMode)

	// Create the TransactionManager
	tm, err := NewTransactionManager(connString)
	if err != nil {
		t.Fatalf("Failed to create TransactionManager: %v", err)
	}

	// Create a test transaction
	testTransaction := &Transaction{
		UserId:           "testuser",
		AccountId:        "testacc",
		AccountName:      "Test Account",
		Amount:           decimal.NewFromFloat(100.50),
		Currency:         "USD",
		AuthorizedDate:   time.Now().UTC().Truncate(time.Second), // Truncate to remove sub-second precision
		MerchantName:     "Test Merchant",
		PaymentChannel:   "online",
		PrimaryCategory:  "Shopping",
		DetailedCategory: "Electronics",
		ConfidenceLevel:  "high",
	}

	// Test storing the transaction
	ctx := context.Background()
	err = tm.StoreTransactions(ctx, []*Transaction{testTransaction})
	if err != nil {
		t.Fatalf("Failed to store test transaction: %v", err)
	}

	// Test retrieving transactions
	startDate := testTransaction.AuthorizedDate.AddDate(0, 0, -1) // 1 day before
	endDate := testTransaction.AuthorizedDate.AddDate(0, 0, 1)    // 1 day after
	transactions, err := tm.GetTransactions(ctx, "testuser", &startDate, &endDate, 10, 0)
	if err != nil {
		t.Fatalf("Failed to retrieve transactions: %v", err)
	}

	// Check if we retrieved the transaction we just stored
	if len(transactions) == 0 {
		t.Fatalf("No transactions retrieved")
	}

	retrievedTransaction := transactions[0]

	// Detailed comparison
	t.Logf("Stored Transaction: %+v", testTransaction)
	t.Logf("Retrieved Transaction: %+v", retrievedTransaction)

	if retrievedTransaction.UserId != testTransaction.UserId {
		t.Errorf("UserId mismatch: got %s, want %s", retrievedTransaction.UserId, testTransaction.UserId)
	}
	if retrievedTransaction.AccountId != testTransaction.AccountId {
		t.Errorf("AccountId mismatch: got %s, want %s", retrievedTransaction.AccountId, testTransaction.AccountId)
	}
	if !retrievedTransaction.Amount.Equal(testTransaction.Amount) {
		t.Errorf("Amount mismatch: got %s, want %s", retrievedTransaction.Amount, testTransaction.Amount)
	}
	if !retrievedTransaction.AuthorizedDate.Equal(testTransaction.AuthorizedDate) {
		t.Errorf("AuthorizedDate mismatch: got %v, want %v", retrievedTransaction.AuthorizedDate, testTransaction.AuthorizedDate)
	}
	if retrievedTransaction.MerchantName != testTransaction.MerchantName {
		t.Errorf("MerchantName mismatch: got %s, want %s", retrievedTransaction.MerchantName, testTransaction.MerchantName)
	}
	if retrievedTransaction.PaymentChannel != testTransaction.PaymentChannel {
		t.Errorf("PaymentChannel mismatch: got %s, want %s", retrievedTransaction.PaymentChannel, testTransaction.PaymentChannel)
	}
	if retrievedTransaction.PrimaryCategory != testTransaction.PrimaryCategory {
		t.Errorf("PrimaryCategory mismatch: got %s, want %s", retrievedTransaction.PrimaryCategory, testTransaction.PrimaryCategory)
	}
	if retrievedTransaction.DetailedCategory != testTransaction.DetailedCategory {
		t.Errorf("DetailedCategory mismatch: got %s, want %s", retrievedTransaction.DetailedCategory, testTransaction.DetailedCategory)
	}
	if retrievedTransaction.ConfidenceLevel != testTransaction.ConfidenceLevel {
		t.Errorf("ConfidenceLevel mismatch: got %s, want %s", retrievedTransaction.ConfidenceLevel, testTransaction.ConfidenceLevel)
	}

	// Test deleting transactions
	deletedCount, err := tm.DeleteTransactions(ctx, "testuser", &startDate, &endDate)
	if err != nil {
		t.Fatalf("Failed to delete transactions: %v", err)
	}
	if deletedCount == 0 {
		t.Errorf("No transactions were deleted")
	}

	// Verify deletion
	transactions, err = tm.GetTransactions(ctx, "testuser", &startDate, &endDate, 10, 0)
	if err != nil {
		t.Fatalf("Failed to retrieve transactions after deletion: %v", err)
	}
	if len(transactions) != 0 {
		t.Errorf("Expected 0 transactions after deletion, got %d", len(transactions))
	}

	t.Logf("Successfully completed all database operations")
}
