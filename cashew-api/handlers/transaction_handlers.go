package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brianykl/cashew/cashew-api/db"
	"github.com/shopspring/decimal"
)

var TokenManager db.TokenManager

var TransactionManager db.TransactionManager

type TransactionsRequest struct {
	UserId      string
	AccessToken string
}

type PlaidResponse struct {
	Accounts []Account          `json:"accounts"`
	Added    []PlaidTransaction `json:"added"`
}

type Account struct {
	AccountID string `json:"account_id"`
	Name      string `json:"name"`
}

type PlaidTransaction struct {
	AccountID               string  `json:"account_id"`
	Amount                  float64 `json:"amount"`
	ISOCurrencyCode         string  `json:"iso_currency_code"`
	AuthorizedDate          string  `json:"authorized_date"`
	Date                    string  `json:"date"`
	MerchantName            string  `json:"merchant_name"`
	PaymentChannel          string  `json:"payment_channel"`
	PersonalFinanceCategory struct {
		Primary         string `json:"primary"`
		Detailed        string `json:"detailed"`
		ConfidenceLevel string `json:"confidence_level"`
	} `json:"personal_finance_category"`
}

func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var req TransactionsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"client_id":    os.Getenv("CLIENT_ID"),
		"secret":       os.Getenv("SANDBOX_SECRET"),
		"access_token": req.AccessToken,
		"count":        250,
	})
	if err != nil {
		http.Error(w, "Failed to create request body", http.StatusInternalServerError)
		return
	}
	resp, err := http.Post("https://sandbox.plaid.com/transactions/sync", "application/json", bytes.NewBuffer(requestBody))
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Plaid", http.StatusInternalServerError)
		return
	}
	log.Printf(string(body))
	transactions, _ := parseTransactions(req.UserId, body)
	// t := transactions[0]
	// log.Printf("Transaction: AccountId: %s, AccountName: %s, Amount: %s %s, "+
	// 	"AuthorizedDate: %s, MerchantName: %s, PaymentChannel: %s, "+
	// 	"PrimaryCategory: %s, DetailedCategory: %s, ConfidenceLevel: %s",
	// 	t.AccountId, t.AccountName,
	// 	t.Amount.String(), t.Currency,
	// 	t.AuthorizedDate.Format(time.RFC3339),
	// 	t.MerchantName, t.PaymentChannel,
	// 	t.PrimaryCategory, t.DetailedCategory, t.ConfidenceLevel)
	TransactionManager.StoreTransactions(r.Context(), transactions)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := struct {
		Transactions []*db.Transaction `json:"transactions"`
	}{
		Transactions: transactions,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func parseTransactions(user_id string, body []byte) ([]*db.Transaction, error) {
	var plaidResp PlaidResponse
	err := json.Unmarshal(body, &plaidResp)
	if err != nil {
		return nil, err
	}

	accountMap := make(map[string]string)
	for _, account := range plaidResp.Accounts {
		accountMap[account.AccountID] = account.Name
	}

	var transactions []*db.Transaction
	for _, pt := range plaidResp.Added {

		authorizedDate, err := time.Parse("2006-01-02", pt.AuthorizedDate)
		if err != nil {

			authorizedDate, err = time.Parse("2006-01-02", pt.Date)
			if err != nil {

				continue
			}
		}

		transaction := &db.Transaction{
			UserId:           user_id,
			AccountId:        pt.AccountID,
			AccountName:      accountMap[pt.AccountID],
			Amount:           decimal.NewFromFloat(pt.Amount),
			Currency:         pt.ISOCurrencyCode,
			AuthorizedDate:   authorizedDate,
			MerchantName:     pt.MerchantName,
			PaymentChannel:   pt.PaymentChannel,
			PrimaryCategory:  pt.PersonalFinanceCategory.Primary,
			DetailedCategory: pt.PersonalFinanceCategory.Detailed,
			ConfidenceLevel:  pt.PersonalFinanceCategory.ConfidenceLevel,
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
