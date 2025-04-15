// Package logicalnames provides constants for entity and field logical names
// used when interacting with Microsoft Dataverse through the Web API.
package logicalnames

// Constants representing entity and field logical names in Microsoft Dataverse.
// These are used when constructing OData queries and processing API responses.
const (
	TableAccount           = "account"
	TableAccountResource   = "accounts"
	TableContactSingular   = "contact"
	TableContactResource   = "contacts"
	ColumnAccountId        = "accountid"
	ColumnAccountName      = "name"
	ColumnAccountCity      = "address1_city"
	ColumnContactId        = "contactid"
	ColumnContactFirstName = "firstname"
	ColumnContactLastName  = "lastname"
	ColumnContactEmail     = "emailaddress1"
)
