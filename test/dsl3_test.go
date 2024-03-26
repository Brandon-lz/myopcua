package test

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TokenType int

const (
	NUMBER TokenType = iota
	PLUS
	MINUS
	TIMES
	DIVIDE
	ASSIGN
	SEMICOLON
	LPAREN
	RPAREN
	COMMA
	VAR
	IF
	THEN
	END
	GREATER
	LESS
	EOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	reader *bufio.Reader
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{reader: bufio.NewReader(r)}
}

func (l *Lexer) NextToken() Token {
	var ch byte
	var err error
	buf := strings.Builder{}
	continueMatch := false
	for {
		if !continueMatch {
			ch, err = l.reader.ReadByte()
			if err != nil {
				return Token{Type: EOF, Value: ""}
			}
			continueMatch = false
		}
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			return Token{Type: EOF, Value: ""}
		}
		switch ch {
		case '+':
			return Token{Type: PLUS, Value: "+"}
		case '-':
			return Token{Type: MINUS, Value: "-"}
		case '*':
			return Token{Type: TIMES, Value: "*"}
		case '/':
			return Token{Type: DIVIDE, Value: "/"}
		case '=':
			return Token{Type: ASSIGN, Value: "="}
		case '(':
			return Token{Type: LPAREN, Value: "("}
		case ')':
			return Token{Type: RPAREN, Value: ")"}
		case ',':
			return Token{Type: COMMA, Value: ","}
		case '>':
			return Token{Type: GREATER, Value: ">"}
		case '<':
			return Token{Type: LESS, Value: "<"}
		case ';':
			return Token{Type: SEMICOLON, Value: ";"}
		case 'I':
			buf.WriteByte(ch)
			ch, err = l.reader.ReadByte()
			if err != nil {
				return Token{Type: EOF, Value: ""}
			}
			if ch == 'F' {
				buf.WriteByte(ch)
				if buf.String() == "IF" {
					return Token{Type: IF, Value: "IF"}
				}else{
					continueMatch = true

				}
			} else {
				
			}
		case 'T':
			ch, err = l.reader.ReadByte()
			if err != nil {
				return Token{Type: EOF, Value: ""}
			}
			if ch == 'H' {
				ch, err = l.reader.ReadByte()
				if err != nil {
					return Token{Type: EOF, Value: ""}
				}
				if ch == 'E' {
					ch, err = l.reader.ReadByte()
					if err != nil {
						return Token{Type: EOF, Value: ""}
					}
					if ch == 'N' {
						return Token{Type: THEN, Value: "THEN"}
					} else {
						l.reader.UnreadByte()
					}
				} else {
					l.reader.UnreadByte()
				}
			} else {
				l.reader.UnreadByte()
			}
		case 'E':
			ch, err = l.reader.ReadByte()
			if err != nil {
				return Token{Type: EOF, Value: ""}
			}
			if ch == 'N' {
				ch, err = l.reader.ReadByte()
				if err != nil {
					return Token{Type: EOF, Value: ""}
				}
				if ch == 'D' {
					return Token{Type: END, Value: "END"}
				} else {
					l.reader.UnreadByte()
				}
			} else {
				l.reader.UnreadByte()
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // number
			buf.WriteByte(ch)
			for {
				ch, err = l.reader.ReadByte()
				if err != nil {
					break
				}
				if ch < '0' || ch > '9' && ch != '.' {
					l.reader.UnreadByte()
					return Token{Type: NUMBER, Value: buf.String()}
				} else {
					buf.WriteByte(ch)
				}
			}
		// variable
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'U', 'V', 'W', 'X', 'Y', 'Z': // variable
			buf := strings.Builder{}
			buf.WriteByte(ch)
			for {
				ch, err = l.reader.ReadByte()
				if err != nil {
					break
				}
				if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') {
					l.reader.UnreadByte()
					return Token{Type: VAR, Value: buf.String()}
				} else {
					buf.WriteByte(ch)
				}
			}
		default:
			fmt.Println("unknown token")
			return Token{Type: EOF, Value: ""}

		}
	}
}

func TestLexer(t *testing.T) {
	require := require.New(t)
	lexer := NewLexer(strings.NewReader(`
		IF a > 0 THEN
			b = a + 1;
		END
	`))
	token := lexer.NextToken()
	require.Equal(IF, token.Type)
	require.Equal("IF", token.Value)
	token = lexer.NextToken()
	require.Equal(VAR, token.Type)
	require.Equal("a", token.Value)
	token = lexer.NextToken()
	require.Equal(GREATER, token.Type)
	require.Equal(">", token.Value)
	token = lexer.NextToken()
	require.Equal(NUMBER, token.Type)
	require.Equal("0", token.Value)
	token = lexer.NextToken()
	require.Equal(THEN, token.Type)
	require.Equal("THEN", token.Value)
	token = lexer.NextToken()
	require.Equal(VAR, token.Type)
	require.Equal("b", token.Value)
	token = lexer.NextToken()
	require.Equal(ASSIGN, token.Type)
	require.Equal("=", token.Value)
	token = lexer.NextToken()
	require.Equal(VAR, token.Type)
	require.Equal("a", token.Value)
	token = lexer.NextToken()
	require.Equal(PLUS, token.Type)
	require.Equal("+", token.Value)
	token = lexer.NextToken()
	require.Equal(NUMBER, token.Type)
	require.Equal("1", token.Value)
	token = lexer.NextToken()
	require.Equal(SEMICOLON, token.Type)
	require.Equal(";", token.Value)
	token = lexer.NextToken()
	require.Equal(END, token.Type)
	require.Equal("END", token.Value)
	token = lexer.NextToken()
	require.Equal(EOF, token.Type)
	require.Equal("", token.Value)

}
