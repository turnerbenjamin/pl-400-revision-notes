package app

import (
	"github.com/turnerbenjamin/go_odata/constants/confirmoption"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func NewConfirmationScreen(msg string) (view.Screen, error) {

	menu, err := view.NewMenuComponent([]string{
		string(confirmoption.Yes),
		string(confirmoption.No),
	})

	if err != nil {
		return nil, err
	}

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent("Confirmation Required", colours.Purple),
		view.NewTextComponent(msg),
		menu,
	},
	)
}
