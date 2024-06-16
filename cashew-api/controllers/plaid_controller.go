package controllers

import (
	plaid_services "cashew-api/plaid"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type PlaidController struct {
	beego.Controller
}

func (c *PlaidController) GenerateLinkToken() {
	// userID := "unique-user-id" // Replace with actual user ID logic
	linkToken, err := plaid_services.CreateLinkToken()
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data["json"] = map[string]string{"link_token": linkToken, "brian": "brian"}
	c.ServeJSON()
}
