package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type AccountsRequest struct {
	UserId string `json:"user_id"`
}

func AccountsHandler(w http.ResponseWriter, r *http.Request) {
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
			Accounts []struct {
				Name string `json:"name"`
			} `json:"accounts"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "failed to parse Plaid response", http.StatusInternalServerError)
			return
		}

		accountNames := make([]string, len(data.Accounts))
		for i, account := range data.Accounts {
			accountNames[i] = account.Name
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string][]string{"account_names": accountNames})
	}

}
