package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

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
}

type ExchangeResponse struct {
	AccessToken string `json:"access_token"`
	ItemID      string `json:"item_id"`
}

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

	log.Printf("plaid response: %s", string(body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body) // do not include the access token in the response lol, we gotta keep that in the backend
}
