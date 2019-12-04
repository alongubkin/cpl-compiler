package parser

// DataType represents the primitive data types available in CPL.
type DataType int

const (
	// Unknown primitive data type.
	Unknown DataType = iota
	// Float means the data type is a float.
	Float DataType = 1
	// Integer means the data type is an integer.
	Integer DataType = 2
)

// Operator represents a boolean or arithmatic operator in CPL.
type Operator int

const (
	// Add (+) two or more numbers
	Add Operator = iota
	// Subtract (-) two or more numbers
	Subtract
	// Multiply (*) two or more numbers
	Multiply
	// Divide (/) two or more numbers
	Divide
)

// Node represents a node in the CPL abstract syntax tree.
type Node interface {
	// node is unexported to ensure implementations of Node
	// can only originate in this package.
	node()
}

// Program represents the root node of a CPL program.
type Program struct {
	Declarations []Declaration
	Statements   []Statement
}

// Declaration of one or more variables.
type Declaration struct {
	Names []string
	Type  DataType
}

// Statement represents a single command in CPL.
type Statement interface {
	Node
	// statement is unexported to ensure implementations of Statement
	// can only originate in this package.
	statement()
}

// AssignmentStatement represents a command for assigning a value to a variable,
// e.g: x = 5;
type AssignmentStatement struct {
	Variable string
	Value    Expression

	// If the assignment doesn't contain static_cast<>, then CastType will be Unknown.
	// Otherwise, CastType will contain the type to cast to.
	CastType DataType
}

// Expression is a combination of numbers, variables and operators that
// can be evaluated to a value.
type Expression interface {
	Node
	// expression is unexported to ensure implementations of Expression
	// can only originate in this package.
	expression()
}

// VariableExpression is an expression that contains a single variable.
type VariableExpression struct {
	Variable string
}

// NumberLiteral is an expression that contains a single constant number.
type NumberLiteral struct {
	Value float64
}

// ArithmeticExpression is an expression that contains a +, -, *, / operator.
type ArithmeticExpression struct {
	LHS      Expression
	Operator Operator
	RHS      Expression
}

func (*Program) node()              {}
func (*Declaration) node()          {}
func (*AssignmentStatement) node()  {}
func (*VariableExpression) node()   {}
func (*NumberLiteral) node()        {}
func (*ArithmeticExpression) node() {}

func (*AssignmentStatement) statement() {}

func (*VariableExpression) expression()   {}
func (*NumberLiteral) expression()        {}
func (*ArithmeticExpression) expression() {}
