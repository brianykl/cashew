package controllers

import (
	"cashew-api/conf"
	"context"
	"net/http"

	beelog "github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/oauth2"
)

type OAuthController struct {
	beego.Controller
	oauthConfig conf.OAuthConfig
}

func NewOAuthController(cfg conf.OAuthConfig) *OAuthController {
	return &OAuthController{
		oauthConfig: cfg,
	}
}

var oauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/oauth/callback", // my callback url
	ClientID:     "myclientid",
	ClientSecret: "myclientsecret",
	Scopes:       []string{"email"}, // scopes i want to access
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https", // plaid authorization url
		TokenURL: "https", // plaid token url
	},
}

func (c *OAuthController) Redirect() {
	state := generateState()
	c.SetSession("oauthState", state)
	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Ctx.Redirect(http.StatusFound, url)
}

func (c *OAuthController) Callback() {
	// retrieve the state from the session and the request
	sessionState := c.GetSession("oauthState").(string)
	requestState := c.GetString("state")

	// validate the state
	if sessionState != requestState {
		beelog.Error("invalid OAuth state")
		c.Abort("400")
		return
	}

	code := c.GetString("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		beelog.Error("failed to exchange token: ", err)
		c.Abort("500")
		return
	}

	c.Ctx.WriteString("OAuth callback recieved: " + token.AccessToken)
}

func generateState() string {
	return "" // randomly generate this
}
