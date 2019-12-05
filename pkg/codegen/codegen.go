package codegen

import (
	"bufio"
	"fmt"
	"io"

	"github.com/alongubkin/cpl-compiler/pkg/parser"
)

type CodeGenerator struct {
	Errors         []Error
	output         *bufio.Writer
	Variables      map[string]parser.DataType
	temporaryIndex int
}

type Expression struct {
	Code string
	Type parser.DataType
}

// NewCodeGenerator returns a new instance of CodeGenerator.
func NewCodeGenerator(output io.Writer) *CodeGenerator {
	return &CodeGenerator{
		Errors:         []Error{},
		output:         bufio.NewWriterSize(output, 1),
		Variables:      map[string]parser.DataType{},
		temporaryIndex: 0,
	}
}

// CodegenProgram generates code for a CPL program.
func (c *CodeGenerator) CodegenProgram(node parser.Program) {
	// Go over variable declarations
	for _, declaration := range node.Declarations {
		for _, name := range declaration.Names {
			if _, exists := c.Variables[name]; exists {
				c.Errors = append(c.Errors, Error{Message: fmt.Sprintf("Variable %s already defined.", name)})
				continue
			}

			c.Variables[name] = declaration.Type
		}
	}

	c.CodegenStatement(node.StatementsBlock)
}

// CodegenStatement generates code for a CPL statement.
func (c *CodeGenerator) CodegenStatement(node parser.Statement) {
	switch s := node.(type) {
	case *parser.StatementsBlock:
		c.CodegenStatementsBlock(s)
	}
}

// CodegenStatementsBlock generates code for a statements block.
func (c *CodeGenerator) CodegenStatementsBlock(node *parser.StatementsBlock) {
	for _, statement := range node.Statements {
		c.CodegenStatement(statement)
	}
}

// CodegenExpression generates code for a CPL expression.
func (c *CodeGenerator) CodegenExpression(node parser.Expression) *Expression {
	switch s := node.(type) {
	case *parser.ArithmeticExpression:
		return c.CodegenArithmeticExpression(s)
	case *parser.VariableExpression:
		return c.CodegenVariableExpression(s)
	case *parser.IntLiteral:
		return c.CodegenIntLiteral(s)
	case *parser.FloatLiteral:
		return c.CodegenFloatLiteral(s)
	}

	return nil
}

// CodegenArithmeticExpression generates code for an arithmetic expression.
func (c *CodeGenerator) CodegenArithmeticExpression(node *parser.ArithmeticExpression) *Expression {
	lhs := c.CodegenExpression(node.LHS)
	rhs := c.CodegenExpression(node.RHS)
	if lhs == nil || rhs == nil {
		return nil
	}

	result := &Expression{
		Code: c.getNewTemporary(),
		Type: calculateExpressionType(lhs.Type, rhs.Type),
	}

	switch node.Operator {
	case parser.Add:
		if result.Type == parser.Integer {
			c.output.WriteString(fmt.Sprintf("IADD %s %s %s\n", result.Code, lhs.Code, rhs.Code))
		} else if result.Type == parser.Float {
			c.output.WriteString(fmt.Sprintf("RADD %s %s %s\n", result.Code, lhs.Code, rhs.Code))
		}

	case parser.Subtract:
		if result.Type == parser.Integer {
			c.output.WriteString(fmt.Sprintf("ISUB %s %s %s\n", result.Code, lhs.Code, rhs.Code))
		} else if result.Type == parser.Float {
			c.output.WriteString(fmt.Sprintf("RSUB %s %s %s\n", result.Code, lhs.Code, rhs.Code))
		}

	case parser.Multiply:
		if result.Type == parser.Integer {
			c.output.WriteString(fmt.Sprintf("IMLT %s %s %s\n", result.Code, lhs.Code, rhs.Code))
		} else if result.Type == parser.Float {
			c.output.WriteString(fmt.Sprintf("RMLT %s %s %s\n", result.Code, lhs.Code, rhs.Code))
		}

	case parser.Divide:
		if result.Type == parser.Integer {
			c.output.WriteString(fmt.Sprintf("IDIV %s %s %s\n", result.Code, lhs.Code, rhs.Code))
		} else if result.Type == parser.Float {
			c.output.WriteString(fmt.Sprintf("RDIV %s %s %s\n", result.Code, lhs.Code, rhs.Code))
		}
	}

	return result
}

// CodegenVariableExpression generates code for a variable expression.
func (c *CodeGenerator) CodegenVariableExpression(node *parser.VariableExpression) *Expression {
	// Make sure the variable is defined.
	if _, exists := c.Variables[node.Variable]; !exists {
		c.Errors = append(c.Errors, Error{Message: fmt.Sprintf(
			"Undefined variable %s.", node.Variable)})
		return nil
	}

	return &Expression{Code: node.Variable, Type: c.Variables[node.Variable]}
}

// CodegenIntLiteral generates code for an integer literal.
func (c *CodeGenerator) CodegenIntLiteral(node *parser.IntLiteral) *Expression {
	return &Expression{
		Code: fmt.Sprintf("%d", node.Value),
		Type: parser.Integer,
	}
}

// CodegenFloatLiteral generates code for an float literal.
func (c *CodeGenerator) CodegenFloatLiteral(node *parser.FloatLiteral) *Expression {
	return &Expression{
		Code: fmt.Sprintf("%f", node.Value),
		Type: parser.Float,
	}
}

func (c *CodeGenerator) getNewTemporary() string {
	c.temporaryIndex++
	return fmt.Sprintf("$t%d", c.temporaryIndex)
}

func calculateExpressionType(types ...parser.DataType) parser.DataType {
	for _, t := range types {
		if t == parser.Float {
			return parser.Float
		}
	}

	return parser.Integer
}
