package jsLoader

import (
	"fmt"
)

type ast struct {
	t        nodeType
	value    string
	children []ast
}

const (
	f_ACCEPT_NO_FUNCTION_CALL = 1 << 0
	f_ACCEPT_NO_IN            = 1 << 1
)

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

func makeNode(t nodeType, value string, children ...ast) ast {
	return ast{t, value, children}
}

var INVALID_NODE = ast{t: n_INVALID}

type parser struct {
	tokens []token
	i      int
	t      token

	node  ast
	flags int
}

func (p *parser) testFlag(f int) bool {
	return (p.flags & f) != 0
}

func (p *parser) tNext() {
	p.i++
	if p.i < len(p.tokens) {
		p.t = p.tokens[p.i]
	}
}

func (p *parser) acceptT(tTypes ...tokenType) bool {
	for _, tType := range tTypes {
		if p.t.tType == tType {
			p.tNext()
			return true
		}
	}
	return false
}

type parseFunc func(*parser) ast

func (p *parser) acceptF(f parseFunc) bool {
	prevIndex := p.i

	n := f(p)
	if n.t == n_INVALID {
		p.i = prevIndex
		p.t = p.tokens[p.i]
		return false
	}

	p.node = n
	return true
}

func (p *parser) takeNode() ast {
	return p.node
}

func (p *parser) takeLexeme() string {
	return p.tokens[p.i-1].lexeme
}

func (p *parser) expectF(f parseFunc) ast {
	if !p.acceptF(f) {
		fmt.Printf("Parsing error, %v", p.t)
		panic("error")
	}
	return p.takeNode()
}

func program(p *parser) ast {
	statements := []ast{}
	for !p.acceptT(tEND_OF_INPUT) {
		statements = append(statements, p.expectF(statement))
	}
	return makeNode(n_PROGRAM, "", statements...)
}

func statement(p *parser) ast {
	if p.acceptT(tSEMI) {
		return makeNode(n_EMPTY_STATEMENT, "")
	} else if p.acceptF(doWhileStatement) ||
		p.acceptF(whileStatement) ||
		p.acceptF(ifStatement) ||
		p.acceptF(forInOfStatement) ||
		p.acceptF(forStatement) ||
		p.acceptF(exportStatement) ||
		p.acceptF(importStatement) ||
		p.acceptF(functionDeclaration) ||
		p.acceptF(expressionStatement) {
		return p.takeNode()
	}

	return INVALID_NODE
}

func importStatement(p *parser) ast {
	if !p.acceptT(tIMPORT) {
		return INVALID_NODE
	}

	var vars, all, path ast
	all.t = n_IMPORT_ALL
	path.t = n_IMPORT_PATH

	if p.acceptT(tMULT) {
		p.tNext()
		if p.takeLexeme() != "as" {
			return INVALID_NODE
		}

		p.expectF(identifier)
		all.value = p.takeLexeme()

		p.tNext()
		if p.takeLexeme() != "from" {
			return INVALID_NODE
		}
	} else {
		if p.acceptT(tNAME) {
			name := makeNode(n_IMPORT_NAME, "default")
			alias := makeNode(n_IMPORT_ALIAS, p.takeLexeme())

			varNode := makeNode(n_IMPORT_VAR, "", name, alias)
			vars.children = append(vars.children, varNode)

			if p.acceptT(tCOMMA) {
				if p.acceptT(tMULT) {
					p.tNext()
					if p.takeLexeme() != "as" {
						return INVALID_NODE
					}

					p.expectF(identifier)
					all.value = p.takeLexeme()

					p.tNext()
					if p.takeLexeme() != "from" {
						return INVALID_NODE
					}
				}
			} else {
				p.tNext()
				if p.takeLexeme() != "from" {
					return INVALID_NODE
				}
			}
		}

		if p.acceptT(tCURLY_LEFT) {
			for !p.acceptT(tCURLY_RIGHT) {
				if p.acceptT(tNAME, tDEFAULT) {
					name := makeNode(n_IMPORT_NAME, p.takeLexeme())
					alias := makeNode(n_IMPORT_ALIAS, p.takeLexeme())

					if p.acceptT(tNAME) && p.takeLexeme() == "as" {
						p.tNext()
						alias = makeNode(n_IMPORT_ALIAS, p.takeLexeme())
					}

					varNode := makeNode(n_IMPORT_VAR, "", name, alias)
					vars.children = append(vars.children, varNode)
				}

				if !p.acceptT(tCOMMA) {
					if !p.acceptT(tCURLY_RIGHT) {
						return INVALID_NODE
					}
					break
				}
			}

			p.tNext()
			if p.takeLexeme() != "from" {
				return INVALID_NODE
			}
		}
	}

	if !p.acceptT(tSTRING) {
		return INVALID_NODE
	}
	path.value = p.takeLexeme()
	if !p.acceptT(tSEMI) {
		return INVALID_NODE
	}

	return makeNode(n_IMPORT_STATEMENT, "", vars, all, path)
}

func expressionStatement(p *parser) ast {
	if !(p.acceptF(varDeclaration) ||
		p.acceptF(sequenceExpression)) {
		return INVALID_NODE
	}
	if !p.acceptT(tSEMI) {
		return INVALID_NODE
	}
	return makeNode(n_EXPRESSION_STATEMENT, "", p.takeNode())
}

func varDeclaration(p *parser) ast {
	if !p.acceptT(tVAR, tCONST, tLET) {
		return INVALID_NODE
	}
	kind := p.takeLexeme()

	declarators := []ast{}
	for ok := true; ok; ok = p.acceptT(tCOMMA) {
		declarators = append(declarators, p.expectF(declarator))
	}
	return makeNode(n_VARIABLE_DECLARATION, kind, declarators...)
}

func declarator(p *parser) ast {
	if !p.acceptF(identifier) &&
		!p.acceptF(objectPattern) {
		return INVALID_NODE
	}
	var name, value ast
	name = p.takeNode()

	if p.acceptT(tASSIGN) {
		value = p.expectF(assignmentExpression)
	}

	return makeNode(n_DECLARATOR, "", name, value)
}

func expression(p *parser) ast {
	return INVALID_NODE
}

func sequenceExpression(p *parser) ast {
	if !p.acceptF(assignmentExpression) {
		return INVALID_NODE
	}
	first := p.takeNode()

	if p.acceptT(tCOMMA) {
		items := []ast{}

		for ok := true; ok; ok = p.acceptT(tCOMMA) {
			items = append(items, p.expectF(assignmentExpression))
		}

		return makeNode(n_SEQUENCE_EXPRESSION, "", items...)
	}

	return first
}

func assignmentExpression(p *parser) ast {
	if !p.acceptF(conditionalExpression) {
		return INVALID_NODE
	}
	left := p.takeNode()

	if p.acceptT(
		tASSIGN, tPLUS_ASSIGN, tMINUS_ASSIGN, tMULT_ASSIGN, tDIV_ASSIGN,
		tBITWISE_SHIFT_LEFT_ASSIGN, tBITWISE_SHIFT_RIGHT_ASSIGN,
		tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN, tBITWISE_AND_ASSIGN,
		tBITWISE_OR_ASSIGN, tBITWISE_XOR_ASSIGN,
	) {
		right := p.expectF(assignmentExpression)
		return makeNode(n_ASSIGNMENT_EXPRESSION, "", left, right)
	}

	return left
}

func conditionalExpression(p *parser) ast {
	if !p.acceptF(binaryExpression) {
		return INVALID_NODE
	}
	cond := p.takeNode()

	if p.acceptT(tQUESTION) {
		conseq := p.expectF(conditionalExpression)
		if !p.acceptT(tCOLON) {
			return INVALID_NODE
		}
		altern := p.expectF(conditionalExpression)

		return makeNode(
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
		nn := makeNode(n_BINARY_EXPRESSION, t.lexeme, left, right)

		outputStack = outputStack[:len(outputStack)-2]
		outputStack = append(outputStack, nn)
		root = &nn
	}

	for {
		if !p.acceptF(prefixUnaryExpression) {
			return INVALID_NODE
		}
		expr := p.takeNode()

		outputStack = append(
			outputStack, expr,
		)

		op, ok := operatorTable[p.t.tType]
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
		opStack = append(opStack, p.t)
		p.tNext()
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
	if p.acceptT(
		tNOT, tBITWISE_NOT, tPLUS, tMINUS,
		tINC, tDEC, tTYPEOF, tVOID, tDELETE,
	) {
		op := p.takeLexeme()
		if !p.acceptF(postfixUnaryExpression) {
			return INVALID_NODE
		}
		value := p.takeNode()
		return makeNode(
			n_PREFIX_UNARY_EXPRESSION, op, value,
		)
	}

	if !p.acceptF(postfixUnaryExpression) {
		return INVALID_NODE
	}
	return p.takeNode()
}

func postfixUnaryExpression(p *parser) ast {
	if !p.acceptF(functionCallOrMemberExpression) {
		return INVALID_NODE
	}
	value := p.takeNode()

	if p.acceptT(tINC, tDEC) {
		return makeNode(n_POSTFIX_UNARY_EXPRESSION, p.takeLexeme(), value)
	}

	return value
}

func functionCallOrMemberExpression(p *parser) ast {
	if !p.acceptF(constructorCall) {
		return INVALID_NODE
	}
	funcName := p.takeNode()

	for {
		if !p.testFlag(f_ACCEPT_NO_FUNCTION_CALL) && p.acceptF(functionArgs) {
			call := makeNode(n_FUNCTION_CALL, "", funcName, p.takeNode())
			funcName = call
		} else {
			var property ast

			if p.acceptT(tDOT) {
				if p.acceptT(tNAME) {
					property = makeNode(n_PROPERTY_NAME, p.takeLexeme())
				} else if p.acceptT(tBRACKET_LEFT) {
					property = makeNode(
						n_CALCULATED_PROPERTY_NAME, "", p.expectF(expression),
					)
				}
			}

			if property.t == n_EMPTY {
				return funcName
			}

			funcName = makeNode(n_MEMBER_EXPRESSION, "", funcName, property)
		}
	}
}

func functionArgs(p *parser) ast {
	if !p.acceptT(tPAREN_LEFT) {
		return INVALID_NODE
	}

	args := []ast{}

	for !p.acceptT(tPAREN_RIGHT) {
		if p.acceptF(assignmentExpression) {
			args = append(args, p.takeNode())
		}

		if !p.acceptT(tCOMMA) {
			if !p.acceptT(tPAREN_RIGHT) {
				return INVALID_NODE
			}
			break
		}
	}

	return makeNode(n_FUNCTION_ARGS, "", args...)
}

func constructorCall(p *parser) ast {
	if p.acceptT(tNEW) {
		var name, args ast

		p.flags = f_ACCEPT_NO_FUNCTION_CALL
		name = p.expectF(functionCallOrMemberExpression)
		p.flags = 0

		if p.acceptF(functionArgs) {
			args = p.takeNode()
		}

		return makeNode(n_CONSTRUCTOR_CALL, "", name, args)
	}

	if !p.acceptF(atom) {
		return INVALID_NODE
	}
	return p.takeNode()
}

func atom(p *parser) ast {
	if p.acceptF(objectLiteral) ||
		p.acceptF(otherLiteral) ||
		p.acceptF(lambdaExpression) ||
		p.acceptF(parenExpression) ||
		p.acceptF(regexpLiteral) ||
		p.acceptF(arrayLiteral) ||
		p.acceptF(identifier) ||
		p.acceptF(functionExpression) {
		return p.takeNode()
	}

	return INVALID_NODE
}

func parenExpression(p *parser) ast {
	if !p.acceptT(tPAREN_LEFT) {
		return INVALID_NODE
	}
	if p.acceptF(sequenceExpression) {
		if !p.acceptT(tPAREN_RIGHT) {
			return INVALID_NODE
		}
		return makeNode(n_PAREN_EXPRESSION, "", p.takeNode())
	}
	return INVALID_NODE
}

func lambdaExpression(p *parser) ast {
	var args, body ast

	if p.acceptF(identifier) {
		args = makeNode(n_FUNCTION_ARGS, "", p.takeNode())
	} else if p.acceptF(functionParameters) {
		args = p.takeNode()
	}

	if !p.acceptT(tLAMBDA) {
		return INVALID_NODE
	}

	if p.acceptF(blockStatement) {
		body = p.takeNode()
	} else {
		body = p.expectF(assignmentExpression)
	}
	return makeNode(n_LAMBDA_EXPRESSION, "", args, body)
}

func regexpLiteral(p *parser) ast {
	if !p.acceptT(tDIV) {
		return INVALID_NODE
	}

	var body, flags ast

	bodyValue := ""
	for {
		if p.acceptT(tDIV) && p.tokens[p.i-2].tType != tESCAPE {
			break
		}
		p.tNext()
		bodyValue += p.takeLexeme()
	}

	body = makeNode(n_REGEXP_BODY, bodyValue)
	if p.acceptT(tNAME) {
		flags = makeNode(n_REGEXP_FLAGS, p.takeLexeme())
	}

	return makeNode(n_REGEXP_LITERAL, "", body, flags)
}

func arrayLiteral(p *parser) ast {
	if !p.acceptT(tBRACKET_LEFT) {
		return INVALID_NODE
	}

	items := []ast{}
	for !p.acceptT(tBRACKET_RIGHT) {
		if p.acceptF(assignmentExpression) {
			items = append(items, p.takeNode())
		}

		if !p.acceptT(tCOMMA) {
			if !p.acceptT(tBRACKET_RIGHT) {
				return INVALID_NODE
			}
			break
		}
	}

	return makeNode(n_ARRAY_LITERAL, "", items...)
}

func objectLiteral(p *parser) ast {
	if !p.acceptT(tCURLY_LEFT) {
		return INVALID_NODE
	}

	props := []ast{}
	for !p.acceptT(tCURLY_RIGHT) {
		if p.acceptF(objectPropertyName) {
			var key, value ast
			key = p.takeNode()

			if p.acceptT(tCOLON) {
				value = p.expectF(assignmentExpression)
			}

			props = append(props, makeNode(n_OBJECT_PROPERTY, "", key, value))
		}

		if !p.acceptT(tCOMMA) {
			if !p.acceptT(tCURLY_RIGHT) {
				return INVALID_NODE
			}
			break
		}
	}

	return makeNode(n_OBJECT_LITERAL, "", props...)
}

func otherLiteral(p *parser) ast {
	if p.acceptT(tSTRING) {
		return makeNode(n_STRING_LITERAL, p.takeLexeme())
	} else if p.acceptT(tNUMBER, tHEX) {
		return makeNode(n_NUMBER_LITERAL, p.takeLexeme())
	} else if p.acceptT(tTRUE, tFALSE) {
		return makeNode(n_BOOL_LITERAL, p.takeLexeme())
	} else if p.acceptT(tNULL) {
		return makeNode(n_NULL, p.takeLexeme())
	} else if p.acceptT(tUNDEFINED) {
		return makeNode(n_UNDEFINED, p.takeLexeme())
	}
	return INVALID_NODE
}

func identifier(p *parser) ast {
	if p.acceptT(tNAME, tTHIS) {
		return makeNode(n_IDENTIFIER, p.takeLexeme())
	}
	return INVALID_NODE
}

func arrayPattern(p *parser) ast {
	return INVALID_NODE
}

func objectPattern(p *parser) ast {
	if !p.acceptT(tCURLY_LEFT) {
		return INVALID_NODE
	}

	props := []ast{}
	for !p.acceptT(tCURLY_RIGHT) {
		if p.acceptF(objectPropertyName) {
			var key, value ast

			key = p.takeNode()
			if p.acceptT(tCOLON) {
				value = p.expectF(assignmentPattern)
			}

			props = append(props, makeNode(n_OBJECT_PROPERTY, "", key, value))
		}

		if !p.acceptT(tCOMMA) {
			if !p.acceptT(tCURLY_RIGHT) {
				return INVALID_NODE
			}
			break
		}
	}

	return makeNode(n_OBJECT_PATTERN, "", props...)
}

func objectPropertyName(p *parser) ast {
	if p.acceptF(identifier) ||
		p.acceptF(otherLiteral) {
		return p.takeNode()
	}

	return INVALID_NODE
}

func assignmentPattern(p *parser) ast {
	var left, right ast

	if p.acceptF(objectPattern) ||
		p.acceptF(identifier) {
		left = p.takeNode()
	}

	if p.acceptT(tASSIGN) {
		right = p.expectF(conditionalExpression)
		return makeNode(n_ASSIGNMENT_PATTERN, "", left, right)
	}

	return left
}

func functionParameter(p *parser) ast {
	if p.acceptT(tSPREAD) {
		if p.acceptF(objectPattern) ||
			p.acceptF(arrayPattern) ||
			p.acceptF(identifier) {
			return makeNode(n_REST_ELEMENT, "", p.takeNode())
		}
		return p.expectF(identifier)
	}

	if p.acceptF(assignmentPattern) {
		return p.takeNode()
	}

	return INVALID_NODE
}

func functionParameters(p *parser) ast {
	if !p.acceptT(tPAREN_LEFT) {
		return INVALID_NODE
	}

	paramsSlice := []ast{}
	for !p.acceptT(tPAREN_RIGHT) {
		if p.acceptF(functionParameter) {
			paramsSlice = append(paramsSlice, p.takeNode())
		}

		if !p.acceptT(tCOMMA) {
			if !p.acceptT(tPAREN_RIGHT) {
				return INVALID_NODE
			}
			break
		}
	}

	return makeNode(n_FUNCTION_PARAMETERS, "", paramsSlice...)
}

func blockStatement(p *parser) ast {
	if !p.acceptT(tCURLY_LEFT) {
		return INVALID_NODE
	}
	statements := []ast{}
	for !p.acceptT(tCURLY_RIGHT) {
		statements = append(statements, p.expectF(statement))
	}

	return makeNode(n_BLOCK_STATEMENT, "", statements...)
}

func functionExpression(p *parser) ast {
	if !p.acceptT(tFUNCTION) {
		return INVALID_NODE
	}

	var name, params, body ast
	if p.acceptF(identifier) {
		name = p.takeNode()
	}

	params = p.expectF(functionParameters)
	body = p.expectF(blockStatement)

	return makeNode(n_FUNCTION_EXPRESSION, "", name, params, body)
}

func functionDeclaration(p *parser) ast {
	if !p.acceptT(tFUNCTION) {
		return INVALID_NODE
	}

	var name, params, body ast

	name = p.expectF(identifier)
	params = p.expectF(functionParameters)
	body = p.expectF(blockStatement)

	return makeNode(n_FUNCTION_DECLARATION, "", name, params, body)
}

func whileStatement(p *parser) ast {
	if !(p.acceptT(tWHILE) && p.acceptT(tPAREN_LEFT)) {
		return INVALID_NODE
	}
	var cond, body ast
	cond = p.expectF(sequenceExpression)
	if !p.acceptT(tPAREN_RIGHT) {
		return INVALID_NODE
	}
	body = p.expectF(statement)

	return makeNode(n_WHILE_STATEMENT, "", cond, body)
}

func doWhileStatement(p *parser) ast {
	if !p.acceptT(tDO) {
		return INVALID_NODE
	}
	var cond, body ast
	body = p.expectF(statement)

	if !(p.acceptT(tWHILE) && p.acceptT(tPAREN_LEFT)) {
		return INVALID_NODE
	}
	cond = p.expectF(sequenceExpression)
	if !p.acceptT(tPAREN_RIGHT) {
		return INVALID_NODE
	}

	return makeNode(n_DO_WHILE_STATEMENT, "", cond, body)
}

func forInOfStatement(p *parser) ast {
	if !(p.acceptT(tFOR) && p.acceptT(tPAREN_LEFT)) {
		return INVALID_NODE
	}
	var left, right, body ast

	if p.acceptF(varDeclaration) ||
		p.acceptF(assignmentPattern) {
		left = p.takeNode()
	} else {
		return INVALID_NODE
	}

	p.tNext()
	kind := p.t.tType

	right = p.expectF(sequenceExpression)

	if !p.acceptT(tPAREN_RIGHT) {
		return INVALID_NODE
	}

	body = p.expectF(statement)

	if kind == tIN {
		return makeNode(n_FOR_IN_STATEMENT, "", left, right, body)
	}
	return makeNode(n_FOR_OF_STATEMENT, "", left, right, body)
}

func forStatement(p *parser) ast {
	if !(p.acceptT(tFOR) && p.acceptT(tPAREN_LEFT)) {
		return INVALID_NODE
	}
	var init, test, final, body ast

	if p.acceptT(tSEMI) {
		init = makeNode(n_EMPTY_EXPRESSION, "")
	} else if p.acceptF(varDeclaration) ||
		p.acceptF(sequenceExpression) {
		init = p.takeNode()

		if !p.acceptT(tSEMI) {
			return INVALID_NODE
		}
	} else {
		return INVALID_NODE
	}

	if p.acceptT(tSEMI) {
		test = makeNode(n_EMPTY_EXPRESSION, "")
	} else {
		test = p.expectF(sequenceExpression)

		if !p.acceptT(tSEMI) {
			return INVALID_NODE
		}
	}

	if p.acceptT(tPAREN_RIGHT) {
		final = makeNode(n_EMPTY_EXPRESSION, "")
	} else {
		final = p.expectF(sequenceExpression)

		if !p.acceptT(tPAREN_RIGHT) {
			return INVALID_NODE
		}
	}

	body = p.expectF(statement)
	return makeNode(n_FOR_STATEMENT, "", init, test, final, body)
}

func ifStatement(p *parser) ast {
	if !(p.acceptT(tIF) &&
		p.acceptT(tPAREN_LEFT)) {
		return INVALID_NODE
	}
	var cond, conseq, alt ast
	cond = p.expectF(sequenceExpression)
	if !(p.acceptT(tPAREN_RIGHT)) {
		return INVALID_NODE
	}
	conseq = p.expectF(statement)

	if p.acceptT(tELSE) {
		alt = p.expectF(statement)
	}

	return makeNode(n_IF_STATEMENT, "", cond, conseq, alt)
}

func exportStatement(p *parser) ast {
	if !p.acceptT(tEXPORT) {
		return INVALID_NODE
	}

	var vars, declaration, path ast
	vars.t = n_EXPORT_VARS
	declaration.t = n_EXPORT_DECLARATION
	path.t = n_EXPORT_PATH

	if p.acceptT(tDEFAULT) {
		var name, alias ast

		alias = makeNode(n_EXPORT_ALIAS, "default")
		if p.acceptF(functionExpression) {
			fe := p.takeNode()
			feName := fe.children[0]

			if feName.t != n_EMPTY {
				declaration = fe
				name = feName
			} else {
				name = fe
			}
		} else {
			p.expectF(sequenceExpression)
			name = p.takeNode()
		}

		ev := makeNode(n_EXPORT_VAR, "", name, alias)
		vars.children = append(vars.children, ev)
		if !p.acceptT(tSEMI) {
			return INVALID_NODE
		}

	} else if p.acceptT(tCURLY_LEFT) {
		for !p.acceptT(tCURLY_RIGHT) {
			if p.acceptF(identifier) {
				name := p.takeNode()
				alias := name

				if p.acceptT(tNAME) && p.takeLexeme() == "as" {
					p.tNext()
					alias = makeNode(n_EXPORT_ALIAS, p.takeLexeme())
				}
				ev := makeNode(n_EXPORT_VAR, "", name, alias)
				vars.children = append(vars.children, ev)
			}

			if !p.acceptT(tCOMMA) {
				if !p.acceptT(tCURLY_RIGHT) {
					return INVALID_NODE
				}
				break
			}
		}

		if p.acceptT(tNAME) && p.takeLexeme() == "from" {
			if !p.acceptT(tSTRING) {
				return INVALID_NODE
			}
			path = makeNode(n_EXPORT_PATH, p.takeLexeme())
		}
		if !p.acceptT(tSEMI) {
			return INVALID_NODE
		}

	} else if p.acceptF(varDeclaration) {
		declaration = p.takeNode()
		for _, d := range declaration.children[0].children {
			name := d.children[0]
			alias := d.children[0]
			ev := makeNode(n_EXPORT_VAR, "", name, alias)

			vars.children = append(vars.children, ev)
		}

	} else if p.acceptF(functionDeclaration) {
		fs := p.takeNode()
		declaration = fs
		name := fs.children[0]
		alias := name
		ev := makeNode(n_EXPORT_VAR, "", name, alias)
		vars.children = append(vars.children, ev)

	} else if p.acceptT(tMULT) {
		if !p.acceptT(tNAME) {
			return INVALID_NODE
		}
		if p.takeLexeme() != "from" {
			return INVALID_NODE
		}
		if !p.acceptT(tSTRING) {
			return INVALID_NODE
		}
		path = makeNode(n_EXPORT_PATH, p.takeLexeme())
		vars.t = n_EXPORT_ALL
		if !p.acceptT(tSEMI) {
			return INVALID_NODE
		}
	}

	result := makeNode(n_EXPORT_STATEMENT, "", vars, declaration, path)
	return result
}
