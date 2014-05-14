package kingpin

import "strings"

type TokenType int

// Token types.
const (
	TokenShort TokenType = iota
	TokenLong
	TokenArg
	TokenEOF
)

var (
	TokenEOFMarker = Token{TokenEOF, ""}
)

type Token struct {
	Type  TokenType
	Value string
}

func (t *Token) IsFlag() bool {
	return t.Type == TokenShort || t.Type == TokenLong
}

func (t *Token) IsEOF() bool {
	return t.Type == TokenEOF
}

func (t *Token) String() string {
	switch t.Type {
	case TokenShort:
		return "-" + t.Value
	case TokenLong:
		return "--" + t.Value
	case TokenArg:
		return t.Value
	case TokenEOF:
		return "<EOF>"
	default:
		panic("unhandled type")
	}
}

type Tokens []*Token

func (t Tokens) String() string {
	out := []string{}
	for _, tok := range t {
		out = append(out, tok.String())
	}
	return "Tokens{" + strings.Join(out, ", ") + "}"
}

func (t Tokens) Next() (*Token, Tokens) {
	if len(t) == 0 {
		return &TokenEOFMarker, nil
	}
	return t[0], t[1:]
}

func (t Tokens) Return(token *Token) Tokens {
	if token.Type == TokenEOF {
		return t
	}
	return append(Tokens{token}, t...)
}

func (t Tokens) Peek() *Token {
	if len(t) == 0 {
		return &TokenEOFMarker
	}
	return t[0]
}

func Tokenize(args []string) (tokens Tokens) {
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg[2:], "=", 2)
			tokens = append(tokens, &Token{TokenLong, parts[0]})
			if len(parts) == 2 {
				tokens = append(tokens, &Token{TokenArg, parts[1]})
			}
		} else if strings.HasPrefix(arg, "-") {
			for _, a := range arg[1:] {
				tokens = append(tokens, &Token{TokenShort, string(a)})
			}
		} else {
			tokens = append(tokens, &Token{TokenArg, arg})
		}
	}
	return
}
