package parser

import (
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kashav/fsql/query"
	"github.com/kashav/fsql/tokenizer"
)

func Run(input string) (*query.Query, error) {
	return (&parser{}).parse(input)
}

type parser struct {
	tokenizer *tokenizer.Tokenizer
	current   *tokenizer.Token
	expected  tokenizer.TokenType
}

func (p *parser) parse(input string) (*query.Query, error) {
	q := query.NewQuery()
	p.tokenizer = tokenizer.NewTokenizer(input)
	if err := p.parseSelectClause(q); err != nil {
		return nil, err
	}
	if err := p.parseFromClause(q); err != nil {
		return nil, err
	}
	if err := p.parseWhereClause(q); err != nil {
		return nil, err
	}
	if err := p.parseLimitClause(q); err != nil {
		return nil, err
	}
	return q, nil
}

func (p *parser) parseSelectClause(q *query.Query) error {
	var showAll = true
	if p.expect(tokenizer.Select) == nil {
		if p.current == nil || p.current.Type == tokenizer.Identifier {
			showAll = false
		} else if p.current.Type == tokenizer.From || p.current.Type == tokenizer.Where {
			showAll = true
		} else {
			return p.currentError()
		}
	} else if current := p.expect(tokenizer.Identifier); current != nil {
		p.current = current
		showAll = false
	}

	if showAll {
		q.Attributes = allAttributes
	} else if err := p.parseAttrs(&q.Attributes, &q.Modifiers); err != nil {
		return err
	}

	return nil
}

func (p *parser) parseFromClause(q *query.Query) error {
	if p.expect(tokenizer.From) == nil {
		err := p.currentError()
		if p.expect(tokenizer.Identifier) != nil {
			return err
		}
		q.Sources["include"] = append(q.Sources["include"], ".")
		return nil
	}

	if err := p.parseSourceList(&q.Sources, &q.SourceAliases); err != nil {
		return err
	}

	u, err := user.Current()
	if err != nil {
		return err
	}
	for _, sourceType := range []string{"include", "exclude"} {
		for i, src := range q.Sources[sourceType] {
			if strings.Contains(src, "~") {
				q.Sources[sourceType][i] = filepath.Join(u.HomeDir, src[1:])
			}
		}
	}

	return nil
}

func (p *parser) parseWhereClause(q *query.Query) error {
	if p.expect(tokenizer.Where) == nil {
		err := p.currentError()
		if p.expect(tokenizer.Identifier) == nil {
			return nil
		}
		return err
	}
	root, err := p.parseConditionTree()
	if err != nil {
		return err
	}
	q.ConditionTree = root

	return nil
}

func (p *parser) parseLimitClause(q *query.Query) error {
	if p.expect(tokenizer.Limit) == nil {
		return nil
	}
	
	current := p.expect(tokenizer.Identifier)
	if current == nil {
		return p.currentError()
	}
	
	limit, err := strconv.Atoi(current.Raw)
	if err != nil {
		return err
	}
	
	q.Limit = limit
	return nil
}

func (p *parser) expect(t tokenizer.TokenType) *tokenizer.Token {
	p.expected = t

	if p.current == nil {
		p.current = p.tokenizer.Next()
	}

	if p.current != nil && p.current.Type == t {
		tok := p.current
		p.current = nil
		return tok
	}

	return nil
}
