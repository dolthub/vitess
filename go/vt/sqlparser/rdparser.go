package sqlparser

import (
	"context"
	"fmt"
)

type parser struct {
	tok    *Tokenizer
	curOk  bool
	curId  int
	cur    []byte
	peekOk bool
	peekId int
	_peek  []byte
}

func (p *parser) parse(ctx context.Context, s string, options ParserOptions) (ret Statement, err error) {
	defer func() {
		if mes := recover(); mes != nil {
			_, ok := mes.(parseErr)
			if !ok {
				err = fmt.Errorf("panic encountered while parsing: %s", mes)
				return
			}
			ret, err = ParseWithOptions(ctx, s, options)
		}
	}()
	// get next token
	p.tok = NewStringTokenizer(s)
	if options.AnsiQuotes {
		p.tok = NewStringTokenizerForAnsiQuotes(s)
	}

	if prePlan, ok := p.statement(ctx); ok {
		return prePlan, nil
	}

	return ParseWithOptions(ctx, s, options)
}

type parseErr struct {
	str string
}

func (p *parser) fail(s string) {
	panic(parseErr{s})
}

func (p *parser) next() (int, []byte) {
	if p.peekOk {
		p.peekOk = false
		p.curId, p.cur = p.peekId, p._peek
	}
	p.curOk = true
	p.curId, p.cur = p.tok.Scan()
	return p.curId, p.cur
}

func (p *parser) peek() (int, []byte) {
	if !p.peekOk {
		p.peekOk = true
		p.peekId, p._peek = p.tok.Scan()
	}
	return p.peekId, p._peek
}
