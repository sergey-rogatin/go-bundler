package jsLoader

import (
	"fmt"
)

/* TODO:

return, yield and other constructions that do not permit newline

check what keywords are allowed to be variable names and object property names

static analysis:
	replace imported vars in whole module
	tree-shaking ? read more about it


test invalid syntax catching
*/

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

func newNode(t nodeType, value string, children ...ast) ast {
	return ast{t, value, children}
}

var INVALID_NODE = ast{t: n_INVALID}

const (
	f_ACCEPT_NO_FUNCTION_CALL = 1 << 0
	f_ACCEPT_NO_IN            = 1 << 1
	f_GENERATOR_FUNCTION      = 1 << 2
	f_ASYNC_FUNCTION          = 1 << 3
)

type parser struct {
	tokens []token
	i      int

	flags    int
	lastNode ast
}

func (p *parser) testFlag(f int) bool {
	return (p.flags & f) != 0
}

func (p *parser) skip(skipNewline bool) int {
	i := p.i
	for {

		if p.tokens[i].tType == t_END_OF_INPUT {
			break
		} else if p.tokens[i].tType == t_NEWLINE && skipNewline {
			i++
		} else if p.tokens[i].tType == t_SPACE {
			i++
		} else if p.tokens[i].tType == t_DIV &&
			(p.tokens[i+1].tType == t_DIV || p.tokens[i+1].tType == t_DIV_ASSIGN) {
			for p.tokens[i].tType != t_NEWLINE && p.tokens[i].tType != t_END_OF_INPUT {
				i++
			}
		} else if p.tokens[i].tType == t_DIV &&
			(p.tokens[i+1].tType == t_MULT || p.tokens[i+1].tType == t_EXP) {
			i += 2
			for !((p.tokens[i].tType == t_MULT || p.tokens[i].tType == t_EXP) &&
				p.tokens[i+1].tType == t_DIV) {
				i++
			}
			i += 2
		} else {
			break
		}
	}

	return i
}

func (p *parser) getNextT() token {
	i := p.skip(true)
	return p.tokens[i]
}

func (p *parser) skipT() {
	p.i = p.skip(true)
	p.i++
}

func (p *parser) acceptT(tTypes ...tokenType) bool {
	tok := p.getNextT()

	for _, tType := range tTypes {
		if tok.tType == tType {
			p.skipT()
			return true
		}
	}

	return false
}

func (p *parser) checkASI(tTypes ...tokenType) {
	i := p.skip(false)

	for _, tType := range tTypes {
		if tType == t_SEMI {
			if p.tokens[i].tType == t_NEWLINE ||
				p.tokens[i].tType == t_CURLY_RIGHT {
				return
			}
			if p.tokens[i].tType == t_END_OF_INPUT {
				return
			}
		}
	}

	panic(&parsingError{p.i, p.tokens})
}

func (p *parser) expectT(tTypes ...tokenType) {
	if !p.acceptT(tTypes...) {
		p.i = p.skip(false)
		p.checkASI(tTypes...)
	}
}

type parseFunc func(*parser) ast

func (p *parser) acceptF(f parseFunc) bool {
	prevIndex := p.i

	n := f(p)
	if n.t == n_INVALID {
		p.i = prevIndex
		return false
	}

	p.lastNode = n
	return true
}
func (p *parser) expectF(f parseFunc) ast {
	if !p.acceptF(f) {
		p.checkASI(t_SEMI)
	}
	return p.getNode()
}

func (p *parser) getNode() ast {
	return p.lastNode
}

func (p *parser) getLexeme() string {
	return p.tokens[p.i-1].lexeme
}

type parsingError struct {
	index  int
	tokens []token
}

func (pe *parsingError) Error() string {
	tok := pe.tokens[pe.index]

	start := pe.index
	for start > 0 && pe.tokens[start-1].tType != t_NEWLINE {
		start--
	}

	end := pe.index
	for end < len(pe.tokens)-1 && pe.tokens[end].tType != t_NEWLINE {
		end++
	}

	resLine0 := fmt.Sprintf("Unexpected token '%s' at %v:%v\n", tok.lexeme, tok.line, tok.column)
	resStart := fmt.Sprintf("%04d", tok.line) + " "
	resLine1 := resStart
	for i := start; i < end; i++ {
		resLine1 += pe.tokens[i].lexeme
	}
	resLine1 += "\n"

	resLine2 := ""
	for i := int32(0); i <= tok.column+int32(len(resStart)); i++ {
		resLine1 += " "
	}
	resLine2 += "^"

	return resLine0 + resLine1 + resLine2
}

func parseTokens(tokens []token) (res ast, err *parsingError) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(*parsingError)
			return
		}
	}()

	p := parser{tokens: tokens}
	res = program(&p)
	return
}

func program(p *parser) ast {
	statements := []ast{}
	for !p.acceptT(t_END_OF_INPUT) {
		statements = append(statements, p.expectF(statement))
	}
	return newNode(n_PROGRAM, "", statements...)
}

func statement(p *parser) ast {
	startIndex := p.i
	if p.acceptT(t_NAME) {
		labelName := p.getLexeme()
		if !p.acceptT(t_COLON) {
			p.i = startIndex
		} else {
			st := p.expectF(statement)
			return newNode(n_LABEL, labelName, st)
		}
	}

	if p.acceptT(t_SEMI) {
		return newNode(n_EMPTY_STATEMENT, "")
	}

	if p.acceptF(throwStatement) ||
		p.acceptF(switchStatement) ||
		p.acceptF(tryCatchStatement) ||
		p.acceptF(declarationStatement) ||
		p.acceptF(classStatement) ||
		p.acceptF(blockStatement) ||
		p.acceptF(doWhileStatement) ||
		p.acceptF(whileStatement) ||
		p.acceptF(ifStatement) ||
		p.acceptF(forInOfStatement) ||
		p.acceptF(forStatement) ||
		p.acceptF(exportStatement) ||
		p.acceptF(importStatement) ||
		p.acceptF(functionExpression(true)) ||
		p.acceptF(withStatement) ||
		p.acceptF(expressionStatement) {
		return p.getNode()
	}

	p.expectT(t_END_OF_INPUT)
	return INVALID_NODE
}

func withStatement(p *parser) ast {
	if !p.acceptT(t_WITH) {
		return INVALID_NODE
	}
	p.expectT(t_PAREN_LEFT)
	expr := p.expectF(sequenceExpression)
	p.expectT(t_PAREN_RIGHT)

	st := p.expectF(statement)

	return newNode(n_WITH_STATEMENT, "", expr, st)
}

func throwStatement(p *parser) ast {
	if !p.acceptT(t_THROW) {
		return INVALID_NODE
	}
	expr := p.expectF(sequenceExpression)
	p.expectT(t_SEMI)
	return newNode(n_THROW_STATEMENT, "", expr)
}

func switchCase(p *parser) ast {
	var cond ast
	statements := []ast{}

	if p.acceptT(t_CASE) {
		cond = p.expectF(sequenceExpression)
	} else if p.acceptT(t_DEFAULT) {

	} else {
		return INVALID_NODE
	}

	p.expectT(t_COLON)

	for !(p.getNextT().tType == t_CASE ||
		p.getNextT().tType == t_DEFAULT ||
		p.getNextT().tType == t_CURLY_RIGHT) {
		statements = append(statements, p.expectF(statement))
	}

	return newNode(
		n_SWITCH_CASE, "", cond,
		newNode(n_MULTISTATEMENT, "", statements...),
	)
}

func switchStatement(p *parser) ast {
	if !p.acceptT(t_SWITCH) {
		return INVALID_NODE
	}

	var cond ast

	p.expectT(t_PAREN_LEFT)
	cond = p.expectF(sequenceExpression)
	p.expectT(t_PAREN_RIGHT)
	p.expectT(t_CURLY_LEFT)

	cases := []ast{}
	for !p.acceptT(t_CURLY_RIGHT) {
		cases = append(cases, p.expectF(switchCase))
	}

	return newNode(
		n_SWITCH_STATEMENT, "", cond,
		newNode(n_MULTISTATEMENT, "", cases...),
	)
}

func tryCatchStatement(p *parser) ast {
	if !p.acceptT(t_TRY) {
		return INVALID_NODE
	}

	var try, errorID, catch, finally ast
	try = p.expectF(blockStatement)

	if p.acceptT(t_CATCH) {
		p.expectT(t_PAREN_LEFT)
		errorID = p.expectF(identifier)
		p.expectT(t_PAREN_RIGHT)
		catch = p.expectF(blockStatement)
	}

	if p.acceptT(t_FINALLY) {
		finally = p.expectF(blockStatement)
	}

	return newNode(n_TRY_CATCH_STATEMENT, "", try, errorID, catch, finally)
}

func importStatement(p *parser) ast {
	if !p.acceptT(t_IMPORT) {
		return INVALID_NODE
	}

	var vars, all, path ast
	all.t = n_IMPORT_ALL
	path.t = n_STRING_LITERAL

	if p.acceptT(t_MULT) {
		p.skipT()
		if p.getLexeme() != "as" {
			p.checkASI(t_NAME)
		}

		p.expectF(identifier)
		all.value = p.getLexeme()

		p.skipT()
		if p.getLexeme() != "from" {
			p.checkASI(t_NAME)
		}
	} else {
		if p.acceptT(t_NAME) {
			name := newNode(n_IMPORT_NAME, "default")
			alias := newNode(n_IMPORT_ALIAS, p.getLexeme())

			varNode := newNode(n_IMPORT_VAR, "", name, alias)
			vars.children = append(vars.children, varNode)

			if p.acceptT(t_COMMA) {
				if p.acceptT(t_MULT) {
					p.skipT()
					if p.getLexeme() != "as" {
						p.checkASI(t_NAME)
					}

					p.expectF(identifier)
					all.value = p.getLexeme()

					p.skipT()
					if p.getLexeme() != "from" {
						p.checkASI(t_NAME)
					}
				}
			} else {
				p.skipT()
				if p.getLexeme() != "from" {
					p.checkASI(t_NAME)
				}
			}
		}

		if p.acceptT(t_CURLY_LEFT) {
			for !p.acceptT(t_CURLY_RIGHT) {
				if p.acceptT(t_NAME, t_DEFAULT) {
					name := newNode(n_IMPORT_NAME, p.getLexeme())
					alias := newNode(n_IMPORT_ALIAS, p.getLexeme())

					if p.acceptT(t_NAME) {
						if p.getLexeme() == "as" {
							p.skipT()
							alias = newNode(n_IMPORT_ALIAS, p.getLexeme())
						} else {
							p.i--
							p.expectT(t_CURLY_RIGHT)
						}
					}

					varNode := newNode(n_IMPORT_VAR, "", name, alias)
					vars.children = append(vars.children, varNode)
				}

				if !p.acceptT(t_COMMA) {
					p.expectT(t_CURLY_RIGHT)
					break
				}
			}

			p.skipT()
			if p.getLexeme() != "from" {
				p.checkASI(t_NAME)
			}
		}
	}

	path.value = p.expectF(stringLiteral).value
	p.expectT(t_SEMI)

	return newNode(n_IMPORT_STATEMENT, "", vars, all, path)
}

func declarationStatement(p *parser) ast {
	if !p.acceptF(varDeclaration) {
		return INVALID_NODE
	}
	p.expectT(t_SEMI)
	return newNode(n_DECLARATION_STATEMENT, "", p.getNode())
}

func expressionStatement(p *parser) ast {
	if p.acceptT(t_RETURN) {
		if p.acceptF(sequenceExpression) {
			return newNode(n_RETURN_STATEMENT, "", p.getNode())
		}
		return newNode(n_RETURN_STATEMENT, "")
	}

	if p.acceptT(t_BREAK, t_CONTINUE) {
		var label ast
		word := p.getLexeme()
		if p.acceptF(identifier) {
			label = p.getNode()
		}
		res := newNode(n_CONTROL_STATEMENT, word, label)
		p.expectT(t_SEMI)
		return res
	}
	if p.acceptT(t_DEBUGGER) {
		res := newNode(
			n_EXPRESSION_STATEMENT, "",
			newNode(n_CONTROL_WORD, p.getLexeme()),
		)
		p.expectT(t_SEMI)
		return res
	}

	if !(p.acceptF(varDeclaration) ||
		p.acceptF(sequenceExpression)) {
		return INVALID_NODE
	}
	p.expectT(t_SEMI)
	return newNode(n_EXPRESSION_STATEMENT, "", p.getNode())
}

func varDeclaration(p *parser) ast {
	if !p.acceptT(t_VAR, t_CONST, t_LET) {
		return INVALID_NODE
	}
	kind := p.getLexeme()

	declarators := []ast{}
	for ok := true; ok; ok = p.acceptT(t_COMMA) {
		declarators = append(declarators, p.expectF(declarator))
	}
	return newNode(n_VARIABLE_DECLARATION, kind, declarators...)
}

func declarator(p *parser) ast {
	if !p.acceptF(identifier) &&
		!p.acceptF(objectPattern) {
		return INVALID_NODE
	}
	var name, value ast
	name = p.getNode()

	if p.acceptT(t_ASSIGN) {
		value = p.expectF(assignmentExpression)
	}

	return newNode(n_DECLARATOR, "", name, value)
}

func sequenceExpression(p *parser) ast {
	if !p.acceptF(assignmentExpression) {
		return INVALID_NODE
	}
	first := p.getNode()

	if p.acceptT(t_COMMA) {
		items := []ast{first}

		for ok := true; ok; ok = p.acceptT(t_COMMA) {
			items = append(items, p.expectF(assignmentExpression))
		}

		return newNode(n_SEQUENCE_EXPRESSION, "", items...)
	}

	return first
}

func assignmentExpression(p *parser) ast {
	if p.testFlag(f_GENERATOR_FUNCTION) {
		if p.acceptT(t_YIELD) {
			value := ""
			if p.acceptT(t_MULT) {
				value = "*"
			}

			var expr ast

			if p.getNextT().tType != t_SEMI {
				expr = p.expectF(assignmentExpression)
			}

			return newNode(n_YIELD_EXPRESSION, value, expr)
		}
	}

	prevIndex := p.i

	if p.acceptF(objectPattern) ||
		p.acceptF(arrayPattern) {
		if p.acceptT(t_ASSIGN) {
			return newNode(
				n_ASSIGNMENT_EXPRESSION, "=", p.getNode(),
				p.expectF(assignmentExpression),
			)
		}
		p.i = prevIndex
	}

	if !p.acceptF(conditionalExpression) {
		return INVALID_NODE
	}
	left := p.getNode()

	if p.acceptT(
		t_ASSIGN, t_PLUS_ASSIGN, t_MINUS_ASSIGN, t_MULT_ASSIGN, t_DIV_ASSIGN,
		t_BITWISE_SHIFT_LEFT_ASSIGN, t_BITWISE_SHIFT_RIGHT_ASSIGN,
		t_BITWISE_SHIFT_RIGHT_ZERO_ASSIGN, t_BITWISE_AND_ASSIGN,
		t_BITWISE_OR_ASSIGN, t_BITWISE_XOR_ASSIGN,
	) {
		kind := p.getLexeme()
		right := p.expectF(assignmentExpression)
		return newNode(n_ASSIGNMENT_EXPRESSION, kind, left, right)
	}

	return left
}

func conditionalExpression(p *parser) ast {
	if !p.acceptF(binaryExpression) {
		return INVALID_NODE
	}
	cond := p.getNode()

	if p.acceptT(t_QUESTION) {
		conseq := p.expectF(conditionalExpression)
		p.expectT(t_COLON)
		altern := p.expectF(conditionalExpression)

		return newNode(
			n_CONDITIONAL_EXPRESSION, "", cond, conseq, altern,
		)
	}

	return cond
}

func binaryExpression(p *parser) ast {
	opStack := make([]token, 0)
	outputStack := make([]ast, 0)
	var root *ast

	addNode := func(t token) {
		right := outputStack[len(outputStack)-1]
		left := outputStack[len(outputStack)-2]
		nn := newNode(n_BINARY_EXPRESSION, t.lexeme, left, right)

		outputStack = outputStack[:len(outputStack)-2]
		outputStack = append(outputStack, nn)
		root = &nn
	}

	for {
		if !p.acceptF(prefixUnaryExpression) {
			return INVALID_NODE
		}
		expr := p.getNode()

		outputStack = append(
			outputStack, expr,
		)

		op, ok := operatorTable[p.getNextT().tType]
		if !ok {
			break
		}

		for {
			if len(opStack) > 0 {
				opToken := opStack[len(opStack)-1]

				opFromStack := operatorTable[opToken.tType]
				if opFromStack.precedence > op.precedence ||
					(opFromStack.precedence == op.precedence && !op.isRightAssociative) {
					addNode(opToken)
					opStack = opStack[:len(opStack)-1]
				} else {
					break
				}
			} else {
				break
			}
		}
		opStack = append(opStack, p.getNextT())
		p.skipT()
	}

	for i := len(opStack) - 1; i >= 0; i-- {
		addNode(opStack[i])
	}

	if root != nil {
		return *root
	}
	return outputStack[len(outputStack)-1]
}

func prefixUnaryExpression(p *parser) ast {
	if p.testFlag(f_ASYNC_FUNCTION) {
		if p.acceptT(t_AWAIT) {
			return newNode(n_AWAIT_EXPRESSION, "", p.expectF(prefixUnaryExpression))
		}
	}

	if p.acceptT(
		t_NOT, t_BITWISE_NOT, t_PLUS, t_MINUS,
		t_INC, t_DEC, t_TYPEOF, t_VOID, t_DELETE,
	) {
		op := p.getLexeme()
		if !p.acceptF(prefixUnaryExpression) {
			return INVALID_NODE
		}
		value := p.getNode()
		return newNode(
			n_PREFIX_UNARY_EXPRESSION, op, value,
		)
	}

	if !p.acceptF(postfixUnaryExpression) {
		return INVALID_NODE
	}
	return p.getNode()
}

func postfixUnaryExpression(p *parser) ast {
	if !p.acceptF(functionCallOrMemberExpression(false)) {
		return INVALID_NODE
	}
	value := p.getNode()

	if p.acceptT(t_INC, t_DEC) {
		return newNode(n_POSTFIX_UNARY_EXPRESSION, p.getLexeme(), value)
	}

	return value
}

func functionCallOrMemberExpression(noFuncCall bool) func(*parser) ast {
	return func(p *parser) ast {
		if !p.acceptF(constructorCall) {
			return INVALID_NODE
		}
		funcName := p.getNode()

		for {
			if !noFuncCall && p.acceptF(functionArgs) {
				call := newNode(n_FUNCTION_CALL, "", funcName, p.getNode())
				funcName = call
			} else {
				var property ast

				if p.acceptT(t_DOT) {
					if p.acceptF(identifier) {
						property = p.getNode()
						property.t = n_OBJECT_MEMBER
					} else {
						next := p.getNextT()
						if isValidPropertyName(next.lexeme) {
							property = newNode(n_OBJECT_MEMBER, next.lexeme)
							p.skipT()
						}
					}
				} else if p.acceptT(t_BRACKET_LEFT) {
					property = newNode(
						n_CALCULATED_PROPERTY_NAME, "", p.expectF(sequenceExpression),
					)
					p.expectT(t_BRACKET_RIGHT)
				}

				if property.t == n_EMPTY {
					p.flags &= ^f_ACCEPT_NO_FUNCTION_CALL

					if p.acceptF(templateLiteral) {
						return newNode(
							n_TAGGED_LITERAL, "", funcName,
							p.getNode(),
						)
					}
					return funcName
				}

				funcName = newNode(n_MEMBER_EXPRESSION, "", funcName, property)
			}
		}
	}
}

func functionArgs(p *parser) ast {
	if !p.acceptT(t_PAREN_LEFT) {
		return INVALID_NODE
	}

	args := []ast{}

	for !p.acceptT(t_PAREN_RIGHT) {
		if p.acceptF(assignmentExpression) {
			args = append(args, p.getNode())
		}

		if !p.acceptT(t_COMMA) {
			p.expectT(t_PAREN_RIGHT)
			break
		}
	}

	return newNode(n_FUNCTION_ARGS, "", args...)
}

func constructorCall(p *parser) ast {
	if p.acceptT(t_NEW) {
		var name, args ast

		name = p.expectF(functionCallOrMemberExpression(true))

		if p.acceptF(functionArgs) {
			args = p.getNode()
		}

		return newNode(n_CONSTRUCTOR_CALL, "", name, args)
	}

	if !p.acceptF(atom) {
		return INVALID_NODE
	}
	return p.getNode()
}

func atom(p *parser) ast {
	if p.acceptF(classExpression) ||
		p.acceptF(objectLiteral) ||
		p.acceptF(otherLiteral) ||
		p.acceptF(lambdaExpression) ||
		p.acceptF(parenExpression) ||
		p.acceptF(regexpLiteral) ||
		p.acceptF(arrayLiteral) ||
		p.acceptF(identifier) ||
		p.acceptF(functionExpression(false)) {
		return p.getNode()
	}

	return INVALID_NODE
}

func parenExpression(p *parser) ast {
	if !p.acceptT(t_PAREN_LEFT) {
		return INVALID_NODE
	}
	expr := p.expectF(sequenceExpression)
	p.expectT(t_PAREN_RIGHT)
	return newNode(n_PAREN_EXPRESSION, "", expr)
}

func lambdaExpression(p *parser) ast {
	var args, body ast

	if p.acceptF(identifier) {
		args = newNode(n_FUNCTION_PARAMETERS, "", p.getNode())
	} else if p.acceptF(functionParameters) {
		args = p.getNode()
	}

	if !p.acceptT(t_LAMBDA) {
		return INVALID_NODE
	}

	if p.acceptF(blockStatement) {
		body = p.getNode()
	} else {
		body = p.expectF(assignmentExpression)
	}
	return newNode(n_LAMBDA_EXPRESSION, "", args, body)
}

func regexpLiteral(p *parser) ast {
	if !p.acceptT(t_DIV) {
		return INVALID_NODE
	}

	var body, flags ast

	bodyValue := ""
	for {
		if p.tokens[p.i].tType == t_BRACKET_LEFT {
			for p.tokens[p.i].tType != t_BRACKET_RIGHT {
				bodyValue += p.tokens[p.i].lexeme
				p.i++
			}
			bodyValue += p.tokens[p.i].lexeme
			p.i++
			continue
		}
		if p.tokens[p.i].tType == t_ESCAPE {
			bodyValue += p.tokens[p.i].lexeme
			bodyValue += p.tokens[p.i+1].lexeme
			p.i += 2
			continue
		}
		if p.tokens[p.i].tType == t_DIV {
			p.i++
			break
		}
		bodyValue += p.tokens[p.i].lexeme
		p.i++
	}

	body = newNode(n_REGEXP_BODY, bodyValue)
	if p.acceptT(t_NAME) {
		flags = newNode(n_REGEXP_FLAGS, p.getLexeme())
	}

	return newNode(n_REGEXP_LITERAL, "", body, flags)
}

func arrayLiteral(p *parser) ast {
	if !p.acceptT(t_BRACKET_LEFT) {
		return INVALID_NODE
	}

	items := []ast{}
	for !p.acceptT(t_BRACKET_RIGHT) {
		if p.acceptT(t_SPREAD) {
			items = append(
				items, newNode(
					n_SPREAD_ELEMENT, "", p.expectF(assignmentExpression),
				),
			)
		} else if p.acceptF(assignmentExpression) {
			items = append(items, p.getNode())
		} else {
			items = append(items, newNode(n_EMPTY_EXPRESSION, ""))
		}

		if !p.acceptT(t_COMMA) {
			p.expectT(t_BRACKET_RIGHT)
			break
		}
	}

	return newNode(n_ARRAY_LITERAL, "", items...)
}

func objectProperty(p *parser) ast {
	var kind string
	var key, value ast

	if p.acceptT(t_SPREAD) {
		return newNode(
			n_SPREAD_ELEMENT, "", p.expectF(assignmentExpression),
		)
	}

	key = p.expectF(propertyKey)
	if key.value == "set" || key.value == "get" {
		if p.getNextT().tType != t_COLON && p.getNextT().tType != t_PAREN_LEFT {
			kind = key.value
			key = p.expectF(propertyKey)
		}
	}

	if p.getNextT().tType == t_PAREN_LEFT {
		params := p.expectF(functionParameters)
		body := p.expectF(blockStatement)
		return newNode(n_OBJECT_METHOD, kind, key, params, body)
	}

	if p.acceptT(t_COLON) {
		value = p.expectF(assignmentExpression)
	}

	if value.t == n_EMPTY {
		value = key
		value.t = n_IDENTIFIER
	}

	return newNode(n_OBJECT_PROPERTY, kind, key, value)
}

func objectLiteral(p *parser) ast {
	if !p.acceptT(t_CURLY_LEFT) {
		return INVALID_NODE
	}

	props := []ast{}
	for !p.acceptT(t_CURLY_RIGHT) {
		props = append(props, p.expectF(objectProperty))

		if !p.acceptT(t_COMMA) {
			p.expectT(t_CURLY_RIGHT)
			break
		}
	}

	return newNode(n_OBJECT_LITERAL, "", props...)
}

func stringLiteral(p *parser) ast {
	if !p.acceptT(t_STRING_QUOTE) {
		return INVALID_NODE
	}
	firstQuote := p.getLexeme()
	value := ""
	for {
		if p.tokens[p.i].tType == t_ESCAPE {
			value += p.tokens[p.i].lexeme
			value += p.tokens[p.i+1].lexeme
			p.i += 2
			continue
		}
		if p.tokens[p.i].lexeme == firstQuote {
			break
		}
		value += p.tokens[p.i].lexeme
		p.i++
	}
	p.i++

	if firstQuote == "'" {
		return newNode(n_STRING_LITERAL, value)
	}
	return newNode(n_DOUBLE_QUOTE_STRING_LITERAL, value)
}

func templateLiteral(p *parser) ast {
	if !p.acceptT(t_TEMPLATE_LITERAL_QUOTE) {
		return INVALID_NODE
	}
	firstQuote := p.getLexeme()
	value := ""
	for {
		if p.tokens[p.i].tType == t_ESCAPE {
			value += p.tokens[p.i].lexeme
			value += p.tokens[p.i+1].lexeme
			p.i += 2
			continue
		}
		if p.tokens[p.i].lexeme == firstQuote {
			break
		}
		value += p.tokens[p.i].lexeme
		p.i++
	}
	p.i++

	return newNode(n_TEMPLATE_LITERAL, value)
}

func otherLiteral(p *parser) ast {
	if p.acceptF(stringLiteral) {
		return p.getNode()
	} else if p.acceptT(t_NUMBER) {
		return newNode(n_NUMBER_LITERAL, p.getLexeme())
	} else if p.acceptT(t_TRUE, t_FALSE) {
		return newNode(n_BOOL_LITERAL, p.getLexeme())
	} else if p.acceptT(t_NULL) {
		return newNode(n_NULL, p.getLexeme())
	} else if p.acceptT(t_UNDEFINED) {
		return newNode(n_UNDEFINED, p.getLexeme())
	} else if p.acceptF(templateLiteral) {
		return p.getNode()
	}
	return INVALID_NODE
}

func identifier(p *parser) ast {
	if p.acceptT(t_NAME, t_THIS, t_OF) {
		return newNode(n_IDENTIFIER, p.getLexeme())
	}
	return INVALID_NODE
}

func arrayPattern(p *parser) ast {
	if !p.acceptT(t_BRACKET_LEFT) {
		return INVALID_NODE
	}

	items := []ast{}
	for !p.acceptT(t_BRACKET_RIGHT) {
		if p.acceptT(t_SPREAD) {
			items = append(items, newNode(
				n_SPREAD_ELEMENT, "", p.expectF(identifier),
			))
			continue
		}

		if p.acceptF(assignmentPattern) {
			items = append(items, p.getNode())
		} else if p.getNextT().tType == t_COMMA {
			items = append(items, newNode(n_EMPTY, ""))
		} else {
			return INVALID_NODE
		}

		if !p.acceptT(t_COMMA) {
			if !p.acceptT(t_BRACKET_RIGHT) {
				return INVALID_NODE
			}
			break
		}
	}

	return newNode(n_ARRAY_PATTERN, "", items...)
}

func objectPattern(p *parser) ast {
	if !p.acceptT(t_CURLY_LEFT) {
		return INVALID_NODE
	}

	props := []ast{}
	for !p.acceptT(t_CURLY_RIGHT) {
		if p.acceptT(t_SPREAD) {
			props = append(props, newNode(
				n_REST_ELEMENT, "", p.expectF(identifier),
			))
			continue
		}

		if p.acceptF(propertyKey) {
			var key, value ast

			key = p.getNode()
			if p.acceptT(t_COLON) {
				if p.acceptF(assignmentPattern) {
					value = p.getNode()
				} else {
					return INVALID_NODE
				}
			}

			props = append(props, newNode(n_OBJECT_PROPERTY, "", key, value))
		}

		if !p.acceptT(t_COMMA) {
			if !p.acceptT(t_CURLY_RIGHT) {
				return INVALID_NODE
			}
			break
		}
	}

	return newNode(n_OBJECT_PATTERN, "", props...)
}

func isValidPropertyName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if !isLetter(name[0]) {
		return false
	}
	for _, symbol := range name {
		if !isLetter(byte(symbol)) && !isNumber(byte(symbol)) {
			return false
		}
	}
	return true
}

func propertyKey(p *parser) ast {
	if p.acceptT(t_BRACKET_LEFT) {
		res := p.expectF(sequenceExpression)
		p.expectT(t_BRACKET_RIGHT)
		return newNode(n_CALCULATED_OBJECT_KEY, "", res)
	}

	if p.acceptF(identifier) {
		return newNode(n_OBJECT_KEY, p.getLexeme())
	}

	if p.acceptF(otherLiteral) {
		return newNode(n_NON_IDENTIFIER_OBJECT_KEY, "", p.getNode())
	}

	if isValidPropertyName(p.getNextT().lexeme) {
		p.skipT()
		return newNode(n_NON_IDENTIFIER_OBJECT_KEY, "", newNode(n_IDENTIFIER, p.getLexeme()))
	}

	return INVALID_NODE
}

func assignmentPattern(p *parser) ast {
	var left, right ast

	if p.acceptF(objectPattern) ||
		p.acceptF(identifier) ||
		p.acceptF(arrayPattern) {
		left = p.getNode()
	} else {
		return INVALID_NODE
	}

	if p.acceptT(t_ASSIGN) {
		right = p.expectF(conditionalExpression)
		return newNode(n_ASSIGNMENT_PATTERN, "", left, right)
	}

	return left
}

func functionParameter(p *parser) ast {
	if p.acceptT(t_SPREAD) {
		if p.acceptF(objectPattern) ||
			p.acceptF(arrayPattern) ||
			p.acceptF(identifier) {
			return newNode(n_REST_ELEMENT, "", p.getNode())
		}
		res := p.expectF(identifier)
		return res
	}

	if p.acceptF(assignmentPattern) {
		return p.getNode()
	}

	return INVALID_NODE
}

func functionParameters(p *parser) ast {
	if !p.acceptT(t_PAREN_LEFT) {
		return INVALID_NODE
	}

	paramsSlice := []ast{}
	for !p.acceptT(t_PAREN_RIGHT) {
		if p.acceptF(functionParameter) {
			paramsSlice = append(paramsSlice, p.getNode())
		}

		if !p.acceptT(t_COMMA) {
			if !p.acceptT(t_PAREN_RIGHT) {
				return INVALID_NODE
			}
			break
		}
	}

	return newNode(n_FUNCTION_PARAMETERS, "", paramsSlice...)
}

func blockStatement(p *parser) ast {
	if !p.acceptT(t_CURLY_LEFT) {
		return INVALID_NODE
	}
	statements := []ast{}
	for !p.acceptT(t_CURLY_RIGHT) {
		statements = append(statements, p.expectF(statement))
	}

	return newNode(n_BLOCK_STATEMENT, "", statements...)
}

func functionExpression(isDeclaration bool) func(*parser) ast {
	return func(p *parser) ast {
		isAsync := false
		if p.acceptT(t_ASYNC) {
			isAsync = true
		}
		if !p.acceptT(t_FUNCTION) {
			return INVALID_NODE
		}
		typ := ""
		if p.acceptT(t_MULT) {
			typ = "*"
		}

		var name, params, body ast
		if isDeclaration {
			name = p.expectF(identifier)
		} else if p.acceptF(identifier) {
			name = p.getNode()
		}

		params = p.expectF(functionParameters)

		if isAsync {
			p.flags |= f_ASYNC_FUNCTION
		}
		if typ == "*" {
			p.flags |= f_GENERATOR_FUNCTION
		}
		body = p.expectF(blockStatement)
		p.flags &= ^f_GENERATOR_FUNCTION
		p.flags &= ^f_ASYNC_FUNCTION

		funcExpr := newNode(n_FUNCTION_EXPRESSION, typ, name, params, body)
		if isAsync {
			return newNode(n_ASYNC_FUNCTION, "", funcExpr)
		}
		return funcExpr
	}
}

func whileStatement(p *parser) ast {
	if !(p.acceptT(t_WHILE) && p.acceptT(t_PAREN_LEFT)) {
		return INVALID_NODE
	}
	var cond, body ast
	cond = p.expectF(sequenceExpression)
	p.expectT(t_PAREN_RIGHT)
	body = p.expectF(statement)

	return newNode(n_WHILE_STATEMENT, "", cond, body)
}

func doWhileStatement(p *parser) ast {
	if !p.acceptT(t_DO) {
		return INVALID_NODE
	}
	var cond, body ast
	body = p.expectF(statement)

	if !(p.acceptT(t_WHILE) && p.acceptT(t_PAREN_LEFT)) {
		return INVALID_NODE
	}
	cond = p.expectF(sequenceExpression)
	p.expectT(t_PAREN_RIGHT)

	return newNode(n_DO_WHILE_STATEMENT, "", cond, body)
}

func forInOfStatement(p *parser) ast {
	if !(p.acceptT(t_FOR) && p.acceptT(t_PAREN_LEFT)) {
		return INVALID_NODE
	}
	var left, right, body ast

	if p.acceptF(varDeclaration) ||
		p.acceptF(assignmentPattern) {
		left = p.getNode()
	} else {
		return INVALID_NODE
	}

	if !p.acceptT(t_IN) && !p.acceptT(t_OF) {
		return INVALID_NODE
	}
	kind := p.getLexeme()

	right = p.expectF(sequenceExpression)

	if !p.acceptT(t_PAREN_RIGHT) {
		return INVALID_NODE
	}

	body = p.expectF(statement)

	if kind == "in" {
		return newNode(n_FOR_IN_STATEMENT, "", left, right, body)
	}
	return newNode(n_FOR_OF_STATEMENT, "", left, right, body)
}

func forStatement(p *parser) ast {
	if !(p.acceptT(t_FOR) && p.acceptT(t_PAREN_LEFT)) {
		return INVALID_NODE
	}
	var init, test, final, body ast

	if p.acceptT(t_SEMI) {
		init = newNode(n_EMPTY_EXPRESSION, "")
	} else if p.acceptF(varDeclaration) ||
		p.acceptF(sequenceExpression) {
		init = p.getNode()

		p.expectT(t_SEMI)
	} else {
		p.expectT(t_SEMI)
	}

	if p.acceptT(t_SEMI) {
		test = newNode(n_EMPTY_EXPRESSION, "")
	} else {
		test = p.expectF(sequenceExpression)

		p.expectT(t_SEMI)
	}

	if p.acceptT(t_PAREN_RIGHT) {
		final = newNode(n_EMPTY_EXPRESSION, "")
	} else {
		final = p.expectF(sequenceExpression)

		p.expectT(t_PAREN_RIGHT)
	}

	body = p.expectF(statement)
	return newNode(n_FOR_STATEMENT, "", init, test, final, body)
}

func ifStatement(p *parser) ast {
	if !(p.acceptT(t_IF) &&
		p.acceptT(t_PAREN_LEFT)) {
		return INVALID_NODE
	}
	var cond, conseq, alt ast
	cond = p.expectF(sequenceExpression)
	if !(p.acceptT(t_PAREN_RIGHT)) {
		return INVALID_NODE
	}
	conseq = p.expectF(statement)

	if p.acceptT(t_ELSE) {
		alt = p.expectF(statement)
	}

	return newNode(n_IF_STATEMENT, "", cond, conseq, alt)
}

func exportStatement(p *parser) ast {
	if !p.acceptT(t_EXPORT) {
		return INVALID_NODE
	}

	var vars, declaration, path ast
	vars.t = n_EXPORT_VARS
	declaration.t = n_EXPORT_DECLARATION
	path.t = n_STRING_LITERAL

	if p.acceptT(t_DEFAULT) {
		var name, alias ast

		alias = newNode(n_EXPORT_ALIAS, "default")
		if p.acceptF(functionExpression(false)) {
			fe := p.getNode()
			feName := fe.children[0]

			if feName.t != n_EMPTY {
				declaration = fe
				name = feName
			} else {
				name = fe
			}
		} else if p.acceptF(classExpression) {
			ce := p.getNode()
			ceName := ce.children[0]
			if ceName.t != n_EMPTY {
				declaration = ce
				name = ceName
			} else {
				name = ce
			}
		} else {
			p.expectF(sequenceExpression)
			name = p.getNode()
		}

		ev := newNode(n_EXPORT_VAR, "", name, alias)
		vars.children = append(vars.children, ev)
		p.expectT(t_SEMI)

	} else if p.acceptT(t_CURLY_LEFT) {
		for !p.acceptT(t_CURLY_RIGHT) {
			if p.acceptF(identifier) {
				name := p.getNode()
				alias := name

				if p.acceptT(t_NAME) && p.getLexeme() == "as" {
					p.skipT()
					alias = newNode(n_EXPORT_ALIAS, p.getLexeme())
				}
				ev := newNode(n_EXPORT_VAR, "", name, alias)
				vars.children = append(vars.children, ev)
			}

			if !p.acceptT(t_COMMA) {
				p.expectT(t_CURLY_RIGHT)
				break
			}
		}

		if p.acceptT(t_NAME) && p.getLexeme() == "from" {
			path = p.expectF(stringLiteral)
		}
		p.expectT(t_SEMI)

	} else if p.acceptF(declarationStatement) {
		declaration = p.getNode()
		for _, d := range declaration.children[0].children {
			name := d.children[0]
			alias := d.children[0]
			ev := newNode(n_EXPORT_VAR, "", name, alias)

			vars.children = append(vars.children, ev)
		}

	} else if p.acceptF(functionExpression(true)) {
		fs := p.getNode()
		declaration = fs
		name := fs.children[0]
		alias := name
		ev := newNode(n_EXPORT_VAR, "", name, alias)
		vars.children = append(vars.children, ev)

	} else if p.acceptF(classStatement) {
		cs := p.getNode()
		declaration = cs
		name := cs.children[0]
		alias := name
		ev := newNode(n_EXPORT_VAR, "", name, alias)
		vars.children = append(vars.children, ev)

	} else if p.acceptT(t_MULT) {
		p.expectT(t_NAME)
		if p.getLexeme() != "from" {
			return INVALID_NODE
		}
		path = p.expectF(stringLiteral)
		vars.t = n_EXPORT_ALL
		p.expectT(t_SEMI)
	}

	result := newNode(n_EXPORT_STATEMENT, "", vars, declaration, path)
	return result
}

func classBody(p *parser) ast {
	if !p.acceptT(t_CURLY_LEFT) {
		return INVALID_NODE
	}

	props := []ast{}
	for !p.acceptT(t_CURLY_RIGHT) {
		var prop, key, value ast
		kind := ""
		isStatic := false

		if p.acceptT(t_STATIC) {
			isStatic = true
		}

		key = p.expectF(propertyKey)
		if key.value == "set" || key.value == "get" {
			kind = key.value
			key = p.expectF(propertyKey)
		}

		if p.getNextT().tType == t_PAREN_LEFT {
			params := p.expectF(functionParameters)
			body := p.expectF(blockStatement)
			prop = newNode(n_CLASS_METHOD, kind, key, params, body)
		} else {
			if p.acceptT(t_ASSIGN, t_COLON) {
				value = p.expectF(sequenceExpression)
			}
			prop = newNode(n_CLASS_PROPERTY, "", key, value)
		}

		if isStatic {
			props = append(props, newNode(
				n_CLASS_STATIC_PROPERTY, "", prop,
			))
		} else {
			props = append(props, prop)
		}

		if prop.t != n_CLASS_METHOD {
			p.expectT(t_SEMI)
		}
	}

	return newNode(n_CLASS_BODY, "", props...)
}

func classExpression(p *parser) ast {
	if !p.acceptT(t_CLASS) {
		return INVALID_NODE
	}

	var name, extends, body ast
	if p.acceptF(identifier) {
		name = p.getNode()
	}
	if p.acceptT(t_EXTENDS) {
		extends = p.expectF(functionCallOrMemberExpression(false))
	}
	body = p.expectF(classBody)

	return newNode(n_CLASS_EXPRESSION, "", name, extends, body)
}

func classStatement(p *parser) ast {
	if !p.acceptT(t_CLASS) {
		return INVALID_NODE
	}

	var name, extends, body ast
	name = p.expectF(identifier)
	if p.acceptT(t_EXTENDS) {
		extends = p.expectF(functionCallOrMemberExpression(false))
	}
	body = p.expectF(classBody)

	return newNode(n_CLASS_STATEMENT, "", name, extends, body)
}
