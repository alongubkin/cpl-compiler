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
	assert.EqualValues(t, "IADD _t1 5 x", buf.String())
	assert.EqualValues(t, exp, &codegen.Expression{Code: "_t1", Type: parser.Integer})
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
	assert.EqualValues(t, `IADD _t1 10 y
IADD _t2 16 _t1
IADD _t3 _t2 x`, buf.String())
	assert.EqualValues(t, exp, &codegen.Expression{Code: "_t3", Type: parser.Integer})
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
	assert.EqualValues(t, `IMLT _t1 10 y
ITOR _t3 _t1
RSUB _t2 16.500000 _t3
RDIV _t4 _t2 x`, buf.String())
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
	assert.EqualValues(t, `ITOR _t2 10
RMLT _t1 _t2 y
ITOR _t4 16
RSUB _t3 _t4 _t1
RDIV _t5 _t3 x`, buf.String())
}

func TestSimpleAssignment(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer

	c.CodegenStatement(&parser.AssignmentStatement{
		Variable: "x",
		Value:    &parser.IntLiteral{Value: 5},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IASN x 5`, buf.String())
}

func TestFloatToIntAssignment(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer

	c.CodegenStatement(&parser.AssignmentStatement{
		Variable: "x",
		Value:    &parser.FloatLiteral{Value: 5},
	})

	assert.EqualValues(t, []codegen.Error{codegen.Error{
		Message: "Cannot assign float value to int variable x."}}, c.Errors)
	assert.EqualValues(t, ``, buf.String())
}

func TestIntToFloat(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float

	c.CodegenStatement(&parser.AssignmentStatement{
		Variable: "x",
		Value:    &parser.IntLiteral{Value: 5},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `ITOR _t1 5
RASN x _t1`, buf.String())
}

func TestFloatToIntAssignmentWithCast(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer

	c.CodegenStatement(&parser.AssignmentStatement{
		Variable: "x",
		Value:    &parser.FloatLiteral{Value: 5},
		CastType: parser.Integer,
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `RTOI _t1 5.000000
IASN x _t1`, buf.String())
}

func TestFloatByCastToIntAssignment(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer

	c.CodegenStatement(&parser.AssignmentStatement{
		Variable: "x",
		Value:    &parser.IntLiteral{Value: 5},
		CastType: parser.Float,
	})

	assert.EqualValues(t, []codegen.Error{codegen.Error{
		Message: "Cannot assign float value to int variable x."}}, c.Errors)
}

func TestCompareIntegersEquality(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Integer

	c.CodegenBooleanExpression(&parser.CompareBooleanExpression{
		LHS:      &parser.VariableExpression{Variable: "x"},
		Operator: parser.EqualTo,
		RHS:      &parser.VariableExpression{Variable: "y"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IEQL _t1 x y`, buf.String())
}

func TestCompareFloatsInequality(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float
	c.Variables["y"] = parser.Float

	c.CodegenBooleanExpression(&parser.CompareBooleanExpression{
		LHS:      &parser.VariableExpression{Variable: "x"},
		Operator: parser.NotEqualTo,
		RHS:      &parser.VariableExpression{Variable: "y"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `RNQL _t1 x y`, buf.String())
}

func TestCompareIntegerLessThanFloat(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Float

	c.CodegenBooleanExpression(&parser.CompareBooleanExpression{
		LHS:      &parser.VariableExpression{Variable: "x"},
		Operator: parser.LessThan,
		RHS:      &parser.VariableExpression{Variable: "y"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `ITOR _t1 x
RLSS _t2 _t1 y`, buf.String())
}

func TestCompareFloatGreaterThanFloat(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float
	c.Variables["y"] = parser.Integer

	c.CodegenBooleanExpression(&parser.CompareBooleanExpression{
		LHS:      &parser.VariableExpression{Variable: "x"},
		Operator: parser.GreaterThan,
		RHS:      &parser.VariableExpression{Variable: "y"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `ITOR _t1 y
RGRT _t2 x _t1`, buf.String())
}

func TestOrExpression(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Integer

	c.CodegenBooleanExpression(&parser.OrBooleanExpression{
		LHS: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "x"},
			Operator: parser.GreaterThan,
			RHS:      &parser.VariableExpression{Variable: "y"},
		},
		RHS: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "y"},
			Operator: parser.EqualTo,
			RHS:      &parser.VariableExpression{Variable: "x"},
		},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IGRT _t1 x y
IEQL _t2 y x
IADD _t3 _t1 _t2
IGRT _t3 _t3 0`, buf.String())
}

func TestAndExpression(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Integer

	c.CodegenBooleanExpression(&parser.AndBooleanExpression{
		LHS: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "x"},
			Operator: parser.GreaterThan,
			RHS:      &parser.VariableExpression{Variable: "y"},
		},
		RHS: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "y"},
			Operator: parser.EqualTo,
			RHS:      &parser.VariableExpression{Variable: "x"},
		},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IGRT _t1 x y
IEQL _t2 y x
IMLT _t3 _t1 _t2`, buf.String())
}

func TestOrAndExpression(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Integer

	c.CodegenBooleanExpression(&parser.OrBooleanExpression{
		LHS: &parser.AndBooleanExpression{
			LHS: &parser.CompareBooleanExpression{
				LHS:      &parser.VariableExpression{Variable: "x"},
				Operator: parser.GreaterThan,
				RHS:      &parser.VariableExpression{Variable: "y"},
			},
			RHS: &parser.CompareBooleanExpression{
				LHS:      &parser.VariableExpression{Variable: "y"},
				Operator: parser.EqualTo,
				RHS:      &parser.VariableExpression{Variable: "x"},
			},
		},
		RHS: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "y"},
			Operator: parser.NotEqualTo,
			RHS:      &parser.VariableExpression{Variable: "x"},
		}})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IGRT _t1 x y
IEQL _t2 y x
IMLT _t3 _t1 _t2
INQL _t4 y x
IADD _t5 _t3 _t4
IGRT _t5 _t5 0`, buf.String())
}

func TestAndFloatExpression(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Float

	c.CodegenBooleanExpression(&parser.AndBooleanExpression{
		LHS: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "x"},
			Operator: parser.GreaterThan,
			RHS:      &parser.VariableExpression{Variable: "y"},
		},
		RHS: &parser.CompareBooleanExpression{
			LHS:      &parser.VariableExpression{Variable: "y"},
			Operator: parser.EqualTo,
			RHS:      &parser.VariableExpression{Variable: "x"},
		},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `ITOR _t1 x
RGRT _t2 _t1 y
ITOR _t3 x
REQL _t4 y _t3
IMLT _t5 _t2 _t4`, buf.String())
}

func TestNotAndFloatExpression(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Float

	c.CodegenBooleanExpression(&parser.NotBooleanExpression{
		Value: &parser.AndBooleanExpression{
			LHS: &parser.CompareBooleanExpression{
				LHS:      &parser.VariableExpression{Variable: "x"},
				Operator: parser.GreaterThan,
				RHS:      &parser.VariableExpression{Variable: "y"},
			},
			RHS: &parser.CompareBooleanExpression{
				LHS:      &parser.VariableExpression{Variable: "y"},
				Operator: parser.EqualTo,
				RHS:      &parser.VariableExpression{Variable: "x"},
			},
		}})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `ITOR _t1 x
RGRT _t2 _t1 y
ITOR _t3 x
REQL _t4 y _t3
IMLT _t5 _t2 _t4
ISUB _t6 1 _t5`, buf.String())
}

func TestCompareGreaterThanOrEqualTo(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Integer

	c.CodegenBooleanExpression(&parser.CompareBooleanExpression{
		LHS:      &parser.VariableExpression{Variable: "x"},
		Operator: parser.GreaterThanOrEqualTo,
		RHS:      &parser.VariableExpression{Variable: "y"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IEQL _t1 x y
IGRT _t2 x y
IADD _t3 _t1 _t2
IGRT _t3 _t3 0`, buf.String())
}

func TestCompareLessThanOrEqualTo(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Integer

	c.CodegenBooleanExpression(&parser.CompareBooleanExpression{
		LHS:      &parser.VariableExpression{Variable: "x"},
		Operator: parser.LessThenOrEqualTo,
		RHS:      &parser.VariableExpression{Variable: "y"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IEQL _t1 x y
ILSS _t2 x y
IADD _t3 _t1 _t2
IGRT _t3 _t3 0`, buf.String())
}

func TestInputInteger(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer

	c.CodegenStatement(&parser.InputStatement{
		Variable: "x",
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IINP x`, buf.String())
}

func TestInputFloat(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float

	c.CodegenStatement(&parser.InputStatement{
		Variable: "x",
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `RINP x`, buf.String())
}

func TestInputVariableNotExists(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.CodegenStatement(&parser.InputStatement{
		Variable: "x",
	})

	assert.EqualValues(t, []codegen.Error{codegen.Error{Message: "Undefined variable x."}}, c.Errors)
}

func TestOutputInteger(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.CodegenStatement(&parser.OutputStatement{
		Value: &parser.IntLiteral{Value: 5},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IPRT 5`, buf.String())
}

func TestOutputFloat(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.CodegenStatement(&parser.OutputStatement{
		Value: &parser.FloatLiteral{Value: 5},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `RPRT 5.000000`, buf.String())
}

func TestIfElse(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float
	c.Variables["y"] = parser.Float

	c.CodegenStatement(&parser.IfStatement{
		Condition: &parser.CompareBooleanExpression{
			LHS:      &parser.IntLiteral{Value: 0},
			Operator: parser.EqualTo,
			RHS:      &parser.IntLiteral{Value: 1},
		},
		IfBranch:   &parser.InputStatement{Variable: "x"},
		ElseBranch: &parser.InputStatement{Variable: "y"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IEQL _t1 0 1
JMPZ @2 _t1
RINP x
JUMP @1
@2:
RINP y
@1:`, buf.String())
}

func TestIfElseIfElse(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float
	c.Variables["y"] = parser.Float

	c.CodegenStatement(&parser.IfStatement{
		Condition: &parser.CompareBooleanExpression{
			LHS:      &parser.IntLiteral{Value: 0},
			Operator: parser.EqualTo,
			RHS:      &parser.IntLiteral{Value: 1},
		},
		IfBranch: &parser.InputStatement{Variable: "x"},
		ElseBranch: &parser.IfStatement{
			Condition: &parser.CompareBooleanExpression{
				LHS:      &parser.IntLiteral{Value: 0},
				Operator: parser.EqualTo,
				RHS:      &parser.IntLiteral{Value: 1},
			},
			IfBranch:   &parser.InputStatement{Variable: "x"},
			ElseBranch: &parser.InputStatement{Variable: "y"},
		},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IEQL _t1 0 1
JMPZ @2 _t1
RINP x
JUMP @1
@2:
IEQL _t2 0 1
JMPZ @4 _t2
RINP x
JUMP @3
@4:
RINP y
@3:
@1:`, buf.String())
}

func TestBreakStatementNoContext(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.CodegenStatement(&parser.BreakStatement{})
	assert.EqualValues(t, []codegen.Error{codegen.Error{
		Message: "Break statement must be inside a while loop or a switch case."}}, c.Errors)
}

func TestWhileLoop(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float
	c.Variables["y"] = parser.Float

	c.CodegenStatement(&parser.WhileStatement{
		Condition: &parser.CompareBooleanExpression{
			LHS:      &parser.IntLiteral{Value: 0},
			Operator: parser.EqualTo,
			RHS:      &parser.IntLiteral{Value: 1},
		},
		Body: &parser.InputStatement{Variable: "x"},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `@1:
IEQL _t1 0 1
JMPZ @2 _t1
RINP x
JUMP @1
@2:`, buf.String())
}

func TestWhileLoopWithBreak(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float
	c.Variables["y"] = parser.Float

	c.CodegenStatement(&parser.WhileStatement{
		Condition: &parser.CompareBooleanExpression{
			LHS:      &parser.IntLiteral{Value: 0},
			Operator: parser.EqualTo,
			RHS:      &parser.IntLiteral{Value: 1},
		},
		Body: &parser.StatementsBlock{Statements: []parser.Statement{
			&parser.InputStatement{Variable: "x"},
			&parser.BreakStatement{},
		}},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `@1:
IEQL _t1 0 1
JMPZ @2 _t1
RINP x
JUMP @2
JUMP @1
@2:`, buf.String())
}

func TestNestedWhileLoopWithBreak(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Float
	c.Variables["y"] = parser.Float

	c.CodegenStatement(&parser.WhileStatement{
		Condition: &parser.CompareBooleanExpression{
			LHS:      &parser.IntLiteral{Value: 0},
			Operator: parser.EqualTo,
			RHS:      &parser.IntLiteral{Value: 1},
		},
		Body: &parser.StatementsBlock{Statements: []parser.Statement{
			&parser.InputStatement{Variable: "x"},
			&parser.BreakStatement{},
			&parser.WhileStatement{
				Condition: &parser.CompareBooleanExpression{
					LHS:      &parser.IntLiteral{Value: 1},
					Operator: parser.NotEqualTo,
					RHS:      &parser.IntLiteral{Value: 2},
				},
				Body: &parser.StatementsBlock{Statements: []parser.Statement{
					&parser.InputStatement{Variable: "y"},
					&parser.BreakStatement{},
				}},
			},
			&parser.BreakStatement{},
		}},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `@1:
IEQL _t1 0 1
JMPZ @2 _t1
RINP x
JUMP @2
@3:
INQL _t2 1 2
JMPZ @4 _t2
RINP y
JUMP @4
JUMP @3
@4:
JUMP @2
JUMP @1
@2:`, buf.String())
}

func TestSwitchStatement(t *testing.T) {
	buf := new(bytes.Buffer)

	c := codegen.NewCodeGenerator(buf)
	c.Variables["x"] = parser.Integer
	c.Variables["y"] = parser.Float

	c.CodegenStatement(&parser.SwitchStatement{
		Expression: &parser.VariableExpression{Variable: "x"},
		Cases: []parser.SwitchCase{
			parser.SwitchCase{
				Value: 1,
				Statements: []parser.Statement{
					&parser.InputStatement{Variable: "x"},
					&parser.BreakStatement{},
				},
			},
			parser.SwitchCase{
				Value: 2,
				Statements: []parser.Statement{
					&parser.InputStatement{Variable: "y"},
					&parser.BreakStatement{},
				},
			},
		},
		DefaultCase: []parser.Statement{
			&parser.InputStatement{Variable: "x"},
			&parser.BreakStatement{},
		},
	})

	assert.Empty(t, c.Errors)
	assert.EqualValues(t, `IASN _t1 x
IEQL _t2 _t1 1
JMPZ @2 _t2
IINP x
JUMP @1
@2:
IEQL _t3 _t1 2
JMPZ @3 _t3
RINP y
JUMP @1
@3:
IINP x
JUMP @1
@1:`, buf.String())
}
