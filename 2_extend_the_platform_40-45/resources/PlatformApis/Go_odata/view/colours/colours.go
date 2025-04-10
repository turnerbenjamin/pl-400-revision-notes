package colours

type Color string

const (
	PURPLE       Color = "\033[38;5;127m"
	RED          Color = "\033[38;5;196m"
	ORANGE       Color = "\033[38;5;208m"
	GREY         Color = "\033[38;5;238m"
	GREEN        Color = "\033[38;5;120m"
	RESET        Color = "\033[0m"
	SELECTED_ROW Color = "\033[48;5;166m"
)
