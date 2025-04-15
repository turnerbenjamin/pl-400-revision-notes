// Package model provides data structures for working with Dataverse OData API
// responses.
package model

import (
	"fmt"

	"github.com/turnerbenjamin/go_odata/view"
)

// Contact represents a contact entity from Microsoft Dataverse.
// It contains basic information about an individual contact record
// and implements the view.Entity interface.
type Contact struct {
	// Id is the unique identifier for the contact
	Id string `json:"contactid,omitempty"`

	// FirstName is the contact's first name
	FirstName string `json:"firstname"`

	// LastName is the contact's last name
	LastName string `json:"lastname"`

	// Email is the contact's primary email address
	Email string `json:"emailaddress1,omitempty"`
}

// ContactListColumns returns a slice of ListColumn configurations for
// displaying Contact entities in a formatted list.
//
// The returned columns include:
// - Name: A formatted string combining first and last name
// - Email: The contact's email address
//
// Returns:
//   - A slice of ListColumn objects for Contact entities
//   - An error if column creation fails
func ContactListColumns() ([]view.ListColumn[*Contact], error) {
	getName := func(c *Contact) string {
		return fmt.Sprintf("%s %s", c.FirstName, c.LastName)
	}

	nameColumn, err := view.NewListColumn("Name", getName)
	if err != nil {
		return nil, err
	}

	emailCol, err := view.NewListColumn("Email", func(a *Contact) string {
		return a.Email
	})
	if err != nil {
		return nil, err
	}

	return []view.ListColumn[*Contact]{
		nameColumn,
		emailCol,
	}, nil
}

// ID implements the view.Entity interface by returning the contact's unique
// identifier.
func (a *Contact) ID() string {
	return a.Id
}

// Label implements the view.Entity interface by returning a human-readable
// label for the contact, combining the first and last name.
func (a *Contact) Label() string {
	return fmt.Sprintf("%s %s", a.FirstName, a.LastName)
}
