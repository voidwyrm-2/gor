package main

import (
	"fmt"
	"slices"
)

// token types
type tokType string

const (
	DOT tokType = "DOT"

	LPAREN   tokType = "LPAREN"
	RPAREN   tokType = "RPAREN"
	LBRACKET tokType = "LBRACKET"
	RBRACKET tokType = "RBRACKET"
	LBRACE   tokType = "LBRACE"
	RBRACE   tokType = "RBRACE"

	PLUS          = "PLUS"
	HYPHEN        = "HYPHEN"
	ASTERISK      = "ASTERISK"
	FORWARD_SLASH = "FORWARD_SLASH"
	BACK_SLASH    = "BACK_SLASH"
	COLON         = "COLON"
	PERCENT_SIGN  = "PERCENT_SIGN"

	IDENT   tokType = "IDENT"
	NUMBER  tokType = "NUMBER"
	STRING  tokType = "STRING"
	KEYWORD tokType = "KEYWORD"

	EQUALS       tokType = "EQUALS"
	NOT_EQUALS   tokType = "NOT_EQUALS"
	LESSER_THAN  tokType = "LESSER_THAN"
	GREATER_THAN tokType = "GREATER_THAN"

	ASSIGN tokType = "ASSIGN"

	NEWLINE tokType = "NEWLINE"
	COMMENT tokType = "COMMENT"
	EOF     tokType = "EOF"
)

var KEYWORDS = []string{
	"func",
	"con",
	"delete",
	"jumpto",
	"if",
	"elsif",
	"else",
	"use",
}

func isValidForIdent(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

type Token struct {
	Type           tokType
	Lit            string
	Start, End, Ln int
}

func (t Token) Length() int {
	return (t.End - t.Start) + 1
}

func (t Token) IsMathOp() bool {
	return t.Type == PLUS || t.Type == HYPHEN || t.Type == ASTERISK || t.Type == FORWARD_SLASH
}

func (t Token) Istype(_type tokType) bool {
	return t.Type == _type
}

func NewToken(t tokType, literal string, start, end, ln int) Token {
	return Token{t, literal, start, end, ln}
}

func NewNilToken(t tokType, start, end, ln int) Token {
	return Token{t, "", start, end, ln}
}

type Lexer struct {
	text    string
	cchar   rune
	idx, ln int
}

func (l *Lexer) advance() {
	l.idx++

	if l.idx < len(l.text) {
		l.cchar = rune(l.text[l.idx])
	} else {
		l.cchar = rune(-1)
	}

	if l.cchar == '\n' {
		l.ln++
	}
}

func (l Lexer) Lex() ([]Token, error) {
	var tokens []Token

	for l.cchar != -1 {
		//fmt.Println(l.idx, string(l.cchar))
		switch l.cchar {
		case ' ', '\t':
			l.advance()
		case '\n':
			tokens = append(tokens, NewToken(NEWLINE, "\n", l.idx, l.idx, l.ln))
			l.advance()
		case '(':
			tokens = append(tokens, NewToken(LPAREN, "(", l.idx, l.idx, l.ln))
			l.advance()
		case ')':
			tokens = append(tokens, NewToken(RPAREN, ")", l.idx, l.idx, l.ln))
			l.advance()
		case '[':
			tokens = append(tokens, NewToken(LBRACKET, "[", l.idx, l.idx, l.ln))
			l.advance()
		case ']':
			tokens = append(tokens, NewToken(RBRACKET, "]", l.idx, l.idx, l.ln))
			l.advance()
		case '{':
			tokens = append(tokens, NewToken(LBRACE, "{", l.idx, l.idx, l.ln))
			l.advance()
		case '}':
			tokens = append(tokens, NewToken(RBRACE, "}", l.idx, l.idx, l.ln))
			l.advance()
		case '.':
			tokens = append(tokens, NewToken(DOT, ".", l.idx, l.idx, l.ln))
			l.advance()
		case '<':
			l.advance()
			if l.cchar == '-' {
				tokens = append(tokens, NewToken(ASSIGN, "<-", l.idx, l.idx, l.ln))
				l.advance()
			} else {
				tokens = append(tokens, NewToken(LESSER_THAN, "<", l.idx, l.idx, l.ln))
			}
		case '>':
			tokens = append(tokens, NewToken(GREATER_THAN, ">", l.idx, l.idx, l.ln))
			l.advance()
		case '+':
			tokens = append(tokens, NewToken(PLUS, "+", l.idx, l.idx, l.ln))
			l.advance()
		case '-':
			tokens = append(tokens, NewToken(HYPHEN, "-", l.idx, l.idx, l.ln))
			l.advance()
		case '*':
			tokens = append(tokens, NewToken(ASTERISK, "*", l.idx, l.idx, l.ln))
			l.advance()
		case '/':
			tokens = append(tokens, NewToken(FORWARD_SLASH, "/", l.idx, l.idx, l.ln))
			l.advance()
		case '\\':
			tokens = append(tokens, NewToken(BACK_SLASH, "\\", l.idx, l.idx, l.ln))
			l.advance()
		case ':':
			tokens = append(tokens, NewToken(COLON, ":", l.idx, l.idx, l.ln))
			l.advance()
		case '%':
			tokens = append(tokens, NewToken(PERCENT_SIGN, "%", l.idx, l.idx, l.ln))
			l.advance()
		case '?':
			l.advance()
			tokens = append(tokens, l.collectComment())
		case '"':
			l.advance()
			tokens = append(tokens, l.collectString())
		default:
			if l.cchar >= '0' && l.cchar <= '9' {
				tok, err := l.collectNumber()
				if err != nil {
					return []Token{}, err
				}
				tokens = append(tokens, tok)
			} else if isValidForIdent(l.cchar) {
				tokens = append(tokens, l.collectIdent())
			} else {
				return []Token{}, NewGorError(NewNilToken(DOT, l.idx, l.idx, l.ln), fmt.Sprintf("illegal character '%c'", l.cchar))
			}
		}
	}

	tokens = append(tokens, NewNilToken(EOF, l.idx, l.idx, l.ln))
	return tokens, nil
}

func (l *Lexer) collectComment() Token {
	start := l.idx
	comment_str := ""

	for l.cchar == ' ' {
		l.advance()
	}

	for l.cchar != -1 && l.cchar != '\n' {
		comment_str += string(l.cchar)
		l.advance()
	}

	return NewToken(COMMENT, comment_str, start, l.idx-1, l.ln)
}

func (l *Lexer) collectIdent() Token {
	start := l.idx
	ident_str := ""

	for l.cchar != -1 && isValidForIdent(l.cchar) {
		ident_str += string(l.cchar)
		l.advance()
	}

	if slices.Contains(KEYWORDS, ident_str) {
		return NewToken(KEYWORD, ident_str, start, l.idx-1, l.ln)
	}
	return NewToken(IDENT, ident_str, start, l.idx-1, l.ln)
}

func (l *Lexer) collectString() Token {
	start := l.idx
	string_str := ""

	for l.cchar != -1 && l.cchar != '"' {
		string_str += string(l.cchar)
		l.advance()
	}

	l.advance()

	return NewToken(STRING, string_str, start, l.idx-1, l.ln)
}

func (l *Lexer) collectNumber() (Token, error) {
	start := l.idx
	num_str := ""
	hasDot := false

	for l.cchar != -1 && (l.cchar >= '0' && l.cchar <= '9' || l.cchar == '.') {
		if l.cchar == '.' {
			if hasDot {
				return Token{}, fmt.Errorf("")
			}
			hasDot = true
		}
		num_str += string(l.cchar)
		l.advance()
	}

	return NewToken(NUMBER, num_str, start, l.idx-1, l.ln), nil
}

func NewLexer(text string) Lexer {
	lexer := Lexer{}

	lexer.text = text
	lexer.ln = 1
	lexer.idx = -1
	lexer.advance()

	return lexer
}
