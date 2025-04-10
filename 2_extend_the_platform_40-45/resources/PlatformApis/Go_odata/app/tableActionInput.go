package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

type tableAction string

const (
	TABLE_ACTION_CREATE tableAction = "Create"
	TABLE_ACTION_LIST   tableAction = "List"
	TABLE_ACTION_BACK   tableAction = "Back"
)

func GetTableActionMenuScreen(tableName string) view.Screen {
	return view.MakeScreen([]view.Component{
		view.BuildTitleComponent(tableName, colours.PURPLE),
		view.BuildTextComponent("Choose an action"),
		view.BuildMenuComponent([]string{string(TABLE_ACTION_CREATE), string(TABLE_ACTION_LIST), string(TABLE_ACTION_BACK)}),
	})
}
