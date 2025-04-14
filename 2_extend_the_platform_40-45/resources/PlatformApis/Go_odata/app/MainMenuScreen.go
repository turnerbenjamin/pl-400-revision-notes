package app

import (
	"github.com/turnerbenjamin/go_odata/constants/mainMenuOption"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func GetMainMenuScreen() (view.Screen, error) {

	menu, err := view.NewMenuComponent([]string{
		string(mainMenuOption.Accounts),
		string(mainMenuOption.Contacts),
		string(mainMenuOption.Exit),
	})

	if err != nil {
		return nil, err
	}

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("Table Selection", colours.Purple),
		view.NewTextComponent("Choose a table"),
		menu,
	},
	)
}
