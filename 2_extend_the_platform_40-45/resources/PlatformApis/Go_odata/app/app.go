package app

import (
	"errors"
	"log"

	"github.com/turnerbenjamin/go_odata/constants/authMode"
	"github.com/turnerbenjamin/go_odata/constants/mainMenuOption"
	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/msal"
	"github.com/turnerbenjamin/go_odata/service"
	"github.com/turnerbenjamin/go_odata/view"
)

var ErrInvalidMenuSelection = errors.New("menu selection input is invalid")

type AppConfig struct {
	ClientId     string
	TenantId     string
	ClientSecret string
	ResourceUrl  string
	APIBaseUrl   string
	Authority    string
	PageLimit    int
}

type App interface {
	Run() error
}
type app struct {
	config              *AppConfig
	dataverseService    service.DataverseService
	accountsService     service.EntityService[*model.Account]
	contactsService     service.EntityService[*model.Contact]
	accountsListColumns []view.ListColumn[*model.Account]
	contactsListColumns []view.ListColumn[*model.Contact]
	ui                  view.UI
}

func Create(config AppConfig) (App, error) {
	a := app{
		config: &config,
	}

	err := a.initAccountListColumns()
	if err != nil {
		return nil, err
	}

	err = a.initContactListColumns()
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (a *app) Run() error {

	ui, err := view.NewConsoleUI()
	if err != nil {
		return err
	}
	a.ui = ui
	defer a.ui.Exit()

	authMode, err := a.getConfigInput()
	if err != nil {
		return err
	}

	err = a.InitialiseDataverseService(authMode)
	if err != nil {
		return err
	}

	return a.startProgramLoop()
}

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

func (a *app) displayMainMenu() (mainMenuOption.MainMenuOption, error) {

	mainMenuScreen, err := GetMainMenuScreen()
	if err != nil {
		return mainMenuOption.Invalid, err
	}

	output, err := a.ui.NavigateTo(mainMenuScreen)
	if err != nil {
		return mainMenuOption.Invalid, err
	}
	return mainMenuOption.MainMenuOption(output.UserInput()), nil
}

func (a *app) displayTable(tableChoice mainMenuOption.MainMenuOption) error {
	switch tableChoice {
	case mainMenuOption.Accounts:
		a.displayAccountsMenu()
	case mainMenuOption.Contacts:
		a.displayContactsMenu()
	}
	return nil
}

func (a *app) displayAccountsMenu() error {
	accountsMenu := entityMenu[*model.Account]{
		ui:          a.ui,
		service:     a.accountsService,
		listColumns: a.accountsListColumns,
		getNewEntity: func() (*model.Account, error) {
			return a.getAccountDetails(nil)
		},
		getUpdatedEntity: func(accountToUpdate *model.Account) (*model.Account, error) {
			return a.getAccountDetails(accountToUpdate)
		},
		entityLabel: "Account",
	}
	return accountsMenu.run()
}

func (a *app) displayContactsMenu() error {
	contactsMenu := entityMenu[*model.Contact]{
		ui:          a.ui,
		service:     a.contactsService,
		listColumns: a.contactsListColumns,
		getNewEntity: func() (*model.Contact, error) {
			return a.getContactDetails(nil)
		},
		getUpdatedEntity: func(accountToUpdate *model.Contact) (*model.Contact, error) {
			return a.getContactDetails(accountToUpdate)
		},
		entityLabel: "Account",
	}
	return contactsMenu.run()
}

func (a *app) getAccountDetails(defaultValues *model.Account) (*model.Account, error) {
	defaultName := ""
	defaultCity := ""

	if defaultValues != nil {
		defaultName = defaultValues.Name
		defaultCity = defaultValues.City
	}

	enterNameOutput, err := a.getScreenOutput(func() (view.Screen, error) {
		return newStringInputScreen("New Account", "Enter account name", "Name", defaultName, true)
	})
	if err != nil {
		return nil, err
	}

	enterCityOutput, err := a.getScreenOutput(func() (view.Screen, error) {
		return newStringInputScreen("New Account", "Enter city", "City", defaultCity, false)
	})
	if err != nil {
		return nil, err
	}

	return &model.Account{
		Name: enterNameOutput.UserInput(),
		City: enterCityOutput.UserInput(),
	}, nil
}

func (a *app) getContactDetails(defaultValues *model.Contact) (*model.Contact, error) {
	defaultFirstName := ""
	defaultLastName := ""
	defaultEmail := ""

	if defaultValues != nil {
		defaultFirstName = defaultValues.FirstName
		defaultLastName = defaultValues.LastName
		defaultEmail = defaultValues.Email
	}

	firstNameOutput, err := a.getScreenOutput(func() (view.Screen, error) {
		return newStringInputScreen("New Contact", "Enter contact first name", "FirstName", defaultFirstName, true)
	})
	if err != nil {
		return nil, err
	}

	lastNameOutput, err := a.getScreenOutput(func() (view.Screen, error) {
		return newStringInputScreen("New Contact", "Enter contact last name", "LastName", defaultLastName, true)
	})
	if err != nil {
		return nil, err
	}

	emailOutput, err := a.getScreenOutput(func() (view.Screen, error) {
		return newStringInputScreen("New Contact", "Enter contact email", "Email", defaultEmail, false)
	})
	if err != nil {
		return nil, err
	}

	return &model.Contact{
		FirstName: firstNameOutput.UserInput(),
		LastName:  lastNameOutput.UserInput(),
		Email:     emailOutput.UserInput(),
	}, nil
}

func (a *app) getScreenOutput(getScreen func() (view.Screen, error)) (view.ScreenOutput, error) {
	s, err := getScreen()
	if err != nil {
		return nil, err
	}
	return a.ui.NavigateTo(s)
}

func (a *app) InitialiseDataverseService(mode authMode.AuthenticationMode) error {

	var getClientFunc func(msal.ClientOptions) (msal.DataverseClient, error)

	switch mode {
	case authMode.Application:
		getClientFunc = msal.GetAppService
	case authMode.User:
		getClientFunc = msal.GetDelegatedService
	default:
		log.Panicf("invalid auth mode: %s", mode)
	}

	client, err := getClientFunc(msal.ClientOptions{
		ClientId:     a.config.ClientId,
		ResourceUrl:  a.config.ResourceUrl,
		Authority:    a.config.Authority,
		ClientSecret: a.config.ClientSecret,
	})
	if err != nil {
		return err
	}

	dataverseService, err := service.NewDataverseService(service.DataverseServiceOptions{
		Client:  client,
		BaseUrl: a.config.APIBaseUrl,
	})

	if err != nil {
		return err
	}

	a.dataverseService = dataverseService
	a.accountsService = service.NewAccountService(dataverseService, a.config.APIBaseUrl, a.config.PageLimit)
	a.contactsService = service.NewContactService(dataverseService, a.config.APIBaseUrl, a.config.PageLimit)
	return nil
}

func (a *app) initAccountListColumns() error {
	columns, err := model.AccountListColumns()
	if err != nil {
		return err
	}
	a.accountsListColumns = columns
	return nil
}

func (a *app) initContactListColumns() error {
	columns, err := model.ContactListColumns()
	if err != nil {
		return err
	}
	a.contactsListColumns = columns
	return nil
}
