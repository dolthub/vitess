package rd

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

// split yacc file into individual functions
// remove comments
// split

type yaccFileContents struct {
	prefix    []string
	goTypes   []goType
	yaccTypes []yaccType
	tokens    []yaccToken
	start     string
	defs      map[string]*def
}

type goType struct {
	name string
	typ  string
}

type yaccType struct {
	name string
	typ  string
}

type yaccToken struct {
	name    string
	yaccTyp string
}

type def struct {
	name  string
	rules *rulePrefix
	done  bool
}

func newDef() *def {
	d := new(def)
	d.rules = new(rulePrefix)
	return d
}

type rule struct {
	name     string
	fields   []string
	set      bool
	body     []string
	usedVars int64
}

type rulePrefix struct {
	prefix   string
	flat     []*rule
	term     []*rule
	pref     []*rulePrefix
	rec      *rulePrefix
	empty    *rule
	usedVars int64
	done     bool
}

func split(infile io.Reader) (yc *yaccFileContents, err error) {
	sc := bufio.NewScanner(infile)
	var buf []string
	var acc bool
	var d *def
	var r *rule
	var nesting int
	yc = new(yaccFileContents)
	yc.defs = make(map[string]*def)
	defer func() {
		if err != nil {
			return
		}
		err = yc.finalize()
	}()
	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line == "/*empty*/" {
			print()
		}

		line = strings.ReplaceAll(line, "%prec", "")

		switch line {
		case "{":
			if r == nil {
				r = new(rule)
			}
			r.set = true

			continue
		case "%{":
			acc = true
			continue
		case "%}":
			yc.prefix = buf
			buf = nil
			acc = false
			continue
		case "%union {":
			acc = true
			continue
		case "},", "})":
			nesting--
		case "}":
			if nesting > 0 {
				nesting--
			} else if d != nil {
				d.rules.addRule(r)
				r = nil
				continue
			} else {
				yc.goTypes, err = parseGoTypes(buf)
				if err != nil {
					return nil, err
				}
				buf = nil
				acc = false
				continue
			}
		case "%%":
			continue
		case "{}":
			if d != nil {
				r := new(rule)
				r.set = true
				r.body = []string{"$$ = nil"}
				d.rules.empty = r
			}
			continue
		}

		if line[0] == '{' && line[len(line)-1] == '}' {
			if r == nil {
				r = new(rule)
			}
			r.set = true
			r.body = append(r.body, strings.TrimSpace(line[1:len(line)-1]))
			d.rules.addRule(r)
			r = nil
			continue
		}

		if acc {
			if line[len(line)-1] == '{' {
				nesting++
			}
			buf = append(buf, line)
		}

		if strings.HasPrefix(line, "%left") {
		} else if strings.HasPrefix(line, "%type") {
			yaccTyp, err := parseYaccType(line)
			if err != nil {
				return nil, err
			}
			yc.yaccTypes = append(yc.yaccTypes, yaccTyp...)
		} else if strings.HasPrefix(line, "%token") {
			//yc.tokens = append(yc.tokens, line)
			continue
		} else if strings.HasPrefix(line, "%right") {

		} else if strings.HasPrefix(line, "%nonassoc") {
		} else if strings.HasPrefix(line, "//") {
			continue
		} else if strings.HasPrefix(line, "%start") {
			yc.start = strings.Split(line, " ")[1]
		} else if strings.HasPrefix(line, "|") && r != nil {
			d.rules.addRule(r)
			r = nil
		} else if strings.HasPrefix(line, "} else") && r != nil {
			nesting--
		}

		if line[len(line)-1] == ':' && !strings.HasPrefix(line, "case ") {
			if r != nil {
				d.rules.addRule(r)
				r = nil
			}
			if d != nil {
				yc.defs[d.name] = d
			}
			d = newDef()
			d.name = line[:len(line)-1]
			continue
		}
		if r != nil {
			if line[len(line)-1] == '{' {
				nesting++
			}
			r.body = append(r.body, line)
		} else if d != nil {
			r = new(rule)
			r.set = true
			if strings.HasSuffix(line, "{}") {
				line = line[:len(line)-2]
				r.name = parseRuleName(line)
				d.rules.addRule(r)
				r = nil
				continue
			}
			r.name = parseRuleName(line)
		}
	}
	if r != nil {
		d.rules.addRule(r)
	}
	if d != nil {
		yc.defs[d.name] = d
	}
	return yc, nil
}

func (yc *yaccFileContents) finalize() error {
	for _, d := range yc.defs {
		if err := d.rules.finalize(d.name, yc.defs); err != nil {
			return err
		}
	}
	return nil
}

func (r *rule) calcUsed() error {
	if len(r.body) == 0 {
		r.usedVars |= 1 << 1
	}
	for i, b := range r.body {
		newB, used, err := normalizeBodyLine(b)
		if err != nil {
			return err
		}
		r.usedVars |= used
		r.body[i] = newB
	}
	return nil
}

func (d *rulePrefix) calcUsed() {
	for _, t := range d.term {
		d.usedVars |= t.usedVars
	}
}

func (d *rulePrefix) addRule(r *rule) {
	d.flat = append(d.flat, r)
}

func (d *rulePrefix) recurseChildren(name string, defs map[string]*def) error {
	// expand/partition children first
	for _, r := range d.flat {
		for _, f := range r.fields {
			if d, ok := defs[f]; ok {
				if err := d.rules.finalize(f, defs); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *rulePrefix) expandConflicts(name string, defs map[string]*def) {
	// fixed point iteration, weasel out conflicts for this prefix level
	// look at first field, store which we've seen

	// token in r, conflicting starting token in r2
	// two different r2's have conflicting starting token
	nextTokens := make(map[string]bool)
	for _, r := range d.flat {
		// toplevel next tokens
		if len(r.fields) == 0 {
			continue
		}
		nextTokens[r.fields[0]] = true
	}
	ignore := make(map[string]bool)
	for {
		var conflicts []*rule
		for _, r := range d.flat {
			// find nested conflicts
			if d, ok := defs[r.name]; ok {
				// next token is a function, go one deeper
				for _, r2 := range d.rules.flat {
					// todo flat needs all fields
					if len(r2.fields) == 0 {
						continue
					}
					if nextTokens[r2.fields[0]] && !ignore[r2.fields[0]] {
						// lookahead conflict
						// advance rule
						conflicts = append(conflicts, r2)
					}
				}
				for _, r2 := range d.rules.flat {
					// catch inter-child-rule conflicts
					// this takes two cycles to unwind
					if len(r2.fields) == 0 {
						continue
					}
					nextTokens[r2.fields[0]] = true
				}
				nextTokens[d.name] = true
			}
		}
		if len(conflicts) == 0 {
			break
		}
		// remove
		for _, r := range conflicts {
			// prevent looping
			ignore[r.fields[0]] = true
		}
		d.flat = append(d.flat, conflicts...)
		conflicts = conflicts[:0]
	}
	return
}

func (d *rulePrefix) finalize(name string, defs map[string]*def) error {
	if d.done {
		return nil
	}
	d.done = true

	for _, r := range d.flat {
		if err := r.calcUsed(); err != nil {
			return err
		}
		r.fields = strings.Fields(r.name)
	}
	d.calcUsed()

	if err := d.recurseChildren(name, defs); err != nil {
		return err
	}

	d.expandConflicts(name, defs)

	if name == "expression" {
		print()
	}

	return d.partition(name)
}

func (r *rule) copy() *rule {
	nr := new(rule)
	nr.fields = make([]string, len(r.fields))
	nr.body = make([]string, len(r.body))
	copy(nr.fields, r.fields)
	copy(nr.body, r.body)
	nr.set = true
	nr.body = r.body
	nr.usedVars = r.usedVars
	nr.name = r.name
	return nr
}

func (d *rulePrefix) addPrefixPartition(name string, rules []*rule) error {
	newRules := make([]*rule, len(rules))
	for i, r := range rules {
		newRules[i] = r.copy()
	}
	p := new(rulePrefix)
	p.prefix = rules[0].fields[0]
	p.flat = newRules
	for _, r := range p.flat {
		r.fields = r.fields[1:]
	}
	if err := p.partition(""); err != nil {
		return err
	}
	if p.prefix == name {
		d.rec = p
	} else {
		d.pref = append(d.pref, p)
	}
	return nil
}

func (d *rulePrefix) partition(name string) error {
	d.calcUsed()

	if len(d.flat) == 1 {
		r := d.flat[0]
		if len(r.fields) == 0 || len(r.fields) == 1 && r.fields[0] == "/*empty*/" {
			r.fields = nil
			d.empty = r
			return nil
		}
	}

	sort.Slice(d.flat, func(i, j int) bool {
		return d.flat[i].name < d.flat[j].name
	})

	var acc []*rule
	for _, r := range d.flat {
		if len(r.fields) == 1 && r.fields[0] == "/*empty*/" {
			r.fields = nil
		}
		if len(r.fields) == 0 {
			// empty rule is special, checked last
			d.empty = r
			continue
		}

		if len(acc) == 0 {
			acc = append(acc, r)
			continue
		}

		match := acc[0].fields[0] == r.fields[0]
		if !match {
			if len(acc) < 2 {
				d.term = append(d.term, acc[0])
			} else {
				if err := d.addPrefixPartition(name, acc); err != nil {
					return err
				}
			}
			acc = acc[:0]
		}
		acc = append(acc, r)
	}
	if len(acc) == 1 {
		d.term = append(d.term, acc[0])
	} else if len(acc) > 0 {
		if err := d.addPrefixPartition(name, acc); err != nil {
			return err
		}
	}

	return nil
}

func parseRuleName(s string) string {
	s = strings.TrimSpace(s)
	if s[0] == '|' {
		s = s[1:]
	}
	s = strings.TrimSpace(s)
	return s
}

func parseGoTypes(in []string) ([]goType, error) {
	var ret []goType
	for _, typ := range in {
		parts := strings.Fields(typ)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid go type: %s", typ)
		}
		ret = append(ret, goType{name: parts[0], typ: parts[1]})
	}
	return ret, nil
}

func parseYaccType(in string) ([]yaccType, error) {
	parts := strings.Split(in, " ")
	if parts[0] != "%type" || len(parts) < 3 {
		return nil, fmt.Errorf("invalid yacc type: %s", in)
	}
	typ := parts[1]
	if typ[0] != '<' || typ[len(typ)-1] != '>' {
		return nil, fmt.Errorf("invalid yacc type: %s", in)
	}
	typ = typ[1 : len(typ)-1]
	var ret []yaccType
	for _, p := range parts[2:] {
		ret = append(ret, yaccType{name: p, typ: typ})
	}
	return ret, nil
}
