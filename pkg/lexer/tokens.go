package lexer

// Token represents a lexical token.
type Token int

// CPL's tokens
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF

	// Symbols
	LPAREN    // (
	RPAREN    // )
	LBRACKET  // {
	RBRACKET  // }
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	EQ        // =

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
