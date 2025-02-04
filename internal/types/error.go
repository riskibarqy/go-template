package types

// Error represents customized error object
type Error struct {
	Path     string
	Message  string
	Error    error
	Type     string
	IsIgnore bool
}
