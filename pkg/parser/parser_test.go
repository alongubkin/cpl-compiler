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
		StatementsBlock: &parser.StatementsBlock{
			Statements: []parser.Statement{},
		},
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
		StatementsBlock: &parser.StatementsBlock{
			Statements: []parser.Statement{},
		},
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
		StatementsBlock: &parser.StatementsBlock{
			Statements: []parser.Statement{},
		},
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
		StatementsBlock: &parser.StatementsBlock{
			Statements: []parser.Statement{},
		},
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
		StatementsBlock: &parser.StatementsBlock{
			Statements: []parser.Statement{},
		},
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
		StatementsBlock: &parser.StatementsBlock{
			Statements: []parser.Statement{},
		},
	}, program)
}

func TestProgramWithAssignmentStatements(t *testing.T) {
	program, errors := parser.Parse("x , y : int; { x = 5 * (y + b); y = static_cast(float)(x + 5); }")
	assert.Empty(t, errors)
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"x", "y"}, Type: parser.Integer},
		},
		StatementsBlock: &parser.StatementsBlock{
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
		},
	}, program)
}

func TestInputStatement(t *testing.T) {
	p := parser.NewParser(strings.NewReader("input(x);"))
	statement := p.ParseStatement()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.InputStatement{Variable: "x"}, statement)
}

func TestOutputStatement(t *testing.T) {
	p := parser.NewParser(strings.NewReader("output(3 + x);"))
	statement := p.ParseStatement()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.OutputStatement{
		Value: &parser.ArithmeticExpression{
			LHS:      &parser.NumberLiteral{Value: 3},
			Operator: parser.Add,
			RHS:      &parser.VariableExpression{Variable: "x"},
		}}, statement)
}

func TestOrAndPrecedence(t *testing.T) {
	p := parser.NewParser(strings.NewReader("x <= 5 || y >= 6 && 3 == 4"))
	expr := p.ParseBooleanExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.OrBooleanExpression{
		LHS: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "x"},
			Operator: parser.LessThenOrEqualTo,
			RHS:      &parser.NumberLiteral{Value: 5},
		},
		RHS: &parser.AndBooleanExpression{
			LHS: &parser.CompareBooleanExpression{
				LHS:      &parser.VariableExpression{Variable: "y"},
				Operator: parser.GreaterThanOrEqualTo,
				RHS:      &parser.NumberLiteral{Value: 6},
			},
			RHS: &parser.CompareBooleanExpression{
				LHS:      &parser.NumberLiteral{Value: 3},
				Operator: parser.EqualTo,
				RHS:      &parser.NumberLiteral{Value: 4},
			},
		},
	}, expr)
}

func TestAndOrPrecedence(t *testing.T) {
	p := parser.NewParser(strings.NewReader("x != 5 && y < 6 || 3 == 4"))
	expr := p.ParseBooleanExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.OrBooleanExpression{
		LHS: &parser.AndBooleanExpression{
			LHS: &parser.CompareBooleanExpression{
				LHS:      &parser.VariableExpression{Variable: "x"},
				Operator: parser.NotEqualTo,
				RHS:      &parser.NumberLiteral{Value: 5},
			},
			RHS: &parser.CompareBooleanExpression{
				LHS:      &parser.VariableExpression{Variable: "y"},
				Operator: parser.LessThan,
				RHS:      &parser.NumberLiteral{Value: 6},
			},
		},
		RHS: &parser.CompareBooleanExpression{
			LHS:      &parser.NumberLiteral{Value: 3},
			Operator: parser.EqualTo,
			RHS:      &parser.NumberLiteral{Value: 4},
		},
	}, expr)
}

func TestNot(t *testing.T) {
	p := parser.NewParser(strings.NewReader("!(x > 5 && y < 6) || 3 == 4"))
	expr := p.ParseBooleanExpression()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.OrBooleanExpression{
		LHS: &parser.NotBooleanExpression{
			Value: &parser.AndBooleanExpression{
				LHS: &parser.CompareBooleanExpression{
					LHS:      &parser.VariableExpression{Variable: "x"},
					Operator: parser.GreaterThan,
					RHS:      &parser.NumberLiteral{Value: 5},
				},
				RHS: &parser.CompareBooleanExpression{
					LHS:      &parser.VariableExpression{Variable: "y"},
					Operator: parser.LessThan,
					RHS:      &parser.NumberLiteral{Value: 6},
				},
			},
		},
		RHS: &parser.CompareBooleanExpression{
			LHS:      &parser.NumberLiteral{Value: 3},
			Operator: parser.EqualTo,
			RHS:      &parser.NumberLiteral{Value: 4},
		},
	}, expr)
}

func TestIfStatement(t *testing.T) {
	p := parser.NewParser(strings.NewReader("if (x == y) input(x); else output(y);"))
	statement := p.ParseStatement()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.IfStatement{
		Condition: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "x"},
			Operator: parser.EqualTo,
			RHS:      &parser.VariableExpression{Variable: "y"},
		},
		IfBranch: &parser.InputStatement{Variable: "x"},
		ElseBranch: &parser.OutputStatement{
			Value: &parser.VariableExpression{Variable: "y"},
		},
	}, statement)
}

func TestElseIfStatement(t *testing.T) {
	p := parser.NewParser(strings.NewReader("if (x == y) { input(x); y = 7; } else if (x == 3) output(y); else t = 6;"))
	statement := p.ParseStatement()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.IfStatement{
		Condition: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "x"},
			Operator: parser.EqualTo,
			RHS:      &parser.VariableExpression{Variable: "y"},
		},
		IfBranch: &parser.StatementsBlock{
			Statements: []parser.Statement{
				&parser.InputStatement{Variable: "x"},
				&parser.AssignmentStatement{
					Variable: "y",
					Value:    &parser.NumberLiteral{Value: 7},
				},
			},
		},
		ElseBranch: &parser.IfStatement{
			Condition: &parser.CompareBooleanExpression{
				LHS:      &parser.VariableExpression{Variable: "x"},
				Operator: parser.EqualTo,
				RHS:      &parser.NumberLiteral{Value: 3},
			},
			IfBranch: &parser.OutputStatement{
				Value: &parser.VariableExpression{Variable: "y"},
			},
			ElseBranch: &parser.AssignmentStatement{
				Variable: "t",
				Value:    &parser.NumberLiteral{Value: 6},
			},
		},
	}, statement)
}

func TestWhileStatement(t *testing.T) {
	p := parser.NewParser(strings.NewReader("while (!(x == y)) { input(x); y = 7; }"))
	statement := p.ParseStatement()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.WhileStatement{
		Condition: &parser.NotBooleanExpression{
			Value: &parser.CompareBooleanExpression{
				LHS:      &parser.VariableExpression{Variable: "x"},
				Operator: parser.EqualTo,
				RHS:      &parser.VariableExpression{Variable: "y"},
			},
		},
		Body: &parser.StatementsBlock{
			Statements: []parser.Statement{
				&parser.InputStatement{Variable: "x"},
				&parser.AssignmentStatement{
					Variable: "y",
					Value:    &parser.NumberLiteral{Value: 7},
				},
			},
		},
	}, statement)
}

func TestBreakStatement(t *testing.T) {
	p := parser.NewParser(strings.NewReader("break;"))
	statement := p.ParseStatement()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.BreakStatement{}, statement)
}

func TestSwitchStatement(t *testing.T) {
	p := parser.NewParser(strings.NewReader(`
		switch (x + y) { 
		case 5: 
			output(x); 
			break;
		case 6: {
			input(y);
			break;
		}
		default: 
			x = y; 
			break;
		}
		`))

	statement := p.ParseStatement()
	assert.Empty(t, p.Errors)
	assert.EqualValues(t, &parser.SwitchStatement{
		Expression: &parser.ArithmeticExpression{
			LHS:      &parser.VariableExpression{Variable: "x"},
			Operator: parser.Add,
			RHS:      &parser.VariableExpression{Variable: "y"},
		},
		Cases: []parser.SwitchCase{
			parser.SwitchCase{
				Value: 5,
				Statements: []parser.Statement{
					&parser.OutputStatement{Value: &parser.VariableExpression{Variable: "x"}},
					&parser.BreakStatement{},
				},
			},
			parser.SwitchCase{
				Value: 6,
				Statements: []parser.Statement{
					&parser.StatementsBlock{Statements: []parser.Statement{
						&parser.InputStatement{Variable: "y"},
						&parser.BreakStatement{},
					}},
				},
			},
		},
		DefaultCase: []parser.Statement{
			&parser.AssignmentStatement{Variable: "x",
				Value: &parser.VariableExpression{Variable: "y"}},
			&parser.BreakStatement{},
		},
	}, statement)
}
