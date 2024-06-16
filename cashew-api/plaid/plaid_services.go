package plaid_services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type LinkTokenCreateRequest struct {
	ClientID       string         `json:"client_id"`
	Secret         string         `json:"secret"`
	User           User           `json:"user"`
	ClientName     string         `json:"client_name"`
	Products       []string       `json:"products"`
	Transactions   Transactions   `json:"transactions"`
	CountryCodes   []string       `json:"country_codes"`
	Language       string         `json:"language"`
	RedirectURI    string         `json:"redirect_uri"`
	AccountFilters AccountFilters `json:"account_filters"`
}

type User struct {
	ClientUserID string `json:"client_user_id"`
	PhoneNumber  string `json:"phone_number"`
}

type Transactions struct {
	DaysRequested int `json:"days_requested"`
}

type AccountFilters struct {
	Depository Depository `json:"depository"`
	Credit     Credit     `json:"credit"`
}

type Depository struct {
	AccountSubtypes []string `json:"account_subtypes"`
}

type Credit struct {
	AccountSubtypes []string `json:"account_subtypes"`
}

type LinkTokenCreateResponse struct {
	LinkToken string `json:"link_token"`
}

func CreateLinkToken() (string, error) {
	err := godotenv.Load("conf/.env.plaid")
	if err != nil {
		log.Fatalf("Error loading .env.plaid file")
	}

	clientID := os.Getenv("PLAID_CLIENT_ID")
	secret := os.Getenv("PLAID_SECRET")
	user := User{
		ClientUserID: time.Now().String(), // maybe switch this to the user id in database
		PhoneNumber:  "+1 647 4690522",
	}
	requestBody := LinkTokenCreateRequest{
		ClientID:     clientID,
		Secret:       secret,
		User:         user,
		ClientName:   "cashew",
		Products:     []string{"transactions"},
		Transactions: Transactions{DaysRequested: 730},
		CountryCodes: []string{"US", "CA"},
		Language:     "en",
		RedirectURI:  "http://localhost:4000/",
		AccountFilters: AccountFilters{
			Depository: Depository{
				AccountSubtypes: []string{"checking", "savings"},
			},
			Credit: Credit{
				AccountSubtypes: []string{"credit card"},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	resp, err := http.Post("https://sandbox.plaid.com/link/token/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("received non-200 response: %v\nResponse body: %s", resp.Status, string(body))
	}

	var result LinkTokenCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return result.LinkToken, nil
}
