package main

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

type atom struct {
	val token
}

func (av atom) String() string {
	return av.val.lexeme
}

type binaryExpression struct {
	left     ast
	right    ast
	operator token
}

func (be binaryExpression) String() string {
	result := ""
	op := operatorTable[be.operator.tType]
	switch t := be.left.(type) {
	case binaryExpression:
		leftPrec := operatorTable[t.operator.tType].precedence
		if leftPrec < op.precedence || (leftPrec == op.precedence && op.isRightAssociative) {
			result += "(" + t.String() + ")"
		} else {
			result += t.String()
		}
	default:
		result += t.String()
	}

	result += be.operator.lexeme

	switch t := be.right.(type) {
	case binaryExpression:
		rightPrec := operatorTable[t.operator.tType].precedence
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

type unaryExpression struct {
	value     ast
	operator  string
	isPostfix bool
}

func (ue unaryExpression) String() string {
	return ue.operator + ue.value.String()
}

type sequenceExpression struct {
	items []ast
}

func (sq sequenceExpression) String() string {
	result := ""
	for i, item := range sq.items {
		result += item.String()
		if i < len(sq.items)-1 {
			result += ","
		}
	}
	return result
}

type conditionalExpression struct {
	test       ast
	consequent ast
	alternate  ast
}

func (ce conditionalExpression) String() string {
	result := fmt.Sprintf("%s?%s:%s", ce.test.String(), ce.consequent.String(), ce.alternate.String())
	return result
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

type memberExpression struct {
	object       ast
	property     ast
	isCalculated bool
}

func (me memberExpression) String() string {
	if !me.isCalculated {
		return me.object.String() + "." + me.property.String()
	}
	return me.object.String() + "[" + me.property.String() + "]"
}

type objectProperty struct {
	key          ast
	value        ast
	isCalculated bool
	isNotName    bool
}

type objectLiteral struct {
	properties []objectProperty
}

func (ol objectLiteral) String() string {
	result := "{"
	for i, prop := range ol.properties {
		if prop.isCalculated {
			result += "["
		}
		result += prop.key.String()
		if prop.isCalculated {
			result += "]"
		}
		if prop.value != nil {
			result += ":" + prop.value.String()
		}
		if i < len(ol.properties)-1 {
			result += ","
		}
	}
	result += "}"
	return result
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

type declaration struct {
	name  string
	value ast
}

type varStatement struct {
	decls   []declaration
	keyword string
}

func (vs varStatement) String() string {
	result := vs.keyword + " "
	for i, decl := range vs.decls {
		result += decl.name
		if decl.value != nil {
			result += "=" + decl.value.String()
		}
		if i < len(vs.decls)-1 {
			result += ","
		}
	}
	result += ";"
	return result
}

type importedVar struct {
	name      string
	pseudonym string
}

type importStatement struct {
	path string
	vars []importedVar
}

func (is importStatement) String() string {
	result := "import{"
	for i, v := range is.vars {
		result += v.name
		if v.pseudonym != "" {
			result += " as " + v.pseudonym
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

type exportedVar struct {
	name      string
	pseudonym string
	value     ast
}

type exportStatement struct {
	exportedVars []exportedVar
}

func (es exportStatement) String() string {
	result := "export "
	for _, exp := range es.exportedVars {
		result += exp.name
		if exp.pseudonym != "" {
			result += " as " + exp.pseudonym
		}
		if exp.value != nil {
			result += " = " + exp.value.String()
		}
	}
	return result + ";"
}

type functionExpression struct {
	isDeclaration bool
	name          string
	args          []string
	body          blockStatement
}

func (fe functionExpression) String() string {
	result := "function " + fe.name + "("
	for i, arg := range fe.args {
		result += arg
		if i < len(fe.args)-1 {
			result += ","
		}
	}
	result += ")" + fe.body.String()
	return result
}

type returnStatement struct {
	value ast
}

func (re returnStatement) String() string {
	result := "return"
	if re.value != nil {
		result += re.value.String()
	}
	result += ";"
	return result
}

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

func parse(src []token) {
	i := 0
	t := src[i]

	next := func() {
		i++
		if i < len(src) {
			t = src[i]
		}
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
		panic(fmt.Sprintf("\nExpected %s, got %s\n", tType, t.tType))
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
		getCond,
		getAssign,
		getYield,
		getSpread,
		getSequence grammar
	)

	getFunctionExpression := func() functionExpression {
		funcExpr := functionExpression{}
		funcExpr.args = make([]string, 0)
		expect(tPAREN_LEFT)
		for !accept(tPAREN_RIGHT) {
			if accept(tNAME) {
				funcExpr.args = append(funcExpr.args, getLexeme())
			}
			if !accept(tCOMMA) {
				expect(tPAREN_RIGHT)
				break
			}
		}
		body := getStatement()
		funcExpr.body = body.(blockStatement)
		return funcExpr
	}

	getObjectKey := func() objectProperty {
		prop := objectProperty{}

		if accept(tBRACKET_LEFT) {
			prop.isCalculated = true
			prop.key = getSequence()
			expect(tBRACKET_RIGHT)
		} else if accept(tNAME) {
			prop.key = atom{getToken()}
		} else if isValidPropertyName(t.lexeme) || t.tType == tNUMBER || t.tType == tSTRING {
			prop.isNotName = true
			next()
			prop.key = atom{getToken()}
		}

		return prop
	}

	getObjectLiteral := func() objectLiteral {
		objLit := objectLiteral{}
		objLit.properties = make([]objectProperty, 0)
		for !accept(tCURLY_RIGHT) {
			prop := getObjectKey()
			if prop.isCalculated || prop.isNotName {
				if src[i].tType == tPAREN_LEFT {
					funcExpr := getFunctionExpression()
					prop.value = funcExpr
				} else {
					expect(tCOLON)
					prop.value = getSpread()
				}
			} else {
				if src[i].tType == tPAREN_LEFT {
					funcExpr := getFunctionExpression()
					prop.value = funcExpr
				} else {
					if accept(tCOLON) {
						prop.value = getSpread()
					}
				}
			}

			objLit.properties = append(objLit.properties, prop)

			if !accept(tCOMMA) {
				expect(tCURLY_RIGHT)
				break
			}
		}
		return objLit
	}

	getArrayLiteral := func() arrayLiteral {
		arrayLit := arrayLiteral{}
		arrayLit.items = make([]ast, 0)
		for !accept(tBRACKET_RIGHT) {
			arrayLit.items = append(arrayLit.items, getSpread())
			if !accept(tCOMMA) {
				expect(tBRACKET_RIGHT)
				break
			}
		}
		return arrayLit
	}

	getAtom = func() ast {
		var result ast
		switch {
		case accept(tTHIS, tNUMBER, tSTRING, tNAME, tTRUE, tFALSE):
			result = atom{getToken()}

		case accept(tPAREN_LEFT):
			result = getSequence()
			expect(tPAREN_RIGHT)

		case accept(tCURLY_LEFT):
			result = getObjectLiteral()

		case accept(tFUNCTION):
			var name string
			if accept(tNAME) {
				name = getLexeme()
			} else {
				name = ""
			}
			funcExpr := getFunctionExpression()
			funcExpr.name = name
			result = funcExpr

		case accept(tBRACKET_LEFT):
			result = getArrayLiteral()

		default:
			panic("Ooos! Unexpected token " + t.tType.String())
		}

		return result
	}

	getMember = func() ast {
		var result ast

		object := getAtom()

		if accept(tDOT) || src[i].tType == tBRACKET_LEFT {
			atomStack := []ast{object}
			isCalcStack := []bool{src[i].tType == tBRACKET_LEFT}

			dotCounter := 0

			for ok := true; ok; ok = accept(tDOT) || src[i].tType == tBRACKET_LEFT {
				objProp := getObjectKey()
				atomStack = append(atomStack, objProp.key)
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

	getCond = func() ast {
		var result ast

		test := getBinaryExp()
		if accept(tQUESTION) {
			condExp := conditionalExpression{}
			condExp.test = test
			condExp.consequent = getCond()
			expect(tCOLON)
			condExp.alternate = getCond()
			result = condExp
		} else {
			result = test
		}

		return result
	}

	getAssign = func() ast {
		var result ast

		// TODO dont assign to literals
		left := getCond()
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

	getStatement = func() ast {
		var result ast

		switch {
		case accept(tEXPORT):
			exSt := exportStatement{}
			exSt.exportedVars = make([]exportedVar, 0)
			if t.tType == tVAR || t.tType == tCONST || t.tType == tLET {
				decl := getStatement()
				varSt := decl.(varStatement)
				for _, declarator := range varSt.decls {
					expVar := exportedVar{}
					expVar.name = declarator.name
					expVar.pseudonym = ""
					expVar.value = declarator.value
					exSt.exportedVars = append(exSt.exportedVars, expVar)
				}
			} else if accept(tDEFAULT) {
				expVar := exportedVar{}
				expVar.name = "default"
				if accept(tFUNCTION) {
					var name string
					if accept(tNAME) {
						name = getLexeme()
					} else {
						name = ""
					}
					decl := getFunctionExpression()
					decl.name = name
					expVar.value = decl
				} else if accept(tCLASS) {

				} else {
					expVar.pseudonym = ""
					expVar.value = getSpread()
					exSt.exportedVars = append(exSt.exportedVars, expVar)
				}
				exSt.exportedVars = append(exSt.exportedVars, expVar)
			}
			result = exSt

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
			keyword := getLexeme()
			varSt := varStatement{make([]declaration, 0), keyword}
			for ok := true; ok; ok = accept(tCOMMA) {
				expect(tNAME)
				decl := declaration{}
				decl.name = getLexeme()
				if accept(tASSIGN) {
					decl.value = getSpread()
				}
				varSt.decls = append(varSt.decls, decl)
			}
			result = varSt
			expect(tSEMI)

		// // tIMPORT [tNAME tFROM] [tCURLY_LEFT [tDEFAULT tAS tNAME] [tNAME {tCOMMA tNAME}] [tCOMMA] tCURLY_RIGHT tFROM] tSTRING;
		case accept(tIMPORT):
			impSt := importStatement{}
			impSt.vars = make([]importedVar, 0)
			// import { foo, bar as i } from "module-name";
			if accept(tNAME) {
				defVar := importedVar{}
				defVar.pseudonym = getLexeme()
				defVar.name = "default"

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
					extVar.name = getLexeme()
					if accept(tAS) {
						expect(tNAME)
						extVar.pseudonym = getLexeme()
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
			name := getLexeme()
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

		default:
			result = getSequence()
			expect(tSEMI)
		}
		return result
	}

	getProgram := func() ast {
		result := program{}
		result.statements = make([]ast, 0)

		for !accept(tEND_OF_INPUT) {
			result.statements = append(result.statements, getStatement())
		}

		return result
	}

	result := getProgram()

	fmt.Println(result)
}
