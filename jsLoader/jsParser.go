package jsLoader

// type parserState struct {
// 	sourceTokens                 []token
// 	index                        int
// 	tok                          token
// 	globalFlagIsInsideForIn      bool
// 	globalFlagPatternParsingFail bool
// }

// func parseTokens(localSrc []token) (ast, error) {
// 	p := parserState{
// 		sourceTokens: localSrc,
// 		index:        0,
// 		tok:          localSrc[0],
// 	}

// 	return parseProgram(&p)
// }

// const (
// 	p_UNEXPECTED_TOKEN = iota
// 	p_WRONG_ASSIGNMENT
// 	p_UNNAMED_FUNCTION
// )

// type parsingError struct {
// 	kind    int
// 	tok     token
// 	context []token
// }

// func (pe parsingError) Error() string {
// 	switch pe.kind {
// 	case p_UNEXPECTED_TOKEN:
// 		return fmt.Sprintf(
// 			"Unexpected token \"%s\" at %v:%v, %s",
// 			pe.tok.lexeme, pe.tok.line, pe.tok.column, pe.context,
// 		)

// 	case p_UNNAMED_FUNCTION:
// 		return fmt.Sprintf(
// 			"Unnamed function declaration at %v:%v",
// 			pe.tok.line, pe.tok.column,
// 		)

// 	case p_WRONG_ASSIGNMENT:
// 		return fmt.Sprintf(
// 			"Invalid left-hand side in assignment at %v:%v",
// 			pe.tok.line, pe.tok.column,
// 		)

// 	default:
// 		return "Unknown parser error"
// 	}
// }

// // helper functions for parsing

// func next(p *parserState) {
// 	p.index++
// 	if p.index < len(p.sourceTokens) {
// 		p.tok = p.sourceTokens[p.index]
// 	}
// }

// func backtrack(p *parserState, backIndex int) {
// 	p.index = backIndex
// 	p.tok = p.sourceTokens[p.index]
// }

// func test(p *parserState, tTypes ...tokenType) bool {
// 	i := p.index
// 	for p.sourceTokens[i].tType == tNEWLINE {
// 		i++
// 	}
// 	for _, tType := range tTypes {
// 		if p.sourceTokens[i].tType == tType {
// 			return true
// 		}
// 	}
// 	return false
// }

// func getNoNewline(p *parserState) token {
// 	i := p.index
// 	for p.sourceTokens[i].tType == tNEWLINE {
// 		i++
// 	}
// 	return p.sourceTokens[i]
// }

// func accept(p *parserState, tTypes ...tokenType) bool {
// 	prevPos := p.index
// 	for p.tok.tType == tNEWLINE {
// 		next(p)
// 	}
// 	if tTypes == nil {
// 		next(p)
// 		return true
// 	}
// 	for _, tType := range tTypes {
// 		if p.tok.tType == tType {
// 			next(p)
// 			return true
// 		}
// 	}
// 	backtrack(p, prevPos)
// 	return false
// }

// func checkASI(p *parserState, tType tokenType) {
// 	if tType == tSEMI {
// 		if p.tok.tType == tNEWLINE {
// 			next(p)
// 			return
// 		}
// 		if p.tok.tType == tCURLY_RIGHT || p.tok.tType == tEND_OF_INPUT {
// 			return
// 		}
// 	}
// 	panic(parsingError{
// 		p_UNEXPECTED_TOKEN, p.tok,
// 		nil,
// 	})
// }

// func expect(p *parserState, tType tokenType) {
// 	if accept(p, tType) {
// 		return
// 	}
// 	checkASI(p, tType)
// }

// func getLexeme(p *parserState) string {
// 	return p.sourceTokens[p.index-1].lexeme
// }

// func getToken(p *parserState) token {
// 	return p.sourceTokens[p.index-1]
// }

// func isValidPropertyName(name string) bool {
// 	if len(name) == 0 {
// 		return false
// 	}
// 	if !isLetter(name[0]) {
// 		return false
// 	}
// 	for _, symbol := range name {
// 		if !isLetter(byte(symbol)) && !isNumber(byte(symbol)) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func parseProgram(p *parserState) (program ast, err error) {
// 	err = nil

// 	defer func() {
// 		if r := recover(); r != nil {
// 			if e, ok := r.(parsingError); ok {
// 				err = e
// 			}
// 		}
// 	}()

// 	statements := []ast{}
// 	for !accept(p, tEND_OF_INPUT) {
// 		statements = append(statements, parseStatement(p))
// 	}

// 	program = makeNode(g_PROGRAM_STATEMENT, "", statements...)
// 	return
// }

// func parseTryCatchStatement(p *parserState) ast {
// 	var try, catch, finally, catchValue ast

// 	expect(p, tCURLY_LEFT)
// 	try = parseBlockStatement(p)
// 	if accept(p, tCATCH) {
// 		expect(p, tPAREN_LEFT)
// 		catchValue = parseExpression(p)
// 		expect(p, tPAREN_RIGHT)

// 		expect(p, tCURLY_LEFT)
// 		catch = parseBlockStatement(p)
// 	}
// 	if accept(p, tFINALLY) {
// 		expect(p, tCURLY_LEFT)
// 		finally = parseBlockStatement(p)
// 	}

// 	return makeNode(
// 		g_TRY_CATCH_STATEMENT, "", try, catchValue, catch, finally,
// 	)
// }

// func parseStatement(p *parserState) ast {
// 	startPos := p.index
// 	if accept(p, tNAME) {
// 		markerName := getLexeme(p)
// 		if accept(p, tCOLON) {
// 			return makeNode(g_MARKER, markerName)
// 		}
// 		backtrack(p, startPos)
// 	}

// 	switch {
// 	case accept(p, tCLASS):
// 		return makeNode(g_CLASS_STATEMENT, "", parseClassExpression(p))
// 	case accept(p, tTHROW):
// 		return makeNode(g_THROW_STATEMENT, "", parseExpression(p))
// 	case accept(p, tTRY):
// 		return parseTryCatchStatement(p)
// 	case accept(p, tVAR, tCONST, tLET):
// 		return parseDeclarationStatement(p)
// 	case accept(p, tRETURN):
// 		return parseReturnStatement(p)
// 	case accept(p, tIMPORT):
// 		return parseImportStatement(p)
// 	case accept(p, tFUNCTION):
// 		return parseFunctionDeclaration(p)
// 	case accept(p, tDO):
// 		return parseDoWhileStatement(p)
// 	case accept(p, tWHILE):
// 		return parseWhileStatement(p)
// 	case accept(p, tFOR):
// 		return parseForStatement(p)
// 	case accept(p, tCURLY_LEFT):
// 		return parseBlockStatement(p)
// 	case accept(p, tIF):
// 		return parseIfStatement(p)
// 	case accept(p, tEXPORT):
// 		return parseExportStatement(p)
// 	case accept(p, tSWITCH):
// 		return parseSwitchStatement(p)
// 	case accept(p, tSEMI):
// 		return makeNode(g_EMPTY_STATEMENT, ";")
// 	default:
// 		return parseExpressionStatement(p)
// 	}
// }

// func parseExportStatement(p *parserState) ast {
// 	declaration := makeNode(g_EXPORT_DECLARATION, "")
// 	path := makeNode(g_EXPORT_PATH, "")
// 	vars := []ast{}
// 	flags := 0

// 	if accept(p, tDEFAULT) {

// 		var name ast
// 		alias := makeNode(g_EXPORT_ALIAS, "default")

// 		if accept(p, tCLASS) {
// 			ce := parseClassExpression(p)
// 			if ce.children[0].value != "" {
// 				declaration = ce
// 				name = ce.children[0]
// 			} else {
// 				name = ce
// 			}
// 		} else if accept(p, tFUNCTION) {
// 			fe := parseFunctionExpression(p)
// 			if fe.value != "" {
// 				declaration = fe
// 				name = makeNode(g_EXPORT_NAME, fe.value)
// 			} else {
// 				name = fe
// 			}
// 		} else {
// 			name = parseExpression(p)
// 		}

// 		ev := makeNode(g_EXPORT_VAR, "", name, alias)
// 		vars = append(vars, ev)
// 		expect(p, tSEMI)

// 	} else if accept(p, tCURLY_LEFT) {
// 		for !accept(p, tCURLY_RIGHT) {
// 			if accept(p, tNAME) {
// 				name := makeNode(g_EXPORT_NAME, getLexeme(p))
// 				alias := name

// 				if accept(p, tNAME) && getLexeme(p) == "as" {
// 					next(p)
// 					alias = makeNode(g_EXPORT_ALIAS, getLexeme(p))
// 				}
// 				ev := makeNode(g_EXPORT_VAR, "", name, alias)
// 				vars = append(vars, ev)
// 			}

// 			if !accept(p, tCOMMA) {
// 				expect(p, tCURLY_RIGHT)
// 				break
// 			}
// 		}

// 		if accept(p, tNAME) && getLexeme(p) == "from" {
// 			expect(p, tSTRING)
// 			path = makeNode(g_EXPORT_PATH, getLexeme(p))
// 		}
// 		expect(p, tSEMI)

// 	} else if accept(p, tVAR, tLET, tCONST) {
// 		declaration = parseDeclarationStatement(p)
// 		for _, d := range declaration.children[0].children {
// 			name := d.children[0]
// 			alias := d.children[0]
// 			ev := makeNode(g_EXPORT_VAR, "", name, alias)

// 			vars = append(vars, ev)
// 		}

// 	} else if accept(p, tFUNCTION) {
// 		fs := parseFunctionDeclaration(p)
// 		declaration = fs
// 		name := makeNode(g_EXPORT_NAME, fs.value)
// 		alias := name
// 		ev := makeNode(g_EXPORT_VAR, "", name, alias)
// 		vars = append(vars, ev)

// 	} else if accept(p, tCLASS) {
// 		ce := parseClassExpression(p)
// 		declaration = ce
// 		name := ce.children[0]
// 		alias := name
// 		ev := makeNode(g_EXPORT_VAR, "", name, alias)
// 		vars = append(vars, ev)

// 	} else if accept(p, tMULT) {
// 		expect(p, tNAME)
// 		if getLexeme(p) != "from" {
// 			checkASI(p, tSEMI)
// 		}
// 		expect(p, tSTRING)
// 		path = makeNode(g_EXPORT_PATH, getLexeme(p))
// 		flags = flags | f_EXPORT_ALL
// 		expect(p, tSEMI)
// 	}
// 	varsNode := makeNode(g_EXPORT_VARS, "", vars...)

// 	result := makeNode(g_EXPORT_STATEMENT, "", varsNode, declaration, path)
// 	result.flags = flags
// 	return result
// }

// func parseSwitchStatement(p *parserState) ast {
// 	expect(p, tPAREN_LEFT)
// 	condition := parseExpression(p)
// 	expect(p, tPAREN_RIGHT)

// 	cases := []ast{}
// 	expect(p, tCURLY_LEFT)
// 	for !accept(p, tCURLY_RIGHT) {
// 		if accept(p, tCASE) {
// 			caseTest := makeNode(g_SWITCH_CASE_TEST, "", parseExpression(p))
// 			expect(p, tCOLON)
// 			caseStatements := []ast{}
// 			for !test(p, tDEFAULT, tCASE, tCURLY_RIGHT) {
// 				caseStatements = append(caseStatements, parseStatement(p))
// 			}
// 			statementNode := makeNode(
// 				g_SWITCH_CASE_STATEMENTS, "", caseStatements...,
// 			)
// 			switchCase := makeNode(
// 				g_SWITCH_CASE, "", caseTest, statementNode,
// 			)
// 			cases = append(cases, switchCase)
// 		}

// 		if accept(p, tDEFAULT) {
// 			expect(p, tCOLON)
// 			caseStatements := []ast{}
// 			for !test(p, tDEFAULT, tCASE, tCURLY_RIGHT) {
// 				caseStatements = append(caseStatements, parseStatement(p))
// 			}
// 			defaultCase := makeNode(g_SWITCH_DEFAULT, "", caseStatements...)
// 			cases = append(cases, defaultCase)
// 		}
// 	}

// 	casesNode := makeNode(g_SWITCH_CASES, "", cases...)
// 	return makeNode(g_SWITCH_STATEMENT, "", condition, casesNode)
// }

// func parseDeclarationStatement(p *parserState) ast {
// 	decl := makeNode(
// 		g_DECLARATION_STATEMENT, "", parseDeclarationExpression(p),
// 	)
// 	expect(p, tSEMI)
// 	return decl
// }

// func parseDeclarationExpression(p *parserState) ast {
// 	keyword := getLexeme(p)

// 	ve := makeNode(g_DECLARATION_EXPRESSION, keyword, []ast{}...)
// 	for ok := true; ok; ok = accept(p, tCOMMA) {
// 		ve.children = append(ve.children, parseDeclarator(p))
// 	}

// 	return ve
// }

// func parseForStatement(p *parserState) ast {
// 	expect(p, tPAREN_LEFT)
// 	var init ast
// 	if accept(p, tVAR, tLET, tCONST) {
// 		init = parseDeclarationExpression(p)
// 	} else if test(p, tSEMI) {
// 		init = makeNode(g_EMPTY_EXPRESSION, "")
// 	} else {
// 		startPos := p.index
// 		init = parseSequence(p)

// 		// we accidentally parsed for in loop
// 		if accept(p, tPAREN_RIGHT) {
// 			backtrack(p, startPos)
// 			p.globalFlagIsInsideForIn = true
// 			init = parseSequence(p)
// 			p.globalFlagIsInsideForIn = false
// 		}
// 	}

// 	if accept(p, tNAME) && getLexeme(p) == "of" {
// 		left := init
// 		right := parseExpression(p)
// 		expect(p, tPAREN_RIGHT)
// 		body := parseStatement(p)
// 		return makeNode(g_FOR_OF_STATEMENT, "", left, right, body)
// 	} else if accept(p, tIN) {
// 		left := init
// 		right := parseExpression(p)
// 		expect(p, tPAREN_RIGHT)
// 		body := parseStatement(p)
// 		return makeNode(g_FOR_IN_STATEMENT, "", left, right, body)
// 	} else {
// 		var test, final ast

// 		expect(p, tSEMI)
// 		if accept(p, tSEMI) {
// 			test = makeNode(g_EMPTY_EXPRESSION, "")
// 		} else {
// 			test = parseExpression(p)
// 			expect(p, tSEMI)
// 		}

// 		if accept(p, tPAREN_RIGHT) {
// 			final = makeNode(g_EMPTY_EXPRESSION, "")
// 		} else {
// 			final = parseExpression(p)
// 			expect(p, tPAREN_RIGHT)
// 		}

// 		body := parseStatement(p)

// 		return makeNode(g_FOR_STATEMENT, "", init, test, final, body)
// 	}
// }

// func parseIfStatement(p *parserState) ast {
// 	expect(p, tPAREN_LEFT)
// 	test := parseExpression(p)
// 	expect(p, tPAREN_RIGHT)
// 	body := parseStatement(p)
// 	if accept(p, tELSE) {
// 		alternate := parseStatement(p)
// 		return makeNode(g_IF_STATEMENT, "", test, body, alternate)
// 	}

// 	return makeNode(g_IF_STATEMENT, "", test, body)
// }

// func parseWhileStatement(p *parserState) ast {
// 	expect(p, tPAREN_LEFT)
// 	test := parseExpression(p)
// 	expect(p, tPAREN_RIGHT)
// 	body := parseStatement(p)

// 	return makeNode(g_WHILE_STATEMENT, "", test, body)
// }

// func parseDoWhileStatement(p *parserState) ast {
// 	body := parseStatement(p)
// 	expect(p, tWHILE)
// 	expect(p, tPAREN_LEFT)
// 	test := parseExpression(p)
// 	expect(p, tPAREN_RIGHT)

// 	return makeNode(g_DO_WHILE_STATEMENT, "", test, body)
// }

// func parseFunctionDeclaration(p *parserState) ast {
// 	expect(p, tNAME)
// 	name := getLexeme(p)

// 	expect(p, tPAREN_LEFT)
// 	params := parseFunctionParameters(p)
// 	expect(p, tCURLY_LEFT)
// 	body := parseBlockStatement(p)

// 	return makeNode(g_FUNCTION_DECLARATION, name, params, body)
// }

// func parseExpressionStatement(p *parserState) ast {
// 	if accept(p, tBREAK) {
// 		return makeNode(g_BREAK_STATEMENT, "")
// 	} else if accept(p, tCONTINUE) {
// 		return makeNode(g_CONTINUE_STATEMENT, "")
// 	} else if accept(p, tDEBUGGER) {
// 		return makeNode(g_DEBUGGER_STATEMENT, "")
// 	}
// 	expr := makeNode(g_EXPRESSION_STATEMENT, "", parseExpression(p))
// 	expect(p, tSEMI)
// 	return expr
// }

// func parseReturnStatement(p *parserState) ast {
// 	if accept(p, tSEMI) {
// 		return makeNode(g_RETURN_STATEMENT, "")
// 	}
// 	expr := parseExpression(p)
// 	return makeNode(g_RETURN_STATEMENT, "", expr)
// }

// func parseExpression(p *parserState) ast {
// 	return parseSequence(p)
// }

// func parseSequence(p *parserState) ast {
// 	firstItem := parseYield(p)

// 	children := []ast{firstItem}
// 	for accept(p, tCOMMA) {
// 		children = append(children, parseYield(p))
// 	}

// 	if len(children) > 1 {
// 		return makeNode(g_SEQUENCE_EXPRESSION, ",", children...)
// 	}
// 	return firstItem
// }

// func parseYield(p *parserState) ast {
// 	if accept(p, tYIELD) {
// 		return makeNode(g_EXPRESSION, "yield", parseYield(p))
// 	}
// 	return parseAssignment(p)
// }

// func parseAssignment(p *parserState) ast {
// 	var left ast
// 	startPos := p.index

// 	if accept(p, tCURLY_LEFT) {
// 		left = parseObjectPattern(p)
// 	} else if accept(p, tBRACKET_LEFT) {
// 		left = parseArrayPattern(p)
// 	}

// 	if left.t == g_INVALID_GRAMMAR {
// 		backtrack(p, startPos)
// 	} else if left.t != g_UNDEFINED {
// 		if accept(p, tASSIGN) {
// 			right := parseAssignment(p)
// 			return makeNode(g_EXPRESSION, "=", left, right)
// 		}
// 		backtrack(p, startPos)
// 	}

// 	left = parseConditional(p)

// 	if accept(p,
// 		tASSIGN, tPLUS_ASSIGN, tMINUS_ASSIGN, tMULT_ASSIGN, tDIV_ASSIGN,
// 		tBITWISE_SHIFT_LEFT_ASSIGN, tBITWISE_SHIFT_RIGHT_ASSIGN,
// 		tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN, tBITWISE_AND_ASSIGN,
// 		tBITWISE_OR_ASSIGN, tBITWISE_XOR_ASSIGN,
// 	) {
// 		op := getLexeme(p)
// 		right := parseAssignment(p)
// 		return makeNode(g_EXPRESSION, op, left, right)
// 	}

// 	return left
// }

// func parseAssignmentPattern(p *parserState) ast {
// 	var left ast

// 	if accept(p, tBRACKET_LEFT) {
// 		left = parseArrayPattern(p)
// 	} else if accept(p, tCURLY_LEFT) {
// 		left = parseObjectPattern(p)
// 	} else if accept(p, tNAME) {
// 		left = makeNode(g_OBJECT_LITERAL_PROPERTY_NAME, getLexeme(p))
// 	} else {
// 		return makeNode(g_INVALID_GRAMMAR, "")
// 	}

// 	if accept(p, tASSIGN) {
// 		right := parseYield(p)
// 		return makeNode(g_ASSIGNMENT_PATTERN, "=", left, right)
// 	}

// 	return left
// }

// func parseArrayPattern(p *parserState) ast {
// 	properties := []ast{}

// 	for !accept(p, tBRACKET_RIGHT) {
// 		if accept(p, tCOMMA) {
// 			properties = append(properties, makeNode(g_EMPTY_EXPRESSION, ""))
// 			continue
// 		}

// 		prop := parseAssignmentPattern(p)
// 		if prop.t == g_INVALID_GRAMMAR {
// 			return prop
// 		}
// 		properties = append(properties, prop)

// 		if !accept(p, tCOMMA) {
// 			if !accept(p, tBRACKET_RIGHT) {
// 				return makeNode(g_INVALID_GRAMMAR, "")
// 			}
// 			break
// 		}
// 	}

// 	return makeNode(g_ARRAY_PATTERN, "", properties...)
// }

// func parseClassExpression(p *parserState) ast {
// 	var name, extends, body ast

// 	if accept(p, tNAME) {
// 		name = makeNode(g_VARIABLE_NAME, getLexeme(p))
// 	}

// 	if accept(p, tEXTENDS) {
// 		extends = parseYield(p)
// 	}

// 	expect(p, tCURLY_LEFT)

// 	props := []ast{}

// 	for !accept(p, tCURLY_RIGHT) {
// 		isStatic := false

// 		if accept(p, tSTATIC) {
// 			isStatic = true
// 		}

// 		if accept(p, tNAME) {
// 			name := makeNode(g_OBJECT_LITERAL_PROPERTY_NAME, getLexeme(p))
// 			var prop ast

// 			if test(p, tPAREN_LEFT) {
// 				prop = makeNode(g_CLASS_METHOD, "", name, parseMemberFunction(p))
// 			} else if accept(p, tASSIGN) {
// 				prop = makeNode(g_CLASS_PROPERTY, "", name, parseExpression(p))
// 			} else {
// 				prop = makeNode(g_CLASS_PROPERTY, "", name)
// 			}

// 			if isStatic {
// 				prop = makeNode(g_CLASS_STATIC_PROP, "", prop)
// 			}
// 			props = append(props, prop)

// 			expect(p, tSEMI)
// 		}
// 	}
// 	body = makeNode(g_BLOCK_STATEMENT, "", props...)

// 	return makeNode(g_CLASS_EXPRESSION, "", name, extends, body)
// }

// func parseObjectPattern(p *parserState) ast {
// 	properties := []ast{}
// 	for !accept(p, tCURLY_RIGHT) {
// 		if accept(p, tSPREAD) {
// 			if !accept(p, tNAME) {
// 				return makeNode(g_INVALID_GRAMMAR, "")
// 			}
// 			name := makeNode(g_OBJECT_LITERAL_PROPERTY_NAME, getLexeme(p))
// 			properties = append(
// 				properties, makeNode(g_REST_EXPRESSION, "", name),
// 			)
// 		} else if accept(p, tNAME) {
// 			prop := makeNode(g_OBJECT_PROPERTY, "", []ast{}...)
// 			key := makeNode(g_OBJECT_LITERAL_PROPERTY_NAME, getLexeme(p))
// 			prop.children = append(prop.children, key)

// 			if accept(p, tCOLON) {
// 				val := parseAssignmentPattern(p)
// 				if val.t == g_INVALID_GRAMMAR {
// 					return val
// 				}
// 				prop.children = append(prop.children, val)
// 			}

// 			properties = append(properties, prop)
// 		}

// 		if !accept(p, tCOMMA) {
// 			if !accept(p, tCURLY_RIGHT) {
// 				return makeNode(g_INVALID_GRAMMAR, "")
// 			}
// 			break
// 		}
// 	}

// 	return makeNode(g_OBJECT_PATTERN, "", properties...)
// }

// func parseConditional(p *parserState) ast {
// 	test := parseBinary(p)

// 	if accept(p, tQUESTION) {
// 		consequent := parseConditional(p)
// 		expect(p, tCOLON)
// 		alternate := parseConditional(p)
// 		return makeNode(g_CONDITIONAL_EXPRESSION, "?", test, consequent, alternate)
// 	}

// 	return test
// }

// func parseBinary(p *parserState) ast {
// 	opStack := make([]token, 0)
// 	outputStack := make([]ast, 0)
// 	var root *ast

// 	addNode := func(t token) {
// 		right := outputStack[len(outputStack)-1]
// 		left := outputStack[len(outputStack)-2]
// 		nn := makeNode(g_EXPRESSION, t.lexeme, left, right)

// 		outputStack = outputStack[:len(outputStack)-2]
// 		outputStack = append(outputStack, nn)
// 		root = &nn
// 	}

// 	for {
// 		outputStack = append(outputStack, parsePrefixUnary(p))

// 		op, ok := operatorTable[p.tok.tType]
// 		if !ok || (p.globalFlagIsInsideForIn && p.tok.tType == tIN) {
// 			break
// 		}

// 		for {
// 			if len(opStack) > 0 {
// 				opToken := opStack[len(opStack)-1]

// 				opFromStack := operatorTable[opToken.tType]
// 				if opFromStack.precedence > op.precedence || (opFromStack.precedence == op.precedence && !op.isRightAssociative) {
// 					addNode(opToken)
// 					opStack = opStack[:len(opStack)-1]
// 				} else {
// 					break
// 				}
// 			} else {
// 				break
// 			}
// 		}
// 		opStack = append(opStack, p.tok)
// 		next(p)
// 	}

// 	for i := len(opStack) - 1; i >= 0; i-- {
// 		addNode(opStack[i])
// 	}

// 	if root != nil {
// 		return *root
// 	}
// 	return outputStack[len(outputStack)-1]
// }

// func parsePrefixUnary(p *parserState) ast {
// 	if accept(p,
// 		tNOT, tBITWISE_NOT, tPLUS, tMINUS,
// 		tINC, tDEC, tTYPEOF, tVOID, tDELETE,
// 	) {
// 		return makeNode(
// 			g_UNARY_PREFIX_EXPRESSION, getLexeme(p), parsePrefixUnary(p),
// 		)
// 	}

// 	return parsePostfixUnary(p)
// }

// func parsePostfixUnary(p *parserState) ast {
// 	value := parseFunctionCallOrMember(p, false)

// 	if accept(p, tINC, tDEC) {
// 		return makeNode(g_UNARY_POSTFIX_EXPRESSION, getLexeme(p), value)
// 	}

// 	return value
// }

// func parseFunctionArgs(p *parserState) ast {
// 	args := []ast{}

// 	for !accept(p, tPAREN_RIGHT) {
// 		args = append(args, parseExpression(p))

// 		if !accept(p, tCOMMA) {
// 			expect(p, tPAREN_RIGHT)
// 			break
// 		}
// 	}

// 	argsNode := makeNode(g_FUNCTION_ARGS, "", args...)

// 	return argsNode
// }

// func parseFunctionCallOrMember(p *parserState, onlyMember bool) ast {
// 	funcName := parseConstructorCall(p)

// 	for {
// 		if !onlyMember && accept(p, tPAREN_LEFT) {
// 			argsNode := parseFunctionArgs(p)
// 			n := makeNode(g_FUNCTION_CALL, "", funcName, argsNode)

// 			funcName = n
// 		} else {
// 			var property ast

// 			if accept(p, tDOT) {
// 				expect(p, tNAME)
// 				property = makeNode(g_VARIABLE_NAME, getLexeme(p))
// 			} else if accept(p, tBRACKET_LEFT) {
// 				property = parseCalculatedPropertyName(p)
// 			} else {
// 				break
// 			}
// 			object := funcName

// 			me := makeNode(g_MEMBER_EXPRESSION, "", object, property)
// 			funcName = me
// 		}
// 	}

// 	return funcName
// }

// func parseConstructorCall(p *parserState) ast {
// 	if accept(p, tNEW) {
// 		name := parseFunctionCallOrMember(p, true)
// 		if accept(p, tPAREN_LEFT) {
// 			return makeNode(
// 				g_CONSTRUCTOR_CALL, "", name, parseFunctionArgs(p),
// 			)
// 		}
// 		return makeNode(g_CONSTRUCTOR_CALL, "", name)
// 	}

// 	return parseAtom(p)
// }

// func parseAtom(p *parserState) ast {
// 	switch {
// 	case accept(p, tCLASS):
// 		return parseClassExpression(p)
// 	case accept(p, tDIV):
// 		return parseRegexp(p)
// 	case accept(p, tHEX):
// 		return makeNode(g_HEX_LITERAL, getLexeme(p))
// 	case accept(p, tPAREN_LEFT):
// 		return parseParensOrLambda(p)
// 	case accept(p, tCURLY_LEFT):
// 		return parseObjectLiteral(p)
// 	case accept(p, tFUNCTION):
// 		return parseFunctionExpression(p)
// 	case accept(p, tBRACKET_LEFT):
// 		return parseArrayLiteral(p)
// 	case accept(p, tNAME):
// 		return parseLambdaOrName(p)
// 	case accept(p, tNUMBER):
// 		return makeNode(g_NUMBER_LITERAL, getLexeme(p))
// 	case accept(p, tSTRING):
// 		return makeNode(g_STRING_LITERAL, getLexeme(p))
// 	case accept(p, tNULL):
// 		return makeNode(g_NULL, getLexeme(p))
// 	case accept(p, tUNDEFINED):
// 		return makeNode(g_UNDEFINED, getLexeme(p))
// 	case accept(p, tTRUE, tFALSE):
// 		return makeNode(g_BOOL_LITERAL, getLexeme(p))
// 	case accept(p, tTHIS):
// 		return makeNode(g_THIS, getLexeme(p))

// 	default:
// 		checkASI(p, tSEMI)
// 	}
// 	return ast{}
// }

// func parseRegexp(p *parserState) ast {
// 	value := "/"
// 	for {
// 		if accept(p, tDIV) {
// 			if p.sourceTokens[p.index-2].tType != tESCAPE {
// 				break
// 			} else {
// 				value += getLexeme(p)
// 			}
// 		}
// 		next(p)
// 		value += getLexeme(p)
// 	}
// 	value += "/"
// 	if accept(p, tNAME) {
// 		value += getLexeme(p)
// 	}

// 	return makeNode(g_REGEXP_LITERAL, value)
// }

// func parseParensOrLambda(p *parserState) ast {
// 	prevPos := p.index

// 	params := parseFunctionParameters(p)

// 	if accept(p, tLAMBDA) {
// 		var body ast
// 		if accept(p, tCURLY_LEFT) {
// 			body = parseBlockStatement(p)
// 		} else {
// 			body = parseYield(p)
// 		}
// 		return makeNode(g_LAMBDA_EXPRESSION, "", params, body)
// 	}

// 	backtrack(p, prevPos)

// 	var value ast
// 	if !accept(p, tPAREN_RIGHT) {
// 		value = parseExpression(p)
// 		expect(p, tPAREN_RIGHT)
// 	}

// 	return makeNode(g_PARENS_EXPRESSION, "", value)
// }

// func parseFunctionParameters(p *parserState) ast {
// 	startPos := p.index

// 	result := ast{g_FUNCTION_PARAMETERS, "", nil, 0}
// 	params := []ast{}
// 	for !accept(p, tPAREN_RIGHT) {
// 		params = append(params, parseFunctionParameter(p))

// 		if !accept(p, tCOMMA) {
// 			if !accept(p, tPAREN_RIGHT) {
// 				backtrack(p, startPos)
// 				return result
// 			}
// 			break
// 		}
// 	}
// 	result.children = params
// 	return result
// }

// func parseFunctionExpression(p *parserState) ast {
// 	name := ""
// 	if accept(p, tNAME) {
// 		name = getLexeme(p)
// 	}

// 	expect(p, tPAREN_LEFT)
// 	params := parseFunctionParameters(p)
// 	expect(p, tCURLY_LEFT)
// 	body := parseBlockStatement(p)

// 	return makeNode(g_FUNCTION_EXPRESSION, name, params, body)
// }

// func parseObjectLiteral(p *parserState) ast {
// 	props := []ast{}

// 	for !accept(p, tCURLY_RIGHT) {
// 		prop := ast{g_OBJECT_PROPERTY, "", []ast{}, 0}
// 		var key, value ast

// 		if accept(p, tSPREAD) {
// 			expr := parseYield(p)
// 			prop = makeNode(g_SPREAD_EXPRESSION, "", expr)
// 		} else {
// 			if getNoNewline(p).lexeme == "get" {
// 				accept(p, tNAME)
// 				prop.value = getLexeme(p)
// 			} else if getNoNewline(p).lexeme == "set" {
// 				accept(p, tNAME)
// 				prop.value = getLexeme(p)
// 			}

// 			if accept(p, tNAME) {
// 				key = makeNode(g_OBJECT_LITERAL_PROPERTY_NAME, getLexeme(p))
// 			} else if accept(p, tBRACKET_LEFT) {
// 				key = parseCalculatedPropertyName(p)
// 			} else if isValidPropertyName(getNoNewline(p).lexeme) || test(p, tNUMBER, tSTRING) {
// 				accept(p)
// 				key = makeNode(g_VALID_PROPERTY_NAME, getLexeme(p))
// 			}
// 			prop.children = append(prop.children, key)

// 			if test(p, tPAREN_LEFT) {
// 				value = parseMemberFunction(p)
// 				prop.children = append(prop.children, value)
// 			} else if accept(p, tCOLON) {
// 				value = parseYield(p)
// 				prop.children = append(prop.children, value)
// 			}
// 		}

// 		props = append(props, prop)
// 		if !accept(p, tCOMMA) {
// 			expect(p, tCURLY_RIGHT)
// 			break
// 		}
// 	}

// 	return makeNode(g_OBJECT_LITERAL, "", props...)
// }

// func parseMemberFunction(p *parserState) ast {
// 	f := parseFunctionExpression(p)
// 	f.t = g_MEMBER_FUNCTION
// 	return f
// }

// func parseArrayLiteral(p *parserState) ast {
// 	values := []ast{}

// 	for !accept(p, tBRACKET_RIGHT) {
// 		if accept(p, tCOMMA) {
// 			values = append(values, makeNode(g_EMPTY_EXPRESSION, ""))
// 			continue
// 		}
// 		if accept(p, tSPREAD) {
// 			values = append(
// 				values,
// 				makeNode(g_SPREAD_EXPRESSION, "", parseYield(p)),
// 			)
// 		} else {
// 			values = append(values, parseYield(p))
// 		}

// 		if !accept(p, tCOMMA) {
// 			expect(p, tBRACKET_RIGHT)
// 			break
// 		}
// 	}

// 	return makeNode(g_ARRAY_LITERAL, "", values...)
// }

// func parseLambdaOrName(p *parserState) ast {
// 	firstParamStr := getLexeme(p)

// 	if accept(p, tLAMBDA) {
// 		var body ast
// 		if accept(p, tCURLY_LEFT) {
// 			body = parseBlockStatement(p)
// 		} else {
// 			body = parseYield(p)
// 		}
// 		params := makeNode(
// 			g_FUNCTION_PARAMETERS, "",
// 			makeNode(g_VARIABLE_NAME, firstParamStr),
// 		)
// 		return makeNode(g_LAMBDA_EXPRESSION, "", params, body)
// 	}

// 	return makeNode(g_VARIABLE_NAME, firstParamStr)
// }

// func parseBlockStatement(p *parserState) ast {
// 	statements := []ast{}
// 	for !accept(p, tCURLY_RIGHT) {
// 		statements = append(statements, parseStatement(p))
// 	}
// 	return makeNode(g_BLOCK_STATEMENT, "", statements...)
// }

// func parseFunctionParameter(p *parserState) ast {
// 	if accept(p, tSPREAD) {
// 		var expr ast
// 		if accept(p, tBRACKET_LEFT) {
// 			expr = parseArrayPattern(p)
// 		} else if accept(p, tCURLY_LEFT) {
// 			expr = parseObjectPattern(p)
// 		} else if accept(p, tNAME) {
// 			expr = makeNode(g_VARIABLE_NAME, getLexeme(p))
// 		}
// 		return makeNode(g_REST_EXPRESSION, "", expr)
// 	}
// 	n := parseAssignmentPattern(p)
// 	return n
// }

// func parseDeclarator(p *parserState) ast {
// 	var left ast
// 	if accept(p, tCURLY_LEFT) {
// 		left = parseObjectPattern(p)
// 	} else if accept(p, tNAME) {
// 		left = makeNode(g_VARIABLE_NAME, getLexeme(p))
// 	} else {
// 		return left
// 	}

// 	if accept(p, tASSIGN) {
// 		right := parseAssignment(p)
// 		return makeNode(g_DECLARATOR, "", left, right)
// 	}

// 	return makeNode(g_DECLARATOR, "", left)
// }

// func parseCalculatedPropertyName(p *parserState) ast {
// 	value := parseExpression(p)
// 	expect(p, tBRACKET_RIGHT)
// 	return makeNode(g_CALCULATED_PROPERTY_NAME, "", value)
// }

// func parseImportStatement(p *parserState) ast {
// 	vars := ast{g_IMPORT_VARS, "", []ast{}, 0}
// 	all := ast{t: g_IMPORT_ALL}
// 	path := ast{t: g_IMPORT_PATH}

// 	if accept(p, tMULT) {
// 		expect(p, tNAME)
// 		if getLexeme(p) != "as" {
// 			checkASI(p, tSEMI)
// 		}
// 		expect(p, tNAME)
// 		all.value = getLexeme(p)

// 		expect(p, tNAME)
// 		if getLexeme(p) != "from" {
// 			checkASI(p, tSEMI)
// 		}
// 	} else {
// 		if accept(p, tNAME) {
// 			name := makeNode(g_IMPORT_NAME, "default")
// 			alias := makeNode(g_IMPORT_ALIAS, getLexeme(p))

// 			varNode := makeNode(g_IMPORT_VAR, "", name, alias)
// 			vars.children = append(vars.children, varNode)

// 			if accept(p, tCOMMA) {
// 				if accept(p, tMULT) {
// 					expect(p, tNAME)
// 					if getLexeme(p) != "as" {
// 						checkASI(p, tSEMI)
// 					}

// 					expect(p, tNAME)
// 					all.value = getLexeme(p)

// 					expect(p, tNAME)
// 					if getLexeme(p) != "from" {
// 						checkASI(p, tSEMI)
// 					}
// 				}
// 			} else {
// 				expect(p, tNAME)
// 				if getLexeme(p) != "from" {
// 					checkASI(p, tSEMI)
// 				}
// 			}
// 		}

// 		if accept(p, tCURLY_LEFT) {
// 			for !accept(p, tCURLY_RIGHT) {
// 				if accept(p, tNAME, tDEFAULT) {
// 					name := makeNode(g_IMPORT_NAME, getLexeme(p))
// 					alias := makeNode(g_IMPORT_ALIAS, getLexeme(p))

// 					if accept(p, tNAME) && getLexeme(p) == "as" {
// 						expect(p, tNAME)
// 						alias = makeNode(g_IMPORT_ALIAS, getLexeme(p))
// 					}

// 					varNode := makeNode(g_IMPORT_VAR, "", name, alias)
// 					vars.children = append(vars.children, varNode)
// 				}

// 				if !accept(p, tCOMMA) {
// 					expect(p, tCURLY_RIGHT)
// 					break
// 				}
// 			}

// 			expect(p, tNAME)
// 			if getLexeme(p) != "from" {
// 				checkASI(p, tSEMI)
// 			}
// 		}
// 	}

// 	expect(p, tSTRING)
// 	path.value = getLexeme(p)
// 	expect(p, tSEMI)

// 	return makeNode(g_IMPORT_STATEMENT, "", vars, all, path)
// }
