package jsLoader

import (
	"fmt"
)

type nodeType int

var nodeToString = []string{
	"n_EMPTY",
	"n_INVALID",
	"n_PROGRAM",
	"n_IDENTIFIER",
	"n_EXPRESSION_STATEMENT",
	"n_FUNCTION_DECLARATION",
	"n_BLOCK_STATEMENT",
	"n_FUNCTION_EXPRESSION",
	"n_FUNCTION_PARAMETERS",
	"n_VARIABLE_DECLARATION",
	"n_DECLARATOR",
	"n_STRING_LITERAL",
	"n_NUMBER_LITERAL",
	"n_BOOL_LITERAL",
	"n_UNDEFINED",
	"n_NULL",
	"n_OBJECT_LITERAL",
	"n_ARRAY_LITERAL",
	"n_REGEXP_LITERAL",
	"n_ASSIGNMENT_PATTERN",
	"n_OBJECT_PATTERN",
	"n_ARRAY_PATTERN",
	"n_OBJECT_PROPERTY",
	"n_REST_ELEMENT",
	"n_SEQUENCE_EXPRESSION",
	"n_ASSIGNMENT_EXPRESSION",
	"n_CONDITIONAL_EXPRESSION",
	"n_BINARY_EXPRESSION",
	"n_PREFIX_UNARY_EXPRESSION",
	"n_POSTFIX_UNARY_EXPRESSION",
	"n_MEMBER_EXPRESSION",
	"n_PROPERTY_NAME",
	"n_CALCULATED_PROPERTY_NAME",
	"n_FUNCTION_CALL",
	"n_FUNCTION_ARGS",
	"n_CONSTRUCTOR_CALL",
}

func (n nodeType) String() string {
	return nodeToString[int(n)]
}

const (
	n_EMPTY nodeType = iota
	n_INVALID
	n_PROGRAM
	n_IDENTIFIER
	n_EXPRESSION_STATEMENT
	n_FUNCTION_DECLARATION
	n_BLOCK_STATEMENT

	n_FUNCTION_EXPRESSION
	n_FUNCTION_PARAMETERS

	n_VARIABLE_DECLARATION
	n_DECLARATOR

	n_STRING_LITERAL
	n_NUMBER_LITERAL
	n_BOOL_LITERAL
	n_UNDEFINED
	n_NULL
	n_OBJECT_LITERAL
	n_ARRAY_LITERAL
	n_REGEXP_LITERAL

	n_ASSIGNMENT_PATTERN
	n_OBJECT_PATTERN
	n_ARRAY_PATTERN
	n_OBJECT_PROPERTY

	n_REST_ELEMENT

	n_SEQUENCE_EXPRESSION
	n_ASSIGNMENT_EXPRESSION
	n_CONDITIONAL_EXPRESSION
	n_BINARY_EXPRESSION
	n_PREFIX_UNARY_EXPRESSION
	n_POSTFIX_UNARY_EXPRESSION
	n_MEMBER_EXPRESSION
	n_PROPERTY_NAME
	n_CALCULATED_PROPERTY_NAME

	n_FUNCTION_CALL
	n_FUNCTION_ARGS

	n_CONSTRUCTOR_CALL
)

type ast struct {
	t        nodeType
	value    string
	children []ast
	flags    int
}

const (
	f_ACCEPT_NO_FUNCTION_CALL = 1 << 0
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
	return ast{t, value, children, 0}
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

func (p *parser) expectF(f parseFunc) {
	if !p.acceptF(f) {
		fmt.Printf("Parsing error, %v", p.t)
		panic("error")
	}
}

func program(p *parser) ast {
	statements := []ast{}
	for !p.acceptT(tEND_OF_INPUT) {
		p.expectF(statement)
		statements = append(statements, p.takeNode())
	}
	return makeNode(n_PROGRAM, "", statements...)
}

func statement(p *parser) ast {
	if p.acceptF(varDeclaration) ||
		p.acceptF(functionDeclaration) ||
		p.acceptF(expressionStatement) {
		return p.takeNode()
	}
	return INVALID_NODE
}

func expressionStatement(p *parser) ast {
	if !p.acceptF(sequenceExpression) {
		return INVALID_NODE
	}
	if !p.acceptT(tSEMI) {
		return INVALID_NODE
	}
	return makeNode(n_EXPRESSION_STATEMENT, "", p.takeNode())
}

func varDeclaration(p *parser) ast {
	if !p.acceptT(tVAR) {
		return INVALID_NODE
	}
	kind := p.takeLexeme()

	declarators := []ast{}
	for ok := true; ok; ok = p.acceptT(tCOMMA) {
		p.expectF(declarator)
		declarators = append(declarators, p.takeNode())
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
		p.expectF(assignmentExpression)
		value = p.takeNode()
	}

	return makeNode(n_DECLARATOR, "", name, value)
}

func expression(p *parser) ast {
	return INVALID_NODE
}

func sequenceExpression(p *parser) ast {
	p.expectF(assignmentExpression)
	first := p.takeNode()

	if p.acceptT(tCOMMA) {
		items := []ast{}

		for ok := true; ok; ok = p.acceptT(tCOMMA) {
			p.expectF(assignmentExpression)
			items = append(items, p.takeNode())
		}

		return makeNode(n_SEQUENCE_EXPRESSION, "", items...)
	}

	return first
}

func assignmentExpression(p *parser) ast {
	p.expectF(conditionalExpression)
	left := p.takeNode()

	if p.acceptT(
		tASSIGN, tPLUS_ASSIGN, tMINUS_ASSIGN, tMULT_ASSIGN, tDIV_ASSIGN,
		tBITWISE_SHIFT_LEFT_ASSIGN, tBITWISE_SHIFT_RIGHT_ASSIGN,
		tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN, tBITWISE_AND_ASSIGN,
		tBITWISE_OR_ASSIGN, tBITWISE_XOR_ASSIGN,
	) {
		p.expectF(assignmentExpression)
		right := p.takeNode()
		return makeNode(n_ASSIGNMENT_EXPRESSION, "", left, right)
	}

	return left
}

func conditionalExpression(p *parser) ast {
	p.expectF(binaryExpression)
	cond := p.takeNode()

	if p.acceptT(tQUESTION) {
		p.expectF(conditionalExpression)
		conseq := p.takeNode()
		if !p.acceptT(tCOLON) {
			return INVALID_NODE
		}
		p.expectF(conditionalExpression)
		altern := p.takeNode()

		return makeNode(n_CONDITIONAL_EXPRESSION, "", cond, conseq, altern)
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
		p.expectF(prefixUnaryExpression)
		outputStack = append(outputStack, p.takeNode())

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
		p.expectF(postfixUnaryExpression)
		value := p.takeNode()
		return makeNode(
			n_PREFIX_UNARY_EXPRESSION, op, value,
		)
	}

	p.expectF(postfixUnaryExpression)
	return p.takeNode()
}

func postfixUnaryExpression(p *parser) ast {
	p.expectF(functionCallOrMemberExpression)
	value := p.takeNode()

	if p.acceptT(tINC, tDEC) {
		return makeNode(n_POSTFIX_UNARY_EXPRESSION, p.takeLexeme(), value)
	}

	return value
}

func functionCallOrMemberExpression(p *parser) ast {
	p.expectF(constructorCall)
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
					p.expectF(expression)
					property = makeNode(n_CALCULATED_PROPERTY_NAME, "", p.takeNode())
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
		p.expectF(functionCallOrMemberExpression)
		p.flags = 0

		name = p.takeNode()
		if p.acceptF(functionArgs) {
			args = p.takeNode()
		}

		return makeNode(n_CONSTRUCTOR_CALL, "", name, args)
	}

	p.expectF(atom)
	return p.takeNode()
}

func atom(p *parser) ast {
	// if p.fAccept(identifier) ||
	// 	p.fAccept(regexpLiteral) ||
	// 	p.fAccept(hexLiteral) ||
	// 	p.fAccept(lambda) ||
	// 	p.fAccept(objectLiteral) ||
	// 	p.fAccept(functionExpression) ||
	// 	p.fAccept(arrayLiteral) ||
	// 	p.fAccept(otherLiteral) { // number null string undefined true false this
	// 	return p.takeNode()
	// }

	if p.acceptF(arrayLiteral) ||
		p.acceptF(objectLiteral) ||
		p.acceptF(otherLiteral) ||
		p.acceptF(identifier) {
		return p.takeNode()
	}

	return INVALID_NODE
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
				p.expectF(assignmentExpression)
				value = p.takeNode()
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
	} else if p.acceptT(tNUMBER) {
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
				p.expectF(assignmentPattern)
				value = p.takeNode()
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
		p.expectF(conditionalExpression)
		right = p.takeNode()
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
		p.expectF(identifier)
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
		p.expectF(statement)
		statements = append(statements, p.takeNode())
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

	p.expectF(functionParameters)
	params = p.takeNode()
	p.expectF(blockStatement)
	body = p.takeNode()

	return makeNode(n_FUNCTION_EXPRESSION, "", name, params, body)
}

func functionDeclaration(p *parser) ast {
	if !p.acceptT(tFUNCTION) {
		return INVALID_NODE
	}

	var name, params, body ast
	p.expectF(identifier)
	name = p.takeNode()

	p.expectF(functionParameters)
	params = p.takeNode()
	p.expectF(blockStatement)
	body = p.takeNode()

	return makeNode(n_FUNCTION_DECLARATION, "", name, params, body)
}
