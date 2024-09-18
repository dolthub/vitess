package rd

import (
	"fmt"
	"io"
	"strings"
)

// generator splits yacc file
// and then codegen's recursive descent
// create mappings for expression types
// gen prefix
// gen all of the functions, using lookups for function name and types

type recursiveGen struct {
	inf       io.Reader
	b         *strings.Builder
	outf      io.Writer
	yacc      *yaccFileContents
	goTypes   map[string]string
	funcExprs map[string]string
}

func newRecursiveGen(inf io.Reader, outf io.Writer) *recursiveGen {
	return &recursiveGen{
		inf:       inf,
		outf:      outf,
		b:         &strings.Builder{},
		goTypes:   make(map[string]string),
		funcExprs: make(map[string]string),
	}
}

func (g *recursiveGen) init() (err error) {
	g.yacc, err = split(g.inf)
	if err != nil {
		return err
	}
	if g.goTypes == nil {
		g.goTypes = make(map[string]string)
	}
	if g.funcExprs == nil {
		g.funcExprs = make(map[string]string)
	}
	for _, typ := range g.yacc.goTypes {
		g.goTypes[typ.name] = typ.typ
	}
	for _, d := range g.yacc.yaccTypes {
		goTyp, ok := g.goTypes[d.typ]
		if !ok {
			return fmt.Errorf("function expression not found: %s", d.name)
		}
		g.funcExprs[d.name] = goTyp
	}
	return nil
}

func (g *recursiveGen) gen() error {
	if g.yacc == nil {
		return fmt.Errorf("missing parsed yacc contents")
	}
	// prefix
	if err := g.genPrefix(); err != nil {
		return err
	}
	if err := g.genStart(); err != nil {
		return err
	}
	// each function
	for _, d := range g.yacc.defs {
		if err := g.genFunc(d); err != nil {
			return err
		}
	}
	return nil
}

func (g *recursiveGen) genPrefix() error {
	for _, line := range g.yacc.prefix {
		fmt.Fprintln(g.b, line)
	}
	return nil
}

func (g *recursiveGen) genStart() error {
	start := g.yacc.start
	if start == "" {
		return fmt.Errorf("start function not found")
	}
	fmt.Fprintf(g.b, "func (p *parser) parse() (statement, bool) {\n")
	fmt.Fprintf(g.b, "  return %s()\n", g.yacc.start)
	fmt.Fprintf(g.b, "}\n\n")
	return nil
}

func (g *recursiveGen) genFunc(d *def) error {
	fmt.Fprintf(g.b, "func (p *parser) %s() (Expr, bool) {\n", d.name)
	fmt.Fprintf(g.b, "  id, tok := p.peek()\n")
	for i, r := range d.rules {
		parts := strings.Fields(r.name)
		for j, p := range parts {
			var cmp string
			if p == "openb" {
				cmp = "'('"
			} else if p == "closeb" {
				cmp = "')'"
			} else if _, ok := g.funcExprs[p]; ok {
			} else {
				cmp = p
			}
			if j == 0 {
				if i == 0 {
					fmt.Fprintf(g.b, "  if ")
				} else {
					fmt.Fprintf(g.b, "} else if ")
				}
				if cmp != "" {
					fmt.Fprintf(g.b, "id == %s {\n", cmp)
					fmt.Fprintf(g.b, "    // %s\n", r.name)
					fmt.Fprintf(g.b, "    var1, _ := p.next()\n")
				} else {
					fmt.Fprintf(g.b, "var1, ok := p.%s(); ok {\n", p)
					fmt.Fprintf(g.b, "    // %s\n", r.name)
				}
				continue
			}
			if cmp != "" {
				fmt.Fprintf(g.b, "    var%d, tok := p.next()\n", j+1)
				fmt.Fprintf(g.b, "    if var%d != %s {\n", j+1, cmp)
				fmt.Fprintf(g.b, "      p.fail(\"expected: '%s: %s <%s>', found: '\" + string(tok) + \"'\")\n", p, strings.Join(parts[:j], " "), cmp)
				fmt.Fprintf(g.b, "    }\n")
			} else if _, ok := g.funcExprs[p]; ok {
				fmt.Fprintf(g.b, "    var%d, ok := p.%s()\n", j+1, p)
				fmt.Fprintf(g.b, "    if !ok {\n")
				fmt.Fprintf(g.b, "      p.fail(\"expected: '%s: %s <%s>', found: 'string(tok)'\")\n", p, strings.Join(parts[:j], " "), p)
				fmt.Fprintf(g.b, "    }\n")
			}
		}
		//success, return
		for _, r := range r.body {
			fmt.Fprintf(g.b, "    %s\n", r)
		}
		fmt.Fprintf(g.b, "  }\n")
	}
	fmt.Fprintf(g.b, "\n  return nil, false\n}\n")
	return nil
}
