package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

type tableMode string

const (
	TABLE_MODE_ACCOUNT tableMode = "Account"
	TABLE_MODE_CONTACT tableMode = "Contact"
	TABLE_MODE_EXIT    tableMode = "Exit"
)

func GetTableMenuScreen() view.Screen {
	return view.MakeScreen([]view.Component{
		view.BuildTitleComponent("Table Selection", colours.PURPLE),
		view.BuildTextComponent("Choose a table"),
		view.BuildMenuComponent([]string{string(TABLE_MODE_ACCOUNT), string(TABLE_MODE_CONTACT), string(TABLE_MODE_EXIT)}),
	})
}
