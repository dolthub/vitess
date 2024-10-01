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
	for _, d := range g.yacc.defs {
		if _, ok := g.funcExprs[d.name]; !ok {
			g.funcExprs[d.name] = "[]byte"
		}
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
	for _, n := range g.yacc.defNames {
		d := g.yacc.defs[n]
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
		g.funcExprs[d.name] = "[]byte"
	}

	fmt.Fprintf(g.b, "func (p *parser) %s(yylex *Tokenizer) (%s, bool) {\n", d.name, typ)
	fmt.Fprintf(g.b, "  var ret %s\n", typ)
	fmt.Fprintf(g.b, "  var matched bool\n")

	g.dfsRuleGen(d.name, d.rules, 1, true, "  ")

	//var emptyRule string
	//firstRule := true
	// recursive can have empty rule
	// recursive structure:
	//  - leaf rules
	//  - if leaf rule nil return nil
	//  - if leaf not nil try to match next characters

	//dfsRuleGen(d.rules, true, "  ")

	//if emptyRule != "" && len(d.rules) == 1 {
	//	fmt.Fprint(g.b, emptyRule)
	//	fmt.Fprintf(g.b, "}\n\n")
	//
	//	return nil
	//}
	//if len(d.selfRec) > 0 {
	//	fmt.Fprintf(g.b, "    if ret == nil {\n")
	//	fmt.Fprintf(g.b, "    	return ret, false\n")
	//	fmt.Fprintf(g.b, "    }\n")
	//	fmt.Fprintf(g.b, "    var1 := ret\n")
	//
	//}

	//fmt.Fprintf(g.b, "  }\n")
	//
	//if emptyRule != "" {
	//	fmt.Fprint(g.b, emptyRule)
	//	fmt.Fprintf(g.b, "  return ret, true\n}\n\n")
	//} else {
	//	fmt.Fprintf(g.b, "  return ret, false\n}\n\n")

	//defRet := "nil"
	//if typ == "string" {
	//	defRet = "\"\""
	//} else if typ == "int" {
	//	defRet = "0"
	//} else if strings.HasPrefix(typ, "*") {
	//} else if typ == "Expr" {
	//} else {
	//	defRet = typ + "{}"
	//}
	//fmt.Fprintf(g.b, "  return %s, false\n}\n\n", defRet)
	//}
	return nil
}

func (g *recursiveGen) dfsRuleGen(defName string, r *rulePrefix, fid int, first bool, indent string) error {
	var nesting int
	if r.prefix != "" {
		match, err := g.nestMatch(r.prefix, fid, first, r.usedVars, indent)
		if err != nil {
			return err
		}
		indent += "  "
		fid++
		first = true
		nesting++
		fmt.Fprintf(g.b, match)
	}

	if r.prefix == "value_expression" {
		print()
	}
	// try to match terminal rules, easiest to distinguish
	for _, r := range r.term {
		// match terminal rule
		match, err := g.termRuleMatch(defName, r, fid, first, r.usedVars, indent)
		if err != nil {
			return err
		}
		first = false
		fmt.Fprintf(g.b, match)
	}

	// nested rule clusters
	for _, r := range r.pref {
		err := g.dfsRuleGen(defName, r, fid, first, indent)
		if err != nil {
			return err
		}
		first = false
	}

	if r.empty != nil {
		// returns before recursive
		if first {
			// no other rules in block
			fmt.Fprintf(g.b, "%sreturn ret, true\n", indent)
		} else {
			fmt.Fprintf(g.b, " else {\n")
			fmt.Fprintf(g.b, "%s  return ret, true\n", indent)
			fmt.Fprintf(g.b, "%s}\n", indent)
		}
	}

	// recursive is special
	// field match is different, because it's whether we've matched current scope
	if r.rec != nil {
		// if no match, we don't recurse
		fmt.Fprintf(g.b, "\n")
		fmt.Fprintf(g.b, "%sif !matched {\n", indent)
		fmt.Fprintf(g.b, "%s  return ret, false\n", indent)
		fmt.Fprintf(g.b, "%s}\n", indent)
		// otherwise |ret| becomes |var1| for the left recursion
		// rule body
		if setIncludes(r.rec.usedVars, 1) {
			fmt.Fprintf(g.b, "%svar1 := ret\n", indent)
		}

		first = true
		r.rec.prefix = ""
		err := g.dfsRuleGen(defName, r.rec, fid+1, first, indent)
		if err != nil {
			return err
		}
	}

	for i := 0; i < nesting; i++ {
		indent = indent[2:]
		fmt.Fprintf(g.b, "\n%s}", indent)
	}
	if fid == 1 {
		//only top level rule closes function body
		fmt.Fprintf(g.b, "\n%sreturn ret, matched\n", indent)
		indent = indent[2:]
		fmt.Fprintf(g.b, "%s}\n\n", indent)
	}
	return nil
}

func (g *recursiveGen) termRuleMatch(defName string, r *rule, fid int, first bool, usedVars int64, indent string) (string, error) {
	var b strings.Builder

	// nest the first token
	match, err := g.nestMatch(r.fields[0], fid, first, r.usedVars, indent)
	if err != nil {
		return "", err
	}
	fmt.Fprint(&b, match)
	fid++

	// rest of tokens have to match
	var okDefined bool
	for _, f := range r.fields[1:] {
		match, okDefined, err = g.termMatch(defName, f, r.name, fid, first, okDefined, r.usedVars, indent+"  ")
		if err != nil {
			return "", err
		}
		fid++
		fmt.Fprint(&b, match)
	}

	if len(r.body) == 0 {
		fmt.Fprintf(&b, "  %sret = var1\n", indent)

	}
	for _, l := range r.body {
		fmt.Fprintf(&b, "  %s%s\n", indent, l)
	}

	// if all tokens match, set |matched| var
	fmt.Fprintf(&b, "  %smatched = true\n", indent)

	fmt.Fprintf(&b, "%s}", indent)

	return b.String(), nil
}

func (g *recursiveGen) nestMatch(f string, fid int, first bool, usedVars int64, indent string) (string, error) {
	var b strings.Builder
	var cmp string
	if f == "openb" {
		cmp = "'('"
	} else if f == "closeb" {
		cmp = "')'"
	} else if f == "}" {
		cmp = "'}'"
	} else if f == "{" {
		cmp = "'{'"
	} else if f == "'%'" {
		cmp = "'%%'"
	} else if _, ok := g.funcExprs[f]; ok {
	} else {
		cmp = f
	}

	if first {
		fmt.Fprintf(&b, "%sif ", indent)
	} else {
		fmt.Fprintf(&b, " else if ")
	}

	if cmp != "" {
		fmt.Fprintf(&b, "id, _ := p.peek(); id == %s {\n", cmp)
		if setIncludes(usedVars, fid) {
			fmt.Fprintf(&b, "%s  _, var%d := p.next()\n", indent, fid)
		} else {
			fmt.Fprintf(&b, "%s  p.next()\n", indent)
		}
	} else {
		if setIncludes(usedVars, fid) {
			fmt.Fprintf(&b, "var%d, ok := p.%s(yylex); ok {\n", fid, f)
		} else {
			fmt.Fprintf(&b, "_, ok := p.%s(yylex); ok {\n", f)
		}
	}
	return b.String(), nil
}

func (g *recursiveGen) termMatch(defName, f string, ruleFields string, fid int, first, okDefined bool, usedVars int64, indent string) (string, bool, error) {
	var b strings.Builder
	var cmp string
	if f == "openb" {
		cmp = "'('"
	} else if f == "closeb" {
		cmp = "')'"
	} else if f == "}" {
		cmp = "'}'"
	} else if f == "{" {
		cmp = "'{'"
	} else if f == "'%'" {
		cmp = "'%%%'"
	} else if _, ok := g.funcExprs[f]; ok {
	} else {
		cmp = f
	}

	if cmp != "" {
		fmt.Fprintf(&b, "%sid%d, var%d := p.next()\n", indent, fid, fid)
		fmt.Fprintf(&b, "%sif id%d != %s {\n", indent, fid, cmp)
		fmt.Fprintf(&b, "%s  p.fail(\"rule: '%s->%s', field %d expected '<%s>', found: '\" + string(var%d) + \"'\")\n", indent, defName, ruleFields, fid, f, fid)
		fmt.Fprintf(&b, "%s}\n", indent)
	} else if _, ok := g.funcExprs[f]; ok {
		fmt.Fprintf(&b, "%s_, tok%d := p.peek()\n", indent, fid)
		if setIncludes(usedVars, fid) {
			fmt.Fprintf(&b, "%svar%d, ok := p.%s(yylex)\n", indent, fid, f)
			okDefined = true
		} else if okDefined {
			fmt.Fprintf(&b, "%s_, ok = p.%s(yylex)\n", indent, f)
		} else {
			fmt.Fprintf(&b, "%s_, ok := p.%s(yylex)\n", indent, f)
			okDefined = true
		}
		fmt.Fprintf(&b, "%sif !ok {\n", indent)
		fmt.Fprintf(&b, "%s  p.fail(\"rule: '%s->%s', field %d expected '<%s>', found: '\"+string(tok%d)+\"'\")\n", indent, defName, ruleFields, fid, f, fid)
		fmt.Fprintf(&b, "%s}\n", indent)
	}
	return b.String(), okDefined, nil
}

//func (g *recursiveGen) genRule(r *rule, indent string) error {
//	if strings.Contains(r.name, "/*empty*/") {
//		return nil
//	}
//
//	// preprocess body
//	bb := strings.Builder{}
//	var usedVars int64
//	for _, line := range r.body {
//		line, vars, err := normalizeBodyLine(line)
//		if err != nil {
//			return err
//		}
//		usedVars |= vars
//		fmt.Fprintf(&bb, "    %s%s\n", indent, line)
//	}
//	fmt.Fprintf(&bb, "    return ret, true\n")
//
//	parts := strings.Fields(r.name)
//	if len(parts) == 0 {
//		emptyRule = bb.String()
//		continue
//	}
//
//	var okDefined bool
//	for j, p := range parts {
//		var cmp string
//		if p == "openb" {
//			cmp = "'('"
//		} else if p == "closeb" {
//			cmp = "')'"
//		} else if p == "}" {
//			cmp = "'}'"
//		} else if p == "{" {
//			cmp = "'{'"
//		} else if _, ok := g.funcExprs[p]; ok {
//		} else {
//			cmp = p
//		}
//
//		if j == 0 {
//			if firstRule {
//				fmt.Fprintf(g.b, "  if ")
//				firstRule = false
//			} else {
//				fmt.Fprintf(g.b, "  } else if ")
//			}
//			if len(parts) == 1 && len(r.body) == 0 {
//				if cmp != "" {
//					fmt.Fprintf(g.b, "id, var1 := p.peek(); id == %s {\n", cmp)
//					fmt.Fprintf(g.b, "    ret = var1\n")
//				} else {
//					fmt.Fprintf(g.b, "ret, ok := p.%s(yylex); ok {\n", p)
//				}
//			} else if cmp != "" {
//				fmt.Fprintf(g.b, "id, _ := p.peek(); id == %s {\n", cmp)
//				fmt.Fprintf(g.b, "    // %s\n", r.name)
//				if setIncludes(usedVars, 1) {
//					//if isLit(cmp) {
//					fmt.Fprintf(g.b, "    _, var1 := p.next()\n")
//					//} else {
//					//	fmt.Fprintf(g.b, "    var1, _ := p.next()\n")
//					//}
//				} else {
//					fmt.Fprintf(g.b, "    p.next()\n")
//				}
//			} else {
//				if setIncludes(usedVars, 1) {
//					fmt.Fprintf(g.b, "var1, ok := p.%s(yylex); ok {\n", p)
//				} else {
//					fmt.Fprintf(g.b, "_, ok := p.%s(yylex); ok {\n", p)
//				}
//				fmt.Fprintf(g.b, "    // %s\n", r.name)
//			}
//			continue
//		}
//		if cmp != "" {
//			fmt.Fprintf(g.b, "    id%d, var%d := p.next()\n", j+1, j+1)
//			fmt.Fprintf(g.b, "    if id%d != %s {\n", j+1, cmp)
//			fmt.Fprintf(g.b, "      p.fail(\"expected: '%s: %s <%s>', found: '\" + string(var%d) + \"'\")\n", p, strings.Join(parts[:j], " "), cmp, j+1)
//			fmt.Fprintf(g.b, "    }\n")
//		} else if _, ok := g.funcExprs[p]; ok {
//			fmt.Fprintf(g.b, "    _, tok%d := p.peek()\n", j+1)
//			if setIncludes(usedVars, j+1) {
//				fmt.Fprintf(g.b, "    var%d, ok := p.%s(yylex)\n", j+1, p)
//				okDefined = true
//			} else if okDefined {
//				fmt.Fprintf(g.b, "    _, ok = p.%s(yylex)\n", p)
//			} else {
//				fmt.Fprintf(g.b, "    _, ok := p.%s(yylex)\n", p)
//				okDefined = true
//			}
//			fmt.Fprintf(g.b, "    if !ok {\n")
//			fmt.Fprintf(g.b, "      p.fail(\"expected: '%s: %s <%s>', found: '\"+string(tok%d)+\"'\")\n", p, strings.Join(parts[:j], " "), p, j+1)
//			fmt.Fprintf(g.b, "    }\n")
//		}
//	}
//	//success, return
//	fmt.Fprint(g.b, bb.String())
//}

var variableRe = regexp.MustCompile("(New\\w*\\()*\\$([1-6]+[0-9]*|[1-9])")

func normalizeBodyLine(r string) (string, int64, error) {
	r = strings.ReplaceAll(r, "%", "%%")
	r = strings.ReplaceAll(r, "$$ =", "ret =")
	r = strings.ReplaceAll(r, "return 1", "return ret, false")

	var variables int64
	r = strings.ReplaceAll(r, "$$", "ret")
	match := variableRe.FindAllStringSubmatchIndex(r, 1)
	for len(match) > 0 {
		m := match[0]
		start, end := m[4], m[5]

		_int64, err := strconv.ParseInt(r[start:end], 10, 64)
		if err != nil {
			return "", 0, fmt.Errorf("failed to parse variable string: %s", r[start:end])
		}
		if _int64 >= 64 {
			return "", 0, fmt.Errorf("variable reference too big: %d", _int64)
		}
		r = r[:start-1] + "var" + r[start:]
		variables |= (1 << _int64)
		match = variableRe.FindAllStringSubmatchIndex(r, 1)
	}

	//r = strings.ReplaceAll(r, "NewIntVal(var", "NewIntVal(tok")
	//r = strings.ReplaceAll(r, "NewStrVal(var", "NewStrVal(tok")
	//r = strings.ReplaceAll(r, "NewValArg(var", "NewValArg(tok")
	//r = strings.ReplaceAll(r, "NewListArg(var", "NewListArg(tok")
	//r = strings.ReplaceAll(r, "NewBitVal(var", "NewBitVal(tok")
	//r = strings.ReplaceAll(r, "NewHexVal(var", "NewHexVal(tok")

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

func isLit(p string) bool {
	switch p {
	case "STRING", "INTEGRAL", "VALUE_ARG", "LIST_ARG", "BIT_LITERAL", "HEX", "FLOAT":
		return true
	default:
		return false
	}
}
