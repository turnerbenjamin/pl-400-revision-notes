package app

import (
	"github.com/turnerbenjamin/go_odata/view"
	"github.com/turnerbenjamin/go_odata/view/colours"
)

func getStringInputScreen(title, text, propertyName string, isRequired bool) view.Screen {
	return view.MakeScreen([]view.Component{
		view.BuildTitleComponent(title, colours.PURPLE),
		view.BuildTextComponent(text),
		view.BuildStringInputComponent(propertyName, isRequired),
	})
}
