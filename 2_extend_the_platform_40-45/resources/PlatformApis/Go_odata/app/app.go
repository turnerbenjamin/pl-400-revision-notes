// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	"fmt"
	"net/url"

	authMode "github.com/turnerbenjamin/go_odata/constants/authmode"
	logicalNames "github.com/turnerbenjamin/go_odata/constants/logicalnames"
	mainMenuOption "github.com/turnerbenjamin/go_odata/constants/mainmenuoption"
	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/msal"
	"github.com/turnerbenjamin/go_odata/service"
	"github.com/turnerbenjamin/go_odata/view"
)

// AppConfig contains the configuration settings required for the application
// to connect to and interact with Dataverse.
type AppConfig struct {
	ClientID     string // EntraID client ID for authentication
	TenantID     string // EntraID tenant ID
	ClientSecret string // EntraID client secret for application authentication
	ResourceURL  string // Resource URL for accessing web api
	APIBaseURL   string // Base URL for the Dataverse API
	Authority    string // Authority URL for authentication
	PageLimit    int    // Maximum number of records to retrieve per page
}

// App defines the interface for the OData client application.
type App interface {
	// Run starts the application, handles authentication, and begins the UI
	// loop.
	// Returns an error if the application fails to start or encounters an
	// unrecoverable error.
	Run() error
}

// app is the concrete implementation of the App interface that manages services
// UI components, and application flow.
type app struct {
	config              *AppConfig
	accountsService     service.EntityService[*model.Account]
	contactsService     service.EntityService[*model.Contact]
	accountsListColumns []view.ListColumn[*model.Account]
	contactsListColumns []view.ListColumn[*model.Contact]
	ui                  view.UI
}

// NewApp creates a new instance of the application with the provided
// configuration.
// It initializes entity list columns and returns the App interface or an
// error.
func NewApp(config AppConfig) (App, error) {
	a := app{
		config: &config,
	}

	err := a.initEntityListColumns()
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// Run starts the application, initializing the UI, handling authentication,
// setting up services, and starting the main program loop.
func (a *app) Run() error {

	ui, err := view.NewConsoleUI()
	if err != nil {
		return err
	}
	a.ui = ui
	defer a.ui.Exit()

	authMode, err := a.getConfigInput()
	if err != nil {
		return a.displayErrorScreen(err)
	}

	ds, err := a.newDataverseService(authMode)
	if err != nil {
		return a.displayErrorScreen(err)
	}

	err = a.initialiseEntityServices(ds)
	if err != nil {
		return a.displayErrorScreen(err)
	}

	return a.startProgramLoop()
}

// getConfigInput displays the configuration screen and returns the selected
// authentication mode or an error.
func (a *app) getConfigInput() (authMode.AuthenticationMode, error) {
	configScreen, err := GetConfigScreen()
	if err != nil {
		return authMode.Invalid, err
	}

	output, err := a.ui.NavigateTo(configScreen)
	if err != nil {
		return authMode.Invalid, err
	}

	return authMode.AuthenticationMode(output.UserInput()), nil
}

// startProgramLoop displays the main menu and handles the main application flow
// until the user chooses to exit.
func (a *app) startProgramLoop() error {
	mainMenuChoice, err := a.displayMainMenu()
	if err != nil {
		return err
	}

	for mainMenuChoice != mainMenuOption.Exit {
		a.displayTable(mainMenuChoice)
		mainMenuChoice, err = a.displayMainMenu()
		if err != nil {
			return err
		}
	}
	return nil
}

// displayMainMenu shows the main menu screen and returns the selected option.
func (a *app) displayMainMenu() (mainMenuOption.MainMenuOption, error) {

	mainMenuScreen, err := newMainMenuScreen()
	if err != nil {
		return mainMenuOption.Invalid, err
	}

	output, err := a.ui.NavigateTo(mainMenuScreen)
	if err != nil {
		return mainMenuOption.Invalid, err
	}
	return mainMenuOption.MainMenuOption(output.UserInput()), nil
}

// displayTable shows the appropriate entity table based on the user's menu
// choice.
func (a *app) displayTable(tableChoice mainMenuOption.MainMenuOption) error {
	switch tableChoice {
	case mainMenuOption.Accounts:
		return a.displayAccountsMenu()
	case mainMenuOption.Contacts:
		return a.displayContactsMenu()
	}
	return nil
}

// displayAccountsMenu shows the accounts menu and handles account-related
// operations.
func (a *app) displayAccountsMenu() error {
	accountsMenu := entityMenu[*model.Account]{
		ui:          a.ui,
		service:     a.accountsService,
		listColumns: a.accountsListColumns,
		getNewEntity: func() (*model.Account, error) {
			defaultValues := model.Account{}
			return getEntityDetails(
				&defaultValues,
				"New Account",
				accountPropertyPrompts,
				a.getScreenOutput)
		},
		getUpdatedEntity: func(accountToUpdate *model.Account) (*model.Account, error) {
			return getEntityDetails(
				accountToUpdate,
				"New Account",
				accountPropertyPrompts,
				a.getScreenOutput)
		},
		entityLabel: "Account",
	}
	return accountsMenu.run()
}

// displayContactsMenu shows the contacts menu and handles contact-related
// operations.
func (a *app) displayContactsMenu() error {
	contactsMenu := entityMenu[*model.Contact]{
		ui:          a.ui,
		service:     a.contactsService,
		listColumns: a.contactsListColumns,
		getNewEntity: func() (*model.Contact, error) {
			defaultValues := model.Contact{}
			return getEntityDetails(
				&defaultValues,
				"New Contact",
				contactPropertyPrompts,
				a.getScreenOutput)
		},
		getUpdatedEntity: func(contactToUpdate *model.Contact) (*model.Contact, error) {
			return getEntityDetails(
				contactToUpdate,
				"New Contact",
				contactPropertyPrompts,
				a.getScreenOutput)
		},
		entityLabel: "Contact",
	}
	return contactsMenu.run()
}

// getScreenOutput is a helper method that abstracts the process of displaying a
// screen and retrieving its output.
func (a *app) getScreenOutput(getScreen func() (view.Screen, error)) (view.ScreenOutput, error) {
	s, err := getScreen()
	if err != nil {
		return nil, err
	}
	return a.ui.NavigateTo(s)
}

// newDataverseService creates a new Dataverse service using the specified
// authentication mode. It configures the appropriate client based on the mode
// and tests the connection.
func (a *app) newDataverseService(mode authMode.AuthenticationMode) (service.DataverseService, error) {
	var getClientFunc func(msal.ClientOptions) (msal.DataverseClient, error)

	switch mode {
	case authMode.Application:
		getClientFunc = msal.GetAppService
	case authMode.User:
		getClientFunc = msal.GetDelegatedService
	default:
		return nil, fmt.Errorf("invalid auth mode: %s", mode)
	}

	client, err := getClientFunc(msal.ClientOptions{
		ClientID:     a.config.ClientID,
		ResourceURL:  a.config.ResourceURL,
		Authority:    a.config.Authority,
		ClientSecret: a.config.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	dataverseService, err := service.NewDataverseService(service.DataverseServiceOptions{
		Client: client,
	})
	if err != nil {
		return nil, err
	}
	err = dataverseService.TestConnection()
	if err != nil {
		return nil, err
	}

	return dataverseService, nil
}

// initialiseEntityServices sets up the Account and Contact services with the
// provided Dataverse service.
func (a *app) initialiseEntityServices(dataverseService service.DataverseService) error {
	baseURL, err := url.Parse(a.config.APIBaseURL)
	if err != nil {
		return err
	}

	err = a.initAccountsService(dataverseService, baseURL)
	if err != nil {
		return err
	}

	return a.initContactsService(dataverseService, baseURL)
}

// initAccountsService initializes the service for working with account
// entities.
func (a *app) initAccountsService(dataverseService service.DataverseService, baseURL *url.URL) error {
	accountServiceOptions := service.EntityServiceOptions{
		DataverseService: dataverseService,
		BaseUrl:          baseURL,
		PageLimit:        a.config.PageLimit,
		ResourcePath:     logicalNames.TableAccountResource,
		SearchFields: []string{
			logicalNames.ColumnAccountName,
			logicalNames.ColumnAccountCity,
		},
		SelectsFields: []string{
			logicalNames.ColumnAccountId,
			logicalNames.ColumnAccountName,
			logicalNames.ColumnAccountCity,
		},
	}

	a.accountsService = service.NewEntityService[*model.Account](accountServiceOptions)
	return nil
}

// initContactsService initializes the service for working with contact
// entities.
func (a *app) initContactsService(dataverseService service.DataverseService, baseURL *url.URL) error {
	contactServiceOptions := service.EntityServiceOptions{
		DataverseService: dataverseService,
		BaseUrl:          baseURL,
		PageLimit:        a.config.PageLimit,
		ResourcePath:     logicalNames.TableContactResource,
		SearchFields: []string{
			logicalNames.ColumnContactFirstName,
			logicalNames.ColumnContactLastName,
			logicalNames.ColumnContactEmail,
		},
		SelectsFields: []string{
			logicalNames.ColumnContactId,
			logicalNames.ColumnContactFirstName,
			logicalNames.ColumnContactLastName,
			logicalNames.ColumnContactEmail,
		},
	}

	a.contactsService = service.NewEntityService[*model.Contact](contactServiceOptions)
	return nil
}

// initEntityListColumns initializes the column definitions for both account and
// contact lists.
func (a *app) initEntityListColumns() error {
	err := a.initAccountListColumns()
	if err != nil {
		return err
	}
	return a.initContactListColumns()
}

// initAccountListColumns initializes the column definitions for the account
// list view.
func (a *app) initAccountListColumns() error {
	columns, err := model.AccountListColumns()
	if err != nil {
		return err
	}
	a.accountsListColumns = columns
	return nil
}

// initContactListColumns initializes the column definitions for the contact
// list view.
func (a *app) initContactListColumns() error {
	columns, err := model.ContactListColumns()
	if err != nil {
		return err
	}
	a.contactsListColumns = columns
	return nil
}

// displayErrorScreen shows an error message to the user.
// Returns the original error if displaying the error screen fails,
// otherwise returns nil to indicate the error was displayed successfully.
func (a *app) displayErrorScreen(originalError error) error {
	if a.ui == nil {
		return originalError
	}

	es, err := newErrorScreen(originalError.Error())
	if err != nil {
		return originalError
	}
	_, err = a.ui.NavigateTo(es)
	if err != nil {
		return originalError
	}
	a.ui.Exit()
	return nil
}
