package app

import (
	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/service"
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

type listControl struct {
	label string
	value string
	key   rune
}

func (lc listControl) GetLabel() string {
	return lc.label
}

func (lc listControl) GetValue() string {
	return lc.value
}

func (lc listControl) GetKey() rune {
	return lc.key
}

func GetAccountListScreen(results service.EntityList[model.Account]) view.Screen {

	controls := []view.ListControl{
		listControl{
			label: "Update selected",
			value: "Update",
			key:   'u',
		},
		listControl{
			label: "Delete selected",
			value: "Delete",
			key:   'd',
		},
		listControl{
			label: "Back to account menu",
			value: "back",
			key:   'b',
		},
	}

	lc := view.CreateListColumns[model.Account]()
	lc = lc.WithColumn("Name", func(a model.Account) string {
		return a.Name
	})
	lc = lc.WithColumn("City", func(a model.Account) string {
		return a.City
	})

	return view.MakeScreen([]view.Component{
		view.BuildTitleComponent("ACCOUNTS", colours.PURPLE),
		view.BuildListComponent(controls, lc, results.GetData(), results.HasNext(), false),
	})
}
