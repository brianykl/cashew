package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/brianykl/cashew-api/handlers"
)

func EnsureValidToken() func(next http.HandlerFunc) http.HandlerFunc {
	issuer := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatalf("failed to parse the issuer url: %v", err)
	}
	audience := os.Getenv("AUTH0_AUDIENCE")

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuer,
		[]string{audience},
	)
	if err != nil {
		panic(err)
	}

	middleware := jwtmiddleware.New(jwtValidator.ValidateToken)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.CheckJWT(next).ServeHTTP
	}
}

func main() {
	http.HandleFunc("/link", handlers.LinkHandler)
	http.HandleFunc("/protected/exchange", EnsureValidToken()(http.HandlerFunc(handlers.ExchangeHandler)))

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
