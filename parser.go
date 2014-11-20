package gohaml

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

type hamlParser struct {
}

func (self *hamlParser) parse(input string) (output *tree, err error) {
	output = newTree()

	var currentNode inode
	var node inode
	lastSpaceChar := '\000'
	line := 1
	j := 0
	for i, r := range input {
		if r == '\n' {
			line += 1
			node, err, lastSpaceChar = parseLeadingSpace(input[j:i], lastSpaceChar, line)
			if err != nil {
				return
			}
			if node != nil && !node.nil() {
				putNodeInPlace(currentNode, node, output)
				currentNode = node
			}
			j = i + 1
		}
	}
	node, err, lastSpaceChar = parseLeadingSpace(input[j:], lastSpaceChar, line)
	if err != nil {
		return
	}
	if node != nil && !node.nil() {
		putNodeInPlace(currentNode, node, output)
	}
	return
}

func putNodeInPlace(cn inode, node inode, t *tree) {
	if node == nil || node.nil() {
		return
	}
	if cn == nil || cn.nil() {
		//t.nodes.Push(node)
		t.nodes = append(t.nodes, node)
	} else if node.indentLevel() < cn.indentLevel() {
		for cn = cn.parent(); cn != nil && node.indentLevel() < cn.indentLevel(); cn = cn.parent() {
		}
		putNodeInPlace(cn, node, t)
	} else if node.indentLevel() == cn.indentLevel() && cn.parent() != nil {
		cn.parent().addChild(node)
	} else if node.indentLevel() == cn.indentLevel() {
		//t.nodes.Push(node)
		t.nodes = append(t.nodes, node)
	} else if node.indentLevel() > cn.indentLevel() {
		cn.addChild(node)
	}
}

var parser hamlParser

func parseLeadingSpace(input string, lastSpaceChar rune, line int) (output inode, err error, spaceChar rune) {
	node := new(node)
	for i, r := range input {
		switch {
		case r == '-':
			output = parseCode(input[i+1:], node, line)
		case r == '%':
			output, err = parseTag(input[i+1:], node, true, line)
		case r == '#':
			output, err = parseId(input[i+1:], node, line)
		case r == '.':
			output, err = parseClass(input[i+1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i+1:]), node, line)
		case r == '\\':
			output = parseRemainder(input[i+1:], node, line)
		case !unicode.IsSpace(r):
			output = parseRemainder(input[i:], node, line)
		case unicode.IsSpace(r):
			if lastSpaceChar > 0 && r != lastSpaceChar {
				from, to := "space", "tab"
				if lastSpaceChar == 9 {
					from = "tab"
					to = "space"
				}
				msg := fmt.Sprintf("Syntax error on line %d: Inconsistent spacing in document changed from %s to %s characters.\n", line, from, to)
				//err = os.NewError(msg)
				err = errors.New(msg)
			} else {
				lastSpaceChar = r
			}
		}
		if nil != err {
			break
		}
		if nil != output {
			output.setIndentLevel(i)
			break
		}
	}
	spaceChar = lastSpaceChar
	return
}

func parseKey(input string, n *node, line int) (output inode) {
	if input[len(input)-1] == '<' {
		n = parseNoNewline("", n, line)
		n.setRemainder(input[0:len(input)-1], true)
	} else {
		n.setRemainder(input, true)
		output = n
	}
	output = n
	return
}

func parseTag(input string, node *node, newTag bool, line int) (output inode, err error) {
	if 0 == len(input) && newTag {
		//err = os.NewError(fmt.Sprintf("Syntax error on line %d: Invalid tag: %s.\n", line, input))
		err = errors.New(fmt.Sprintf("Syntax error on line %d: Invalid tag: %s.\n", line, input))
		return
	}
	for i, r := range input {
		switch {
		case r == '.':
			output, err = parseClass(input[i+1:], node, line)
		case r == '#':
			output, err = parseId(input[i+1:], node, line)
		case r == '{':
			output, err = parseAttributes(tl(input[i+1:]), node, line)
		case r == '<':
			output = parseNoNewline(input[i+1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i+1:]), node, line)
		case r == '/':
			output = parseAutoclose("", node, line)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i+1:], node, line)
		}
		if nil != err {
			break
		}
		if nil != output {
			node._name = input[0:i]
			break
		}
	}
	if nil == output {
		node._name = input
		output = node
	}
	return
}

func parseAutoclose(input string, node *node, line int) (output inode) {
	node._autoclose = true
	output = node
	return
}

func parseAttributes(input string, node *node, line int) (output inode, err error) {
	const (
		t_end = iota
		t_ident
		t_string
		t_ws
		t_close
		t_sym
		t_comma
		t_rocket
		t_err
	)

	err = nil
	pos := 0

	errout := func(s string) {
		err = fmt.Errorf("parse error line %d character %d, <<%s>>", line, pos, s)
		panic(err)
	}

	lex := func() (int, string) {
		idset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-."
		start := pos

		if pos < len(input) {
			if strings.ContainsAny(string(input[pos]), " \t") {
				pos += 1
				for pos < len(input) && strings.ContainsAny(string(input[pos]), " \t") {
					pos += 1
				}

				return t_ws, input[start:pos]
			}

			if input[pos] == '"' {
				pos += 1
				for pos < len(input) && input[pos] != '"'  { 
					pos += 1
				}
				pos += 1

				return t_string, input[start:pos]
			}
			
			if strings.ContainsAny(string(input[pos]), idset) {
				pos += 1
				for pos < len(input) && strings.ContainsAny(string(input[pos]), idset) {
					pos += 1
				}

				return t_ident, input[start:pos]
			}

			switch input[pos] {
			case '}':
				pos += 1
				return t_close, "}"
			case ':':
				pos += 1
				return t_sym, ":"
			case ',':
				pos += 1
				return t_comma, ","
			case '=':
				pos += 1
				if pos < len(input) && input[pos] == '>' {
					pos += 1
					return t_rocket, "=>"
				}
				fallthrough
			default:
				errout(string(input[pos]))
				return t_err, input[start:pos]
			}

		} else {
			return t_end, ""
		}
	}
		
	parseKey := func() string {
		c, s := lex()
		if c == t_ws {
			// skip leading
			c, s = lex()
		}

		if c != t_sym {
			errout(s)
			return ""
		}

		// i am dubious about the value of the lookup-key feature
		c, s = lex()		
		if c == t_string {
			return fmt.Sprintf(":%s", s[1:len(s)-1])
		} else if c == t_ident {
			return fmt.Sprintf(":%s", s)
		}
		
		errout(s)
		return ""
	}
	
	parseValue := func() string {
		c, s := lex()
		if c == t_ws {
			// skip leading
			c, s = lex()
		}

		if c == t_string {
			return s
		} else if c == t_ident {
			return s
		}
		
		errout(s)
		return ""
	}

	parseAttr := func() (key, value string) {
		k := parseKey()
		if err != nil {
			return
		}

		c, s := lex()
		if c == t_ws {
			c, s = lex()
		}
		if c == t_rocket {
			value := parseValue()
			return k, value
		} else {
			errout(s)
			return
		}
	}

	parseAttrs := func() {
		for {
   		key, value := parseAttr()
   		if err != nil {
   			return
   		} else {
				node.addAttr(key, value)
			}

   		c, s := lex()
			if c == t_ws {
				c, s = lex()
			}
			
   		if c == t_comma {
				continue
			} else if c == t_close {
				break
			} else {
				errout(s)
				return
			}
		}
	}

	parseAttrs()

	output, err = parseTag(input[pos:], node, false, line)
	return
}

func oldBrokenParseAttributes(input string, node *node, line int) (output inode, err error) {
	inKey := true
	inRocket := false

	keyEnd, attrStart := 0, 0
	for i, r := range input {
		if inKey && (r == '=' || unicode.IsSpace(r)) {
			inKey = false
			inRocket = true
			keyEnd = i
		} else if inRocket && r != '>' && r != '=' && r != '}' && !unicode.IsSpace(r) {
			inRocket = false
			attrStart = i
		} else if r == ',' {
			node.addAttr(t(input[0:keyEnd]), t(input[attrStart:i]))
			output, err = parseAttributes(tl(input[i+1:]), node, line)
			break
		} else if r == '}' {
			if attrStart == 0 {
				msg := fmt.Sprintf("Syntax error on line %d: Attribute requires a value.\n", line)
				//err = os.NewError(msg)
				err = errors.New(msg)
				return
			}
			if inKey {
				msg := fmt.Sprintf("Syntax error on line %d: Attribute requires a rocket and value.\n", line)
				err = errors.New(msg)
				return
			}
			attrValue := t(input[attrStart:i])
			node.addAttr(input[0:keyEnd], attrValue)
			output, _ = parseTag(input[i+1:], node, false, line)
			break
		}
	}
	if nil == output {
		msg := fmt.Sprintf("Syntax error on line %d: Attributes must have closing '}'.\n", line)
		err = errors.New(msg)
	}
	return
}

func parseId(input string, node *node, line int) (output inode, err error) {
	defer func() {
		if nil == output {
			msg := fmt.Sprintf("Syntax error on line %d: Illegal element: classes and ids must have values.\n", line)
			err = errors.New(msg)
		}
	}()
	if len(input) == 0 {
		return
	}
	for i, r := range input {
		if r == '.' || r == '=' || r == '{' || unicode.IsSpace(r) {
			if i == 0 {
				return
			}
			node.addAttrNoLookup("id", input[0:i])
		}
		switch {
		case r == '.':
			output, _ = parseClass(input[i+1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i+1:]), node, line)
		case r == '{':
			output, err = parseAttributes(tl(input[i+1:]), node, line)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i+1:], node, line)
		}
		if nil != output {
			break
		}
	}
	if nil == output {
		output = node
		node.addAttrNoLookup("id", input)
	}
	return
}

func parseClass(input string, node *node, line int) (output inode, err error) {
	defer func() {
		if nil == output {
			msg := fmt.Sprintf("Syntax error on line %d: Illegal element: classes and ids must have values.\n", line)
			err = errors.New(msg)
		}
	}()
	if len(input) == 0 {
		return
	}
	for i, r := range input {
		if r == '{' || r == '.' || r == '=' || unicode.IsSpace(r) {
			if i == 0 {
				return
			}
			node.addAttrNoLookup("class", input[0:i])
		}
		switch {
		case r == '{':
			output, err = parseAttributes(tl(input[i+1:]), node, line)
		case r == '.':
			output, err = parseClass(input[i+1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i+1:]), node, line)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i+1:], node, line)
		}
		if nil != output {
			break
		}
	}
	if nil == output {
		node.addAttrNoLookup("class", input)
		output = node
	}
	return
}

func parseRemainder(input string, node *node, line int) (output inode) {
	if input[len(input)-1] == '<' {
		node = parseNoNewline("", node, line)
		node._remainder.value = input[0 : len(input)-1]
		node._remainder.needsResolution = false
	} else {
		node._remainder.value = input
		node._remainder.needsResolution = false
	}
	output = node
	return
}

func parseNoNewline(input string, node *node, line int) (output *node) {
	node.setNoNewline(true)
	output = node
	return
}

func t(input string) (output string) {
	output = strings.Trim(input, " 	")
	return
}

func tl(input string) (output string) {
	output = strings.TrimLeft(input, " 	")
	return
}

func parseCode(input string, node inode, line int) (output inode) {
	l.init(strings.NewReader(input))

	success := yyParse(l)
	if success != 0 {
		fmt.Fprintf(os.Stderr, "Did not recognize %s", input)
	}
	output = Output
	return
}

// var s scanner.Scanner

type Lexer struct {
	s *scanner.Scanner
}

var l = &Lexer{new(scanner.Scanner)}

func (l *Lexer) init(reader *strings.Reader) {
	l.s.Init(reader)
}

func (l *Lexer) Lex(v *yySymType) (output int) {
	i := l.s.Scan()

	switch i {
	case scanner.Ident:
		switch l.s.TokenText() {
		case "for":
			output = FOR
		case "range":
			output = RANGE
		default:
			output = IDENT
		}
		v.s = l.s.TokenText()
	case scanner.String, scanner.RawString:
		output = ATOM
		text := l.s.TokenText()
		v.i = text[1 : len(text)-1]
	case scanner.Int:
		output = ATOM
		v.i, _ = strconv.Atoi(l.s.TokenText())
	case scanner.Float:
		output = ATOM
		v.i, _ = strconv.ParseFloat(l.s.TokenText(), 64)
	case scanner.EOF:
		output = 0
	default:
		output = int(i)
	}
	return
}

func (l *Lexer) Error(e string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", e)
}
