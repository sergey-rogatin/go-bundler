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

type varStatement struct {
	decls            []declarator
	keyword          string
	isInForStatement bool
}

func (vs varStatement) String() string {
	result := vs.keyword + " "
	for i, decl := range vs.decls {
		result += decl.name.String()
		if decl.value != nil {
			result += "=" + decl.value.String()
		}
		if i < len(vs.decls)-1 {
			result += ","
		}
	}
	if !vs.isInForStatement {
		result += ";"
	}
	return result
}

type importedVar struct {
	name      atom
	pseudonym atom
}

type importStatement struct {
	path string
	vars []importedVar
}

func (is importStatement) String() string {
	result := "import{"
	for i, v := range is.vars {
		result += v.name.String()
		if v.pseudonym.String() != "" {
			result += " as " + v.pseudonym.String()
		}
		if i < len(is.vars)-1 {
			result += ","
		}
	}
	result += "}from" + is.path + ";"
	return result
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

type whileStatement struct {
	test ast
	body ast
}

func (ws whileStatement) String() string {
	result := "while(" + ws.test.String() + ")" + ws.body.String()
	return result
}

type exportedVar struct {
	name      atom
	pseudonym atom
	value     ast
}

type exportStatement struct {
	exportedVars  []exportedVar
	isDeclaration bool
}

func (es exportStatement) String() string {
	result := "export "
	for _, exp := range es.exportedVars {
		result += exp.name.String()
		if exp.pseudonym.String() != "" {
			result += " as " + exp.pseudonym.String()
		}
		if exp.value != nil {
			result += " = " + exp.value.String()
		}
	}
	return result + ";"
}

type objectPattern struct {
	properties []objectProperty
}

func (op objectPattern) String() string {
	result := "{"
	for i, prop := range op.properties {
		result += prop.key.String()
		if prop.value != nil {
			result += ":" + prop.value.String()
		}

		if i < len(op.properties)-1 {
			result += ","
		}
	}
	result += "}"
	return result
}

type assignmentPattern struct {
	left  atom
	right ast
}

func (ap assignmentPattern) String() string {
	result := ap.left.String() + "=" + ap.right.String()
	return result
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

type expressionStatement struct {
	expr ast
}

func (es expressionStatement) String() string {
	if es.expr != nil {
		return es.expr.String() + ";"
	}
	return ";"
}

type ast interface {
	String() string
}

type operatorInfo struct {
	precedence         int
	isRightAssociative bool
}

var operatorTable = map[tokenType]operatorInfo{
	tOR:                       operatorInfo{5, false},
	tAND:                      operatorInfo{6, false},
	tBITWISE_OR:               operatorInfo{7, false},
	tBITWISE_XOR:              operatorInfo{8, false},
	tBITWISE_AND:              operatorInfo{9, false},
	tEQUALS:                   operatorInfo{10, false},
	tNOT_EQUALS:               operatorInfo{10, false},
	tEQUALS_STRICT:            operatorInfo{10, false},
	tNOT_EQUALS_STRICT:        operatorInfo{10, false},
	tLESS:                     operatorInfo{11, false},
	tLESS_EQUALS:              operatorInfo{11, false},
	tGREATER:                  operatorInfo{11, false},
	tGREATER_EQUALS:           operatorInfo{11, false},
	tIN:                       operatorInfo{11, false},
	tINSTANCEOF:               operatorInfo{11, false},
	tBITWISE_SHIFT_LEFT:       operatorInfo{12, false},
	tBITWISE_SHIFT_RIGHT:      operatorInfo{12, false},
	tBITWISE_SHIFT_RIGHT_ZERO: operatorInfo{12, false},
	tPLUS:  operatorInfo{13, false},
	tMINUS: operatorInfo{13, false},
	tMULT:  operatorInfo{14, false},
	tDIV:   operatorInfo{14, false},
	tMOD:   operatorInfo{14, false},
	tEXP:   operatorInfo{15, true},
}

var src []token

var i = 0
var t = src[i]

func next() {
	i++
	if i < len(src) {
		t = src[i]
	}
}

func backtrack(backIndex int) {
	i = backIndex
	t = src[i]
}

func accept(tTypes ...tokenType) bool {
	for _, tType := range tTypes {
		if t.tType == tType {
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
	panic(fmt.Sprintf("\nExpected %s, got %s\n", tType, t))
}

func getLexeme() string {
	return src[i-1].lexeme
}

func getToken() token {
	return src[i-1]
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
	} else if s, ok := getImportStatement(); ok {
		st = s
	} else if s, ok := getExportStatement(); ok {
		st = s
	} else if s, ok := getBlockStatement(); ok {
		st = s
	} else if s, ok := getIfStatement(); ok {
		st = s
	} else if s, ok := getVarStatement(); ok {
		st = s
	} else if s, ok := getFunctionStatement(); ok {
		st = s
	} else if s, ok := getReturnStatement(); ok {
		st = s
	} else if s, ok := getContinueStatement(); ok {
		st = s
	}
	return st
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
	result := "for (" + fs.init.String() + ";"
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
	fs.test = getSequence()
	expect(tSEMI)
	fs.final = getSequence()
	expect(tPAREN_RIGHT)
	fs.body = getStatement()

	return fs, true
}

type declarator struct {
	name  atom
	value ast
}

func (d declarator) String() string {
	result := d.name.String()
	if d.value != nil {
		result += "=" + d.value.String()
	}
	return result
}

func getDeclarator() (declarator, bool) {
	d := declarator{}

	if accept(tNAME) {
		d.name = atom{getToken()}

		if accept(tASSIGN) {
			d.value = getYield()
		}
		return d, true
	}

	return d, false
}

type varDeclaration struct {
	declarations []declarator
	keyword      atom
}

func (vd varDeclaration) String() string {
	result := vd.keyword.String()
	for _, d := range vd.declarations {
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
			panic("Unexpected token " + t.String())
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
	se.items = []ast{}
	for ok := true; ok; ok = accept(tCOMMA) {
		se.items = append(se.items, getYield())
	}
	if len(se.items) > 0 {
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
		if rightPrec <= op.precedence {
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
	// TODO dont assign to literals
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

		op, ok := operatorTable[t.tType]
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
		opStack = append(opStack, t)
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
	value := getFunctionCall()

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

func getFunctionCall() ast {
	funcName := getConstructorCall()

	if accept(tPAREN_LEFT) {
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
		return fc
	}
	return funcName
}

func getConstructorCall() ast {
	if accept(tNEW) {
		fc := functionCall{}
		fc.name = getConstructorCall()
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

	return getMember()
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

func getMember() ast {
	object := getAtom()

	me := memberExpression{}
	if accept(tDOT) || t.tType == tBRACKET_LEFT {
		prop := getPropertyName()
		me.object = getMember()
		me.property = prop
	} else {
		return object
	}

	me.object = object
	return me
}

type atom struct {
	val token
}

func (av atom) String() string {
	return av.val.lexeme
}

func getAtom() ast {
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
	} else if accept(tNAME, tNUMBER, tSTRING, tFALSE, tNULL, tUNDEFINED, tTHIS) {
		return atom{getToken()}
	} else {
		panic("Unexpected token " + t.String())
	}
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
	startPos := i

	le := lambdaExpression{}
	if accept(tNAME) {
		le.args = []declarator{declarator{name: atom{getToken()}}}
	} else if args, ok := getFunctionArgs(); ok {
		le.args = args
	}

	if le.args != nil {
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
	startPos := i
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
			expect(tPAREN_RIGHT)
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
		}

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
	result += fe.name.String()

	for i, arg := range fe.args {
		result += arg.String()
		if i < len(fe.args)-1 {
			result += ","
		}
	}
	return result
}

func getFunctionExpression(isMember bool) (functionExpression, bool) {
	fe := functionExpression{}
	if isMember {
		if t.tType != tPAREN_LEFT {
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
	} else if isValidPropertyName(t.lexeme) || t.tType == tNUMBER || t.tType == tSTRING {
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
	result += " of " + fs.right.String() + ")"
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

/*







 */

func parse(src []token) program {
	i := 0
	t := src[i]

	next := func() {
		i++
		if i < len(src) {
			t = src[i]
		}
	}

	backtrack := func(backIndex int) {
		i = backIndex
		t = src[i]
	}

	accept := func(tTypes ...tokenType) bool {
		for _, tType := range tTypes {
			if t.tType == tType {
				next()
				return true
			}
		}
		return false
	}

	expect := func(tType tokenType) bool {
		if accept(tType) {
			return true
		}
		panic(fmt.Sprintf("\nExpected %s, got %s\n", tType, t))
	}

	getLexeme := func() string {
		return src[i-1].lexeme
	}

	getToken := func() token {
		return src[i-1]
	}

	type grammar func() ast

	var (
		getStatement,
		getAtom,
		getMember,
		getConstructorCall,
		getFunctionCall,
		getPostfixUnary,
		getPrefixUnary,
		getBinaryExp,
		getTernary,
		getAssign,
		getYield,
		getSpread,
		getSequence grammar
	)

	getMember = func() ast {
		var result ast

		object := getAtom()

		if accept(tDOT) || src[i].tType == tBRACKET_LEFT {
			atomStack := []ast{object}
			isCalcStack := []bool{src[i].tType == tBRACKET_LEFT}

			dotCounter := 0

			for ok := true; ok; ok = accept(tDOT) || src[i].tType == tBRACKET_LEFT {
				objProp := getFunctionCall()
				atomStack = append(atomStack, objProp)
				isCalcStack = append(isCalcStack, src[i].tType == tBRACKET_LEFT)
				dotCounter++
			}

			for i := 0; i < dotCounter; i++ {
				newMe := memberExpression{}
				newMe.object = atomStack[i]
				newMe.property = atomStack[i+1]
				newMe.isCalculated = isCalcStack[i]
				atomStack[i+1] = newMe
			}

			result = atomStack[dotCounter]
		} else {
			result = object
		}

		return result
	}

	getConstructorCall = func() ast {
		var result ast

		if accept(tNEW) {
			fc := functionCall{}
			fc.name = getConstructorCall()
			fc.isConstructor = true
			fc.args = make([]ast, 0)
			if accept(tPAREN_LEFT) {
				for !accept(tPAREN_RIGHT) {
					fc.args = append(fc.args, getSpread())
					if !accept(tCOMMA) {
						expect(tPAREN_RIGHT)
						break
					}
				}
			}
			result = fc
		} else {
			result = getMember()
		}

		return result
	}

	getFunctionCall = func() ast {
		var result ast

		funcName := getConstructorCall()
		switch {
		case accept(tPAREN_LEFT):
			fc := functionCall{}
			fc.name = funcName
			fc.args = make([]ast, 0)
			for !accept(tPAREN_RIGHT) {
				fc.args = append(fc.args, getSpread())
				if !accept(tCOMMA) {
					expect(tPAREN_RIGHT)
					break
				}
			}
			result = fc

		default:
			result = funcName
		}

		return result
	}

	getPostfixUnary = func() ast {
		var result ast

		value := getFunctionCall()
		if accept(tINC, tDEC) {
			unExp := unaryExpression{}
			unExp.isPostfix = true
			unExp.operator = getLexeme()
			unExp.value = value
			result = unExp
		} else {
			result = value
		}

		return result
	}

	getPrefixUnary = func() ast {
		var result ast

		if accept(
			tNOT, tBITWISE_NOT, tPLUS, tMINUS,
			tINC, tDEC, tTYPEOF, tVOID, tDELETE,
		) {
			unExp := unaryExpression{}
			unExp.operator = getLexeme()
			unExp.value = getPostfixUnary()
			result = unExp
		} else {
			result = getPostfixUnary()
		}

		return result
	}

	getBinaryExp = func() ast {
		opStack := make([]token, 0)
		outputStack := make([]ast, 0)

		var root *binaryExpression

		addNode := func(t token) {
			newNode := binaryExpression{}
			newNode.right = outputStack[len(outputStack)-1]
			newNode.left = outputStack[len(outputStack)-2]
			newNode.operator = t

			outputStack = outputStack[:len(outputStack)-2]
			outputStack = append(outputStack, newNode)
			root = &newNode
		}

		for {
			outputStack = append(outputStack, getPrefixUnary())

			op, ok := operatorTable[t.tType]
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
			opStack = append(opStack, t)
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

	getTernary = func() ast {
		var result ast

		test := getBinaryExp()
		if accept(tQUESTION) {
			condExp := conditionalExpression{}
			condExp.test = test
			condExp.consequent = getTernary()
			expect(tCOLON)
			condExp.alternate = getTernary()
			result = condExp
		} else {
			result = test
		}

		return result
	}

	getAssign = func() ast {
		var result ast

		// TODO dont assign to literals
		left := getTernary()
		if accept(
			tASSIGN, tPLUS_ASSIGN, tMINUS_ASSIGN, tMULT_ASSIGN, tDIV_ASSIGN,
			tBITWISE_SHIFT_LEFT_ASSIGN, tBITWISE_SHIFT_RIGHT_ASSIGN, tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN,
			tBITWISE_AND_ASSIGN, tBITWISE_OR_ASSIGN, tBITWISE_XOR_ASSIGN,
		) {
			binExpr := binaryExpression{}
			binExpr.operator = src[i-1]
			binExpr.left = left
			binExpr.right = getAssign()
			result = binExpr
		} else {
			result = left
		}

		return result
	}

	// TODO:
	// yield
	// spread
	getYield = func() ast {
		return getAssign()
	}
	getSpread = func() ast {
		return getYield()
	}

	getSequence = func() ast {
		var result ast

		firstInSeq := getSpread()
		if accept(tCOMMA) {
			seqExp := sequenceExpression{[]ast{firstInSeq}}
			for ok := true; ok; ok = accept(tCOMMA) {
				seqExp.items = append(seqExp.items, getSpread())
			}
			result = seqExp
		} else {
			result = firstInSeq
		}

		return result
	}

	getVarDeclaration := func() varStatement {
		keyword := getLexeme()
		varSt := varStatement{make([]declaration, 0), keyword, false}
		for ok := true; ok; ok = accept(tCOMMA) {
			decl := declaration{}

			// descructuring declaration
			if accept(tCURLY_LEFT) {
				objPat := objectPattern{}
				objPat.properties = make([]objectProperty, 0)
				for !accept(tCURLY_RIGHT) {
					prop := getObjectKey()
					if accept(tCOLON) {
						expect(tNAME)
						valueName := getToken()
						if accept(tASSIGN) {
							assP := assignmentPattern{}
							assP.left = atom{valueName}
							assP.right = getSpread()
							prop.value = assP
						} else {
							prop.value = atom{valueName}
						}
					}
					objPat.properties = append(objPat.properties, prop)

					if !accept(tCOMMA) {
						expect(tCURLY_RIGHT)
						break
					}
				}
				decl.name = objPat
			} else {
				expect(tNAME)
				decl.name = atom{getToken()}
			}

			if accept(tASSIGN) {
				decl.value = getSpread()
			}
			varSt.decls = append(varSt.decls, decl)
		}
		return varSt
	}

	getStatement = func() ast {
		var result ast

		switch {

		case accept(tFOR):
			expect(tPAREN_LEFT)

			var init ast
			if accept(tVAR, tCONST, tLET) {
				decl := getVarDeclaration()
				decl.isInForStatement = true
				init = decl
			} else {
				init = getSequence()
			}

			if accept(tIN, tOF) {
				fio := forInOfStatement{}
				if vs, ok := init.(varStatement); ok {
					if len(vs.decls) > 1 {
						panic("Invalid left-hand expression \"" + vs.String() + "\"")
					}
				}
				fio.value = init
				fio.keyword = getLexeme()
				fio.iterable = getSpread()
				expect(tPAREN_RIGHT)
				fio.body = getStatement()
				result = fio
			} else {
				fs := forStatement{}
				fs.init = init
				expect(tSEMI)
				fs.test = getSequence()
				expect(tSEMI)
				if !accept(tPAREN_RIGHT) {
					fs.final = getSequence()
					expect(tPAREN_RIGHT)
				} else {
					fs.final = atom{}
				}
				fs.body = getStatement()
				result = fs
			}

		case accept(tSEMI):
			result = expressionStatement{nil}

		case accept(tEXPORT):
			es := exportStatement{}
			es.exportedVars = make([]exportedVar, 0)

			if accept(tFUNCTION) {
				es.isDeclaration = true
				expVar := exportedVar{}
				name := atom{token{tType: tNAME, lexeme: ""}}
				if accept(tNAME) {
					name = atom{getToken()}
				}
				decl := getFunctionExpression()
				decl.name = name
				expVar.value = decl
				expVar.name = name
				es.exportedVars = append(es.exportedVars, expVar)
			} else if t.tType == tVAR || t.tType == tCONST || t.tType == tLET {
				es.isDeclaration = true
				decl := getStatement()
				varSt := decl.(varStatement)
				for _, declarator := range varSt.decls {
					expVar := exportedVar{}
					expVar.name = declarator.name.(atom)
					expVar.pseudonym = atom{token{tType: tNAME, lexeme: ""}}
					expVar.value = declarator.value
					es.exportedVars = append(es.exportedVars, expVar)
				}
			} else if accept(tDEFAULT) {
				expVar := exportedVar{}
				expVar.name = atom{token{tType: tNAME, lexeme: "default"}}
				if accept(tFUNCTION) {
					es.isDeclaration = true
					name := atom{token{tType: tNAME, lexeme: ""}}
					if accept(tNAME) {
						name = atom{getToken()}
					}
					decl := getFunctionExpression()
					decl.name = name
					expVar.value = decl
				} else if accept(tCLASS) {

				} else {
					expVar.pseudonym = atom{token{tType: tNAME, lexeme: ""}}
					expVar.value = getSpread()
				}
				es.exportedVars = append(es.exportedVars, expVar)
				expect(tSEMI)
			}
			result = es

		case accept(tWHILE):
			wlSt := whileStatement{}
			expect(tPAREN_LEFT)
			wlSt.test = getSequence()
			expect(tPAREN_RIGHT)
			wlSt.body = getStatement()
			result = wlSt

		case accept(tCURLY_LEFT):
			blSt := blockStatement{}
			blSt.statements = make([]ast, 0)
			for !accept(tCURLY_RIGHT) {
				blSt.statements = append(blSt.statements, getStatement())
			}
			result = blSt

		case accept(tIF):
			ifSt := ifStatement{}
			expect(tPAREN_LEFT)
			ifSt.test = getSequence()
			expect(tPAREN_RIGHT)
			ifSt.consequent = getStatement()
			if accept(tELSE) {
				ifSt.alternate = getStatement()
			}
			result = ifSt

		// tVAR tNAME [tEQUALS add] {tCOMMA tNAME [tEQUALS add]} tSEMI
		case accept(tVAR) || accept(tCONST) || accept(tLET):
			result = getVarDeclaration()
			expect(tSEMI)

		// // tIMPORT [tNAME tFROM] [tCURLY_LEFT [tDEFAULT tAS tNAME] [tNAME {tCOMMA tNAME}] [tCOMMA] tCURLY_RIGHT tFROM] tSTRING;
		case accept(tIMPORT):
			impSt := importStatement{}
			impSt.vars = make([]importedVar, 0)
			// import { foo, bar as i } from "module-name";
			if accept(tNAME) {
				defVar := importedVar{}
				defVar.pseudonym = atom{getToken()}
				defVar.name = atom{token{lexeme: "default", tType: tNAME}}

				impSt.vars = append(impSt.vars, defVar)
				if src[i+1].tType == tCURLY_LEFT {
					expect(tCOMMA)
				} else {
					expect(tFROM)
				}
			}
			if accept(tCURLY_LEFT) {
				for accept(tNAME) || accept(tDEFAULT) {
					extVar := importedVar{}
					extVar.name = atom{getToken()}
					if accept(tAS) {
						expect(tNAME)
						extVar.pseudonym = atom{getToken()}
					} else {
						extVar.pseudonym = atom{getToken()}
					}
					impSt.vars = append(impSt.vars, extVar)
					if src[i+2].tType == tNAME || src[i+2].tType == tDEFAULT {
						expect(tCOMMA)
					} else {
						accept(tCOMMA)
					}
				}

				expect(tCURLY_RIGHT)
				expect(tFROM)
			}
			expect(tSTRING)
			impSt.path = getLexeme()
			result = impSt
			expect(tSEMI)

		case accept(tFUNCTION):
			expect(tNAME)
			name := atom{getToken()}
			funcDecl := getFunctionExpression()
			funcDecl.name = name
			result = funcDecl

		case accept(tRETURN):
			rs := returnStatement{}
			if !accept(tSEMI) {
				rs.value = getSequence()
				expect(tSEMI)
			}
			result = rs

		case accept(tCONTINUE):
			result = expressionStatement{atom{getToken()}}
			expect(tSEMI)

		default:
			result = expressionStatement{getSequence()}
			expect(tSEMI)
		}
		return result
	}

	getProgram := func() program {
		result := program{}
		result.statements = make([]ast, 0)

		for !accept(tEND_OF_INPUT) {
			result.statements = append(result.statements, getStatement())
		}

		return result
	}

	return getProgram()
}
