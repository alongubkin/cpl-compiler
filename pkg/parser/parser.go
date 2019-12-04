package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alongubkin/cpl-compiler/pkg/lexer"
)

// Parser represents a CPL parser.
type Parser struct {
	Errors    []ParseError
	scanner   *lexer.Scanner
	lookahead lexer.Token
}

// NewParser returns a new instance of Parser.
func NewParser(reader io.Reader) *Parser {
	scanner := lexer.NewScanner(reader)
	return &Parser{
		Errors:    []ParseError{},
		scanner:   scanner,
		lookahead: scanner.Scan(),
	}
}

// Parse parses a CPL program and returns its AST representation.
func Parse(s string) (*Program, []ParseError) {
	parser := NewParser(strings.NewReader(s))
	return parser.ParseProgram(), parser.Errors
}

func (p *Parser) matchToken(tokenTypes ...lexer.TokenType) (*lexer.Token, bool) {
	for _, tokType := range tokenTypes {
		if tokType == p.lookahead.TokenType {
			token := p.lookahead
			p.lookahead = p.scanner.Scan()
			return &token, true
		}
	}

	return &p.lookahead, false
}

func (p *Parser) match(tokenTypes ...lexer.TokenType) (*lexer.Token, bool) {
	// Try to find the requested token.
	if token, ok := p.matchToken(tokenTypes...); ok {
		return token, true
	}

	nextRealToken := p.lookahead

	// If no such token was found, skip tokens until a correct one was found.
	skips := 0
	for {
		token, ok := p.matchToken(tokenTypes...)
		if ok {
			// A token was found! Continue parsing from here.
			// TODO: Add support for token types
			p.addError(newParseError(nextRealToken.Lexeme, []string{}, nextRealToken.Position))
			return token, true
		} else if token.TokenType == lexer.EOF {
			break
		}

		// Skip token.
		p.lookahead = p.scanner.Scan()
		skips++
	}

	// We reached EOF and no token was found. Backtrack to the current token.
	for i := 0; i < skips; i++ {
		p.scanner.Unscan()
	}

	// Revert the lookahead to original one.
	p.lookahead = nextRealToken // p.scanner.Scan()
	return &nextRealToken, false
}

// ParseProgram parses a CPL program and returns a Program AST object.
// 	program -> declarations stmt_block
func (p *Parser) ParseProgram() *Program {
	program := &Program{}

	// Parse declarations.
	program.Declarations = p.ParseDeclarations()

	// Parse statements.
	program.Statements = p.ParseStatementsBlock()

	// Make sure there's an EOF at the end of the file.
	if token, ok := p.match(lexer.EOF); !ok {
		p.addError(newParseError(token.Lexeme, []string{"EOF"}, token.Position))
	}

	return program
}

// ParseDeclarations parses a list of declarations and returns a Declaration AST array.
// 	declarations -> declaration declarations | ε
func (p *Parser) ParseDeclarations() []Declaration {
	declarations := []Declaration{}
	for p.lookahead.TokenType == lexer.ID {
		declarations = append(declarations, *p.ParseDeclaration())
	}

	return declarations
}

// ParseDeclaration parses a declaration and returns a Declaration AST object.
// 	declaration -> idlist ':' type ';'
func (p *Parser) ParseDeclaration() *Declaration {
	declaration := &Declaration{}
	declaration.Names = p.ParseIDList()

	if token, ok := p.match(lexer.COLON); !ok {
		p.addError(newParseError(token.Lexeme, []string{":"}, token.Position))
	}

	declaration.Type = p.ParseType()

	if token, ok := p.match(lexer.SEMICOLON); !ok {
		p.addError(newParseError(token.Lexeme, []string{";"}, token.Position))
	}

	return declaration
}

// ParseType parses a type returns it as a DataType.
// 	type -> INT | FLOAT
func (p *Parser) ParseType() DataType {
	token, ok := p.match(lexer.INT, lexer.FLOAT)
	if !ok {
		p.addError(newParseError(token.Lexeme, []string{"int", "float"}, token.Position))
		return Unknown
	}

	switch token.TokenType {
	case lexer.INT:
		return Integer
	case lexer.FLOAT:
		return Float
	}

	return Unknown
}

// ParseIDList parses a list of IDs and returns a string array.
// 	idlist -> ID idlist'
// 	idlist' -> ',' ID idlist' | ε
func (p *Parser) ParseIDList() []string {
	names := []string{}

	// Parse the first name
	if token, ok := p.match(lexer.ID); ok {
		names = append(names, token.Lexeme)
	} else {
		p.addError(newParseError(token.Lexeme, []string{"ID"}, token.Position))
	}

	// Parse other names if exist
	for p.lookahead.TokenType == lexer.COMMA {
		p.match(lexer.COMMA)

		if token, ok := p.match(lexer.ID); ok {
			names = append(names, token.Lexeme)
		} else {
			p.addError(newParseError(token.Lexeme, []string{"ID"}, token.Position))
		}
	}

	return names
}

// ParseStatement parses a CPL statement.
//	stmt -> assignment_stmt | input_stmt | output_stmt | if_stmt | while_stmt
//		| switch_stmt | break_stmt | stmt_block
func (p *Parser) ParseStatement() Statement {
	switch p.lookahead.TokenType {
	case lexer.ID:
		return p.ParseAssignmentStatement()
	}

	return nil
}

// ParseAssignmentStatement parses a CPL assignment statement.
// 	assignment_stmt -> ID '=' assignment_stmt'
// 	assignment_stmt' -> expression ';'
//   	| STATIC_CAST '(' type ')' '(' expression ')' ';
func (p *Parser) ParseAssignmentStatement() *AssignmentStatement {
	result := &AssignmentStatement{}

	if token, ok := p.match(lexer.ID); ok {
		result.Variable = token.Lexeme
	} else {
		p.addError(newParseError(token.Lexeme, []string{"ID"}, token.Position))
	}

	if token, ok := p.match(lexer.EQUALS); !ok {
		p.addError(newParseError(token.Lexeme, []string{"ID"}, token.Position))
	}

	// Parse static_cast(type) if exists
	if p.lookahead.TokenType == lexer.STATICCAST {
		p.match(lexer.STATICCAST)

		if token, ok := p.match(lexer.LPAREN); !ok {
			p.addError(newParseError(token.Lexeme, []string{"("}, token.Position))
		}

		result.CastType = p.ParseType()

		if token, ok := p.match(lexer.RPAREN); !ok {
			p.addError(newParseError(token.Lexeme, []string{")"}, token.Position))
		}
	}

	// Parse expression
	result.Value = p.ParseExpression()

	if token, ok := p.match(lexer.SEMICOLON); !ok {
		p.addError(newParseError(token.Lexeme, []string{";"}, token.Position))
	}

	return result
}

// ParseStatementsBlock parses a block of statements.
//	stmt_block -> '{' stmtlist '}'
func (p *Parser) ParseStatementsBlock() []Statement {
	// Parse {
	startBlock := false
	startBlockToken, startBlock := p.match(lexer.LBRACKET)
	if !startBlock {
		p.addError(newParseError(startBlockToken.Lexeme,
			[]string{"{"}, startBlockToken.Position))
	}

	statements := p.ParseStatements()

	// Parse }
	// Only show an error for the } if there was a {
	if token, ok := p.match(lexer.RBRACKET); !ok && startBlock {
		p.addError(newParseError(token.Lexeme, []string{"}"}, token.Position))
	}

	return statements
}

// ParseStatements parses zero or more statements.
//	stmtlist -> stmt stmtlist | ε
func (p *Parser) ParseStatements() []Statement {
	statements := []Statement{}
	for {
		statement := p.ParseStatement()
		if statement == nil {
			break
		}

		statements = append(statements, statement)
	}

	return statements
}

// ParseExpression parses expressions that might contain additions or subtractions.
// 	expression -> term expression'
//  expression' -> ADDOP term expression' | ε
func (p *Parser) ParseExpression() Expression {
	result := p.ParseTerm()
	for p.lookahead.TokenType == lexer.ADDOP {
		var operator Operator
		switch token, _ := p.match(lexer.ADDOP); token.Lexeme {
		case "+":
			operator = Add
		case "-":
			operator = Subtract
		}

		rhs := p.ParseTerm()
		result = &ArithmeticExpression{LHS: result, RHS: rhs, Operator: operator}
	}

	return result
}

// ParseTerm parses expressions that might contain multipications or divisions.
// 	term -> factor term'
// 	term' -> MULOP factor term' | ε
func (p *Parser) ParseTerm() Expression {
	result := p.ParseFactor()
	for p.lookahead.TokenType == lexer.MULOP {
		var operator Operator
		switch token, _ := p.match(lexer.MULOP); token.Lexeme {
		case "*":
			operator = Multiply
		case "/":
			operator = Divide
		}

		result = &ArithmeticExpression{LHS: result, RHS: p.ParseFactor(), Operator: operator}
	}

	return result
}

// ParseFactor parses a single variable, single constant number or (...some expr...).
// 	factor -> '(' expression ')' | ID | NUM
func (p *Parser) ParseFactor() Expression {
	switch p.lookahead.TokenType {
	case lexer.LPAREN:
		p.match(lexer.LPAREN)

		expr := p.ParseExpression()

		if token, ok := p.match(lexer.RPAREN); !ok {
			p.addError(newParseError(token.Lexeme, []string{")"}, token.Position))
		}

		return expr

	case lexer.ID:
		token, _ := p.match(lexer.ID)
		return &VariableExpression{Variable: token.Lexeme}

	case lexer.NUM:
		token, _ := p.match(lexer.NUM)

		value, err := strconv.ParseFloat(token.Lexeme, 64)
		if err != nil {
			p.addError(ParseError{Message: fmt.Sprintf("%s is not number", token.Lexeme)})
		}

		return &NumberLiteral{Value: value}

	default:
		p.addError(newParseError(p.lookahead.Lexeme, []string{"(", "ID", "NUM"},
			p.lookahead.Position))
		return nil
	}
}

func (p *Parser) addError(e ParseError) {
	for _, err := range p.Errors {
		if err.Pos == e.Pos {
			return
		}
	}

	p.Errors = append(p.Errors, e)
}
