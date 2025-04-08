package accountsService

import (
	"fmt"
	"log"

	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/msal"
)

var selects = "?$select=accountid,name,address1_city"

type accountsOdataService struct {
	dataverseService msal.DataverseService
}

func NewAccountsOdataService(dataverseService msal.DataverseService) AccountsService {
	return &accountsOdataService{dataverseService: dataverseService}
}

func (s *accountsOdataService) Create(account *model.Account) {

	endpoint := "accounts" + selects
	res, err := s.dataverseService.Post(endpoint, account.ToJSON())
	if err != nil {
		log.Println("Failed to create request:", err)
	}
	if res.IsSuccessful {
		log.Println("Account created successfully")
		account := model.NewAccountFromJson(*res.Body)
		log.Printf("Id: %s, Name: %s, City: %s", account.Id, account.Name, account.City)
	} else {
		log.Printf("Account creation failed with status %d (%s)\n", res.StatusCode, string(*res.Body))
	}

}

func (s *accountsOdataService) RetrieveByName(name string) {
	endpoint := fmt.Sprintf("accounts(name='%s')%s", name, selects)
	res, err := s.dataverseService.Get(endpoint)
	if err != nil {
		log.Println("Failed to create request:", err)
	}
	if res.IsSuccessful {
		account := model.NewAccountFromJson(*res.Body)
		log.Println("Account retrieved successfully")
		log.Printf("Id: %s, Name: %s, City: %s", account.Id, account.Name, account.City)
	} else {
		log.Println("Unable to retrieve account")
	}

}
