package app

import (
	"strings"

	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func getErrorScreen(msg string) view.Screen {
	fm := strings.ReplaceAll(msg, ". ", ".\n")

	return view.MakeScreen([]view.Component{
		view.BuildTitleComponent("ERROR", colours.RED),
		view.BuildTextComponent(fm),
		view.BuildAnyKeyToContinueComponent(),
	})
}
