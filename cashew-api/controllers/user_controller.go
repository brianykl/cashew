package controllers

import (
	"context"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/brianykl/cashew/services/user/client"
)

type UserController struct {
	beego.Controller
}

type UserData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserController) Post() {
	var userData UserData
	err := c.Ctx.Input.Bind(&userData, "body")
	if err != nil {
		c.Abort("400")
		return
	}

	userClient, err := client.NewUserClient("address-of-user-service:port")
	if err != nil {
		c.Abort("500")
		return
	}

	ctx := context.Background()
	response, err := userClient.CreateUser(ctx, userData.Name, userData.Email, userData.Password)
	if err != nil {
		c.Abort("500")
		return
	}

	c.Data["json"] = map[string]interface{}{
		"UserID":   response.UserId,
		"Email":    response.Email,
		"Name":     response.Name,
		"Password": response.Password,
	}
	c.ServeJSON()
}
