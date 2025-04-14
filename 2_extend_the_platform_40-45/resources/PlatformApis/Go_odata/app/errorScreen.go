package app

import (
	"strings"

	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func getErrorScreen(msg string) (view.Screen, error) {
	fm := strings.ReplaceAll(msg, ". ", ".\n")

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("ERROR", colours.Red),
		view.NewTextComponent(fm),
		view.NewAnyKeyToContinueComponent(),
	})
}
