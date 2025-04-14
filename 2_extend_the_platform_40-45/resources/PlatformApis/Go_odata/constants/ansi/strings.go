package ansi

const (
	CursorHide = "\033[?25l"
	CursorShow = "\033[?25h"
	CursorHome = "\033[H"

	ClearScreen     = "\033[2J"
	ClearScrollback = "\033[3J"
	ClearToEnd      = "\033[J"

	ClearAll  = "\033[H\033[2J\033[3J"
	ResetView = "\033[H\033[J"
)
