package app

import (
	"github.com/turnerbenjamin/go_odata/constants/main_menu_option"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func GetMainMenuScreen() (view.Screen, error) {

	menu, err := view.NewMenuComponent([]string{
		string(main_menu_option.Accounts),
		string(main_menu_option.Contacts),
		string(main_menu_option.Exit),
	})

	if err != nil {
		return nil, err
	}

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("Table Selection", colours.Purple),
		view.NewTextComponent("Choose a table"),
		menu})
}
