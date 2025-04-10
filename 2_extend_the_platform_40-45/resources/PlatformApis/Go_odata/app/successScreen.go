package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func getSuccessScreen(msg string) view.Screen {
	return view.MakeScreen([]view.Component{
		view.BuildTitleComponent("SUCCESS", colours.GREEN),
		view.BuildTextComponent(msg),
		view.BuildAnyKeyToContinueComponent(),
	})
}
