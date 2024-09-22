package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type AccountsRequest struct {
	UserId string `json:"user_id"`
}

type AccountName struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Balances Balances `json:"balances"`
}

type Balances struct {
	Available *float64 `json:"available"`
}

type AccountResponse struct {
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	AvailableBalance float64 `json:"available_balance"`
}

func AccountsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("PULLING ACCOUNTS")
	var req AccountsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	accessTokens, err := TokenManager.GetTokens(req.UserId)
	if err != nil {
		http.Error(w, "did not find access tokens linked to user", http.StatusBadRequest)
		return
	}

	var allAccounts []AccountName

	for _, accessToken := range accessTokens {
		requestBody, err := json.Marshal(map[string]interface{}{
			"client_id":    os.Getenv("CLIENT_ID"),
			"secret":       os.Getenv("SANDBOX_SECRET"),
			"access_token": accessToken,
		})
		resp, err := http.Post("https://sandbox.plaid.com/accounts/get", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			http.Error(w, "Failed to read response from Plaid", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var data struct {
			Accounts []AccountName `json:"accounts"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "failed to parse Plaid response", http.StatusInternalServerError)
			return
		}

		allAccounts = append(allAccounts, data.Accounts...)
	}

	var responseAccounts []AccountResponse

	for _, account := range allAccounts {
		availableBalance := 0.0
		if account.Balances.Available != nil {
			availableBalance = *account.Balances.Available
		}

		responseAccount := AccountResponse{
			Name:             account.Name,
			Type:             account.Type,
			AvailableBalance: availableBalance,
		}
		responseAccounts = append(responseAccounts, responseAccount)
	}
	log.Printf("%v", responseAccounts)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseAccounts)
}
