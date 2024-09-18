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

	"github.com/joho/godotenv"
)

type PlaidTokenRequest struct {
	ClientID     string   `json:"client_id"`
	Secret       string   `json:"secret"`
	ClientName   string   `json:"client_name"`
	User         User     `json:"user"`
	Products     []string `json:"products"`
	CountryCodes []string `json:"country_codes"`
	Language     string   `json:"language"`
	// Webhook      string   `json:"webhook"`
	RedirectURI string `json:"redirect_uri"`
}

// User represents the user object in the Plaid API request
type User struct {
	ClientUserID string `json:"client_user_id"`
}

func LinkHandler(w http.ResponseWriter, r *http.Request) {
	// Prepare the request payload

	err := godotenv.Load(".env.local")
	if err != nil {
		http.Error(w, "failed to load environment variables", http.StatusInternalServerError)
		return
	}
	plaidRequest := PlaidTokenRequest{
		ClientID:     os.Getenv("CLIENT_ID"),
		Secret:       os.Getenv("SANDBOX_SECRET"),
		ClientName:   "cashew",
		User:         User{ClientUserID: "unique_user_id"},
		Products:     []string{"auth"},
		CountryCodes: []string{"US"},
		Language:     "en",
		RedirectURI:  "http://localhost:8080/callback",
	}

	// Convert the struct to JSON
	requestBody, err := json.Marshal(plaidRequest)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Make the POST request to Plaid
	resp, err := http.Post("https://sandbox.plaid.com/link/token/create", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Failed to make request to Plaid", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response from Plaid
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Plaid", http.StatusInternalServerError)
		return
	}

	// Write the Plaid response back to the client
	log.Printf("plaid response: %s", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

type ExchangeRequest struct {
	PublicToken string `json:"public_token"`
	UserId      string `json:"user_id"`
}

type ExchangeResponse struct {
	AccessToken string `json:"access_token"`
	ItemID      string `json:"item_id"`
}

var TokenManager db.TokenManager

var TransactionManager db.TransactionManager

// example response:
//
//	{
//		"access_token": "access-sandbox-c5cf65ec-b58f-4fe1-8a91-f8b4bf383355",
//		"item_id": "X7yj8NJ8PnFQkr3zENzdIGDDrL9lmeFd9gN8W",
//		"request_id": "4arMsEAbVi64dYA"
//	  }
func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	var req ExchangeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	requestBody, _ := json.Marshal(map[string]string{
		"client_id":    os.Getenv("CLIENT_ID"),
		"secret":       os.Getenv("SANDBOX_SECRET"),
		"public_token": req.PublicToken,
	})

	resp, _ := http.Post("https://sandbox.plaid.com/item/public_token/exchange", "application/json", bytes.NewBuffer(requestBody))
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response from Plaid", http.StatusInternalServerError)
		return
	}

	var exchangeResp struct {
		AccessToken string `json:"access_token"`
		ItemID      string `json:"item_id"`
		RequestID   string `json:"request_id"`
	}

	err = json.Unmarshal(body, &exchangeResp)

	log.Printf("plaid response: %s", string(body))

	if resp.StatusCode != 200 {
		w.Write(body)
	} else {
		err := TokenManager.StoreToken(req.UserId, exchangeResp.AccessToken, 30*24*time.Hour)
		if err != nil {
			http.Error(w, "Failed to store token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		log.Printf("token storage: success")
	}
}

type PrevConnectionRequest struct {
	UserId string `json:"user_id"`
}

func PrevConnectionHandler(w http.ResponseWriter, r *http.Request) {
	var req ExchangeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	access_token, err := TokenManager.GetToken(req.UserId)
	if err != nil {
		http.Error(w, "Error fetching token", http.StatusBadRequest)
		log.Printf("Error fetching token %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	log.Printf("token fetch: success")
	log.Printf(access_token)
}

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
