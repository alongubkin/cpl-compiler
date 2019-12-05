package codegen_test

import (
	"bytes"
	"testing"

	"github.com/alongubkin/cpl-compiler/pkg/codegen"
	"github.com/alongubkin/cpl-compiler/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestCodegenAddExpression(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer

	exp := c.CodegenExpression(&parser.ArithmeticExpression{
		LHS:      &parser.IntLiteral{Value: 5},
		Operator: parser.Add,
		RHS:      &parser.VariableExpression{Variable: "x"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, "IADD $t1 5 x", buf.String())
	assert.EqualValues(t, exp, &codegen.Expression{Code: "$t1", Type: parser.Integer})
}

func TestCodegenAddExpressionVariableNotExists(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.CodegenExpression(&parser.ArithmeticExpression{
		LHS:      &parser.IntLiteral{Value: 5},
		Operator: parser.Add,
		RHS:      &parser.VariableExpression{Variable: "x"},
	})

	assert.EqualValues(t, []codegen.Error{codegen.Error{Message: "Undefined variable x."}}, c.Errors)
	assert.EqualValues(t, "", buf.String())
}

func TestCodegenComplexAddExpression(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Integer

	exp := c.CodegenExpression(&parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS:      &parser.IntLiteral{Value: 16},
			Operator: parser.Add,
			RHS: &parser.ArithmeticExpression{
				LHS:      &parser.IntLiteral{Value: 10},
				Operator: parser.Add,
				RHS:      &parser.VariableExpression{Variable: "y"},
			},
		},
		Operator: parser.Add,
		RHS:      &parser.VariableExpression{Variable: "x"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IADD $t1 10 y
IADD $t2 16 $t1
IADD $t3 $t2 x`, buf.String())
	assert.EqualValues(t, exp, &codegen.Expression{Code: "$t3", Type: parser.Integer})
}

func TestCodegenComplexExpression(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float
	c.Variables["y"] = parser.Integer

	c.CodegenExpression(&parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS:      &parser.FloatLiteral{Value: 16.5},
			Operator: parser.Subtract,
			RHS: &parser.ArithmeticExpression{
				LHS:      &parser.IntLiteral{Value: 10},
				Operator: parser.Multiply,
				RHS:      &parser.VariableExpression{Variable: "y"},
			},
		},
		Operator: parser.Divide,
		RHS:      &parser.VariableExpression{Variable: "x"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IMLT $t1 10 y
RSUB $t2 16.500000 $t1
RDIV $t3 $t2 x`, buf.String())
}

func TestVariableType(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float
	c.Variables["y"] = parser.Float

	c.CodegenExpression(&parser.ArithmeticExpression{
		LHS: &parser.ArithmeticExpression{
			LHS:      &parser.IntLiteral{Value: 16},
			Operator: parser.Subtract,
			RHS: &parser.ArithmeticExpression{
				LHS:      &parser.IntLiteral{Value: 10},
				Operator: parser.Multiply,
				RHS:      &parser.VariableExpression{Variable: "y"},
			},
		},
		Operator: parser.Divide,
		RHS:      &parser.VariableExpression{Variable: "x"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `RMLT $t1 10 y
RSUB $t2 16 $t1
RDIV $t3 $t2 x`, buf.String())
}
