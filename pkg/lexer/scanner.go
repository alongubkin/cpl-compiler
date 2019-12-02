package lexer

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// TODO: Comment
const MaxIdentifierLength = 9

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// Skip comments and whitespaces.
	for {
		if ch == '/' {
			ch2 := s.read()
			if ch2 == '*' {
				if err := s.skipUntilEndComment(); err != nil {
					return ILLEGAL, ""
				}
			} else {
				s.unread()
				break
			}
		} else if isWhitespace(ch) {
			s.scanWhitespace()
		} else {
			break
		}

		ch = s.read()
	}

	// If we see a letter then consume as an ID or reserved word.
	if isLetter(ch) {
		s.unread()
		return s.scanIdentifier()
	} else if isDigit(ch) {
		s.unread()
		return s.scanNumber()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""

	case '>', '<':
		ch2 := s.read()
		if ch2 == '=' {
			return RELOP, string(ch) + string(ch2)
		}

		s.unread()
		return RELOP, string(ch)

	case '=':
		ch2 := s.read()
		if ch2 == '=' {
			return RELOP, "=="
		}

		s.unread()
		return EQ, string(ch)

	case '!':
		ch2 := s.read()
		if ch2 == '=' {
			return RELOP, "!="
		}

		s.unread()
		return NOT, string(ch)

	case '|':
		ch2 := s.read()
		if ch2 == '|' {
			return OR, "||"
		}

		s.unread()
		return ILLEGAL, string(ch)

	case '&':
		ch2 := s.read()
		if ch2 == '&' {
			return AND, "&&"
		}

		s.unread()
		return ILLEGAL, string(ch)

	case '+', '-':
		return ADDOP, string(ch)

	case '*', '/':
		return MULOP, string(ch)

	case ';':
		return SEMICOLON, string(ch)

	case '(':
		return LPAREN, string(ch)

	case ')':
		return RPAREN, string(ch)

	case '{':
		return LBRACKET, string(ch)

	case '}':
		return RBRACKET, string(ch)

	case ',':
		return COMMA, string(ch)

	case ':':
		return COLON, string(ch)
	}

	return ILLEGAL, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() {
	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		}
	}
}

// scanIdentifier consumes the current rune and all contiguous identifier runes.
func (s *Scanner) scanIdentifier() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch buf.String() {
	case "break":
		return BREAK, buf.String()
	case "case":
		return CASE, buf.String()
	case "default":
		return DEFAULT, buf.String()
	case "else":
		return ELSE, buf.String()
	case "float":
		return FLOAT, buf.String()
	case "if":
		return IF, buf.String()
	case "input":
		return INPUT, buf.String()
	case "int":
		return INT, buf.String()
	case "output":
		return OUTPUT, buf.String()
	case "switch":
		return SWITCH, buf.String()
	case "while":
		return WHILE, buf.String()
	case "static_cast":
		return STATICCAST, buf.String()
	}

	// Otherwise return as a regular identifier - just need to make sure its length is okay
	// and it doesn't contain an underscore, which is an illegal character for IDs.
	if len(buf.String()) <= MaxIdentifierLength && !strings.ContainsRune(buf.String(), '_') {
		return ID, buf.String()
	}

	return ILLEGAL, buf.String()
}

// scanNumber consumes a contiguous series of digits.
func (s *Scanner) scanNumber() (tok Token, lit string) {
	var buf bytes.Buffer
	for {
		ch := s.read()
		if !isDigit(ch) && ch != '.' {
			s.unread()
			break
		}
		_, _ = buf.WriteRune(ch)
	}

	return NUM, buf.String()
}

// skipUntilEndComment skips characters until it reaches a '*/' symbol.
func (s *Scanner) skipUntilEndComment() error {
	for {
		if ch := s.read(); ch == '*' {
			// We might be at the end.
		star:
			ch2 := s.read()
			if ch2 == '/' {
				return nil
			} else if ch2 == '*' {
				// We are back in the state machine since we see a star.
				// TODO: Remove goto
				goto star
			} else if ch2 == eof {
				return io.EOF
			}
		} else if ch == eof {
			return io.EOF
		}
	}
}
