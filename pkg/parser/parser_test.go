package parser_test

import (
	"strings"
	"testing"

	"github.com/alongubkin/cpl-compiler/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeclarationOneID(t *testing.T) {
	program, err := parser.Parse("var1 : int;")
	require.NoError(t, err)
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1"}, Type: parser.Integer},
		},
	}, program)
}

func TestDeclarationMultipeIDs(t *testing.T) {
	program, err := parser.Parse("var1, var2, var3 : float;")
	require.NoError(t, err)
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1", "var2", "var3"}, Type: parser.Float},
		},
	}, program)
}

func TestDeclarationInvalidType(t *testing.T) {
	_, err := parser.Parse("var1, var2, var3 : uu;")
	assert.Error(t, err)
}

func TestMultipleDeclarations(t *testing.T) {
	program, err := parser.Parse("var1, var2 : int; var3 : float; var4,var5:int;")
	require.NoError(t, err)
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1", "var2"}, Type: parser.Integer},
			parser.Declaration{Names: []string{"var3"}, Type: parser.Float},
			parser.Declaration{Names: []string{"var4", "var5"}, Type: parser.Integer},
		},
	}, program)
}

func TestAddTwoLiteralsExpression(t *testing.T) {
	expr, err := parser.NewParser(strings.NewReader("1 + 3")).ParseExpression()
	require.NoError(t, err)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS:      &parser.NumberLiteral{Value: 1},
		Operator: parser.Add,
		RHS:      &parser.NumberLiteral{Value: 3},
	}, expr)
}

func TestAddMultipleLiteralsExpression(t *testing.T) {
	expr, err := parser.NewParser(strings.NewReader("1 + 3 + 7 + 10")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("1 - 3")).ParseExpression()
	require.NoError(t, err)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS:      &parser.NumberLiteral{Value: 1},
		Operator: parser.Subtract,
		RHS:      &parser.NumberLiteral{Value: 3},
	}, expr)
}

func TestSubtractMultipleLiteralsExpression(t *testing.T) {
	expr, err := parser.NewParser(strings.NewReader("1 - 3 - 7 - 10")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("1 + 3 - 5")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("1 * 3")).ParseExpression()
	require.NoError(t, err)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS:      &parser.NumberLiteral{Value: 1},
		Operator: parser.Multiply,
		RHS:      &parser.NumberLiteral{Value: 3},
	}, expr)
}

func TestMultiplyMultipleLiteralsExpression(t *testing.T) {
	expr, err := parser.NewParser(strings.NewReader("1 * 3 * 7 * 10")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("1 / 3")).ParseExpression()
	require.NoError(t, err)
	assert.EqualValues(t, &parser.ArithmeticExpression{
		LHS:      &parser.NumberLiteral{Value: 1},
		Operator: parser.Divide,
		RHS:      &parser.NumberLiteral{Value: 3},
	}, expr)
}

func TestDivideMultipleLiteralsExpression(t *testing.T) {
	expr, err := parser.NewParser(strings.NewReader("1 / 3 / 7 / 10")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("1 * 3 / 5")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("1 + 3 * 5")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("1 * 3 + 5")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("(1 + 3) * 5")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("(1 + (3 + 7)) * 5")).ParseExpression()
	require.NoError(t, err)
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
	expr, err := parser.NewParser(strings.NewReader("(x + (y + 7)) / c")).ParseExpression()
	require.NoError(t, err)
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
