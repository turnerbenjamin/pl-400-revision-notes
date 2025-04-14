package table_menu_option

type TableMenuOption string

const (
	Search TableMenuOption = "Search"
	Create TableMenuOption = "Create"
	Update TableMenuOption = "Update"
	Delete TableMenuOption = "Delete"
	Back   TableMenuOption = "Back"
)
