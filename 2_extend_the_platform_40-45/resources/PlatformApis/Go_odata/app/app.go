package app

import (
	"fmt"
	"log"

	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/msal"
	"github.com/turnerbenjamin/go_odata/service"
	"github.com/turnerbenjamin/go_odata/view"
)

type App interface {
	Run()
}
type app struct {
	config           *AppConfig
	dataverseService msal.DataverseService
	accountsService  service.AccountsService
	ui               view.UI
}

func Create(config AppConfig) App {
	return &app{
		config: &config,
	}
}

func (a *app) Run() {

	a.ui = view.NewConsoleUI()
	defer a.ui.Exit()

	authMode := a.getConfigInput()
	a.InitialiseDataverseService(authMode)
	a.startProgramLoop()
}

func (a *app) startProgramLoop() {
	tableChoice := a.getTableChoiceInput()
	for tableChoice != TABLE_MODE_EXIT {
		a.startTableLoop(tableChoice)
		tableChoice = a.getTableChoiceInput()
	}
}

func (a *app) startTableLoop(tableChoice tableMode) {

	tableAction := a.getTableActionInput(tableChoice)
	for tableAction != TABLE_ACTION_BACK {
		a.performTableAction(tableChoice, tableAction)
		tableAction = a.getTableActionInput(tableChoice)
	}
}

func (a *app) getConfigInput() authenticationMode {
	output := a.ui.NavigateTo(GetConfigScreen())
	return authenticationMode(output.GetUserInput())
}

func (a *app) getTableChoiceInput() tableMode {
	output := a.ui.NavigateTo(GetTableMenuScreen())
	return tableMode(output.GetUserInput())
}

func (a *app) getTableActionInput(tableChoice tableMode) tableAction {
	output := a.ui.NavigateTo(GetTableActionMenuScreen(string(tableChoice)))
	return tableAction(output.GetUserInput())
}

func (a *app) performTableAction(table tableMode, action tableAction) {
	switch table {
	case TABLE_MODE_ACCOUNT:
		a.performAccountAction(action)
	case TABLE_MODE_CONTACT:
		a.performContactAction(action)
	default:
		log.Fatalf("Invalid table selection: %s", string(table))
	}
}

func (a *app) performAccountAction(action tableAction) {
	switch action {
	case TABLE_ACTION_CREATE:
		a.createAccountControl()
	case TABLE_ACTION_LIST:
		a.listAccountsControl()
	default:
		log.Fatalf("Invalid table action: %s", string(action))
	}
}

func (a *app) performContactAction(action tableAction) {
	switch action {
	case TABLE_ACTION_CREATE:
	default:
		log.Fatalf("Invalid table action: %s", string(action))
	}
}

func (a *app) createAccountControl() {
	nis := getStringInputScreen("CREATE ACCOUNT", "Enter account name", "Name", true)
	n := a.ui.NavigateTo(nis).GetUserInput()
	cis := getStringInputScreen("CREATE ACCOUNT", "Enter account city", "City", false)
	c := a.ui.NavigateTo(cis).GetUserInput()
	account, err := a.accountsService.Create(&model.Account{Name: n, City: c})

	if err != nil {
		es := getErrorScreen(err.Error())
		a.ui.NavigateTo(es)
		return
	}

	ss := getSuccessScreen(fmt.Sprintf("Created account: %s", account.Name))
	a.ui.NavigateTo(ss)
}

func (a *app) listAccountsControl() {

	res, err := a.accountsService.List("", 1)
	for {
		if err != nil {
			es := getErrorScreen(err.Error())
			a.ui.NavigateTo(es)
			break
		}
		ls := GetAccountListScreen(res)
		a.ui.NavigateTo(ls)
		if !res.HasNext() {
			break
		}
		res, err = res.GetNext()
	}

}

func (a *app) InitialiseDataverseService(authMode authenticationMode) {

	var getServiceFunc func(msal.ClientConfig) (msal.DataverseService, error)

	switch authMode {
	case CONFIG_APPLICATION_MODE:
		getServiceFunc = msal.GetAppService
	case CONFIG_USER_MODE:
		getServiceFunc = msal.GetDelegatedService
	default:
		log.Panicf("invalid auth mode: %s", authMode)
	}

	dvs, err := getServiceFunc(msal.ClientConfig{
		ClientId:     a.config.ClientId,
		ResourceUrl:  a.config.ResourceUrl,
		APIBaseUrl:   a.config.APIBaseUrl,
		Authority:    a.config.Authority,
		ClientSecret: a.config.ClientSecret,
	})

	if err != nil {
		log.Panic(err.Error())
	}

	err = dvs.Connect()
	if err != nil {
		log.Fatalf("Unable to connect: %s", err.Error())
	}
	a.dataverseService = dvs
	a.accountsService = service.NewAccountService(dvs)
}
