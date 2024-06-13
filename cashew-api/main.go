package main

import (
	_ "cashew-api/conf"
	_ "cashew-api/routers"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	context "github.com/beego/beego/v2/server/web/context"
	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func main() {

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	logs.SetLogger(logs.AdapterConsole, `{"level":7, "color":true}`)

	// Set up file logging
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/app.log", "level":7}`)

	// Enable function call depth (optional)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	beego.InsertFilter("*", beego.BeforeRouter, corsMiddleware)
	beego.InsertFilter("/v1/protected/*", beego.BeforeRouter, JWTMiddleware())
	beego.Run()
}

func corsMiddleware(ctx *context.Context) {
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	if ctx.Request.Method == "OPTIONS" {
		ctx.ResponseWriter.WriteHeader(200)
		return
	}
}

func JWTMiddleware() beego.FilterFunc {
	return func(ctx *context.Context) {
		authHeader := ctx.Input.Header("Authorization")
		if authHeader == "" {
			ctx.Abort(http.StatusUnauthorized, "missing authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			ctx.Abort(http.StatusUnauthorized, "Invalid token")
			return
		}
		ctx.Input.SetData("email", claims.Email)
	}
}
