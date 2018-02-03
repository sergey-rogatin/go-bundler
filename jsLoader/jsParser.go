package jsLoader

import "fmt"

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

type ast interface {
	String() string
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

var sourceTokens []token
var index int
var tok token

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

func accept(tTypes ...tokenType) bool {
	for _, tType := range tTypes {
		if tok.tType == tType {
			next()
			return true
		}
	}
	return false
}

func expect(tType tokenType) bool {
	if accept(tType) {
		return true
	}
	panic(fmt.Sprintf("\nExpected %s, got %s\n", tType, tok))
}

func getLexeme() string {
	return sourceTokens[index-1].lexeme
}

func getToken() token {
	return sourceTokens[index-1]
}

//
//
// parsing types
type program struct {
	statements []ast
}

func (ps program) String() string {
	result := ""
	for _, s := range ps.statements {
		result += s.String()
	}
	return result
}

func getProgram() program {
	p := program{}
	p.statements = []ast{}

	for !accept(tEND_OF_INPUT) {
		p.statements = append(p.statements, getStatement())
	}
	return p
}

func getStatement() ast {
	var st ast

	if s, ok := getEmptyStatement(); ok {
		st = s
	} else if s, ok := getForStatement(); ok {
		st = s
	} else if s, ok := getWhileStatement(); ok {
		st = s
	} else if s, ok := getDoWhileStatement(); ok {
		st = s
	} else if s, ok := getBlockStatement(); ok {
		st = s
	} else if s, ok := getIfStatement(); ok {
		st = s
	} else if s, ok := getFunctionStatement(); ok {
		st = s
	} else if s, ok := getReturnStatement(); ok {
		st = s
	} else if s, ok := getImportStatement(); ok {
		st = s
	} else if s, ok := getSwitchStatement(); ok {
		st = s
	} else {
		st = getExpressionStatement()
	}

	// else if s, ok := getExportStatement(); ok {
	// 	st = s
	// }
	return st
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

func getSwitchStatement() (switchStatement, bool) {
	sw := switchStatement{}
	if !accept(tSWITCH) {
		return sw, false
	}

	expect(tPAREN_LEFT)
	sw.condition = getSequence()
	expect(tPAREN_RIGHT)
	expect(tCURLY_LEFT)

	sw.cases = []switchCase{}
	for !accept(tCURLY_RIGHT) {
		if accept(tCASE) {
			c := switchCase{}
			c.condition = getSequence()
			c.statements = []ast{}

			expect(tCOLON)
			for tok.tType != tCASE && tok.tType != tDEFAULT && tok.tType != tCURLY_RIGHT {
				c.statements = append(c.statements, getStatement())
			}
			sw.cases = append(sw.cases, c)

			if accept(tDEFAULT) {
				expect(tCOLON)
				sw.defCase = []ast{}
				for tok.tType != tCASE && tok.tType != tCURLY_RIGHT {
					sw.defCase = append(sw.defCase, getStatement())
				}
			}
		}
	}

	return sw, true
}

type exportedVar struct {
	name  ast
	alias atom
}

type exportStatement struct {
	exportedVars []exportedVar
	declaration  ast
	source       atom
	exportAll    bool
}

func (es exportStatement) String() string {
	result := "export "

	if es.exportAll {
		result += "*"
	} else if len(es.exportedVars) == 1 && es.exportedVars[0].alias.val.lexeme == "default" {
		result += "default "
		if es.declaration != nil {
			result += es.declaration.String()
		} else {
			result += es.exportedVars[0].name.String()
		}
	} else {
		if es.declaration != nil {
			result += es.declaration.String()
		} else {
			result += "{"
			for i, ev := range es.exportedVars {
				result += ev.name.String()
				if ev.alias.val.lexeme != "" {
					result += " as " + ev.alias.val.lexeme
				}
				if i < len(es.exportedVars)-1 {
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

func getExportStatement() (exportStatement, bool) {
	es := exportStatement{}
	if !accept(tEXPORT) {
		return es, false
	}
	es.exportedVars = []exportedVar{}

	if accept(tDEFAULT) {
		ev := exportedVar{}
		ev.alias = atom{makeToken("default")}

		if fe, ok := getFunctionExpression(false); ok {
			if fe.name.val.lexeme != "" {
				es.declaration = fe
				ev.name = fe.name
			} else {
				ev.name = fe
			}
		} else {
			ev.name = getYield()
		}

		expect(tSEMI)
		es.exportedVars = append(es.exportedVars, ev)

	} else if accept(tCURLY_LEFT) {
		for !accept(tCURLY_RIGHT) {
			if accept(tNAME) {
				ev := exportedVar{}
				ev.name = atom{getToken()}
				if accept(tAS) {
					next()
					ev.alias = atom{getToken()}
				}
				es.exportedVars = append(es.exportedVars, ev)
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

	} else if vd, ok := getVarDeclaration(); ok {
		es.declaration = vd
		expect(tSEMI)
		for _, d := range vd.declarations {
			ev := exportedVar{}
			ev.name = d.left
			es.exportedVars = append(es.exportedVars, ev)
		}

	} else if fs, ok := getFunctionStatement(); ok {
		es.declaration = fs
		ev := exportedVar{}
		ev.name = fs.funcExpr.name
		es.exportedVars = append(es.exportedVars, ev)

	} else if accept(tMULT) {
		expect(tFROM)
		expect(tSTRING)
		es.source = atom{getToken()}
		es.exportAll = true
		expect(tSEMI)
	}

	return es, true
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

func getImportStatement() (importStatement, bool) {
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

	return is, true
}

type returnStatement struct {
	value ast
}

func (re returnStatement) String() string {
	result := "return"
	if re.value != nil {
		result += " " + re.value.String()
	}
	result += ";"
	return result
}

func getReturnStatement() (returnStatement, bool) {
	rs := returnStatement{}
	if !accept(tRETURN) {
		return rs, false
	}

	rs.value = getYield()
	return rs, true
}

type functionStatement struct {
	funcExpr functionExpression
}

func (fs functionStatement) String() string {
	return fs.funcExpr.String()
}

func getFunctionStatement() (functionStatement, bool) {
	fs := functionStatement{}

	if fe, ok := getFunctionExpression(false); ok {
		if fe.name.val.lexeme == "" {
			panic("Function declaration must be named! " + getToken().String())
		} else {
			fs.funcExpr = fe
			return fs, true
		}
	}

	return fs, false
}

type expressionStatement struct {
	expr ast
}

func (es expressionStatement) String() string {
	if es.expr != nil {
		return es.expr.String() + ";"
	}
	return ";"
}

func getExpressionStatement() expressionStatement {
	es := expressionStatement{}
	if vd, ok := getVarDeclaration(); ok {
		es.expr = vd
	} else if accept(tCONTINUE, tDEBUGGER, tBREAK) {
		es.expr = atom{getToken()}
	} else {
		es.expr = getSequence()
	}
	expect(tSEMI)
	return es
}

type ifStatement struct {
	test       ast
	consequent ast
	alternate  ast
}

func (is ifStatement) String() string {
	result := "if(" + is.test.String() + ")" + is.consequent.String()
	if is.alternate != nil {
		result += " else " + is.alternate.String()
	}
	return result
}

func getIfStatement() (ifStatement, bool) {
	is := ifStatement{}

	if !accept(tIF) {
		return is, false
	}
	expect(tPAREN_LEFT)
	is.test = getSequence()
	expect(tPAREN_RIGHT)
	is.consequent = getStatement()

	if accept(tELSE) {
		is.alternate = getStatement()
	}

	return is, true
}

type whileStatement struct {
	test ast
	body ast
}

func (ws whileStatement) String() string {
	result := "while(" + ws.test.String() + ")" + ws.body.String()
	return result
}

func getWhileStatement() (whileStatement, bool) {
	ws := whileStatement{}
	if !accept(tWHILE) {
		return ws, false
	}
	expect(tPAREN_LEFT)
	ws.test = getSequence()
	expect(tPAREN_RIGHT)
	ws.body = getStatement()

	return ws, true
}

type doWhileStatement struct {
	test ast
	body ast
}

func (ds doWhileStatement) String() string {
	result := "do " + ds.body.String() + "while" + "(" + ds.test.String() + ")" + ";"
	return result
}

func getDoWhileStatement() (doWhileStatement, bool) {
	ds := doWhileStatement{}

	if !accept(tDO) {
		return ds, false
	}

	ds.body = getStatement()
	expect(tWHILE)
	expect(tPAREN_LEFT)
	ds.test = getSequence()
	expect(tPAREN_RIGHT)

	return ds, true
}

type emptyStatement struct{}

func (es emptyStatement) String() string {
	return ";"
}

func getEmptyStatement() (emptyStatement, bool) {
	es := emptyStatement{}
	if accept(tSEMI) {
		return es, true
	}
	return es, false
}

type forStatement struct {
	init  ast
	test  ast
	final ast
	body  ast
}

func (fs forStatement) String() string {
	result := "for(" + fs.init.String() + ";"
	result += fs.test.String() + ";"
	result += fs.final.String() + ")"
	result += fs.body.String()
	return result
}

func getForStatement() (ast, bool) {
	if !accept(tFOR) {
		return nil, false
	}
	expect(tPAREN_LEFT)

	var init ast
	if vd, ok := getVarDeclaration(); ok {
		init = vd
	} else if tok.tType == tSEMI {
		init = atom{token{lexeme: ""}}
	} else {
		init = getSequence()
	}

	if accept(tOF) {
		return getForOfStatement(init), true
	} else if accept(tIN) {
		return getForInStatement(init), true
	}

	fs := forStatement{}
	fs.init = init
	expect(tSEMI)

	if accept(tSEMI) {
		fs.test = atom{token{lexeme: ""}}
	} else {
		fs.test = getSequence()
		expect(tSEMI)
	}

	if accept(tPAREN_RIGHT) {
		fs.final = atom{token{lexeme: ""}}
	} else {
		fs.final = getSequence()
		expect(tPAREN_RIGHT)
	}
	fs.body = getStatement()

	return fs, true
}

type objectPattern struct {
	properties []objectProperty
}

func (op objectPattern) String() string {
	return objectLiteral.String(objectLiteral(op))
}

func getObjectPattern() (objectPattern, bool) {
	op := objectPattern{}
	if !accept(tCURLY_LEFT) {
		return op, false
	}
	op.properties = []objectProperty{}

	for !accept(tCURLY_RIGHT) {
		prop := objectProperty{}
		prop.key = getPropertyName()

		if accept(tCOLON) {
			prop.value = getAssignment()
		}

		op.properties = append(op.properties, prop)

		if !accept(tCOMMA) {
			expect(tCURLY_RIGHT)
			break
		}
	}

	return op, true
}

type declarator struct {
	left  ast
	value ast
}

func (d declarator) String() string {
	result := d.left.String()
	if d.value != nil {
		result += "=" + d.value.String()
	}
	return result
}

func getDeclarator() (declarator, bool) {
	d := declarator{}

	if accept(tNAME) {
		d.left = atom{getToken()}
	} else if op, ok := getObjectPattern(); ok {
		d.left = op
	} else {
		return d, false
	}

	if accept(tASSIGN) {
		d.value = getAssignment()
	}

	return d, true
}

type varDeclaration struct {
	declarations []declarator
	keyword      atom
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

func getVarDeclaration() (varDeclaration, bool) {
	vd := varDeclaration{}
	if !accept(tVAR, tLET, tCONST) {
		return vd, false
	}
	vd.keyword = atom{getToken()}
	vd.declarations = []declarator{}

	for ok := true; ok; ok = accept(tCOMMA) {
		if d, ok := getDeclarator(); ok {
			vd.declarations = append(vd.declarations, d)
		} else {
			panic("Unexpected token " + tok.String())
		}
	}

	return vd, true
}

type sequenceExpression struct {
	items []ast
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

func getSequence() ast {
	firstItem := getYield()

	se := sequenceExpression{}
	se.items = []ast{firstItem}
	for accept(tCOMMA) {
		se.items = append(se.items, getYield())
	}
	if len(se.items) > 1 {
		return se
	}

	return firstItem
}

type yieldExpression struct {
	value ast
}

func (ye yieldExpression) String() string {
	return "yield " + ye.value.String()
}

func getYield() ast {
	if accept(tYIELD) {
		ye := yieldExpression{}
		ye.value = getYield()
		return ye
	}
	return getAssignment()
}

type binaryExpression struct {
	left     ast
	right    ast
	operator atom
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

func getAssignment() ast {
	if op, ok := getObjectPattern(); ok {
		binExpr := binaryExpression{}
		binExpr.left = op

		if accept(tASSIGN) {
			binExpr.operator = atom{getToken()}
			binExpr.right = getAssignment()
			return binExpr
		}
		return binExpr.left
	}

	left := getConditional()

	if accept(
		tASSIGN, tPLUS_ASSIGN, tMINUS_ASSIGN, tMULT_ASSIGN, tDIV_ASSIGN,
		tBITWISE_SHIFT_LEFT_ASSIGN, tBITWISE_SHIFT_RIGHT_ASSIGN,
		tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN, tBITWISE_AND_ASSIGN,
		tBITWISE_OR_ASSIGN, tBITWISE_XOR_ASSIGN,
	) {
		binExpr := binaryExpression{}
		binExpr.operator = atom{getToken()}
		binExpr.left = left
		binExpr.right = getAssignment()

		return binExpr
	}

	return left
}

type conditionalExpression struct {
	test       ast
	consequent ast
	alternate  ast
}

func (ce conditionalExpression) String() string {
	result := fmt.Sprintf(
		"%s?%s:%s", ce.test.String(), ce.consequent.String(), ce.alternate.String(),
	)
	return result
}

func getConditional() ast {
	test := getBinary()

	if accept(tQUESTION) {
		condExp := conditionalExpression{}
		condExp.test = test
		condExp.consequent = getConditional()
		expect(tCOLON)
		condExp.alternate = getConditional()
		return condExp
	}

	return test
}

func getBinary() ast {
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
		outputStack = append(outputStack, getPrefixUnary())

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

type unaryExpression struct {
	value     ast
	operator  atom
	isPostfix bool
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

func getPrefixUnary() ast {
	if accept(
		tNOT, tBITWISE_NOT, tPLUS, tMINUS,
		tINC, tDEC, tTYPEOF, tVOID, tDELETE,
	) {
		unExp := unaryExpression{}
		unExp.operator = atom{getToken()}
		unExp.value = getPostfixUnary()
		return unExp
	}

	return getPostfixUnary()
}

func getPostfixUnary() ast {
	value := getFunctionCallOrMember(false)

	if accept(tINC, tDEC) {
		unExp := unaryExpression{}
		unExp.isPostfix = true
		unExp.operator = atom{getToken()}
		unExp.value = value
		return unExp
	}

	return value
}

type functionCall struct {
	name          ast
	args          []ast
	isConstructor bool
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

func getFunctionCallOrMember(noFuncCall bool) ast {
	funcName := getConstructorCall()

	for {
		if !noFuncCall && accept(tPAREN_LEFT) {
			fc := functionCall{}
			fc.name = funcName
			fc.args = []ast{}

			for !accept(tPAREN_RIGHT) {
				fc.args = append(fc.args, getYield())
				if !accept(tCOMMA) {
					expect(tPAREN_RIGHT)
					break
				}
			}
			funcName = fc
		} else if accept(tDOT) || tok.tType == tBRACKET_LEFT {
			me := memberExpression{}
			me.object = funcName
			me.property = getPropertyName()
			funcName = me
		} else {
			break
		}
	}

	return funcName
}

func getConstructorCall() ast {
	if accept(tNEW) {
		fc := functionCall{}
		fc.name = getFunctionCallOrMember(true)
		fc.isConstructor = true

		if accept(tPAREN_LEFT) {
			fc.args = []ast{}

			for !accept(tPAREN_RIGHT) {
				fc.args = append(fc.args, getYield())
				if !accept(tCOMMA) {
					expect(tPAREN_RIGHT)
					break
				}
			}
		}

		return fc
	}

	return getAtom()
}

type memberExpression struct {
	object   ast
	property propertyName
}

func (me memberExpression) String() string {
	result := me.object.String()
	if !me.property.isCalculated && !me.property.isNotName {
		result += "."
	}
	result += me.property.String()
	return result
}

type atom struct {
	val token
}

func (av atom) String() string {
	return av.val.lexeme
}

func getAtom() ast {
	//fmt.Println(tok)

	if a, ok := getLambda(); ok {
		return a
	} else if a, ok := getParens(); ok {
		return a
	} else if a, ok := getObjectLiteral(); ok {
		return a
	} else if a, ok := getFunctionExpression(false); ok {
		return a
	} else if a, ok := getArrayLiteral(); ok {
		return a
	} else if accept(tNAME, tNUMBER, tSTRING, tTRUE, tFALSE, tNULL, tUNDEFINED, tTHIS) {
		return atom{getToken()}
	} else {
		expect(tEND_OF_INPUT)
	}
	return nil
}

type lambdaExpression struct {
	args []declarator
	body ast
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

func getLambda() (lambdaExpression, bool) {
	startPos := index

	le := lambdaExpression{}
	if accept(tNAME) {
		le.args = []declarator{declarator{left: atom{getToken()}}}
	} else if args, ok := getFunctionArgs(); ok {
		le.args = args
	}

	if le.args != nil {
		if !accept(tLAMBDA) {
			backtrack(startPos)
			return le, false
		}

		if bs, ok := getBlockStatement(); ok {
			le.body = bs
		} else {
			le.body = getYield()
		}
		return le, true
	}

	backtrack(startPos)
	return le, false
}

func getFunctionArgs() ([]declarator, bool) {
	startPos := index
	if !accept(tPAREN_LEFT) {
		return nil, false
	}

	args := []declarator{}
	for !accept(tPAREN_RIGHT) {
		if d, ok := getDeclarator(); ok {
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

func getParens() (ast, bool) {
	if !accept(tPAREN_LEFT) {
		return nil, false
	}
	res := getSequence()
	expect(tPAREN_RIGHT)
	return res, true
}

type objectProperty struct {
	key   propertyName
	value ast
}

type objectLiteral struct {
	properties []objectProperty
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

func getObjectLiteral() (objectLiteral, bool) {
	ol := objectLiteral{}

	if !accept(tCURLY_LEFT) {
		return ol, false
	}

	ol.properties = []objectProperty{}

	for !accept(tCURLY_RIGHT) {
		prop := objectProperty{}
		prop.key = getPropertyName()

		if accept(tCOLON) {
			prop.value = getYield()
		} else if fe, ok := getFunctionExpression(true); ok {
			prop.value = fe
		} else if prop.key.isCalculated || prop.key.isNotName {
			expect(tEND_OF_INPUT)
		}

		ol.properties = append(ol.properties, prop)
		if !accept(tCOMMA) {
			expect(tCURLY_RIGHT)
			break
		}
	}

	return ol, true
}

type functionExpression struct {
	name     atom
	args     []declarator
	body     blockStatement
	isMember bool
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

func getFunctionExpression(isMember bool) (functionExpression, bool) {
	fe := functionExpression{}
	if isMember {
		if tok.tType != tPAREN_LEFT {
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

	fe.args, _ = getFunctionArgs()
	fe.body, _ = getBlockStatement()

	return fe, true
}

type blockStatement struct {
	statements []ast
}

func (bs blockStatement) String() string {
	result := "{"
	for _, st := range bs.statements {
		result += st.String()
	}
	result += "}"
	return result
}

func getBlockStatement() (blockStatement, bool) {
	bs := blockStatement{}

	if !accept(tCURLY_LEFT) {
		return bs, false
	}

	bs.statements = []ast{}

	for !accept(tCURLY_RIGHT) {
		bs.statements = append(bs.statements, getStatement())
	}

	return bs, true
}

type arrayLiteral struct {
	items []ast
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

func getArrayLiteral() (arrayLiteral, bool) {
	al := arrayLiteral{}
	if !accept(tBRACKET_LEFT) {
		return al, false
	}

	al.items = []ast{}
	for !accept(tBRACKET_RIGHT) {
		al.items = append(al.items, getYield())
		if !accept(tCOMMA) {
			expect(tBRACKET_RIGHT)
			break
		}
	}

	return al, true
}

type propertyName struct {
	value        ast
	isCalculated bool
	isNotName    bool
}

func (pn propertyName) String() string {
	if pn.isCalculated || pn.isNotName {
		return "[" + pn.value.String() + "]"
	}
	return pn.value.String()
}

func getPropertyName() propertyName {
	name := propertyName{}
	if accept(tBRACKET_LEFT) {
		name.isCalculated = true
		name.value = getSequence()
		expect(tBRACKET_RIGHT)
	} else if accept(tNAME) {
		name.value = atom{getToken()}
	} else if isValidPropertyName(tok.lexeme) || tok.tType == tNUMBER || tok.tType == tSTRING {
		name.isNotName = true
		next()
		name.value = atom{getToken()}
	}
	return name
}

type forOfStatement struct {
	left  ast
	right ast
	body  ast
}

func (fs forOfStatement) String() string {
	result := "for(" + fs.left.String()
	result += " of " + fs.right.String() + ")"
	result += fs.body.String()
	return result
}

func getForOfStatement(left ast) forOfStatement {
	fs := forOfStatement{}
	if vs, ok := left.(varDeclaration); ok {
		if len(vs.declarations) > 1 {
			panic("Wrong left side in for-of statement")
		}
	}
	fs.left = left
	fs.right = getYield()
	expect(tPAREN_RIGHT)
	fs.body = getStatement()

	return fs
}

type forInStatement struct {
	left  ast
	right ast
	body  ast
}

func (fs forInStatement) String() string {
	result := "for(" + fs.left.String()
	result += " in " + fs.right.String() + ")"
	result += fs.body.String()
	return result
}

func getForInStatement(left ast) forInStatement {
	fs := forInStatement{}
	if vs, ok := left.(varDeclaration); ok {
		if len(vs.declarations) > 1 {
			panic("Wrong left side in for-of statement")
		}
	}
	fs.left = left
	fs.right = getYield()
	expect(tPAREN_RIGHT)
	fs.body = getStatement()

	return fs
}

func parse(localSrc []token) program {
	sourceTokens = localSrc
	index = 0
	tok = sourceTokens[index]

	return getProgram()
}
