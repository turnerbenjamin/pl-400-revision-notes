// Package model provides data structures for working with Dataverse OData
// API responses.
package model

import (
	"github.com/turnerbenjamin/go_odata/view"
)

// Account represents an account entity from Microsoft Dataverse.
// It contains basic information about a business account record and implements
// the view.Entity interface.
type Account struct {
	// Id is the unique identifier for the account
	Id string `json:"accountid,omitempty"`

	// Name is the account's primary name
	Name string `json:"name"`

	// City is the city of the account's primary address
	City string `json:"address1_city,omitempty"`
}

// AccountListColumns returns a slice of ListColumn configurations
// for displaying Account entities in a formatted list.
//
// The returned columns include:
// - Name: The account's name
// - City: The account's primary address city
//
// Returns:
//   - A slice of ListColumn objects for Account entities
//   - An error if column creation fails
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

// ID implements the view.Entity interface by returning the account's unique
// identifier.
func (a *Account) ID() string {
	return a.Id
}

// Label implements the view.Entity interface by returning a human-readable
// label for the account, which is its name.
func (a *Account) Label() string {
	return a.Name
}
