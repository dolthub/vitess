package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

//go:generate go run ./main.go ../testdata/expr_functions

func main() {
	// split file into definitions
	// each definition gets a symbol
	// parse each definition into a function
	args := os.Args[1:]
	//outfile, err := os.CreateTemp("", "")
	//	//if err != nil {
	//	//	log.Fatal(err)
	//	//}
	outfile, err := os.OpenFile("../../sql.gen.go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("writing to: ", outfile.Name())

	for _, inp := range args {
		bytes, err := os.ReadFile(inp)
		if err != nil {
			log.Fatalf("invalid filename: %s; err=%s", inp, err)
		}
		p := &pgen{lines: strings.Split(string(bytes), "\n")}
		out, err := p.gen()
		if err != nil {
			log.Fatalf("invalid contents: %s; err=%s", inp, err)
		}
		n, err := outfile.WriteString(out)
		if err != nil || n != len(out) {
			log.Fatalf("failed to write results: %s; err=%s", inp, err)
		}
	}
	outfile.Sync()
}

type pgen struct {
	lines []string
	i     int
}

func (p *pgen) nextLine() (string, bool) {
	if p.i >= len(p.lines) {
		return "", false
	}
	l := p.lines[p.i]
	p.i++
	l = strings.TrimSpace(l)
	if len(l) > 2 && l[:2] == "| " {
		l = l[2:]
	}
	if len(l) > 2 && l[:2] == "//" {
		return p.nextLine()
	}

	if len(l) == 0 {
		return p.nextLine()
	}
	return l, true
}

func (p *pgen) gen() (string, error) {
	b := strings.Builder{}
	fmt.Fprintf(&b, "package sqlparser\n\n")

	//l, ok := p.nextLine()
	//if !ok {
	//	return "", fmt.Errorf("function id invalid")
	//}
	//fid := strings.Split(l, ":")
	//if len(fid) < 1 {
	//	return "", fmt.Errorf("function id invalid")
	//}
	//fmt.Fprintf(&b, "func (p *parser) %s() (Expr, bool) {\n", camelize(fid[0]))
	//fmt.Fprintf(&b, "  id, tok := p.peek()\n")
	topIfElse := "  if"
	currentFname := ""
	for {
		def, ok := p.nextLine()
		if !ok {
			break
		}

		if def[len(def)-1] == ':' {
			fid := strings.Split(def, ":")
			if len(fid) < 1 {
				return "", fmt.Errorf("function id invalid")
			}
			if p.i > 1 {
				fmt.Fprintf(&b, "\n  return nil, false\n}\n")
			}
			currentFname = fid[0]
			fmt.Fprintf(&b, "func (p *parser) %s() (Expr, bool) {\n", camelize(fid[0]))
			fmt.Fprintf(&b, "  id, tok := p.peek()\n")
			topIfElse = "  if"
			continue
		}

		open, ok := p.nextLine()
		if !ok || open != "{" {
			return "", fmt.Errorf("missing openb, line %d, %s->%s", p.i, def, open)
		}
		ret, ok := p.nextLine()
		if !ok {
			return "", fmt.Errorf("missing return type")
		}
		close, ok := p.nextLine()
		if !ok || close != "}" {
			return "", fmt.Errorf("missing closeb")
		}

		// progress string for fail message
		// differentiate end token vs function
		// extend reading tokens into variables
		// return
		parts := strings.Split(def, " ")
		for i, p := range parts {
			var cmp string
			if p == "openb" {
				cmp = "'('"
			} else if p == "closeb" {
				cmp = "')'"
			} else if strings.ToLower(p) == p {
			} else {
				cmp = p
			}

			if i == 0 {
				if cmp != "" {
					fmt.Fprintf(&b, "%s id == %s {\n", topIfElse, cmp)
					fmt.Fprintf(&b, "    var1, _ := p.next()\n")
				} else {
					fname := camelize(p)
					fmt.Fprintf(&b, "%s var1, ok := p.%s(); ok {\n", topIfElse, fname)
				}
				topIfElse = " else if"
				continue
			}
			if cmp != "" {
				fmt.Fprintf(&b, "    var%d, tok := p.next()\n", i+1)
				fmt.Fprintf(&b, "    if var%d != %s {\n", i+1, cmp)
				fmt.Fprintf(&b, "      p.fail(\"expected: '%s: %s <%s>', found: '\" + string(tok) + \"'\")\n", currentFname, strings.Join(parts[:i], " "), cmp)
				fmt.Fprintf(&b, "    }\n")
			} else if strings.ToLower(p) == p {
				fname := camelize(p)
				fmt.Fprintf(&b, "    var%d, ok := p.%s()\n", i+1, fname)
				fmt.Fprintf(&b, "    if !ok {\n")
				fmt.Fprintf(&b, "      p.fail(\"expected: '%s: %s <%s>', found: 'string(tok)'\")\n", currentFname, strings.Join(parts[:i], " "), p)
				fmt.Fprintf(&b, "    }\n")
			}
		}
		ret = strings.Replace(ret, "$$ = ", "", -1)
		if ret[0] == '&' {
			ret = ret[1:]
		}
		ret = strings.Replace(ret, "$", "var", -1)
		fmt.Fprintf(&b, "    return &%s, true\n", ret)
		fmt.Fprintf(&b, "  }")
	}
	fmt.Fprintf(&b, "\n  return nil, false\n}\n")

	return b.String(), nil
}

func camelize(name string) string {
	words := strings.Split(name, "_")
	key := strings.ToLower(words[0])
	for _, word := range words[1:] {
		key += strings.Title(word)
	}
	return key
}
