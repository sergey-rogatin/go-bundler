package jsLoader

import "fmt"

var sourceTokens []token
var index int
var tok token

func parseTokens(localSrc []token) (astNode, error) {
	sourceTokens = localSrc
	index = 0
	tok = sourceTokens[index]

	return parseProgram()
}

const (
	p_UNEXPECTED_TOKEN = iota
	p_WRONG_ASSIGNMENT
	p_UNNAMED_FUNCTION
)

type parsingError struct {
	kind    int
	tok     token
	context []token
}

func (pe parsingError) Error() string {
	switch pe.kind {
	case p_UNEXPECTED_TOKEN:
		return fmt.Sprintf(
			"Unexpected token \"%s\" at %v:%v, %s",
			pe.tok.lexeme, pe.tok.line, pe.tok.column, pe.context,
		)

	case p_UNNAMED_FUNCTION:
		return fmt.Sprintf(
			"Unnamed function declaration at %v:%v",
			pe.tok.line, pe.tok.column,
		)

	case p_WRONG_ASSIGNMENT:
		return fmt.Sprintf(
			"Invalid left-hand side in assignment at %v:%v",
			pe.tok.line, pe.tok.column,
		)

	default:
		return "Unknown parser error"
	}
}

// helper functions for parsing

func next() {
	index++
	if index < len(sourceTokens) {
		tok = sourceTokens[index]
	}
}

func backtrack(backIndex int) {
	index = backIndex
	tok = sourceTokens[index]
}

func test(tTypes ...tokenType) bool {
	i := index
	for sourceTokens[i].tType == tNEWLINE {
		i++
	}
	for _, tType := range tTypes {
		if sourceTokens[i].tType == tType {
			return true
		}
	}
	return false
}

func getNoNewline() token {
	i := index
	for sourceTokens[i].tType == tNEWLINE {
		i++
	}
	return sourceTokens[i]
}

func accept(tTypes ...tokenType) bool {
	prevPos := index
	for tok.tType == tNEWLINE {
		next()
	}
	if tTypes == nil {
		next()
		return true
	}
	for _, tType := range tTypes {
		if tok.tType == tType {
			next()
			return true
		}
	}
	backtrack(prevPos)
	return false
}

func checkASI(tType tokenType) {
	if tType == tSEMI {
		if tok.tType == tNEWLINE {
			next()
			return
		}
		if tok.tType == tCURLY_RIGHT || tok.tType == tEND_OF_INPUT {
			return
		}
	}
	panic(parsingError{
		p_UNEXPECTED_TOKEN, tok,
		sourceTokens[index-5 : index+5],
	})
}

func expect(tType tokenType) {
	if accept(tType) {
		return
	}
	checkASI(tType)
}

func getLexeme() string {
	return sourceTokens[index-1].lexeme
}

func getToken() token {
	return sourceTokens[index-1]
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

const (
	f_FUNCTION_PARAM_REST = 1 << 0
	f_EXPORT_ALL          = 1 << 1
)

type astNode struct {
	t        grammarType
	value    string
	children []astNode
	flags    int
}

func (a astNode) String() string {
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

func makeNode(t grammarType, value string, children ...astNode) astNode {
	return astNode{t, value, children, 0}
}

func parseProgram() (program astNode, err error) {
	err = nil

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(parsingError); ok {
				err = e
			}
		}
	}()

	statements := []astNode{}
	for !accept(tEND_OF_INPUT) {
		statements = append(statements, parseStatement())
	}

	program = makeNode(g_PROGRAM_STATEMENT, "", statements...)
	return
}

func parseTryCatchStatement() astNode {
	var try, catch, finally, catchValue astNode

	expect(tCURLY_LEFT)
	try = parseBlockStatement()
	if accept(tCATCH) {
		expect(tPAREN_LEFT)
		catchValue = parseExpression()
		expect(tPAREN_RIGHT)

		expect(tCURLY_LEFT)
		catch = parseBlockStatement()
	}
	if accept(tFINALLY) {
		expect(tCURLY_LEFT)
		finally = parseBlockStatement()
	}

	return makeNode(
		g_TRY_CATCH_STATEMENT, "", try, catchValue, catch, finally,
	)
}

func parseStatement() astNode {
	switch {
	case accept(tTHROW):
		return makeNode(g_THROW_STATEMENT, "", parseExpression())
	case accept(tTRY):
		return parseTryCatchStatement()
	case accept(tVAR, tCONST, tLET):
		return parseDeclarationStatement()
	case accept(tRETURN):
		return parseReturnStatement()
	case accept(tIMPORT):
		return parseImportStatement()
	case accept(tFUNCTION):
		return parseFunctionDeclaration()
	case accept(tDO):
		return parseDoWhileStatement()
	case accept(tWHILE):
		return parseWhileStatement()
	case accept(tFOR):
		return parseForStatement()
	case accept(tCURLY_LEFT):
		return parseBlockStatement()
	case accept(tIF):
		return parseIfStatement()
	case accept(tEXPORT):
		return parseExportStatement()
	case accept(tSWITCH):
		return parseSwitchStatement()
	case accept(tSEMI):
		return makeNode(g_EMPTY_STATEMENT, ";")
	default:
		return parseExpressionStatement()
	}
}

func parseExportStatement() astNode {
	declaration := makeNode(g_EXPORT_DECLARATION, "")
	path := makeNode(g_EXPORT_PATH, "")
	vars := []astNode{}
	flags := 0

	if accept(tDEFAULT) {

		var name astNode
		alias := makeNode(g_EXPORT_ALIAS, "default")

		if accept(tFUNCTION) {
			fe := parseFunctionExpression()
			if fe.value != "" {
				declaration = fe
				name = makeNode(g_EXPORT_NAME, fe.value)
			} else {
				name = fe
			}
		} else {
			name = parseExpression()
		}

		ev := makeNode(g_EXPORT_VAR, "", name, alias)
		vars = append(vars, ev)
		expect(tSEMI)

	} else if accept(tCURLY_LEFT) {
		for !accept(tCURLY_RIGHT) {
			if accept(tNAME) {
				name := makeNode(g_EXPORT_NAME, getLexeme())
				alias := name

				if accept(tAS) {
					next()
					alias = makeNode(g_EXPORT_ALIAS, getLexeme())
				}
				ev := makeNode(g_EXPORT_VAR, "", name, alias)
				vars = append(vars, ev)
			}

			if !accept(tCOMMA) {
				expect(tCURLY_RIGHT)
				break
			}
		}

		if accept(tFROM) {
			expect(tSTRING)
			path = makeNode(g_EXPORT_PATH, getLexeme())
		}
		expect(tSEMI)

	} else if accept(tVAR, tLET, tCONST) {
		declaration = parseDeclarationStatement()
		for _, d := range declaration.children[0].children {
			name := d.children[0]
			alias := d.children[0]
			ev := makeNode(g_EXPORT_VAR, "", name, alias)

			vars = append(vars, ev)
		}

	} else if accept(tFUNCTION) {
		fs := parseFunctionDeclaration()
		declaration = fs
		name := makeNode(g_EXPORT_NAME, fs.value)
		alias := name
		ev := makeNode(g_EXPORT_VAR, "", name, alias)
		vars = append(vars, ev)

	} else if accept(tMULT) {
		expect(tFROM)
		expect(tSTRING)
		path = makeNode(g_EXPORT_PATH, getLexeme())
		flags = flags | f_EXPORT_ALL
		expect(tSEMI)
	}
	varsNode := makeNode(g_EXPORT_VARS, "", vars...)

	result := makeNode(g_EXPORT_STATEMENT, "", varsNode, declaration, path)
	result.flags = flags
	return result
}

func parseSwitchStatement() astNode {
	expect(tPAREN_LEFT)
	condition := parseExpression()
	expect(tPAREN_RIGHT)

	cases := []astNode{}
	expect(tCURLY_LEFT)
	for !accept(tCURLY_RIGHT) {
		if accept(tCASE) {
			caseTest := makeNode(g_SWITCH_CASE_TEST, "", parseExpression())
			expect(tCOLON)
			caseStatements := []astNode{}
			for !test(tDEFAULT, tCASE, tCURLY_RIGHT) {
				caseStatements = append(caseStatements, parseStatement())
			}
			statementNode := makeNode(
				g_SWITCH_CASE_STATEMENTS, "", caseStatements...,
			)
			switchCase := makeNode(
				g_SWITCH_CASE, "", caseTest, statementNode,
			)
			cases = append(cases, switchCase)
		}

		if accept(tDEFAULT) {
			expect(tCOLON)
			caseStatements := []astNode{}
			for !test(tDEFAULT, tCASE, tCURLY_RIGHT) {
				caseStatements = append(caseStatements, parseStatement())
			}
			defaultCase := makeNode(g_SWITCH_DEFAULT, "", caseStatements...)
			cases = append(cases, defaultCase)
		}
	}

	casesNode := makeNode(g_SWITCH_CASES, "", cases...)
	return makeNode(g_SWITCH_STATEMENT, "", condition, casesNode)
}

func parseDeclarationStatement() astNode {
	decl := makeNode(
		g_DECLARATION_STATEMENT, "", parseDeclarationExpression(),
	)
	expect(tSEMI)
	return decl
}

func parseDeclarationExpression() astNode {
	keyword := getLexeme()

	ve := makeNode(g_DECLARATION_EXPRESSION, keyword, []astNode{}...)
	for ok := true; ok; ok = accept(tCOMMA) {
		ve.children = append(ve.children, parseDeclarator())
	}

	return ve
}

func parseForStatement() astNode {
	expect(tPAREN_LEFT)
	var init astNode
	if accept(tVAR, tLET, tCONST) {
		init = parseDeclarationExpression()
	} else if test(tSEMI) {
		init = makeNode(g_EMPTY_EXPRESSION, "")
	} else {
		init = parseSequence()
	}

	if accept(tOF) {
		left := init
		right := parseExpression()
		expect(tPAREN_RIGHT)
		body := parseStatement()
		return makeNode(g_FOR_OF_STATEMENT, "", left, right, body)
	} else if accept(tIN) {
		left := init
		right := parseExpression()
		expect(tPAREN_RIGHT)
		body := parseStatement()
		return makeNode(g_FOR_IN_STATEMENT, "", left, right, body)
	} else {
		var test, final astNode

		expect(tSEMI)
		if accept(tSEMI) {
			test = makeNode(g_EMPTY_EXPRESSION, "")
		} else {
			test = parseExpression()
			expect(tSEMI)
		}

		if accept(tPAREN_RIGHT) {
			final = makeNode(g_EMPTY_EXPRESSION, "")
		} else {
			final = parseExpression()
			expect(tPAREN_RIGHT)
		}

		body := parseStatement()

		return makeNode(g_FOR_STATEMENT, "", init, test, final, body)
	}
}

func parseIfStatement() astNode {
	expect(tPAREN_LEFT)
	test := parseExpression()
	expect(tPAREN_RIGHT)
	body := parseStatement()
	if accept(tELSE) {
		alternate := parseStatement()
		return makeNode(g_IF_STATEMENT, "", test, body, alternate)
	}

	return makeNode(g_IF_STATEMENT, "", test, body)
}

func parseWhileStatement() astNode {
	expect(tPAREN_LEFT)
	test := parseExpression()
	expect(tPAREN_RIGHT)
	body := parseStatement()

	return makeNode(g_WHILE_STATEMENT, "", test, body)
}

func parseDoWhileStatement() astNode {
	body := parseStatement()
	expect(tWHILE)
	expect(tPAREN_LEFT)
	test := parseExpression()
	expect(tPAREN_RIGHT)

	return makeNode(g_DO_WHILE_STATEMENT, "", test, body)
}

func parseFunctionDeclaration() astNode {
	expect(tNAME)
	name := getLexeme()

	expect(tPAREN_LEFT)
	params := parseFunctionParameters()
	expect(tCURLY_LEFT)
	body := parseBlockStatement()

	return makeNode(g_FUNCTION_DECLARATION, name, params, body)
}

func parseExpressionStatement() astNode {
	if accept(tBREAK) {
		return makeNode(g_BREAK_STATEMENT, "")
	} else if accept(tCONTINUE) {
		return makeNode(g_CONTINUE_STATEMENT, "")
	} else if accept(tDEBUGGER) {
		return makeNode(g_DEBUGGER_STATEMENT, "")
	}
	expr := makeNode(g_EXPRESSION_STATEMENT, "", parseExpression())
	expect(tSEMI)
	return expr
}

func parseReturnStatement() astNode {
	if accept(tSEMI) {
		return makeNode(g_RETURN_STATEMENT, "")
	}
	expr := parseExpression()
	return makeNode(g_RETURN_STATEMENT, "", expr)
}

func parseExpression() astNode {
	return parseSequence()
}

func parseSequence() astNode {
	firstItem := parseYield()

	children := []astNode{firstItem}
	for accept(tCOMMA) {
		children = append(children, parseYield())
	}

	if len(children) > 1 {
		return makeNode(g_SEQUENCE_EXPRESSION, ",", children...)
	}
	return firstItem
}

func parseYield() astNode {
	if accept(tYIELD) {
		return makeNode(g_EXPRESSION, "yield", parseYield())
	}
	return parseAssignment()
}

func parseAssignment() astNode {
	// if accept(tCURLY_LEFT) {
	// 	left := parseObjectPattern()

	// 	if accept(tASSIGN) {
	// 		right := parseAssignment()
	// 		return makeNode(g_EXPRESSION, "=", left, right)
	// 	}
	// }

	left := parseConditional()

	if accept(
		tASSIGN, tPLUS_ASSIGN, tMINUS_ASSIGN, tMULT_ASSIGN, tDIV_ASSIGN,
		tBITWISE_SHIFT_LEFT_ASSIGN, tBITWISE_SHIFT_RIGHT_ASSIGN,
		tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN, tBITWISE_AND_ASSIGN,
		tBITWISE_OR_ASSIGN, tBITWISE_XOR_ASSIGN,
	) {
		op := getLexeme()
		right := parseAssignment()
		return makeNode(g_EXPRESSION, op, left, right)
	}

	return left
}

func parseAssignmentPattern() astNode {
	var left astNode
	if accept(tCURLY_LEFT) {
		left = parseObjectPattern()
	} else if accept(tNAME) {
		left = makeNode(g_NAME, getLexeme())
	} else {
		checkASI(tSEMI)
	}

	if accept(tASSIGN) {
		right := parseExpression()
		return makeNode(g_ASSIGNMENT_PATTERN, "=", left, right)
	}

	return left
}

func parseObjectPattern() astNode {
	properties := []astNode{}
	for !accept(tCURLY_RIGHT) {
		if accept(tNAME) {
			prop := makeNode(g_OBJECT_PROPERTY, "", []astNode{}...)
			key := makeNode(g_NAME, getLexeme())
			prop.children = append(prop.children, key)

			if accept(tCOLON) {
				prop.children = append(prop.children, parseAssignmentPattern())
			}

			properties = append(properties, prop)
		}

		if !accept(tCOMMA) {
			expect(tCURLY_RIGHT)
			break
		}
	}

	return makeNode(g_OBJECT_PATTERN, "", properties...)
}

func parseConditional() astNode {
	test := parseBinary()

	if accept(tQUESTION) {
		consequent := parseConditional()
		expect(tCOLON)
		alternate := parseConditional()
		return makeNode(g_CONDITIONAL_EXPRESSION, "?", test, consequent, alternate)
	}

	return test
}

func parseBinary() astNode {
	opStack := make([]token, 0)
	outputStack := make([]astNode, 0)
	var root *astNode

	addNode := func(t token) {
		right := outputStack[len(outputStack)-1]
		left := outputStack[len(outputStack)-2]
		nn := makeNode(g_EXPRESSION, t.lexeme, left, right)

		outputStack = outputStack[:len(outputStack)-2]
		outputStack = append(outputStack, nn)
		root = &nn
	}

	for {
		outputStack = append(outputStack, parsePrefixUnary())

		op, ok := operatorTable[tok.tType]
		if !ok {
			break
		}

		for {
			if len(opStack) > 0 {
				opToken := opStack[len(opStack)-1]

				opFromStack := operatorTable[opToken.tType]
				if opFromStack.precedence > op.precedence || (opFromStack.precedence == op.precedence && !op.isRightAssociative) {
					addNode(opToken)
					opStack = opStack[:len(opStack)-1]
				} else {
					break
				}
			} else {
				break
			}
		}
		opStack = append(opStack, tok)
		next()
	}

	for i := len(opStack) - 1; i >= 0; i-- {
		addNode(opStack[i])
	}

	if root != nil {
		return *root
	}
	return outputStack[len(outputStack)-1]
}

func parsePrefixUnary() astNode {
	if accept(
		tNOT, tBITWISE_NOT, tPLUS, tMINUS,
		tINC, tDEC, tTYPEOF, tVOID, tDELETE,
	) {
		return makeNode(
			g_UNARY_PREFIX_EXPRESSION, getLexeme(), parsePostfixUnary(),
		)
	}

	return parsePostfixUnary()
}

func parsePostfixUnary() astNode {
	value := parseFunctionCallOrMember(false)

	if accept(tINC, tDEC) {
		return makeNode(g_UNARY_POSTFIX_EXPRESSION, getLexeme(), value)
	}

	return value
}

func parseFunctionArgs() astNode {
	args := []astNode{}

	for !accept(tPAREN_RIGHT) {
		args = append(args, parseExpression())

		if !accept(tCOMMA) {
			expect(tPAREN_RIGHT)
			break
		}
	}

	argsNode := makeNode(g_FUNCTION_ARGS, "", args...)

	return argsNode
}

func parseFunctionCallOrMember(onlyMember bool) astNode {
	funcName := parseConstructorCall()

	for {
		if !onlyMember && accept(tPAREN_LEFT) {
			argsNode := parseFunctionArgs()
			n := makeNode(g_FUNCTION_CALL, "", funcName, argsNode)

			funcName = n
		} else {
			var property astNode

			if accept(tDOT) {
				expect(tNAME)
				property = makeNode(g_NAME, getLexeme())
			} else if accept(tBRACKET_LEFT) {
				property = parseCalculatedPropertyName()
			} else {
				break
			}
			object := funcName

			me := makeNode(g_MEMBER_EXPRESSION, "", object, property)
			funcName = me
		}
	}

	return funcName
}

func parseConstructorCall() astNode {
	if accept(tNEW) {
		name := parseFunctionCallOrMember(true)
		if accept(tPAREN_LEFT) {
			return makeNode(
				g_CONSTRUCTOR_CALL, "", name, parseFunctionArgs(),
			)
		}
		return makeNode(g_CONSTRUCTOR_CALL, "", name)
	}

	return parseAtom()
}

func parseAtom() astNode {
	switch {
	case accept(tDIV):
		return parseRegexp()
	case accept(tHEX):
		return makeNode(g_HEX_LITERAL, getLexeme())
	case accept(tPAREN_LEFT):
		return parseParensOrLambda()
	case accept(tCURLY_LEFT):
		return parseObjectLiteral()
	case accept(tFUNCTION):
		return parseFunctionExpression()
	case accept(tBRACKET_LEFT):
		return parseArrayLiteral()
	case accept(tNAME):
		return parseLambdaOrName()
	case accept(tNUMBER):
		return makeNode(g_NUMBER_LITERAL, getLexeme())
	case accept(tSTRING):
		return makeNode(g_STRING_LITERAL, getLexeme())
	case accept(tNULL):
		return makeNode(g_NULL, getLexeme())
	case accept(tUNDEFINED):
		return makeNode(g_UNDEFINED, getLexeme())
	case accept(tTRUE, tFALSE):
		return makeNode(g_BOOL_LITERAL, getLexeme())
	case accept(tTHIS):
		return makeNode(g_THIS, getLexeme())

	default:
		checkASI(tSEMI)
	}
	return astNode{}
}

func parseRegexp() astNode {
	value := "/"
	for {
		if accept(tDIV) {
			if sourceTokens[index-2].tType != tESCAPE {
				break
			} else {
				value += getLexeme()
			}
		}
		next()
		value += getLexeme()
	}
	value += "/"
	if accept(tNAME) {
		value += getLexeme()
	}
	fmt.Println(value)
	return makeNode(g_REGEXP_LITERAL, value)
}

func parseParensOrLambda() astNode {
	prevPos := index

	params := parseFunctionParameters()

	if accept(tLAMBDA) {
		var body astNode
		if accept(tCURLY_LEFT) {
			body = parseBlockStatement()
		} else {
			body = parseYield()
		}
		return makeNode(g_LAMBDA_EXPRESSION, "", params, body)
	}

	backtrack(prevPos)

	var value astNode
	if !accept(tPAREN_RIGHT) {
		value = parseExpression()
		expect(tPAREN_RIGHT)
	}

	return makeNode(g_PARENS_EXPRESSION, "", value)
}

func parseFunctionParameters() astNode {
	startPos := index

	result := astNode{g_FUNCTION_PARAMETERS, "", nil, 0}
	params := []astNode{}
	for !accept(tPAREN_RIGHT) {
		params = append(params, parseFunctionParameter())

		if !accept(tCOMMA) {
			if !accept(tPAREN_RIGHT) {
				backtrack(startPos)
				return result
			}
			break
		}
	}
	result.children = params
	return result
}

func parseFunctionExpression() astNode {
	name := ""
	if accept(tNAME) {
		name = getLexeme()
	}

	expect(tPAREN_LEFT)
	params := parseFunctionParameters()
	expect(tCURLY_LEFT)
	body := parseBlockStatement()

	return makeNode(g_FUNCTION_EXPRESSION, name, params, body)
}

func parseObjectLiteral() astNode {
	props := []astNode{}

	for !accept(tCURLY_RIGHT) {
		prop := astNode{g_OBJECT_PROPERTY, "", []astNode{}, 0}
		var key, value astNode

		if tok.lexeme == "get" {
			next()
			prop.value = getLexeme()
		} else if tok.lexeme == "set" {
			next()
			prop.value = getLexeme()
		}

		if accept(tNAME) {
			key = makeNode(g_NAME, getLexeme())
		} else if accept(tBRACKET_LEFT) {
			key = parseCalculatedPropertyName()
		} else if isValidPropertyName(getNoNewline().lexeme) || test(tNUMBER, tSTRING) {
			accept()
			key = makeNode(g_VALID_PROPERTY_NAME, getLexeme())
		}
		prop.children = append(prop.children, key)

		if test(tPAREN_LEFT) {
			value = parseMemberFunction()
			prop.children = append(prop.children, value)
		} else if accept(tCOLON) {
			value = parseYield()
			prop.children = append(prop.children, value)
		}

		props = append(props, prop)
		if !accept(tCOMMA) {
			expect(tCURLY_RIGHT)
			break
		}
	}

	return makeNode(g_OBJECT_LITERAL, "", props...)
}

func parseMemberFunction() astNode {
	f := parseFunctionExpression()
	f.t = g_MEMBER_FUNCTION
	return f
}

func parseArrayLiteral() astNode {
	values := []astNode{}

	for !accept(tBRACKET_RIGHT) {
		values = append(values, parseYield())

		if !accept(tCOMMA) {
			expect(tBRACKET_RIGHT)
			break
		}
	}

	return makeNode(g_ARRAY_LITERAL, "", values...)
}

func parseLambdaOrName() astNode {
	firstParamStr := getLexeme()

	if accept(tLAMBDA) {
		var body astNode
		if accept(tCURLY_LEFT) {
			body = parseBlockStatement()
		} else {
			body = parseYield()
		}
		params := makeNode(
			g_FUNCTION_PARAMETERS, "",
			makeNode(
				g_FUNCTION_PARAMETER, "", makeNode(g_NAME, firstParamStr),
			),
		)
		return makeNode(g_LAMBDA_EXPRESSION, "", params, body)
	}

	return makeNode(g_NAME, firstParamStr)
}

func parseBlockStatement() astNode {
	statements := []astNode{}
	for !accept(tCURLY_RIGHT) {
		statements = append(statements, parseStatement())
	}
	return makeNode(g_BLOCK_STATEMENT, "", statements...)
}

func parseFunctionParameter() astNode {
	if accept(tSPREAD) {
		var left astNode
		if accept(tCURLY_LEFT) {
			left = parseObjectPattern()
		} else if accept(tNAME) {
			left = makeNode(g_NAME, getLexeme())
		}
		n := makeNode(g_FUNCTION_PARAMETER, "", left)
		n.flags = f_FUNCTION_PARAM_REST
		return n
	}

	n := parseDeclarator()
	n.t = g_FUNCTION_PARAMETER
	return n
}

func parseDeclarator() astNode {
	var left astNode
	if accept(tCURLY_LEFT) {
		left = parseObjectPattern()
	} else if accept(tNAME, tFROM) {
		left = makeNode(g_NAME, getLexeme())
	} else {
		return left
		//checkASI(tSEMI)
	}

	if accept(tASSIGN) {
		right := parseAssignment()
		return makeNode(g_DECLARATOR, "", left, right)
	}

	return makeNode(g_DECLARATOR, "", left)
}

func parseCalculatedPropertyName() astNode {
	value := parseExpression()
	expect(tBRACKET_RIGHT)
	return makeNode(g_CALCULATED_PROPERTY_NAME, "", value)
}

func parseImportStatement() astNode {
	vars := astNode{g_IMPORT_VARS, "", []astNode{}, 0}
	all := astNode{t: g_IMPORT_ALL}
	path := astNode{t: g_IMPORT_PATH}

	if accept(tMULT) {
		expect(tAS)
		expect(tNAME)
		all.value = getLexeme()
		expect(tFROM)
	} else {
		if accept(tNAME) {
			name := makeNode(g_IMPORT_NAME, "default")
			alias := makeNode(g_IMPORT_ALIAS, getLexeme())

			varNode := makeNode(g_IMPORT_VAR, "", name, alias)
			vars.children = append(vars.children, varNode)

			if accept(tCOMMA) {
				if accept(tMULT) {
					expect(tAS)
					expect(tNAME)
					all.value = getLexeme()
					expect(tFROM)
				}
			} else {
				expect(tFROM)
			}
		}

		if accept(tCURLY_LEFT) {
			for !accept(tCURLY_RIGHT) {
				if accept(tNAME, tDEFAULT) {
					name := makeNode(g_IMPORT_NAME, getLexeme())
					alias := makeNode(g_IMPORT_ALIAS, getLexeme())

					if accept(tAS) {
						expect(tNAME)
						alias = makeNode(g_IMPORT_ALIAS, getLexeme())
					}

					varNode := makeNode(g_IMPORT_VAR, "", name, alias)
					vars.children = append(vars.children, varNode)
				}

				if !accept(tCOMMA) {
					expect(tCURLY_RIGHT)
					break
				}
			}

			expect(tFROM)
		}
	}

	expect(tSTRING)
	path.value = getLexeme()
	expect(tSEMI)

	return makeNode(g_IMPORT_STATEMENT, "", vars, all, path)
}
