package app

import (
	"github.com/turnerbenjamin/go_odata/constants/table_menu_option"
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

var entityListControls = []view.ListControl{
	listControl{
		label: "Set/Clear search term",
		value: string(table_menu_option.Search),
		key:   's',
	},
	listControl{
		label: "Create",
		value: string(table_menu_option.Create),
		key:   'c',
	},
	listControl{
		label: "Update",
		value: string(table_menu_option.Update),
		key:   'u',
	},
	listControl{
		label: "Delete",
		value: string(table_menu_option.Delete),
		key:   'd',
	},
	listControl{
		label: "Back to main menu",
		value: string(table_menu_option.Back),
		key:   'b',
	},
}

type listScreenOptions[T view.Entity] struct {
	entityList view.EntityList[T]
	columns    []view.ListColumn[T]
}

func NewEntityListScreen[T view.Entity](entityLabel string, listScreenOptions listScreenOptions[T]) (view.Screen, error) {

	listOptions := view.ListComponentOptions[T]{
		Controls:   entityListControls,
		EntityList: listScreenOptions.entityList,
		Columns:    listScreenOptions.columns,
	}

	listComponent, err := view.BuildListComponent(listOptions)
	if err != nil {
		return nil, err
	}

	return view.MakeScreen([]view.Component{
		view.NewTitleComponent(entityLabel, colours.Purple),
		listComponent,
	})
}
