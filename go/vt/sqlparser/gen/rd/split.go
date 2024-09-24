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
	defs      []*def
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
}

func newDef() *def {
	d := new(def)
	d.rules = new(rulePrefix)
	return d
}

type rule struct {
	name     string
	fields   []string
	start    bool
	body     []string
	usedVars int64
}

type rulePrefix struct {
	prefix   string
	term     []*rule
	pref     []*rulePrefix
	rec      *rulePrefix
	empty    *rule
	usedVars int64
}

func split(infile io.Reader) (*yaccFileContents, error) {
	sc := bufio.NewScanner(infile)
	var buf []string
	var acc bool
	var d *def
	var r *rule
	var err error
	var nesting int
	yc := new(yaccFileContents)
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
			r.start = true

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
				d.rules.term = append(d.rules.term, r)
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
			continue
		}

		if line[0] == '{' && line[len(line)-1] == '}' {
			if r == nil {
				r = new(rule)
			}
			r.body = append(r.body, strings.TrimSpace(line[1:len(line)-1]))
			d.rules.term = append(d.rules.term, r)
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
			d.rules.term = append(d.rules.term, r)
			r = nil
		} else if strings.HasPrefix(line, "} else") && r != nil {
			nesting--
		}

		if line[len(line)-1] == ':' && !strings.HasPrefix(line, "case ") {
			if r != nil {
				d.rules.term = append(d.rules.term, r)
				r = nil
			}
			if d != nil {
				yc.defs = append(yc.defs, d)
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
			if strings.HasSuffix(line, "{}") {
				line = line[:len(line)-2]
				r.name = parseRuleName(line)
				d.rules.term = append(d.rules.term, r)
				r = nil
				continue
			}
			r.name = parseRuleName(line)
		}
	}
	if r != nil {
		d.rules.term = append(d.rules.term, r)
	}
	if d != nil {
		if err := d.finalize(); err != nil {
			return nil, err
		}
		yc.defs = append(yc.defs, d)
	}
	return yc, nil
}

func (d *def) finalize() error {
	sort.Slice(d.rules.term, func(i, j int) bool {
		return d.rules.term[i].name < d.rules.term[j].name
	})
	for _, r := range d.rules.term {
		if err := r.calcUsed(); err != nil {
			return err
		}
		r.fields = strings.Fields(r.name)
	}
	d.rules.partition(d.name)
	return nil
}

func (r *rule) calcUsed() error {
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

func (d *rulePrefix) partition(name string) {
	d.calcUsed()
	j := 0
	for i := 0; i < len(d.term); i++ {
		if len(d.term[i].fields) == 0 {
			// empty rule is special, checked last
			d.empty = d.term[i]
			d.term = append(d.term[:i], d.term[i+1:]...)
			i--
			continue
		}
		if d.term[i].fields[0] != d.term[j].fields[0] {
			if i-j > 1 {
				// leading field doesn't match, but sequence did
				// new partition for the sequence
				// recursively partition the subsequence
				p := new(rulePrefix)
				p.prefix = d.term[j].fields[0]
				p.term = make([]*rule, i-j)
				copy(p.term, d.term[j:i])
				d.term = append(d.term[:j], d.term[i:]...)
				for _, r := range p.term {
					r.fields = r.fields[1:]
				}
				i = j
				p.partition(name)
				if p.prefix == name {
					d.rec = p
				} else {
					d.pref = append(d.pref, p)
				}
			}
		}
	}
	if d.term[j-1].fields[0] == d.term[j].fields[0] && len(d.term)-j > 1 {
		// leading field doesn't match, but sequence did
		// new partition for the sequence
		// recursively partition the subsequence
		p := new(rulePrefix)
		p.prefix = d.term[j].fields[0]
		p.term = make([]*rule, len(d.term)-j)
		copy(p.term, d.term[j:len(d.term)])
		d.term = append(d.term[:j], d.term[len(d.term):]...)
		for _, r := range p.term {
			r.fields = r.fields[1:]
		}
		p.partition(name)
		if p.prefix == name {
			d.rec = p
		} else {
			d.pref = append(d.pref, p)
		}
	}
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
