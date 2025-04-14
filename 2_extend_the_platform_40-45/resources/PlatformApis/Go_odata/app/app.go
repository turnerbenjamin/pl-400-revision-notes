package app

import (
	"errors"
	"log"
	"net/url"

	"github.com/turnerbenjamin/go_odata/constants/auth_mode"
	"github.com/turnerbenjamin/go_odata/constants/logical_names"
	"github.com/turnerbenjamin/go_odata/constants/main_menu_option"
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

func (a *app) getConfigInput() (auth_mode.AuthenticationMode, error) {
	configScreen, err := GetConfigScreen()
	if err != nil {
		return auth_mode.Invalid, err
	}

	output, err := a.ui.NavigateTo(configScreen)
	if err != nil {
		return auth_mode.Invalid, err
	}

	return auth_mode.AuthenticationMode(output.UserInput()), nil
}

func (a *app) startProgramLoop() error {
	mainMenuChoice, err := a.displayMainMenu()
	if err != nil {
		return err
	}

	for mainMenuChoice != main_menu_option.Exit {
		a.displayTable(mainMenuChoice)
		mainMenuChoice, err = a.displayMainMenu()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *app) displayMainMenu() (main_menu_option.MainMenuOption, error) {

	mainMenuScreen, err := GetMainMenuScreen()
	if err != nil {
		return main_menu_option.Invalid, err
	}

	output, err := a.ui.NavigateTo(mainMenuScreen)
	if err != nil {
		return main_menu_option.Invalid, err
	}
	return main_menu_option.MainMenuOption(output.UserInput()), nil
}

func (a *app) displayTable(tableChoice main_menu_option.MainMenuOption) error {
	switch tableChoice {
	case main_menu_option.Accounts:
		a.displayAccountsMenu()
	case main_menu_option.Contacts:
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
		entityLabel: "Contact",
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

func (a *app) InitialiseDataverseService(mode auth_mode.AuthenticationMode) error {

	var getClientFunc func(msal.ClientOptions) (msal.DataverseClient, error)

	switch mode {
	case auth_mode.Application:
		getClientFunc = msal.GetAppService
	case auth_mode.User:
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
		Client: client,
	})

	if err != nil {
		return err
	}

	err = dataverseService.TestConnection()
	if err != nil {
		return err
	}

	a.dataverseService = dataverseService

	baseUrl, err := url.Parse(a.config.APIBaseUrl)
	if err != nil {
		return err
	}

	baseEntityServiceOptions := &service.EntityServiceOptions{
		DataverseService: dataverseService,
		BaseUrl:          baseUrl,
		PageLimit:        a.config.PageLimit,
	}

	accountServiceOptions := *baseEntityServiceOptions
	accountServiceOptions.ResourcePath = logical_names.TableAccountResource
	accountServiceOptions.SearchFields = []string{
		logical_names.ColumnAccountName,
		logical_names.ColumnAccountCity,
	}
	accountServiceOptions.SelectsFields = append(accountServiceOptions.SearchFields, logical_names.ColumnAccountId)
	a.accountsService = service.NewEntityService[*model.Account](accountServiceOptions)

	contactServiceOptions := *baseEntityServiceOptions
	contactServiceOptions.ResourcePath = logical_names.TableContactResource
	contactServiceOptions.SearchFields = []string{
		logical_names.ColumnContactFirstName,
		logical_names.ColumnContactLastName,
		logical_names.ColumnContactEmail,
	}

	contactServiceOptions.SelectsFields = append(contactServiceOptions.SearchFields, logical_names.ColumnContactId)
	a.contactsService = service.NewEntityService[*model.Contact](contactServiceOptions)
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
