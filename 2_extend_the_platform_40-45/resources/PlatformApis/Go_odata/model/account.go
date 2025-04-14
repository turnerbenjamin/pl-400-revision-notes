package model

import (
	"encoding/json"
	"log"

	"github.com/turnerbenjamin/go_odata/view"
)

type Account struct {
	Id   string `json:"accountid,omitempty"`
	Name string `json:"name"`
	City string `json:"address1_city,omitempty"`
}

func (a *Account) ToJSON() []byte {
	body, err := json.Marshal(a)
	if err != nil {
		log.Fatal("failed to serialise account")
	}
	return body
}

func NewAccountFromJson(accountJson []byte) *Account {
	var account Account
	json.Unmarshal(accountJson, &account)
	return &account
}

func AccountListColumns() ([]view.ListColumn[*Account], error) {
	getName := func(a *Account) string {
		return a.Name
	}
	nameColumn, err := view.NewListColumn("Name", getName)

	if err != nil {
		return nil, err
	}

	cityCol, err := view.NewListColumn("City", func(a *Account) string {
		return a.City
	})
	if err != nil {
		return nil, err
	}

	return []view.ListColumn[*Account]{
		nameColumn,
		cityCol,
	}, nil
}

func (a *Account) ID() string {
	return a.Id
}

func (a *Account) Label() string {
	return a.Name
}
