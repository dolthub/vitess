package rd

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
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
	fmt.Fprintln(g.b, "package sqlparser")
	fmt.Fprintln(g.b, "import \"fmt\"")
	fmt.Fprintln(g.b, "import \"strings\"")
	return nil
	for _, line := range g.yacc.prefix {
		fmt.Fprintln(g.b, line)
	}
	return nil
}

func (g *recursiveGen) genStart() error {
	return nil
	start := g.yacc.start
	if start == "" {
		return fmt.Errorf("start function not found")
	}
	fmt.Fprintf(g.b, "func (p *parser) parse(yylex *Tokenizer) (Expr, bool) {\n")
	fmt.Fprintf(g.b, "  return p.%s(yylex)\n", g.yacc.start)
	fmt.Fprintf(g.b, "}\n\n")
	return nil
}

func (g *recursiveGen) genFunc(d *def) error {
	typ, ok := g.funcExprs[d.name]
	if !ok {
		log.Printf("function type not found: %s\n", d.name)
		typ = "[]byte"
	}
	//switch d.name {
	//case "lexer_old_position", "lexer_position", "special_comment_mode":
	//	return nil
	//}
	if d.name == "view_name_list" {
		print()
	}
	fmt.Fprintf(g.b, "func (p *parser) %s(yylex *Tokenizer) (%s, bool) {\n", d.name, typ)
	fmt.Fprintf(g.b, "  var ret %s\n", typ)
	var emptyRule string
	firstRule := true
	for _, r := range d.rules {
		if strings.Contains(r.name, "/*empty*/") {
			continue
		}
		// preprocess body
		bb := strings.Builder{}
		var usedVars int64
		for _, line := range r.body {
			line, vars, err := normalizeBodyLine(line)
			if err != nil {
				return err
			}
			usedVars |= vars
			fmt.Fprintf(&bb, "    %s\n", line)
		}
		fmt.Fprintf(&bb, "    return ret, true\n")

		parts := strings.Fields(r.name)
		if len(parts) == 0 {
			emptyRule = bb.String()
			continue
		}

		var okDefined bool
		for j, p := range parts {
			var cmp string
			if p == "openb" {
				cmp = "'('"
			} else if p == "closeb" {
				cmp = "')'"
			} else if p == "}" {
				cmp = "'}'"
			} else if p == "{" {
				cmp = "'{'"
			} else if _, ok := g.funcExprs[p]; ok {
			} else {
				cmp = p
			}
			if j == 0 {
				if firstRule {
					fmt.Fprintf(g.b, "  if ")
					firstRule = false
				} else {
					fmt.Fprintf(g.b, "  } else if ")
				}
				if len(parts) == 1 && len(r.body) == 0 {
					if cmp != "" {
						fmt.Fprintf(g.b, "id, tok := p.peek(); id == %s {\n", cmp)
						fmt.Fprintf(g.b, "    ret = tok\n")
					} else {
						fmt.Fprintf(g.b, "ret, ok := p.%s(yylex); ok {\n", p)
					}
				} else if cmp != "" {
					fmt.Fprintf(g.b, "id, _ := p.peek(); id == %s {\n", cmp)
					fmt.Fprintf(g.b, "    // %s\n", r.name)
					if setIncludes(usedVars, 1) {
						fmt.Fprintf(g.b, "    var1, _ := p.next()\n")
					} else {
						fmt.Fprintf(g.b, "    p.next()\n")
					}
				} else {
					if setIncludes(usedVars, 1) {
						fmt.Fprintf(g.b, "var1, ok := p.%s(yylex); ok {\n", p)
					} else {
						fmt.Fprintf(g.b, "_, ok := p.%s(yylex); ok {\n", p)
					}
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
				if setIncludes(usedVars, j+1) {
					fmt.Fprintf(g.b, "    var%d, ok := p.%s(yylex)\n", j+1, p)
					okDefined = true
				} else if okDefined {
					fmt.Fprintf(g.b, "    _, ok = p.%s(yylex)\n", p)
				} else {
					fmt.Fprintf(g.b, "    _, ok := p.%s(yylex)\n", p)
					okDefined = true
				}
				fmt.Fprintf(g.b, "    if !ok {\n")
				fmt.Fprintf(g.b, "      p.fail(\"expected: '%s: %s <%s>', found: 'string(tok)'\")\n", p, strings.Join(parts[:j], " "), p)
				fmt.Fprintf(g.b, "    }\n")
			}
		}
		//success, return
		fmt.Fprint(g.b, bb.String())
	}
	if emptyRule != "" && len(d.rules) == 1 {
		fmt.Fprint(g.b, emptyRule)
		fmt.Fprintf(g.b, "}\n\n")

		return nil
	}

	fmt.Fprintf(g.b, "  }\n")

	if emptyRule != "" {
		fmt.Fprint(g.b, emptyRule)
		fmt.Fprintf(g.b, "  return ret, true\n}\n\n")
	} else {
		fmt.Fprintf(g.b, "  return nil, false\n}\n\n")
	}
	return nil
}

var variableRe = regexp.MustCompile("\\$([1-6]+[0-9]*|[1-9])")

func normalizeBodyLine(r string) (string, int64, error) {
	r = strings.ReplaceAll(r, "$$ =", "ret =")
	r = strings.ReplaceAll(r, "return 1", "return nil, false")

	var variables int64
	r = strings.ReplaceAll(r, "$$", "ret")
	match := variableRe.FindAllStringSubmatchIndex(r, 1)
	for len(match) > 0 {
		m := match[0]
		start, end := m[2], m[3]
		_int64, err := strconv.ParseInt(r[start:end], 10, 64)
		if err != nil {
			return "", 0, fmt.Errorf("failed to parse variable string: %s", r[start:end])
		}
		if _int64 >= 64 {
			return "", 0, fmt.Errorf("variable reference too big: %d", _int64)
		}
		variables |= (1 << _int64)
		r = r[:m[0]] + "var" + r[m[0]+1:]
		match = variableRe.FindAllStringSubmatchIndex(r, 1)
	}

	//r = strings.ReplaceAll(r, "$1", "var1")
	//r = strings.ReplaceAll(r, "$2", "var2")
	//r = strings.ReplaceAll(r, "$3", "var3")
	//r = strings.ReplaceAll(r, "$4", "var4")
	//r = strings.ReplaceAll(r, "$5", "var5")
	//r = strings.ReplaceAll(r, "$6", "var6")
	//r = strings.ReplaceAll(r, "$7", "var7")
	//r = strings.ReplaceAll(r, "$8", "var8")
	//r = strings.ReplaceAll(r, "$9", "var9")
	return r, variables, nil
}

func setIncludes(set int64, i int) bool {
	return set&(1<<i) > 0
}
