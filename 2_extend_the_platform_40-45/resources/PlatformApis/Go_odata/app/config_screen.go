package app

import (
	"github.com/turnerbenjamin/go_odata/constants/auth_mode"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func GetConfigScreen() (view.Screen, error) {

	menu, err := view.NewMenuComponent([]string{
		string(auth_mode.Application),
		string(auth_mode.User),
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
