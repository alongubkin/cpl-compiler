package parser

// DataType represents the primitive data types available in CPL.
type DataType int

const (
	// Float means the data type is a float.
	Float DataType = 1
	// Integer means the data type is an integer.
	Integer DataType = 2
)

// Node represents a node in the CPL abstract syntax tree.
type Node interface {
	// node is unexported to ensure implementations of Node
	// can only originate in this package.
	node()
	String() string
}

type Program struct {
	declarations []Declaration
}

type Declaration struct {
	Names []string
	Type  DataType
}

func (*Program) node()     {}
func (*Declaration) node() {}
