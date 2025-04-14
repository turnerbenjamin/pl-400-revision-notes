package authMode

type AuthenticationMode string

const (
	Application AuthenticationMode = "Application"
	User        AuthenticationMode = "User"
	Invalid     AuthenticationMode = "Invalid"
)
