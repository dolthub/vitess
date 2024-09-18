package rd

import (
	"bufio"
	"fmt"
	"io"
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
	rules []*rule
}

type rule struct {
	name  string
	start bool
	body  []string
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
		case "}":
			if nesting > 0 {
				nesting--
			} else if d != nil {
				d.rules = append(d.rules, r)
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
		} else if strings.HasPrefix(line, "%start") {
			yc.start = strings.Split(line, " ")[1]
		} else if strings.HasPrefix(line, "|") && r != nil {
			d.rules = append(d.rules, r)
			r = nil
		}

		if line[len(line)-1] == ':' {
			if d != nil {
				yc.defs = append(yc.defs, d)
			}
			d = new(def)
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
			r.name = parseRuleName(line)
		}
	}
	if d != nil {
		yc.defs = append(yc.defs, d)
	}
	return yc, nil
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
