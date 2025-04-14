package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func newStringInputScreen(title, text, propertyName, value string, isRequired bool) (view.Screen, error) {
	return view.MakeScreen([]view.Component{
		view.NewTitleComponent(title, colours.Purple),
		view.NewTextComponent(text),
		view.NewStringInputComponent(propertyName, value, isRequired),
	})
}
