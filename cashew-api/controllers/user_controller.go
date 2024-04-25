package controllers

import (
	"context"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/brianykl/cashew/services/users/client"
)

type UserController struct {
	beego.Controller
	UserClient *client.UserClient
}

type UserData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserController) Prepare() {
	var err error
	c.UserClient, err = client.NewUserClient("localhost:5001")
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.Body([]byte("failed to connect to user microservice"))
		return
	}
}

func (c *UserController) Post() {
	var inputData struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputData)

	response, err := c.UserClient.CreateUser(context.Background(), inputData.Name, inputData.Email, inputData.Password)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.Ctx.Output.SetStatus(500)
	} else {
		c.Data["json"] = response
	}
	c.ServeJSON()
}

func (c *UserController) GetUser() {
	var inputData struct {
		UserId string `json:"userid"`
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputData)

	response, err := c.UserClient.GetUser(context.Background(), inputData.UserId)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.Ctx.Output.SetStatus(500)
	} else {
		c.Data["json"] = response
	}
	c.ServeJSON()
}

func (c *UserController) UpdateUser() {
	var inputData struct {
		UserId   string `json:"userid"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &inputData)

	response, err := c.UserClient.UpdateUser(
		context.Background(),
		inputData.UserId,
		inputData.Name,
		inputData.Email,
		inputData.Password)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.Ctx.Output.SetStatus(500)
	} else {
		c.Data["json"] = response
	}
	c.ServeJSON()
}
