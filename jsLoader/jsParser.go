package jsLoader

import "fmt"

var sourceTokens []token
var index int
var tok token

func parseTokens(localSrc []token) (program, error) {
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
	kind int
	tok  token
}

func (pe parsingError) Error() string {
	switch pe.kind {
	case p_UNEXPECTED_TOKEN:
		return fmt.Sprintf("Unexpected token \"%s\" at %v:%v", pe.tok.lexeme, pe.tok.line, pe.tok.column)

	case p_UNNAMED_FUNCTION:
		return fmt.Sprintf("Unnamed function declaration at %v:%v", pe.tok.line, pe.tok.column)

	case p_WRONG_ASSIGNMENT:
		return fmt.Sprintf("Invalid left-hand side in assignment at %v:%v", pe.tok.line, pe.tok.column)

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

func accept(tTypes ...tokenType) bool {
	prevPos := index
	for tok.tType == tNEWLINE {
		next()
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
	panic(parsingError{p_UNEXPECTED_TOKEN, tok})
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

// all types below implement ast interface

type ast interface {
	String() string
}

type program struct {
	statements []ast
}

type switchCase struct {
	condition  ast
	statements []ast
}

type switchStatement struct {
	condition ast
	cases     []switchCase
	defCase   []ast
}

type exportedVar struct {
	name  ast
	alias atom
}

type exportStatement struct {
	vars        []exportedVar
	declaration ast
	source      atom
	exportAll   bool
}

type importedVar struct {
	name  atom
	alias atom
}

type importStatement struct {
	path atom
	vars []importedVar
	all  atom
}

type returnStatement struct {
	value ast
}

type functionStatement struct {
	funcExpr functionExpression
}

type expressionStatement struct {
	expr ast
}

type ifStatement struct {
	test       ast
	consequent ast
	alternate  ast
}

type whileStatement struct {
	test ast
	body ast
}

type doWhileStatement struct {
	test ast
	body ast
}

type emptyStatement struct{}

type forStatement struct {
	init  ast
	test  ast
	final ast
	body  ast
}

type forInStatement struct {
	left  ast
	right ast
	body  ast
}

type forOfStatement struct {
	left  ast
	right ast
	body  ast
}

type blockStatement struct {
	statements []ast
}

type objectPattern struct {
	properties []objectProperty
}

type arrayPattern struct {
	properties []declarator
}

type declarator struct {
	left  ast
	value ast
}

type varDeclaration struct {
	declarations []declarator
	keyword      atom
}

type sequenceExpression struct {
	items []ast
}

type yieldExpression struct {
	value ast
}

type binaryExpression struct {
	left     ast
	right    ast
	operator atom
}

type conditionalExpression struct {
	test       ast
	consequent ast
	alternate  ast
}

type unaryExpression struct {
	value     ast
	operator  atom
	isPostfix bool
}

type functionCall struct {
	name          ast
	args          []ast
	isConstructor bool
}

type memberExpression struct {
	object   ast
	property propertyName
}

type lambdaExpression struct {
	args []declarator
	body ast
}

type functionExpression struct {
	name     atom
	args     []declarator
	body     blockStatement
	isMember bool
}

type propertyName struct {
	value        ast
	isCalculated bool
	isNotName    bool
}

type objectProperty struct {
	key   propertyName
	value ast
}

type objectLiteral struct {
	properties []objectProperty
}

type arrayLiteral struct {
	items []ast
}

type parenExpression struct {
	inner ast
}

type atom struct {
	val token
}

type operatorInfo struct {
	precedence         int
	isRightAssociative bool
}

var operatorTable = map[tokenType]operatorInfo{
	tASSIGN:                          operatorInfo{3, true},
	tPLUS_ASSIGN:                     operatorInfo{3, true},
	tMINUS_ASSIGN:                    operatorInfo{3, true},
	tMULT_ASSIGN:                     operatorInfo{3, true},
	tDIV_ASSIGN:                      operatorInfo{3, true},
	tMOD_ASSIGN:                      operatorInfo{3, true},
	tBITWISE_SHIFT_LEFT_ASSIGN:       operatorInfo{3, true},
	tBITWISE_SHIFT_RIGHT_ASSIGN:      operatorInfo{3, true},
	tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN: operatorInfo{3, true},
	tBITWISE_AND_ASSIGN:              operatorInfo{3, true},
	tBITWISE_OR_ASSIGN:               operatorInfo{3, true},
	tBITWISE_XOR_ASSIGN:              operatorInfo{3, true},
	tOR:                              operatorInfo{5, false},
	tAND:                             operatorInfo{6, false},
	tBITWISE_OR:                      operatorInfo{7, false},
	tBITWISE_XOR:                     operatorInfo{8, false},
	tBITWISE_AND:                     operatorInfo{9, false},
	tEQUALS:                          operatorInfo{10, false},
	tNOT_EQUALS:                      operatorInfo{10, false},
	tEQUALS_STRICT:                   operatorInfo{10, false},
	tNOT_EQUALS_STRICT:               operatorInfo{10, false},
	tLESS:                            operatorInfo{11, false},
	tLESS_EQUALS:                     operatorInfo{11, false},
	tGREATER:                         operatorInfo{11, false},
	tGREATER_EQUALS:                  operatorInfo{11, false},
	tINSTANCEOF:                      operatorInfo{11, false},
	tBITWISE_SHIFT_LEFT:              operatorInfo{12, false},
	tBITWISE_SHIFT_RIGHT:             operatorInfo{12, false},
	tBITWISE_SHIFT_RIGHT_ZERO:        operatorInfo{12, false},
	tPLUS:  operatorInfo{13, false},
	tMINUS: operatorInfo{13, false},
	tMULT:  operatorInfo{14, false},
	tDIV:   operatorInfo{14, false},
	tMOD:   operatorInfo{14, false},
	tEXP:   operatorInfo{15, true},
}

func parseProgram() (p program, err error) {
	p = program{}
	p.statements = []ast{}
	err = nil

	defer func() {
		if r := recover(); r != nil {
			if pe, ok := r.(parsingError); ok {
				err = pe
			}
		}
	}()

	for !accept(tEND_OF_INPUT) {
		p.statements = append(p.statements, parseStatement())
	}
	return
}

func parseStatement() ast {
	var st ast

	if s, ok := parseEmptyStatement(); ok {
		st = s
	} else if s, ok := parseForStatement(); ok {
		st = s
	} else if s, ok := parseWhileStatement(); ok {
		st = s
	} else if s, ok := parseDoWhileStatement(); ok {
		st = s
	} else if s, ok := parseBlockStatement(); ok {
		st = s
	} else if s, ok := parseIfStatement(); ok {
		st = s
	} else if s, ok := parseFunctionStatement(); ok {
		st = s
	} else if s, ok := parseReturnStatement(); ok {
		st = s
	} else if s, ok := parseImportStatement(); ok {
		st = s
	} else if s, ok := parseSwitchStatement(); ok {
		st = s
	} else if s, ok := parseExportStatement(); ok {
		st = s
	} else {
		st = parseExpressionStatement()
	}

	return st
}

func parseSwitchStatement() (switchStatement, bool) {
	sw := switchStatement{}
	if !accept(tSWITCH) {
		return sw, false
	}

	expect(tPAREN_LEFT)
	sw.condition = parseSequence()
	expect(tPAREN_RIGHT)
	expect(tCURLY_LEFT)

	sw.cases = []switchCase{}
	for !accept(tCURLY_RIGHT) {
		if accept(tCASE) {
			c := switchCase{}
			c.condition = parseSequence()
			c.statements = []ast{}

			expect(tCOLON)
			for !test(tCASE, tDEFAULT, tCURLY_RIGHT) {
				c.statements = append(c.statements, parseStatement())
			}
			sw.cases = append(sw.cases, c)

			if accept(tDEFAULT) {
				expect(tCOLON)
				sw.defCase = []ast{}
				for !test(tCASE, tCURLY_RIGHT) {
					sw.defCase = append(sw.defCase, parseStatement())
				}
			}
		}
	}

	return sw, true
}

func parseExportStatement() (exportStatement, bool) {
	es := exportStatement{}
	if !accept(tEXPORT) {
		return es, false
	}
	es.vars = []exportedVar{}

	if accept(tDEFAULT) {
		ev := exportedVar{}
		ev.alias = atom{makeToken("default")}

		if fe, ok := parseFunctionExpression(false); ok {
			if fe.name.val.lexeme != "" {
				es.declaration = fe
				ev.name = fe.name
			} else {
				ev.name = fe
			}
		} else {
			ev.name = parseYield()
		}

		expect(tSEMI)
		es.vars = append(es.vars, ev)

	} else if accept(tCURLY_LEFT) {
		for !accept(tCURLY_RIGHT) {
			if accept(tNAME) {
				ev := exportedVar{}
				ev.name = atom{getToken()}
				if accept(tAS) {
					next()
					ev.alias = atom{getToken()}
				}
				es.vars = append(es.vars, ev)
			}

			if !accept(tCOMMA) {
				expect(tCURLY_RIGHT)
				break
			}
		}

		if accept(tFROM) {
			expect(tSTRING)
			es.source = atom{getToken()}
		}
		expect(tSEMI)

	} else if vd, ok := parseVarDeclaration(); ok {
		es.declaration = vd
		expect(tSEMI)
		for _, d := range vd.declarations {
			ev := exportedVar{}
			ev.name = d.left
			es.vars = append(es.vars, ev)
		}

	} else if fs, ok := parseFunctionStatement(); ok {
		es.declaration = fs
		ev := exportedVar{}
		ev.name = fs.funcExpr.name
		es.vars = append(es.vars, ev)

	} else if accept(tMULT) {
		expect(tFROM)
		expect(tSTRING)
		es.source = atom{getToken()}
		es.exportAll = true
		expect(tSEMI)
	}

	return es, true
}

func parseImportStatement() (importStatement, bool) {
	is := importStatement{}
	if !accept(tIMPORT) {
		return is, false
	}
	is.vars = []importedVar{}

	if accept(tMULT) {
		expect(tAS)
		expect(tNAME)
		is.all = atom{getToken()}
		expect(tFROM)
	} else {
		if accept(tNAME) {
			v := importedVar{}
			v.name = atom{makeToken("default")}
			v.alias = atom{getToken()}
			is.vars = append(is.vars, v)

			if accept(tCOMMA) {
				if accept(tMULT) {
					expect(tAS)
					expect(tNAME)
					is.all = atom{getToken()}
					expect(tFROM)
				}
			} else {
				expect(tFROM)
			}
		}

		if accept(tCURLY_LEFT) {
			for !accept(tCURLY_RIGHT) {
				if accept(tNAME, tDEFAULT) {
					v := importedVar{}
					v.name = atom{getToken()}

					if accept(tAS) {
						expect(tNAME)
						v.alias = atom{getToken()}
					}

					is.vars = append(is.vars, v)
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
	is.path = atom{getToken()}
	expect(tSEMI)

	return is, true
}

func parseReturnStatement() (returnStatement, bool) {
	rs := returnStatement{}
	if !accept(tRETURN) {
		return rs, false
	}
	if !accept(tSEMI) {
		rs.value = parseYield()
		expect(tSEMI)
	}
	return rs, true
}

func parseFunctionStatement() (functionStatement, bool) {
	fs := functionStatement{}

	if fe, ok := parseFunctionExpression(false); ok {
		if fe.name.val.lexeme == "" {
			panic(parsingError{p_UNNAMED_FUNCTION, tok})
		} else {
			fs.funcExpr = fe
			return fs, true
		}
	}

	return fs, false
}

func parseExpressionStatement() expressionStatement {
	es := expressionStatement{}
	if vd, ok := parseVarDeclaration(); ok {
		es.expr = vd
	} else if accept(tCONTINUE, tDEBUGGER, tBREAK) {
		es.expr = atom{getToken()}
	} else {
		es.expr = parseSequence()
	}
	expect(tSEMI)
	return es
}

func parseIfStatement() (ifStatement, bool) {
	is := ifStatement{}

	if !accept(tIF) {
		return is, false
	}
	expect(tPAREN_LEFT)
	is.test = parseSequence()
	expect(tPAREN_RIGHT)
	is.consequent = parseStatement()

	if accept(tELSE) {
		is.alternate = parseStatement()
	}

	return is, true
}

func parseWhileStatement() (whileStatement, bool) {
	ws := whileStatement{}
	if !accept(tWHILE) {
		return ws, false
	}
	expect(tPAREN_LEFT)
	ws.test = parseSequence()
	expect(tPAREN_RIGHT)
	ws.body = parseStatement()

	return ws, true
}

func parseDoWhileStatement() (doWhileStatement, bool) {
	ds := doWhileStatement{}

	if !accept(tDO) {
		return ds, false
	}

	ds.body = parseStatement()
	expect(tWHILE)
	expect(tPAREN_LEFT)
	ds.test = parseSequence()
	expect(tPAREN_RIGHT)

	return ds, true
}

func parseEmptyStatement() (emptyStatement, bool) {
	es := emptyStatement{}
	if accept(tSEMI) {
		return es, true
	}
	return es, false
}

func parseForStatement() (ast, bool) {
	if !accept(tFOR) {
		return nil, false
	}
	expect(tPAREN_LEFT)

	var init ast
	if vd, ok := parseVarDeclaration(); ok {
		init = vd
	} else if test(tSEMI) {
		init = atom{token{lexeme: ""}}
	} else {
		init = parseSequence()
	}

	if accept(tOF) {
		return parseForOfStatement(init), true
	} else if accept(tIN) {
		return parseForInStatement(init), true
	}

	fs := forStatement{}
	fs.init = init
	expect(tSEMI)

	if accept(tSEMI) {
		fs.test = atom{token{lexeme: ""}}
	} else {
		fs.test = parseSequence()
		expect(tSEMI)
	}

	if accept(tPAREN_RIGHT) {
		fs.final = atom{token{lexeme: ""}}
	} else {
		fs.final = parseSequence()
		expect(tPAREN_RIGHT)
	}
	fs.body = parseStatement()

	return fs, true
}

func parseObjectPattern() (objectPattern, bool) {
	startPos := index

	op := objectPattern{}
	if !accept(tCURLY_LEFT) {
		return op, false
	}
	op.properties = []objectProperty{}

	for !accept(tCURLY_RIGHT) {
		prop := objectProperty{}
		prop.key = parsePropertyName()

		if accept(tCOLON) {
			prop.value = parseAssignment()
		}

		op.properties = append(op.properties, prop)

		if !accept(tCOMMA) {
			if !accept(tCURLY_RIGHT) {
				backtrack(startPos)
				return op, false
			}
			break
		}
	}

	return op, true
}

func parseArrayPattern() (arrayPattern, bool) {
	ap := arrayPattern{}

	if !accept(tBRACKET_LEFT) {
		return ap, false
	}

	ap.properties = []declarator{}
	for ok := true; ok; ok = accept(tCOMMA) {
		d, _ := parseDeclarator()
		ap.properties = append(ap.properties, d)
	}
	expect(tBRACKET_RIGHT)

	return ap, true
}

func parseDeclarator() (declarator, bool) {
	d := declarator{}

	if accept(tNAME) {
		d.left = atom{getToken()}
	} else if op, ok := parseObjectPattern(); ok {
		d.left = op
	} else {
		return d, false
	}

	if accept(tASSIGN) {
		d.value = parseAssignment()
	}

	return d, true
}

func parseVarDeclaration() (varDeclaration, bool) {
	vd := varDeclaration{}
	if !accept(tVAR, tLET, tCONST) {
		return vd, false
	}
	vd.keyword = atom{getToken()}
	vd.declarations = []declarator{}

	for ok := true; ok; ok = accept(tCOMMA) {
		if d, ok := parseDeclarator(); ok {
			vd.declarations = append(vd.declarations, d)
		} else {
			checkASI(tSEMI)
		}
	}

	return vd, true
}

func parseSequence() ast {
	firstItem := parseYield()

	se := sequenceExpression{}
	se.items = []ast{firstItem}
	for accept(tCOMMA) {
		se.items = append(se.items, parseYield())
	}
	if len(se.items) > 1 {
		return se
	}

	return firstItem
}

func parseYield() ast {
	if accept(tYIELD) {
		ye := yieldExpression{}
		ye.value = parseYield()
		return ye
	}
	return parseAssignment()
}

func parseAssignment() ast {
	if op, ok := parseObjectPattern(); ok {
		binExpr := binaryExpression{}
		binExpr.left = op

		if accept(tASSIGN) {
			binExpr.operator = atom{getToken()}
			binExpr.right = parseAssignment()
			return binExpr
		}
		return binExpr.left
	}

	left := parseConditional()

	if accept(
		tASSIGN, tPLUS_ASSIGN, tMINUS_ASSIGN, tMULT_ASSIGN, tDIV_ASSIGN,
		tBITWISE_SHIFT_LEFT_ASSIGN, tBITWISE_SHIFT_RIGHT_ASSIGN,
		tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN, tBITWISE_AND_ASSIGN,
		tBITWISE_OR_ASSIGN, tBITWISE_XOR_ASSIGN,
	) {
		binExpr := binaryExpression{}
		binExpr.operator = atom{getToken()}
		binExpr.left = left
		binExpr.right = parseAssignment()

		return binExpr
	}

	return left
}

func parseConditional() ast {
	test := parseBinary()

	if accept(tQUESTION) {
		condExp := conditionalExpression{}
		condExp.test = test
		condExp.consequent = parseConditional()
		expect(tCOLON)
		condExp.alternate = parseConditional()
		return condExp
	}

	return test
}

func parseBinary() ast {
	opStack := make([]token, 0)
	outputStack := make([]ast, 0)

	var root *binaryExpression

	addNode := func(t token) {
		newNode := binaryExpression{}
		newNode.right = outputStack[len(outputStack)-1]
		newNode.left = outputStack[len(outputStack)-2]
		newNode.operator = atom{t}

		outputStack = outputStack[:len(outputStack)-2]
		outputStack = append(outputStack, newNode)
		root = &newNode
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

func parsePrefixUnary() ast {
	if accept(
		tNOT, tBITWISE_NOT, tPLUS, tMINUS,
		tINC, tDEC, tTYPEOF, tVOID, tDELETE,
	) {
		unExp := unaryExpression{}
		unExp.operator = atom{getToken()}
		unExp.value = parsePostfixUnary()
		return unExp
	}

	return parsePostfixUnary()
}

func parsePostfixUnary() ast {
	value := parseFunctionCallOrMember(false)

	if accept(tINC, tDEC) {
		unExp := unaryExpression{}
		unExp.isPostfix = true
		unExp.operator = atom{getToken()}
		unExp.value = value
		return unExp
	}

	return value
}

func parseFunctionCallOrMember(noFuncCall bool) ast {
	funcName := parseConstructorCall()

	for {
		if !noFuncCall && accept(tPAREN_LEFT) {
			fc := functionCall{}
			fc.name = funcName
			fc.args = []ast{}

			for !accept(tPAREN_RIGHT) {
				fc.args = append(fc.args, parseYield())
				if !accept(tCOMMA) {
					expect(tPAREN_RIGHT)
					break
				}
			}
			funcName = fc
		} else if accept(tDOT) || test(tBRACKET_LEFT) {
			me := memberExpression{}
			me.object = funcName
			me.property = parsePropertyName()
			funcName = me
		} else {
			break
		}
	}

	return funcName
}

func parseConstructorCall() ast {
	if accept(tNEW) {
		fc := functionCall{}
		fc.name = parseFunctionCallOrMember(true)
		fc.isConstructor = true

		if accept(tPAREN_LEFT) {
			fc.args = []ast{}

			for !accept(tPAREN_RIGHT) {
				fc.args = append(fc.args, parseYield())
				if !accept(tCOMMA) {
					expect(tPAREN_RIGHT)
					break
				}
			}
		}

		return fc
	}

	return parseAtom()
}

func parseAtom() ast {
	if a, ok := parseLambda(); ok {
		return a
	} else if a, ok := parseParens(); ok {
		return a
	} else if a, ok := parseObjectLiteral(); ok {
		return a
	} else if a, ok := parseFunctionExpression(false); ok {
		return a
	} else if a, ok := parseArrayLiteral(); ok {
		return a
	} else if accept(tNAME, tNUMBER, tSTRING, tTRUE, tFALSE, tNULL, tUNDEFINED, tTHIS) {
		return atom{getToken()}
	} else {
		checkASI(tSEMI)
	}
	return nil
}

func parseLambda() (lambdaExpression, bool) {
	startPos := index

	le := lambdaExpression{}
	if accept(tNAME) {
		le.args = []declarator{declarator{left: atom{getToken()}}}
	} else if args, ok := parseFunctionArgs(); ok {
		le.args = args
	}

	if le.args != nil {
		if !accept(tLAMBDA) {
			backtrack(startPos)
			return le, false
		}

		if bs, ok := parseBlockStatement(); ok {
			le.body = bs
		} else {
			le.body = parseYield()
		}
		return le, true
	}

	backtrack(startPos)
	return le, false
}

func parseFunctionArgs() ([]declarator, bool) {
	startPos := index
	if !accept(tPAREN_LEFT) {
		return nil, false
	}

	args := []declarator{}
	for !accept(tPAREN_RIGHT) {
		if d, ok := parseDeclarator(); ok {
			args = append(args, d)
		} else if !accept(tPAREN_RIGHT) {
			backtrack(startPos)
			return nil, false
		}

		if !accept(tCOMMA) {
			if !accept(tPAREN_RIGHT) {
				backtrack(startPos)
				return nil, false
			}
			break
		}
	}

	return args, true
}

func parseParens() (ast, bool) {
	if !accept(tPAREN_LEFT) {
		return nil, false
	}
	res := parseSequence()
	expect(tPAREN_RIGHT)
	return parenExpression{res}, true
}

func parseObjectLiteral() (objectLiteral, bool) {
	ol := objectLiteral{}

	if !accept(tCURLY_LEFT) {
		return ol, false
	}

	ol.properties = []objectProperty{}

	for !accept(tCURLY_RIGHT) {
		prop := objectProperty{}
		prop.key = parsePropertyName()

		if accept(tCOLON) {
			prop.value = parseYield()
		} else if fe, ok := parseFunctionExpression(true); ok {
			prop.value = fe
		} else if prop.key.isCalculated || prop.key.isNotName {
			checkASI(tSEMI)
		}

		ol.properties = append(ol.properties, prop)
		if !accept(tCOMMA) {
			expect(tCURLY_RIGHT)
			break
		}
	}

	return ol, true
}

func parseFunctionExpression(isMember bool) (functionExpression, bool) {
	fe := functionExpression{}
	if isMember {
		if !test(tPAREN_LEFT) {
			return fe, false
		}

		fe.isMember = true
	} else {
		if !accept(tFUNCTION) {
			return fe, false
		}

		if accept(tNAME) {
			fe.name = atom{getToken()}
		}
	}

	fe.args, _ = parseFunctionArgs()
	fe.body, _ = parseBlockStatement()

	return fe, true
}

func parseBlockStatement() (blockStatement, bool) {
	bs := blockStatement{}

	if !accept(tCURLY_LEFT) {
		return bs, false
	}

	bs.statements = []ast{}

	for !accept(tCURLY_RIGHT) {
		bs.statements = append(bs.statements, parseStatement())
	}

	return bs, true
}

func parseArrayLiteral() (arrayLiteral, bool) {
	al := arrayLiteral{}
	if !accept(tBRACKET_LEFT) {
		return al, false
	}

	al.items = []ast{}
	for !accept(tBRACKET_RIGHT) {
		al.items = append(al.items, parseYield())
		if !accept(tCOMMA) {
			expect(tBRACKET_RIGHT)
			break
		}
	}

	return al, true
}

func parsePropertyName() propertyName {
	name := propertyName{}
	if accept(tBRACKET_LEFT) {
		name.isCalculated = true
		name.value = parseSequence()
		expect(tBRACKET_RIGHT)
	} else if accept(tNAME) {
		name.value = atom{getToken()}
	} else if isValidPropertyName(tok.lexeme) || test(tNUMBER, tSTRING) {
		name.isNotName = true
		next()
		name.value = atom{getToken()}
	}
	return name
}

func parseForOfStatement(left ast) forOfStatement {
	fs := forOfStatement{}

	if vs, ok := left.(varDeclaration); ok {
		if len(vs.declarations) > 1 {
			panic(parsingError{p_WRONG_ASSIGNMENT, tok})
		}
	}

	fs.left = left
	fs.right = parseYield()
	expect(tPAREN_RIGHT)
	fs.body = parseStatement()

	return fs
}

func parseForInStatement(left ast) forInStatement {
	fs := forInStatement{}
	if vs, ok := left.(varDeclaration); ok {
		if len(vs.declarations) > 1 {
			panic(parsingError{p_WRONG_ASSIGNMENT, tok})
		}
	}

	fs.left = left
	fs.right = parseYield()
	expect(tPAREN_RIGHT)
	fs.body = parseStatement()

	return fs
}

// ast implementations

func (ps program) String() string {
	result := ""
	for _, s := range ps.statements {
		result += s.String()
	}
	return result
}

func (sw switchStatement) String() string {
	result := "switch(" + sw.condition.String() + "){"
	for _, c := range sw.cases {
		result += "case " + c.condition.String() + ":"
		for _, s := range c.statements {
			result += s.String()
		}
	}
	if sw.defCase != nil {
		result += "default:"
		for _, s := range sw.defCase {
			result += s.String()
		}
	}
	result += "}"

	return result
}

func (es exportStatement) String() string {
	result := "export "

	if es.exportAll {
		result += "*"
	} else if len(es.vars) == 1 && es.vars[0].alias.val.lexeme == "default" {
		result += "default "
		if es.declaration != nil {
			result += es.declaration.String()
		} else {
			result += es.vars[0].name.String()
		}
	} else {
		if es.declaration != nil {
			result += es.declaration.String()
		} else {
			result += "{"
			for i, ev := range es.vars {
				result += ev.name.String()
				if ev.alias.val.lexeme != "" {
					result += " as " + ev.alias.val.lexeme
				}
				if i < len(es.vars)-1 {
					result += ","
				}
			}
			result += "}"
		}
	}

	if es.source.val.lexeme != "" {
		result += " from" + es.source.val.lexeme
	}

	return result + ";"
}

func (is importStatement) String() string {
	result := "import"

	alreadyPrintedIndex := -1
	// print default import first
	for i, v := range is.vars {
		if v.name.val.lexeme == "default" {
			result += " " + v.alias.val.lexeme
			alreadyPrintedIndex = i
			break
		}
	}

	// print * import
	if is.all.val.lexeme != "" {
		if result != "import" {
			result += ","
		}
		result += "*as " + is.all.val.lexeme
	}

	// print all other imports
	if len(is.vars) > 0 && !(len(is.vars) == 1 && alreadyPrintedIndex == 0) {
		if result != "import" {
			result += ","
		}
		result += "{"
		for i, v := range is.vars {
			if i != alreadyPrintedIndex {
				result += v.name.String()
				if v.alias.String() != "" {
					result += " as " + v.alias.String()
				}
				if i < len(is.vars)-1 {
					result += ","
				}
			}
		}
		result += "}"
	} else {
		result += " "
	}
	if result != "import " {
		result += "from"
	}

	result += is.path.String() + ";"
	return result
}

func (re returnStatement) String() string {
	result := "return"
	if re.value != nil {
		result += " " + re.value.String()
	}
	result += ";"
	return result
}

func (fs functionStatement) String() string {
	return fs.funcExpr.String()
}

func (es expressionStatement) String() string {
	if es.expr != nil {
		return es.expr.String() + ";"
	}
	return ";"
}

func (is ifStatement) String() string {
	result := "if(" + is.test.String() + ")" + is.consequent.String()
	if is.alternate != nil {
		result += " else " + is.alternate.String()
	}
	return result
}

func (ws whileStatement) String() string {
	result := "while(" + ws.test.String() + ")" + ws.body.String()
	return result
}

func (ds doWhileStatement) String() string {
	result := "do " + ds.body.String() + "while" + "(" + ds.test.String() + ")" + ";"
	return result
}

func (es emptyStatement) String() string {
	return ";"
}

func (fs forStatement) String() string {
	result := "for(" + fs.init.String() + ";"
	result += fs.test.String() + ";"
	result += fs.final.String() + ")"
	result += fs.body.String()
	return result
}

func (op objectPattern) String() string {
	return objectLiteral.String(objectLiteral(op))
}

func (d declarator) String() string {
	result := ""
	if d.left != nil {
		result += d.left.String()
	}
	if d.value != nil {
		result += "=" + d.value.String()
	}
	return result
}

func (vd varDeclaration) String() string {
	result := vd.keyword.String() + " "
	for i, d := range vd.declarations {
		result += d.String()

		if i < len(vd.declarations)-1 {
			result += ","
		}
	}
	return result
}

func (se sequenceExpression) String() string {
	result := ""
	for i, item := range se.items {
		result += item.String()
		if i < len(se.items)-1 {
			result += ","
		}
	}
	return result
}

func (ye yieldExpression) String() string {
	return "yield " + ye.value.String()
}

func (be binaryExpression) String() string {
	result := ""
	op := operatorTable[be.operator.val.tType]
	switch t := be.left.(type) {
	case binaryExpression:
		leftPrec := operatorTable[t.operator.val.tType].precedence
		if leftPrec < op.precedence || (leftPrec == op.precedence && op.isRightAssociative) {
			result += "(" + t.String() + ")"
		} else {
			result += t.String()
		}
	default:
		result += t.String()
	}

	result += be.operator.String()

	switch t := be.right.(type) {
	case binaryExpression:
		rightPrec := operatorTable[t.operator.val.tType].precedence
		if (!op.isRightAssociative && rightPrec == op.precedence) || op.precedence > rightPrec {
			result += "(" + t.String() + ")"
		} else {
			result += t.String()
		}
	default:
		result += t.String()
	}
	return result
}

func (ce conditionalExpression) String() string {
	result := fmt.Sprintf(
		"%s?%s:%s", ce.test.String(), ce.consequent.String(), ce.alternate.String(),
	)
	return result
}

func (ue unaryExpression) String() string {
	result := ""
	if !ue.isPostfix {
		result += ue.operator.String()
	}
	if _, ok := ue.value.(atom); ok {
		result += ue.value.String()
	} else {
		result += "(" + ue.value.String() + ")"
	}
	if ue.isPostfix {
		result += ue.operator.String()
	}
	return result
}

func (fc functionCall) String() string {
	result := ""
	if fc.isConstructor {
		result += "new "
	}
	result += fc.name.String() + "("
	for i, arg := range fc.args {
		result += arg.String()
		if i < len(fc.args)-1 {
			result += ","
		}
	}
	result += ")"
	return result + ""
}

func (me memberExpression) String() string {
	result := me.object.String()
	if !me.property.isCalculated && !me.property.isNotName {
		result += "."
	}
	result += me.property.String()
	return result
}

func (av atom) String() string {
	return av.val.lexeme
}

func (le lambdaExpression) String() string {
	result := "("
	for i, arg := range le.args {
		result += arg.String()
		if i < len(le.args)-1 {
			result += ","
		}
	}
	result += ")=>"
	result += le.body.String()
	return result
}

func (pe parenExpression) String() string {
	return "(" + pe.inner.String() + ")"
}

func (ol objectLiteral) String() string {
	result := "{"
	for i, prop := range ol.properties {
		result += prop.key.String()

		fe, valueIsFunction := prop.value.(functionExpression)
		if prop.value != nil && (!valueIsFunction || !fe.isMember) {
			result += ":"
		}
		if prop.value != nil {
			result += prop.value.String()
		}
		if i < len(ol.properties)-1 {
			result += ","
		}
	}
	result += "}"
	return result
}

func (fe functionExpression) String() string {
	result := ""
	if !fe.isMember {
		result += "function"
	}
	if fe.name.val.lexeme != "" {
		result += " " + fe.name.String()
	}

	result += "("
	for i, arg := range fe.args {
		result += arg.String()
		if i < len(fe.args)-1 {
			result += ","
		}
	}
	result += ")"
	result += fe.body.String()
	return result
}

func (bs blockStatement) String() string {
	result := "{"
	for _, st := range bs.statements {
		result += st.String()
	}
	result += "}"
	return result
}

func (al arrayLiteral) String() string {
	result := "["
	for i, item := range al.items {
		result += item.String()
		if i < len(al.items)-1 {
			result += ","
		}
	}
	result += "]"
	return result
}

func (pn propertyName) String() string {
	if pn.isCalculated || pn.isNotName {
		return "[" + pn.value.String() + "]"
	}
	return pn.value.String()
}

func (fs forOfStatement) String() string {
	result := "for(" + fs.left.String()
	result += " of " + fs.right.String() + ")"
	result += fs.body.String()
	return result
}

func (fs forInStatement) String() string {
	result := "for(" + fs.left.String()
	result += " in " + fs.right.String() + ")"
	result += fs.body.String()
	return result
}

func (ap arrayPattern) String() string {
	result := "["
	for i, prop := range ap.properties {
		result += prop.String()
		if i < len(ap.properties)-1 {
			result += ","
		}
	}
	result += "]"
	return result
}
