package lexer_test

import (
	"strings"
	"testing"

	"github.com/alongubkin/cpl-compiler/pkg/lexer"
)

func TestScannerOneCharacter(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("a"))
	assertToken(t, s, lexer.ID, "a")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerOneDigit(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("9"))
	assertToken(t, s, lexer.NUM, "9")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerID(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("heLlo"))
	assertToken(t, s, lexer.ID, "heLlo")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerLiterals(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("hello1 \n\n\n 1234 \t\n\t   hhh4h33 111 34"))
	assertToken(t, s, lexer.ID, "hello1")
	assertToken(t, s, lexer.NUM, "1234")
	assertToken(t, s, lexer.ID, "hhh4h33")
	assertToken(t, s, lexer.NUM, "111")
	assertToken(t, s, lexer.NUM, "34")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerKeywords(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(`
		break BreAK case default hello 
		else float    123     if input int a1  x
				output   static_cast   switch while
	`))
	assertToken(t, s, lexer.BREAK, "break")
	assertToken(t, s, lexer.ID, "BreAK")
	assertToken(t, s, lexer.CASE, "case")
	assertToken(t, s, lexer.DEFAULT, "default")
	assertToken(t, s, lexer.ID, "hello")
	assertToken(t, s, lexer.ELSE, "else")
	assertToken(t, s, lexer.FLOAT, "float")
	assertToken(t, s, lexer.NUM, "123")
	assertToken(t, s, lexer.IF, "if")
	assertToken(t, s, lexer.INPUT, "input")
	assertToken(t, s, lexer.INT, "int")
	assertToken(t, s, lexer.ID, "a1")
	assertToken(t, s, lexer.ID, "x")
	assertToken(t, s, lexer.OUTPUT, "output")
	assertToken(t, s, lexer.STATICCAST, "static_cast")
	assertToken(t, s, lexer.SWITCH, "switch")
	assertToken(t, s, lexer.WHILE, "while")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerSymbols(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(`(){,},    :;=`))
	assertToken(t, s, lexer.LPAREN, "(")
	assertToken(t, s, lexer.RPAREN, ")")
	assertToken(t, s, lexer.LBRACKET, "{")
	assertToken(t, s, lexer.COMMA, ",")
	assertToken(t, s, lexer.RBRACKET, "}")
	assertToken(t, s, lexer.COMMA, ",")
	assertToken(t, s, lexer.COLON, ":")
	assertToken(t, s, lexer.SEMICOLON, ";")
	assertToken(t, s, lexer.EQ, "=")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerInvalidIDs(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(`vvvvvvvvvv xx_y 111a`))
	assertToken(t, s, lexer.ILLEGAL, "vvvvvvvvvv")
	assertToken(t, s, lexer.ILLEGAL, "xx_y")
	assertToken(t, s, lexer.NUM, "111")
	assertToken(t, s, lexer.ID, "a")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerDecimalNumbers(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(`123.11 123. .456 0123.001`))
	assertToken(t, s, lexer.NUM, "123.11")
	assertToken(t, s, lexer.NUM, "123.")
	assertToken(t, s, lexer.ILLEGAL, ".")
	assertToken(t, s, lexer.NUM, "456")
	assertToken(t, s, lexer.NUM, "0123.001")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerOperators(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(`< = <= > = >= != ! = = == + - ++ -- * / | | || & & && ! !`))
	assertToken(t, s, lexer.RELOP, "<")
	assertToken(t, s, lexer.EQ, "=")
	assertToken(t, s, lexer.RELOP, "<=")
	assertToken(t, s, lexer.RELOP, ">")
	assertToken(t, s, lexer.EQ, "=")
	assertToken(t, s, lexer.RELOP, ">=")
	assertToken(t, s, lexer.RELOP, "!=")
	assertToken(t, s, lexer.NOT, "!")
	assertToken(t, s, lexer.EQ, "=")
	assertToken(t, s, lexer.EQ, "=")
	assertToken(t, s, lexer.RELOP, "==")
	assertToken(t, s, lexer.ADDOP, "+")
	assertToken(t, s, lexer.ADDOP, "-")
	assertToken(t, s, lexer.ADDOP, "+")
	assertToken(t, s, lexer.ADDOP, "+")
	assertToken(t, s, lexer.ADDOP, "-")
	assertToken(t, s, lexer.ADDOP, "-")
	assertToken(t, s, lexer.MULOP, "*")
	assertToken(t, s, lexer.MULOP, "/")
	assertToken(t, s, lexer.ILLEGAL, "|")
	assertToken(t, s, lexer.ILLEGAL, "|")
	assertToken(t, s, lexer.OR, "||")
	assertToken(t, s, lexer.ILLEGAL, "&")
	assertToken(t, s, lexer.ILLEGAL, "&")
	assertToken(t, s, lexer.AND, "&&")
	assertToken(t, s, lexer.NOT, "!")
	assertToken(t, s, lexer.NOT, "!")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerComments(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(
		`*/ hello /* break ** hello */ while/*asdf**/  /*asdfa*/ break hello 5 / * 4 /*`))
	assertToken(t, s, lexer.MULOP, "*")
	assertToken(t, s, lexer.MULOP, "/")
	assertToken(t, s, lexer.ID, "hello")
	assertToken(t, s, lexer.WHILE, "while")
	assertToken(t, s, lexer.BREAK, "break")
	assertToken(t, s, lexer.ID, "hello")
	assertToken(t, s, lexer.NUM, "5")
	assertToken(t, s, lexer.MULOP, "/")
	assertToken(t, s, lexer.MULOP, "*")
	assertToken(t, s, lexer.NUM, "4")
	assertToken(t, s, lexer.ILLEGAL, "")
	assertToken(t, s, lexer.EOF, "")

}

func TestScannerNestedComments(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader(
		`hello /* /* test */ id1 */ id2`))
	assertToken(t, s, lexer.ID, "hello")
	assertToken(t, s, lexer.ID, "id1")
	assertToken(t, s, lexer.MULOP, "*")
	assertToken(t, s, lexer.MULOP, "/")
	assertToken(t, s, lexer.ID, "id2")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerWhitespace(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("hello\n\t    \n\tbreak\t\t    \t\t test"))
	assertToken(t, s, lexer.ID, "hello")
	assertToken(t, s, lexer.BREAK, "break")
	assertToken(t, s, lexer.ID, "test")
	assertToken(t, s, lexer.EOF, "")
}

func TestScannerNotID(t *testing.T) {
	s := lexer.NewScanner(strings.NewReader("!id ! id1"))
	assertToken(t, s, lexer.NOT, "!")
	assertToken(t, s, lexer.ID, "id")
	assertToken(t, s, lexer.NOT, "!")
	assertToken(t, s, lexer.ID, "id1")
	assertToken(t, s, lexer.EOF, "")
}

func assertToken(t *testing.T, s *lexer.Scanner, tok lexer.Token, lit string) {
	real_tok, real_lit := s.Scan()
	if tok != real_tok {
		t.Errorf("Unexpected token %v (lexeme = %v)", real_tok, real_lit)
		return
	} else if lit != real_lit {
		t.Errorf("Token %v has unexpected lexeme: %v", real_tok, real_lit)
		return
	}

	t.Logf("Read token %v (lexeme %s)", real_tok, real_lit)
}
