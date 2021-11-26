package util

import (
	"github.com/Nerzal/gocloak/v8"
	"golang.org/x/net/context"
	"os"
)

import _ "github.com/joho/godotenv/autoload"

func GetToken() string {
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))
	ctx := context.Background()

	token, err := client.Login(ctx,
		os.Getenv("KEYCLOAK_CLIENT_ID"),
		os.Getenv("KEYCLOAK_SECRET"),
		os.Getenv("KEYCLOAK_REALM"),
		os.Getenv("KEYCLOAK_USERNAME"),
		os.Getenv("KEYCLOAK_PASSWORD"),
	)

	if err != nil {
		panic("Login failed:" + err.Error())
	}

	return token.AccessToken
}
