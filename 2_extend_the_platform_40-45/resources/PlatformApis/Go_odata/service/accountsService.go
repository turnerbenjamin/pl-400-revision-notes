package service

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/msal"
)

var selects = "?$select=accountid,name,address1_city"

type AccountsService interface {
	Create(*model.Account) (*model.Account, error)
	List(string, int) (EntityList[model.Account], error)
	RetrieveByName(name string)
}

type accountsService struct {
	dataverseService msal.DataverseService
}

func NewAccountService(dataverseService msal.DataverseService) AccountsService {
	return &accountsService{dataverseService: dataverseService}
}

func (s *accountsService) Create(account *model.Account) (*model.Account, error) {
	endpoint := "accounts" + selects
	res, err := s.dataverseService.Post(endpoint, account.ToJSON())
	if err != nil {
		return nil, err
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(*res.Body)
		return nil, errors.New(errMsg)
	}
	newAccount := model.NewAccountFromJson(*res.Body)
	return newAccount, nil
}

func (s *accountsService) List(searchTerm string, page int) (EntityList[model.Account], error) {
	filter := ""
	if searchTerm != "" {
		filterValue := fmt.Sprintf("contains(name,'%s') or contains(address1_city,'%s')", searchTerm, searchTerm)
		filter = "&$filter=" + url.QueryEscape(filterValue)
	}
	endpoint := fmt.Sprintf("accounts%s%s", selects, filter)
	res, err := s.dataverseService.GetMany(endpoint)
	if err != nil {
		return nil, err
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(*res.Body)
		return nil, errors.New(string(errMsg))
	}

	gmr := model.NewGetManyResponseFromJson[model.Account](*res.Body)
	return &entityList[model.Account]{
		res: gmr,
		svc: s.dataverseService,
	}, nil
}

func (s *accountsService) RetrieveByName(name string) {
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
