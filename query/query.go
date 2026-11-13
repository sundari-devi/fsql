package query

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrLimitReached = errors.New("query limit reached")

type Query struct {
	Attributes []string
	Modifiers  map[string][]Modifier

	Sources       map[string][]string
	SourceAliases map[string]string

	ConditionTree *ConditionNode

	Limit int
}

func NewQuery() *Query {
	return &Query{
		Attributes: make([]string, 0),
		Modifiers:  make(map[string][]Modifier),
		Sources: map[string][]string{
			"include": make([]string, 0),
			"exclude": make([]string, 0),
		},
		SourceAliases: make(map[string]string),
		ConditionTree: nil,
		Limit:         0,
	}
}

func (q *Query) HasAttribute(attributes ...string) bool {
	for _, attribute := range attributes {
		for _, queryAttribute := range q.Attributes {
			if attribute == queryAttribute {
				return true
			}
		}
	}
	return false
}

func (q *Query) Execute(workFunc interface{}) error {
	seen := map[string]bool{}
	excluder := &regexpExclude{exclusions: q.Sources["exclude"]}
	count := 0

	for _, src := range q.Sources["include"] {
		if strings.ContainsAny(src, "*?") {
			matches, err := filepath.Glob(src)
			if err != nil {
				return err
			}

			for _, match := range matches {
				err = filepath.Walk(match, q.walkFunc(seen, excluder, workFunc, &count))
				if err == ErrLimitReached {
					return nil
				}
				if err != nil {
					return err
				}
			}
			continue
		}

		err := filepath.Walk(src, q.walkFunc(seen, excluder, workFunc, &count))
		if err == ErrLimitReached {
			return nil
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *Query) walkFunc(seen map[string]bool, excluder Excluder,
	workFunc interface{}, count *int) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == "." {
			return nil
		}

		if _, ok := seen[path]; ok {
			return nil
		}
		seen[path] = true

		if excluder.shouldExclude(path) {
			return nil
		}

		if ok, err := q.ConditionTree.evaluateTree(path, info); err != nil {
			return err
		} else if !ok {
			return nil
		}

		results, err := q.applyModifiers(path, info)
		if err != nil {
			return err
		}
		
		workFunc.(func(string, os.FileInfo, map[string]interface{}))(path, info, results)
		
		*count++
		if q.Limit > 0 && *count >= q.Limit {
			return ErrLimitReached
		}

		return nil
	}
}
