package model

import (
	"fmt"

	"github.com/turnerbenjamin/go_odata/view"
)

type Contact struct {
	Id        string `json:"contactid,omitempty"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"emailaddress1,omitempty"`
}

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

func (a *Contact) ID() string {
	return a.Id
}

func (a *Contact) Label() string {
	return fmt.Sprintf("%s %s", a.FirstName, a.LastName)
}
