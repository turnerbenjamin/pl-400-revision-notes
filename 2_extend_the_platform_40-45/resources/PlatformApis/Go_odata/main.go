package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/turnerbenjamin/go_odata/app"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	eu := os.Getenv("ENVIRONMENT_URL")
	tid := os.Getenv("TENANT_ID")

	c := app.AppConfig{
		ClientId:     os.Getenv("CLIENT_ID"),
		TenantId:     tid,
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		ResourceUrl:  eu,
		APIBaseUrl:   eu + "api/data/v9.2/",
		Authority:    fmt.Sprintf("https://login.microsoftonline.com/%s", tid),
		PageLimit:    5,
	}
	a, err := app.Create(c)
	if err != nil {
		log.Fatal(err)
	}

	a.Run()
}
