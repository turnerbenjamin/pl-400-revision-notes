package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/turnerbenjamin/go_odata/accountsService"
	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/msal"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	eu := os.Getenv("ENVIRONMENT_URL")
	tid := os.Getenv("TENANT_ID")

	c := msal.DataverseServiceConfig{
		ClientId:     os.Getenv("CLIENT_ID"),
		TenantId:     tid,
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		ResourceUrl:  eu,
		APIBaseUrl:   eu + "api/data/v9.2/",
		Authority:    fmt.Sprintf("https://login.microsoftonline.com/%s", tid),
	}
	doUseAppAuth := true

	var dataverseService msal.DataverseService

	if doUseAppAuth {
		dataverseService, err = msal.GetAppService(c)
	} else {
		dataverseService, err = msal.GetDelegatedService(c)
	}
	if err != nil {
		log.Fatal(err)
	}

	appService := accountsService.NewAccountsOdataService(dataverseService)
	appService.Create(&model.Account{Name: "Go Account", City: "Peterborough"})
	appService.RetrieveByName("Go Account")

}
