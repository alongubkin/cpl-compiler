package lexer

// Token represents a lexical token.
type TokenType int

// CPL's tokens
const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Symbols
	LPAREN    // (
	RPAREN    // )
	LBRACKET  // {
	RBRACKET  // }
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	EQUALS    // =

	// Keywords
	BREAK
	CASE
	DEFAULT
	ELSE
	FLOAT
	IF
	INPUT
	INT
	OUTPUT
	STATICCAST
	SWITCH
	WHILE

	// Operators
	RELOP // == | != | < | > | >= | <=
	ADDOP // + | -
	MULOP // * | /
	OR    // ||
	AND   // &&
	NOT   // !

	// Literals
	ID
	NUM
)

// Position specifies the line and character position of a token.
// The Column and Line are both zero-based indexes.
type Position struct {
	Line   int
	Column int
}

type Token struct {
	TokenType TokenType
	Lexeme    string
	Position  Position
}
