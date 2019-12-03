package parser

import (
	"io"
	"strings"

	"github.com/alongubkin/cpl-compiler/pkg/lexer"
)

// Parser represents a CPL parser.
type Parser struct {
	scanner *lexer.Scanner
	// lookahead
}

// NewParser returns a new instance of Parser.
func NewParser(reader io.Reader) *Parser {
	return &Parser{scanner: lexer.NewScanner(reader)}
}

// Parse parses a CPL program and returns its AST representation.
func Parse(s string) (*Program, error) {
	return NewParser(strings.NewReader(s)).ParseProgram()
}

// TODO: USE LOOKAHEAD

// ParseProgram parses a CPL program and returns a Program AST object.
func (p *Parser) ParseProgram() (*Program, error) {
	program := &Program{
		declarations: []Declaration{},
	}

	declarations, err := p.ParseDeclarations()
	if err != nil {
		return nil, err
	}

	eof := p.scanner.Scan()
	if eof.TokenType != lexer.EOF {
		return nil, newParseError(eof.Lexeme, []string{"EOF"}, eof.Position)
	}

	return program, nil
}

// ParseDeclaration parses a declaration and returns a Declaration AST object.
func (p *Parser) ParseDeclaration() (*Declaration, error) {
	id := p.scanner.Scan()
	if id.TokenType != lexer.ID {
		return nil, newParseError(id.Lexeme, []string{"ID"}, id.Position)
	}

	declaration := &Declaration{
		Names: []string{id.Lexeme},
	}

	for {
		token := p.scanner.Scan()
		switch token.TokenType {
		case lexer.COMMA:
			id = p.scanner.Scan()
			if id.TokenType != lexer.ID {
				return nil, newParseError(id.Lexeme, []string{"ID"}, id.Position)
			}

			declaration.Names = append(declaration.Names, id.Lexeme)
		case lexer.COLON:
			break

		default:
			return nil, newParseError(token.Lexeme, []string{", or :"}, token.Position)
		}
	}

	typeToken := p.scanner.Scan()
	switch typeToken.TokenType {
	case lexer.INT:
		declaration.Type = Integer
	case lexer.FLOAT:
		declaration.Type = Float
	default:
		return nil, newParseError(typeToken.Lexeme, []string{"int, float"}, typeToken.Position)
	}

	return declaration, nil
}

// // ParseIDList parses a list of IDs and returns a string array.
// func (p *Parser) ParseIDList() ([]string, error) {

// }
