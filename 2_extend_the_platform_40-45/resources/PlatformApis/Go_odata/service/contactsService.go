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

var contactSelects = "$select=contactid,firstname,lastname,emailaddress1"

type contactService struct {
	dataverseService    DataverseService
	getNextResponseFunc func(string) (*model.GetManyResponse[*model.Contact], error)
	baseUrl             string
	pageLimit           int
}

func NewContactService(dataverseService DataverseService, baseUrl string, pageLimit int) EntityService[*model.Contact] {
	return &contactService{
		dataverseService:    dataverseService,
		getNextResponseFunc: buildGetNextResultFunction[*model.Contact](dataverseService, pageLimit),
		baseUrl:             baseUrl,
		pageLimit:           pageLimit,
	}
}

func (s *contactService) Create(contact *model.Contact) (*model.Contact, error) {

	path := s.baseUrl + "contacts"
	req, err := reqbuilder.NewODataReqBuilder(http.MethodPost, path, bytes.NewReader(contact.ToJSON())).
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

	newAccount := model.NewContactFromJson(*res.Body)
	return newAccount, nil
}

func (s *contactService) List(searchTerm string) (view.EntityList[*model.Contact], error) {

	path := s.baseUrl + "contacts"
	rb := reqbuilder.NewODataReqBuilder(http.MethodGet, path, nil).
		AddQueryParam(contactSelects)
	if searchTerm != "" {
		filter := url.QueryEscape(fmt.Sprintf("contains(fullname,'%s') or contains(emailaddress1,'%s')", searchTerm, searchTerm))
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

	gmr := model.NewGetManyResponseFromJson[*model.Contact](*res.Body)

	return model.CreateEntityList(*gmr, s.getNextResponseFunc), nil
}

func (s *contactService) Get(guid string) (*model.Contact, error) {

	path := fmt.Sprintf("%scontacts(%s)", s.baseUrl, guid)

	req, err := reqbuilder.NewODataReqBuilder(http.MethodGet, path, nil).
		AddQueryParam(contactSelects).
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

	account := model.NewContactFromJson(*res.Body)
	return account, nil
}

func (s *contactService) Update(guid string, entityToUpdate *model.Contact) error {
	path := fmt.Sprintf("%scontacts(%s)", s.baseUrl, guid)

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

func (s *contactService) Delete(guid string) error {
	path := fmt.Sprintf("%scontacts(%s)", s.baseUrl, guid)

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

// func buildGetNextResultFunction[T any](s DataverseService, pageLimit int) func(url string) (*model.GetManyResponse[T], error) {
// 	return func(url string) (*model.GetManyResponse[T], error) {
// 		req, err := reqbuilder.NewODataReqBuilder(http.MethodGet, url, nil).Build()
// 		if err != nil {
// 			return nil, err
// 		}

// 		req.Header.Set("Prefer", fmt.Sprintf("odata.maxpagesize=%d", pageLimit))

// 		dr, err := s.Execute(req)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return model.NewGetManyResponseFromJson[T](*dr.Body), nil
// 	}
// }
