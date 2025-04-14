package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func newSuccessScreen(msg string) (view.Screen, error) {
	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("SUCCESS", colours.Green),
		view.NewTextComponent(msg),
		view.NewAnyKeyToContinueComponent(),
	})
}
