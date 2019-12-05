package codegen

// Error represents an error that occurred during code generation.
type Error struct {
	Message string
	// TODO: add position
	// Pos      lexer.Position
}

// Error returns the string representation of the error.
func (e *Error) Error() string {
	return e.Message
}
