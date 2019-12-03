package parser

import (
	"io"
	"strings"

	"github.com/alongubkin/cpl-compiler/pkg/lexer"
)

// Parser represents a CPL parser.
type Parser struct {
	scanner   *lexer.Scanner
	lookahead lexer.Token
}

// NewParser returns a new instance of Parser.
func NewParser(reader io.Reader) *Parser {
	scanner := lexer.NewScanner(reader)
	return &Parser{
		scanner:   scanner,
		lookahead: scanner.Scan(),
	}
}

// Parse parses a CPL program and returns its AST representation.
func Parse(s string) (*Program, error) {
	return NewParser(strings.NewReader(s)).ParseProgram()
}

func (p *Parser) match(tokenTypes ...lexer.TokenType) (*lexer.Token, bool) {
	for _, tokType := range tokenTypes {
		if tokType == p.lookahead.TokenType {
			token := p.lookahead
			p.lookahead = p.scanner.Scan()
			return &token, true
		}
	}

	return &p.lookahead, false
}

// ParseProgram parses a CPL program and returns a Program AST object.
// 	program -> declarations stmt_block
func (p *Parser) ParseProgram() (*Program, error) {
	program := &Program{}

	// Parse declarations.
	if declarations, err := p.ParseDeclarations(); err == nil {
		program.Declarations = declarations
	} else {
		return nil, err
	}

	// Make sure there's an EOF at the end of the file.
	if token, ok := p.match(lexer.EOF); !ok {
		return nil, newParseError(token.Lexeme, []string{"EOF"}, token.Position)
	}

	return program, nil
}

// ParseDeclarations parses a list of declarations and returns a Declaration AST array.
// 	declarations -> declaration declarations | ε
func (p *Parser) ParseDeclarations() ([]Declaration, error) {
	declarations := []Declaration{}
	for p.lookahead.TokenType == lexer.ID {
		if declaration, err := p.ParseDeclaration(); err == nil {
			declarations = append(declarations, *declaration)
		} else {
			return nil, err
		}
	}

	return declarations, nil
}

// ParseDeclaration parses a declaration and returns a Declaration AST object.
// 	declaration -> idlist ':' type ';'
func (p *Parser) ParseDeclaration() (*Declaration, error) {
	declaration := &Declaration{}

	if idlist, err := p.ParseIDList(); err == nil {
		declaration.Names = idlist
	} else {
		return nil, err
	}

	if token, ok := p.match(lexer.COLON); !ok {
		return nil, newParseError(token.Lexeme, []string{":"}, token.Position)
	}

	if datatype, err := p.ParseType(); err == nil {
		declaration.Type = datatype
	} else {
		return nil, err
	}

	if token, ok := p.match(lexer.SEMICOLON); !ok {
		return nil, newParseError(token.Lexeme, []string{";"}, token.Position)
	}

	return declaration, nil
}

// ParseType parses a type returns it as a DataType.
// 	type -> INT | FLOAT
func (p *Parser) ParseType() (DataType, error) {
	token, ok := p.match(lexer.INT, lexer.FLOAT)
	if !ok {
		return Unknown, newParseError(token.Lexeme, []string{"int", "float"}, token.Position)
	}

	switch token.TokenType {
	case lexer.INT:
		return Integer, nil
	case lexer.FLOAT:
		return Float, nil
	default:
		panic("Unknown token type")
	}
}

// ParseIDList parses a list of IDs and returns a string array.
// 	idlist -> ID idlist'
// 	idlist' -> ',' ID idlist' | ε
func (p *Parser) ParseIDList() ([]string, error) {
	names := []string{}

	// Parse the first name
	if token, ok := p.match(lexer.ID); ok {
		names = append(names, token.Lexeme)
	} else {
		return nil, newParseError(token.Lexeme, []string{"ID"}, token.Position)
	}

	// Parse other names if exist
	for p.lookahead.TokenType == lexer.COMMA {
		p.match(lexer.COMMA)

		if token, ok := p.match(lexer.ID); ok {
			names = append(names, token.Lexeme)
		} else {
			return nil, newParseError(token.Lexeme, []string{"ID"}, token.Position)
		}
	}

	return names, nil
}
