package model

import (
	"encoding/json"
	"log"
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
