package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

type authenticationMode string

const (
	CONFIG_APPLICATION_MODE authenticationMode = "Application"
	CONFIG_USER_MODE        authenticationMode = "User"
)

func GetConfigScreen() view.Screen {
	return view.MakeScreen([]view.Component{
		view.BuildTitleComponent("CONFIGURATION", colours.PURPLE),
		view.BuildTextComponent("Select authentication mode"),
		view.BuildMenuComponent([]string{string(CONFIG_APPLICATION_MODE), string(CONFIG_USER_MODE)}),
	})
}
