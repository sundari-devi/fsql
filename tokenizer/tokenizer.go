package tokenizer

import (
	"fmt"
	"strings"
	"unicode"
)

func (t *Tokenizer) All() []Token {
	tokens := []Token{}
	for tok := t.Next(); tok != nil; tok = t.Next() {
		tokens = append(tokens, *tok)
	}
	return tokens
}

func (t *Tokenizer) Next() *Token {
	for unicode.IsSpace(t.current()) {
		t.input = t.input[1:]
	}

	current := t.current()
	if current == -1 {
		return nil
	}

	switch current {
	case '(':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: OpenParen, Raw: "("})
	case ')':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: CloseParen, Raw: ")"})
	case '[':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: OpenBracket, Raw: "["})
	case ']':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: CloseBracket, Raw: "]"})
	case ',':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: Comma, Raw: ","})
	case '-':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: Hyphen, Raw: "-"})
	case '!':
		if t.getRuneAt(1) == '=' {
			t.input = t.input[2:]
			return t.setToken(&Token{Type: NotEquals, Raw: "!="})
		}
		return t.setToken(&Token{Type: ExclamationMark, Raw: "!"})
	case '=':
		t.input = t.input[1:]
		return t.setToken(&Token{Type: Equals, Raw: "="})
	case '>':
		if t.getRuneAt(1) == '=' {
			t.input = t.input[2:]
			return t.setToken(&Token{Type: GreaterThanEquals, Raw: ">="})
		}
		t.input = t.input[1:]
		return t.setToken(&Token{Type: GreaterThan, Raw: ">"})
	case '<':
		if t.getRuneAt(1) == '=' {
			t.input = t.input[2:]
			return t.setToken(&Token{Type: LessThanEquals, Raw: "<="})
		}
		if t.getRuneAt(1) == '>' {
			t.input = t.input[2:]
			return t.setToken(&Token{Type: NotEquals, Raw: "<>"})
		}
		t.input = t.input[1:]
		return t.setToken(&Token{Type: LessThan, Raw: "<"})
	}

	if !t.currentIs(-1, ',', '\'', '"', '`', '(', ')', '[', ']') {
		word := t.readWord()
		tok := &Token{Raw: word}

		switch strings.ToUpper(word) {
		case "SELECT":
			tok.Type = Select
		case "FROM":
			tok.Type = From
		case "WHERE":
			tok.Type = Where
		case "LIMIT":
			tok.Type = Limit
		case "AS":
			tok.Type = As
		case "OR":
			tok.Type = Or
		case "AND":
			tok.Type = And
		case "NOT":
			tok.Type = Not
		case "IN":
			tok.Type = In
		case "IS":
			tok.Type = Is
		case "LIKE":
			tok.Type = Like
		case "REGEXP", "RLIKE":
			tok.Type = RLike
		default:
			tok.Type = Identifier
		}

		if t.getPreviousToken() != nil && t.getPreviousToken().Type == OpenParen &&
			t.getTokenAt(1) != nil && t.getTokenAt(1).Type == In {
			tok.Type = Subquery
			tok.Raw = fmt.Sprintf("%s %s", word, t.readQuery())
		}

		return t.setToken(tok)
	}

	tok := &Token{Type: Unknown, Raw: string(current)}

	if t.currentIs('\'', '"', '`') {
		t.input = t.input[1:]
		tok.Raw = t.readWord() + t.readUntil(current)
		tok.Type = Identifier
	}

	t.input = t.input[1:]
	return t.setToken(tok)
}

func (t *Tokenizer) setToken(token *Token) *Token {
	t.tokens = append(t.tokens, token)
	return token
}

func (t *Tokenizer) getPreviousToken() *Token {
	return t.getTokenAt(0)
}

func (t *Tokenizer) getTokenAt(i int) *Token {
	j := len(t.tokens) - 1 - i
	if j < 0 {
		return nil
	}
	return t.tokens[j]
}

func (t *Tokenizer) current() rune {
	return t.getRuneAt(0)
}

func (t *Tokenizer) getRuneAt(i int) rune {
	if len(t.input) == i {
		return -1
	}
	return t.input[i]
}

func (t *Tokenizer) currentIs(runes ...rune) bool {
	for _, r := range runes {
		if r == t.current() {
			return true
		}
	}
	return false
}

func (t *Tokenizer) readWord() string {
	word := []rune{}
	for {
		if unicode.IsSpace(t.current()) ||
			t.currentIs(-1, ',', '\'', '"', '`', '(', ')', '[', ']') {
			return string(word)
		}
		word = append(word, t.current())
		t.input = t.input[1:]
	}
}

func (t *Tokenizer) readQuery() string {
	var query string
	var count = 1
	for count > 0 {
		for unicode.IsSpace(t.current()) {
			t.input = t.input[1:]
		}
		word := fmt.Sprintf("%s", t.readWord())
		if t.current() == -1 {
			break
		}
		if t.current() == '(' {
			count++
			word = "("
		} else if t.current() == ')' {
			count--
			if count <= 0 {
				query += word
				break
			}
			word = ")"
		} else if t.currentIs('\'', '`') {
			word += string(t.current())
		} else {
			word += " "
		}
		query += word
		t.input = t.input[1:]
	}
	return query
}

func (t *Tokenizer) readUntil(runes ...rune) string {
	var word string
	for !t.currentIs(runes...) {
		for unicode.IsSpace(t.current()) {
			t.input = t.input[1:]
		}
		word = fmt.Sprintf("%s %s", word, t.readWord())
	}
	return word
}
