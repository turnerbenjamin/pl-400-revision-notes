package service

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/reqbuilder"
	"github.com/turnerbenjamin/go_odata/view"
)

var accountSelects = "$select=accountid,name,address1_city"

type accountsService struct {
	dataverseService    DataverseService
	getNextResponseFunc func(string) (*model.GetManyResponse[*model.Account], error)
	baseUrl             string
	pageLimit           int
}

func NewAccountService(dataverseService DataverseService, baseUrl string, pageLimit int) EntityService[*model.Account] {
	return &accountsService{
		dataverseService:    dataverseService,
		getNextResponseFunc: buildGetNextResultFunction[*model.Account](dataverseService, pageLimit),
		baseUrl:             baseUrl,
		pageLimit:           pageLimit,
	}
}

func (s *accountsService) Create(account *model.Account) (*model.Account, error) {

	path := s.baseUrl + "accounts"
	req, err := reqbuilder.NewODataReqBuilder(http.MethodPost, path, bytes.NewReader(account.ToJSON())).
		Build()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	res, err := s.dataverseService.Execute(req)
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

func (s *accountsService) List(searchTerm string) (view.EntityList[*model.Account], error) {

	path := s.baseUrl + "accounts"
	rb := reqbuilder.NewODataReqBuilder(http.MethodGet, path, nil).
		AddQueryParam(accountSelects)
	if searchTerm != "" {
		filter := url.QueryEscape(fmt.Sprintf("contains(name,'%s') or contains(address1_city,'%s')", searchTerm, searchTerm))
		rb.AddQueryParam("$filter=" + filter)
	}

	req, err := rb.Build()
	if err != nil {
		return nil, err
	}

	req.Header.Set("Prefer", fmt.Sprintf("odata.maxpagesize=%d", s.pageLimit))

	res, err := s.dataverseService.Execute(req)
	if err != nil {
		return nil, err
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(*res.Body)
		return nil, errors.New(string(errMsg))
	}

	gmr := model.NewGetManyResponseFromJson[*model.Account](*res.Body)

	return model.CreateEntityList(*gmr, s.getNextResponseFunc), nil
}

func (s *accountsService) Get(guid string) (*model.Account, error) {

	path := fmt.Sprintf("%saccounts(%s)", s.baseUrl, guid)

	req, err := reqbuilder.NewODataReqBuilder(http.MethodGet, path, nil).
		AddQueryParam(accountSelects).
		Build()
	if err != nil {
		return nil, err
	}

	res, err := s.dataverseService.Execute(req)
	if err != nil {
		return nil, err
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(*res.Body)
		return nil, errors.New(errMsg)
	}

	account := model.NewAccountFromJson(*res.Body)
	return account, nil
}

func (s *accountsService) Update(guid string, entityToUpdate *model.Account) error {
	path := fmt.Sprintf("%saccounts(%s)", s.baseUrl, guid)

	req, err := reqbuilder.NewODataReqBuilder(http.MethodPatch, path, bytes.NewReader(entityToUpdate.ToJSON())).
		Build()
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := s.dataverseService.Execute(req)
	if err != nil {
		return err
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(*res.Body)
		return errors.New(errMsg)
	}
	return nil
}

func (s *accountsService) Delete(guid string) error {
	path := fmt.Sprintf("%saccounts(%s)", s.baseUrl, guid)

	req, err := reqbuilder.NewODataReqBuilder(http.MethodDelete, path, nil).
		Build()
	if err != nil {
		return err
	}

	res, err := s.dataverseService.Execute(req)
	if err != nil {
		return err
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(*res.Body)
		return errors.New(errMsg)
	}
	return nil
}

func buildGetNextResultFunction[T any](s DataverseService, pageLimit int) func(url string) (*model.GetManyResponse[T], error) {
	return func(url string) (*model.GetManyResponse[T], error) {
		req, err := reqbuilder.NewODataReqBuilder(http.MethodGet, url, nil).Build()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Prefer", fmt.Sprintf("odata.maxpagesize=%d", pageLimit))

		dr, err := s.Execute(req)
		if err != nil {
			return nil, err
		}
		return model.NewGetManyResponseFromJson[T](*dr.Body), nil
	}
}
