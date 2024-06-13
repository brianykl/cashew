package controllers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/brianykl/cashew/services/users/client"
	jwt "github.com/dgrijalva/jwt-go"
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

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
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

func (c *UserController) Create() {
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

func (c *UserController) Login() {
	var loginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.Unmarshal(c.Ctx.Input.RequestBody, &loginInfo)

	response, err := c.UserClient.VerifyUser(context.Background(), loginInfo.Email, loginInfo.Password)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.Ctx.Output.SetStatus(500)
		c.ServeJSON()
		return
	}

	if !response.IsValid {
		c.Data["json"] = map[string]string{"error": "Invalid credentials"}
		c.Ctx.Output.SetStatus(401)
	} else {
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Email: loginInfo.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			c.Data["json"] = map[string]string{"error": "Failed to generate token"}
			c.Ctx.Output.SetStatus(500)
			c.ServeJSON()
			return
		}

		responseData := map[string]interface{}{
			"user":  loginInfo.Email,
			"token": tokenString,
		}
		responseJSON, _ := json.Marshal(responseData)
		log.Printf("response JSON: %s", responseJSON)
		c.Data["json"] = responseData
		c.Ctx.Output.SetStatus(200)

	}
	c.ServeJSON()
}

// func (c *UserController) Get() {
// 	var loginInfo struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}
// 	json.Unmarshal(c.Ctx.Input.RequestBody, &loginInfo)

// 	response, err := c.UserClient.VerifyUser(context.Background(), loginInfo.Email, loginInfo.Password)
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 		c.Ctx.Output.SetStatus(500)
// 	} else {
// 		c.Data["json"] = response
// 	}
// 	c.ServeJSON()
// }

// func (c *UserController) GetUser() {
// 	var inputData struct {
// 		UserId string `json:"userid"`
// 	}
// 	json.Unmarshal(c.Ctx.Input.RequestBody, &inputData)

// 	response, err := c.UserClient.GetUser(context.Background(), inputData.UserId)
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 		c.Ctx.Output.SetStatus(500)
// 	} else {
// 		c.Data["json"] = response
// 	}
// 	c.ServeJSON()
// }

// func (c *UserController) UpdateUser() {
// 	var inputData struct {
// 		UserId   string `json:"userid"`
// 		Email    string `json:"email"`
// 		Name     string `json:"name"`
// 		Password string `json:"password"`
// 	}
// 	json.Unmarshal(c.Ctx.Input.RequestBody, &inputData)

// 	response, err := c.UserClient.UpdateUser(
// 		context.Background(),
// 		inputData.UserId,
// 		inputData.Name,
// 		inputData.Email,
// 		inputData.Password)
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 		c.Ctx.Output.SetStatus(500)
// 	} else {
// 		c.Data["json"] = response
// 	}
// 	c.ServeJSON()
// }
