package parser_test

import (
	"strings"
	"testing"

	"github.com/alongubkin/cpl-compiler/pkg/lexer"
	"github.com/alongubkin/cpl-compiler/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestEmptyProgram(t *testing.T) {
	program, errors := parser.Parse("{}")
	assert.Empty(t, errors)
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{},
		Statements:   []parser.Statement{},
	}, program)
}

func TestDeclarationOneID(t *testing.T) {
	p := parser.NewParser(strings.NewReader("var1 : int;"))
	declarations := p.ParseDeclarations()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, []parser.Declaration{
		parser.Declaration{Names: []string{"var1"}, Type: parser.Integer},
	}, declarations)
}

func TestDeclarationMultipeIDs(t *testing.T) {
	p := parser.NewParser(strings.NewReader("var1, var2, var3 : float;"))
	declarations := p.ParseDeclarations()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, []parser.Declaration{
		parser.Declaration{Names: []string{"var1", "var2", "var3"}, Type: parser.Float},
	}, declarations)
}

func TestDeclarationInvalidType(t *testing.T) {
	p := parser.NewParser(strings.NewReader("var1, var2, var3 : uu;"))
	declarations := p.ParseDeclarations()
	assert.Len(t, p.Errors, 1)
	assert.EqualValues(t, []parser.Declaration{
		parser.Declaration{Names: []string{"var1", "var2", "var3"}, Type: parser.Unknown},
	}, declarations)
}

func TestMultipleDeclarations(t *testing.T) {
	p := parser.NewParser(strings.NewReader("var1, var2 : int; var3 : float; var4,var5:int;"))
	declarations := p.ParseDeclarations()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, []parser.Declaration{
		parser.Declaration{Names: []string{"var1", "var2"}, Type: parser.Integer},
		parser.Declaration{Names: []string{"var3"}, Type: parser.Float},
		parser.Declaration{Names: []string{"var4", "var5"}, Type: parser.Integer},
	}, declarations)
}

func TestAddTwoLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 + 3"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS:      &parser.NumberLiteral{Value: 1},
		Operator: parser.Add,
		RHS:      &parser.NumberLiteral{Value: 3},
	}, expr)
}

func TestAddMultipleLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 + 3 + 7 + 10"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS: &parser.ArithmeticExpression{
				LHS:      &parser.NumberLiteral{Value: 1},
				Operator: parser.Add,
				RHS:      &parser.NumberLiteral{Value: 3},
			},
			Operator: parser.Add,
			RHS:      &parser.NumberLiteral{Value: 7},
		},
		Operator: parser.Add,
		RHS:      &parser.NumberLiteral{Value: 10},
	}, expr)
}

func TestSubtractTwoLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 - 3"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS:      &parser.NumberLiteral{Value: 1},
		Operator: parser.Subtract,
		RHS:      &parser.NumberLiteral{Value: 3},
	}, expr)
}

func TestSubtractMultipleLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 - 3 - 7 - 10"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS: &parser.ArithmeticExpression{
				LHS:      &parser.NumberLiteral{Value: 1},
				Operator: parser.Subtract,
				RHS:      &parser.NumberLiteral{Value: 3},
			},
			Operator: parser.Subtract,
			RHS:      &parser.NumberLiteral{Value: 7},
		},
		Operator: parser.Subtract,
		RHS:      &parser.NumberLiteral{Value: 10},
	}, expr)
}

func TestAddAndSubtractLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 + 3 - 5"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS:      &parser.NumberLiteral{Value: 1},
			Operator: parser.Add,
			RHS:      &parser.NumberLiteral{Value: 3},
		},
		Operator: parser.Subtract,
		RHS:      &parser.NumberLiteral{Value: 5},
	}, expr)
}

func TestMultiplyTwoLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 * 3"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS:      &parser.NumberLiteral{Value: 1},
		Operator: parser.Multiply,
		RHS:      &parser.NumberLiteral{Value: 3},
	}, expr)
}

func TestMultiplyMultipleLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 * 3 * 7 * 10"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS: &parser.ArithmeticExpression{
				LHS:      &parser.NumberLiteral{Value: 1},
				Operator: parser.Multiply,
				RHS:      &parser.NumberLiteral{Value: 3},
			},
			Operator: parser.Multiply,
			RHS:      &parser.NumberLiteral{Value: 7},
		},
		Operator: parser.Multiply,
		RHS:      &parser.NumberLiteral{Value: 10},
	}, expr)
}

func TestDivideTwoLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 / 3"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS:      &parser.NumberLiteral{Value: 1},
		Operator: parser.Divide,
		RHS:      &parser.NumberLiteral{Value: 3},
	}, expr)
}

func TestDivideMultipleLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 / 3 / 7 / 10"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS: &parser.ArithmeticExpression{
				LHS:      &parser.NumberLiteral{Value: 1},
				Operator: parser.Divide,
				RHS:      &parser.NumberLiteral{Value: 3},
			},
			Operator: parser.Divide,
			RHS:      &parser.NumberLiteral{Value: 7},
		},
		Operator: parser.Divide,
		RHS:      &parser.NumberLiteral{Value: 10},
	}, expr)
}

func TestMultiplyAndDivideLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 * 3 / 5"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS:      &parser.NumberLiteral{Value: 1},
			Operator: parser.Multiply,
			RHS:      &parser.NumberLiteral{Value: 3},
		},
		Operator: parser.Divide,
		RHS:      &parser.NumberLiteral{Value: 5},
	}, expr)
}

func TestMultiplyAndAddLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 + 3 * 5"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS:      &parser.NumberLiteral{Value: 1},
		Operator: parser.Add,
		RHS: &parser.ArithmeticExpression{
			LHS:      &parser.NumberLiteral{Value: 3},
			Operator: parser.Multiply,
			RHS:      &parser.NumberLiteral{Value: 5},
		},
	}, expr)
}

func TestAddAndMultiplyLiteralsExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("1 * 3 + 5"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS:      &parser.NumberLiteral{Value: 1},
			Operator: parser.Multiply,
			RHS:      &parser.NumberLiteral{Value: 3},
		},
		Operator: parser.Add,
		RHS:      &parser.NumberLiteral{Value: 5},
	}, expr)
}

func TestParenthesisExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("(1 + 3) * 5"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS:      &parser.NumberLiteral{Value: 1},
			Operator: parser.Add,
			RHS:      &parser.NumberLiteral{Value: 3},
		},
		Operator: parser.Multiply,
		RHS:      &parser.NumberLiteral{Value: 5},
	}, expr)
}

func TestMultipleParenthesisExpression(t *testing.T) {
	p := parser.NewParser(strings.NewReader("(1 + (3 + 7)) * 5"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS:      &parser.NumberLiteral{Value: 1},
			Operator: parser.Add,
			RHS: &parser.ArithmeticExpression{
				LHS:      &parser.NumberLiteral{Value: 3},
				Operator: parser.Add,
				RHS:      &parser.NumberLiteral{Value: 7},
			},
		},
		Operator: parser.Multiply,
		RHS:      &parser.NumberLiteral{Value: 5},
	}, expr)
}

func TestExpressionWithVariables(t *testing.T) {
	p := parser.NewParser(strings.NewReader("(x + (y + 7)) / c"))
	expr := p.ParseExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS:      &parser.VariableExpression{Variable: "x"},
			Operator: parser.Add,
			RHS: &parser.ArithmeticExpression{
				LHS:      &parser.VariableExpression{Variable: "y"},
				Operator: parser.Add,
				RHS:      &parser.NumberLiteral{Value: 7},
			},
		},
		Operator: parser.Divide,
		RHS:      &parser.VariableExpression{Variable: "c"},
	}, expr)
}

func TestErrorRecoveryOneToken(t *testing.T) {
	program, errors := parser.Parse("var1 : uu int; {}")
	assert.EqualValues(t, errors, []parser.ParseError{
		parser.ParseError{Expected: []string{}, Found: "uu",
			Pos: lexer.Position{Line: 0, Column: 7}},
	})
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1"}, Type: parser.Integer},
		},
		Statements: []parser.Statement{},
	}, program)
}

func TestErrorRecoveryMultipleTokens(t *testing.T) {
	program, errors := parser.Parse("var1 : kk  * * / break hello int;")
	assert.EqualValues(t, errors, []parser.ParseError{
		parser.ParseError{Expected: []string{}, Found: "kk", Pos: lexer.Position{Line: 0, Column: 7}},
		parser.ParseError{Expected: []string{"{"}, Found: "EOF", Pos: lexer.Position{Line: 0, Column: 33}},
	})
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1"}, Type: parser.Integer},
		},
		Statements: []parser.Statement{},
	}, program)
}

func TestErrorRecoveryMultipleTokensAndDeclarations(t *testing.T) {
	program, errors := parser.Parse("var1 : kk  * * / break hello int; var2: x float;")
	assert.EqualValues(t, errors, []parser.ParseError{
		parser.ParseError{Expected: []string{}, Found: "kk", Pos: lexer.Position{Line: 0, Column: 7}},
		parser.ParseError{Expected: []string{}, Found: "x", Pos: lexer.Position{Line: 0, Column: 40}},
		parser.ParseError{Expected: []string{"{"}, Found: "EOF", Pos: lexer.Position{Line: 0, Column: 48}},
	})
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1"}, Type: parser.Integer},
			parser.Declaration{Names: []string{"var2"}, Type: parser.Float},
		},
		Statements: []parser.Statement{},
	}, program)
}

func TestErrorRecoveryEOF(t *testing.T) {
	program, errors := parser.Parse("var1 : int; {")
	assert.EqualValues(t, errors, []parser.ParseError{
		parser.ParseError{Expected: []string{"}"}, Found: "EOF", Pos: lexer.Position{Line: 0, Column: 13}},
	})
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1"}, Type: parser.Integer},
		},
		Statements: []parser.Statement{},
	}, program)
}

func TestErrorRecoveryTwiceWithEOF(t *testing.T) {
	program, errors := parser.Parse("var1 : int {")
	assert.EqualValues(t, errors, []parser.ParseError{
		parser.ParseError{Expected: []string{";"}, Found: "{", Pos: lexer.Position{Line: 0, Column: 11}},
		parser.ParseError{Expected: []string{"}"}, Found: "EOF", Pos: lexer.Position{Line: 0, Column: 12}},
	})
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1"}, Type: parser.Integer},
		},
		Statements: []parser.Statement{},
	}, program)
}

func TestProgramWithAssignmentStatements(t *testing.T) {
	program, errors := parser.Parse("x , y : int; { x = 5 * (y + b); y = static_cast(float)(x + 5); }")
	assert.Empty(t, errors)
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"x", "y"}, Type: parser.Integer},
		},
		Statements: []parser.Statement{
			&parser.AssignmentStatement{Variable: "x", Value: &parser.ArithmeticExpression{
				LHS:      &parser.NumberLiteral{Value: 5},
				Operator: parser.Multiply,
				RHS: &parser.ArithmeticExpression{
					LHS:      &parser.VariableExpression{Variable: "y"},
					Operator: parser.Add,
					RHS:      &parser.VariableExpression{Variable: "b"},
				},
			}},
			&parser.AssignmentStatement{Variable: "y", Value: &parser.ArithmeticExpression{
				LHS:      &parser.VariableExpression{Variable: "x"},
				Operator: parser.Add,
				RHS:      &parser.NumberLiteral{Value: 5},
			}, CastType: parser.Float},
		},
	}, program)
}
