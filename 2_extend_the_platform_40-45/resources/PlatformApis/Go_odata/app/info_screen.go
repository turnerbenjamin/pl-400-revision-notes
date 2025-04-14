package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func newInfoScreen(msg string) (view.Screen, error) {
	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("INFO", colours.Blue),
		view.NewTextComponent(msg),
		view.NewAnyKeyToContinueComponent(),
	})
}
