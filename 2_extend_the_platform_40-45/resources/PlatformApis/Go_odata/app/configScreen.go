package app

import (
	"github.com/turnerbenjamin/go_odata/constants/authMode"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func GetConfigScreen() (view.Screen, error) {

	menu, err := view.NewMenuComponent([]string{
		string(authMode.Application),
		string(authMode.User),
	},
	)

	if err != nil {
		return nil, err
	}

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("CONFIGURATION", colours.Purple),
		view.NewTextComponent("Select authentication mode"),
		menu,
	})
}
