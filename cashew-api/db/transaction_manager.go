package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	UserId           string
	AccountId        string
	AccountName      string
	Amount           decimal.Decimal
	Currency         string
	AuthorizedDate   time.Time
	MerchantName     string
	PaymentChannel   string
	PrimaryCategory  string
	DetailedCategory string
	ConfidenceLevel  string
}

type TransactionManager interface {
	StoreTransactions(ctx context.Context, transactions []*Transaction) error
	GetTransactions(ctx context.Context, userId string, startDate *time.Time, endDate *time.Time, limit int, offset int) ([]*Transaction, error)
	DeleteTransactions(ctx context.Context, userId string, startDate *time.Time, endDate *time.Time) (int64, error)
}

type postgresTransactionManager struct {
	client *sql.DB
}

func NewTransactionManager(connDetails string) (TransactionManager, error) {
	db, err := sql.Open("postgres", connDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %v", err)
	}

	return &postgresTransactionManager{client: db}, nil
}

func (ptm *postgresTransactionManager) StoreTransactions(ctx context.Context, transactions []*Transaction) error {
	tx, err := ptm.client.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction %v", err)
	}
	defer tx.Rollback()

	txValues := make([]string, 0, len(transactions))
	txArgs := make([]interface{}, 0, len(transactions)*11)

	for i, t := range transactions {
		txValues = append(txValues, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*11+1, i*11+2, i*11+3, i*11+4, i*11+5, i*11+6, i*11+7, i*11+8, i*11+9, i*11+10, i*11+11))
		txArgs = append(txArgs,
			t.UserId,
			t.AccountId,
			t.AccountName,
			t.Amount,
			t.Currency,
			t.AuthorizedDate,
			t.MerchantName,
			t.PaymentChannel,
			t.PrimaryCategory,
			t.DetailedCategory,
			t.ConfidenceLevel)
	}

	query := fmt.Sprintf(`
		INSERT INTO transactions 
		(user_id, account_id, account_name, amount, currency, authorized_date, 
		merchant_name, payment_channel, primary_category, detailed_category, confidence_level)
		VALUES %s`,
		strings.Join(txValues, ","))

	_, err = tx.ExecContext(ctx, query, txArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (ptm *postgresTransactionManager) GetTransactions(ctx context.Context, userId string, startDate *time.Time, endDate *time.Time, limit int, offset int) ([]*Transaction, error) {
	query := `
        SELECT user_id, account_id, account_name, amount, currency, authorized_date, 
               merchant_name, payment_channel, primary_category, detailed_category, confidence_level
        FROM transactions
        WHERE user_id = $1 AND authorized_date >= $2 AND authorized_date <= $3
        ORDER BY authorized_date DESC
        LIMIT $4 OFFSET $5`

	rows, err := ptm.client.QueryContext(ctx, query, userId, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions: %v", err)
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		t := &Transaction{}
		err := rows.Scan(
			&t.UserId,
			&t.AccountId,
			&t.AccountName,
			&t.Amount,
			&t.Currency,
			&t.AuthorizedDate,
			&t.MerchantName,
			&t.PaymentChannel,
			&t.PrimaryCategory,
			&t.DetailedCategory,
			&t.ConfidenceLevel,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return transactions, nil
}

func (ptm *postgresTransactionManager) DeleteTransactions(ctx context.Context, userId string, startDate *time.Time, endDate *time.Time) (int64, error) {
	tx, err := ptm.client.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	query := `
        DELETE FROM transactions
        WHERE user_id = $1
          AND ($2::timestamp IS NULL OR authorized_date >= $2)
          AND ($3::timestamp IS NULL OR authorized_date <= $3)`

	var result sql.Result
	result, err = tx.ExecContext(ctx, query, userId, startDate, endDate)

	if err != nil {
		return 0, fmt.Errorf("failed to delete transactions: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting rows affected: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return rowsAffected, nil
}

func (t Transaction) String() string {
	return fmt.Sprintf("Transaction{UserId: %s, AccountId: %s, AccountName: %s, "+
		"Amount: %s %s, AuthorizedDate: %s, MerchantName: %s, PaymentChannel: %s, "+
		"PrimaryCategory: %s, DetailedCategory: %s, ConfidenceLevel: %s}",
		t.UserId, t.AccountId, t.AccountName,
		t.Amount.String(), t.Currency, t.AuthorizedDate.Format(time.RFC3339),
		t.MerchantName, t.PaymentChannel,
		t.PrimaryCategory, t.DetailedCategory, t.ConfidenceLevel)
}
