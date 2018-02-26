package htmlLoader

import (
	"fmt"
	"strings"
)

type nodeType int32

func (n nodeType) String() string {
	return nodeToString[int32(n)]
}

var nodeToString = []string{
	"n_UNDEFINED",
	"n_TAG",
	"n_TAG_NAME",
	"n_TAG_ATTRIBUTES",
	"n_TAG_CHILDREN",
	"n_ATTRIBUTE",
	"n_ATTRIBUTE_VALUE",
	"n_TEXT_CHUNK",
	"n_DOCUMENT",
}

const (
	n_UNDEFINED nodeType = iota
	n_TAG
	n_TAG_NAME
	n_TAG_ATTRIBUTES
	n_TAG_CHILDREN
	n_ATTRIBUTE
	n_ATTRIBUTE_VALUE
	n_TEXT_CHUNK
	n_DOCUMENT
)

type ast struct {
	t        nodeType
	value    string
	children []ast
}

func (a ast) String() string {
	result := "{" + fmt.Sprint(a.t)
	if a.value != "" {
		result += ", " + a.value
	}
	if a.children != nil {
		result += ", " + fmt.Sprint(a.children)
	}
	result += "}"
	return result
}

func strEquals(a, b string) bool {
	if strings.ToLower(a) == strings.ToLower(b) {
		return true
	}
	return false
}

func newNode(t nodeType, value string, children ...ast) ast {
	return ast{t, value, children}
}

type parser struct {
	tokens   []token
	i        int
	lastNode ast
}

func (p *parser) next() {
	p.i++
}

func (p *parser) nextWS() {
	for p.tokens[p.i].tType == t_WHITESPACE {
		p.i++
	}
	p.i++
}

func (p *parser) accept(tType tokenType) bool {
	if p.tokens[p.i].tType == tType {
		p.next()
		return true
	}
	return false
}

func (p *parser) getNextWS(num int) token {
	i := p.i
	for iter := 0; iter < num; iter++ {
		for p.tokens[i].tType == t_WHITESPACE {
			i++
		}
		i++
	}
	return p.tokens[i-1]
}

func (p *parser) acceptWS(tType tokenType) bool {
	for p.tokens[p.i].tType == t_WHITESPACE {
		p.i++
	}
	return p.accept(tType)
}

func (p *parser) expect(tType tokenType) token {
	if !p.accept(tType) {
		panic(p.tokens[p.i])
	}
	return p.tokens[p.i-1]
}

func (p *parser) expectWS(tType tokenType) token {
	for p.tokens[p.i].tType == t_WHITESPACE {
		p.i++
	}
	return p.expect(tType)
}

func parseTokens(tokens []token) ast {
	p := parser{tokens: tokens}
	return document(&p)
}

func document(p *parser) ast {
	children := []ast{}
	for !p.accept(t_END_OF_INPUT) {
		if p.tokens[p.i].tType == t_ANGLE_BRACKET_LEFT {
			children = append(children, tag(p))
		} else {
			children = append(children, textChunk(p))
		}
	}

	return newNode(n_DOCUMENT, "", children...)
}

var voidTags = map[string]bool{
	"area":     true,
	"base":     true,
	"br":       true,
	"col":      true,
	"embed":    true,
	"hr":       true,
	"img":      true,
	"input":    true,
	"keygen":   true,
	"link":     true,
	"menuitem": true,
	"meta":     true,
	"param":    true,
	"source":   true,
	"track":    true,
	"wbr":      true,
	"!doctype": true,
}

func tag(p *parser) ast {
	p.expect(t_ANGLE_BRACKET_LEFT)
	name := strings.ToLower(p.expectWS(t_WORD).lexeme)

	attrs := []ast{}
	for {
		if p.acceptWS(t_ANGLE_BRACKET_RIGHT) {
			if voidTags[name] {
				return newNode(
					n_TAG, name,
					newNode(n_TAG_ATTRIBUTES, "", attrs...),
					newNode(n_TAG_CHILDREN, ""),
				)
			}
			break
		}
		if p.acceptWS(t_SLASH) {
			p.expectWS(t_ANGLE_BRACKET_RIGHT)
			return newNode(
				n_TAG, name,
				newNode(n_TAG_ATTRIBUTES, "", attrs...),
				newNode(n_TAG_CHILDREN, ""),
			)
		}
		attrs = append(attrs, attribute(p))
	}

	children := []ast{}

	for {
		if p.getNextWS(1).tType == t_ANGLE_BRACKET_LEFT {
			if p.getNextWS(2).tType == t_SLASH {
				p.nextWS()
				p.nextWS()
				closeName := p.expectWS(t_WORD).lexeme
				if !strEquals(name, closeName) {
					panic("Closing tag name does not match!")
				}
				p.expectWS(t_ANGLE_BRACKET_RIGHT)
				break
			}
			child := tag(p)
			children = append(children, child)
		} else {
			children = append(children, textChunk(p))
		}
	}

	return newNode(
		n_TAG, name,
		newNode(n_TAG_ATTRIBUTES, "", attrs...),
		newNode(n_TAG_CHILDREN, "", children...),
	)
}

func textChunk(p *parser) ast {
	value := ""
	for p.tokens[p.i].tType != t_ANGLE_BRACKET_LEFT {
		value += p.tokens[p.i].lexeme
		p.next()
	}
	return newNode(n_TEXT_CHUNK, value)
}

func attribute(p *parser) ast {
	name := p.expectWS(t_WORD).lexeme

	var value ast
	if p.acceptWS(t_ASSIGN) {
		value = newNode(n_ATTRIBUTE_VALUE, p.expectWS(t_STRING).lexeme)
	}

	return newNode(n_ATTRIBUTE, name, value)
}
