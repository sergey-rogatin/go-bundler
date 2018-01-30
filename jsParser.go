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

type atomValue struct {
	val string
}

func (av atomValue) String() string {
	return av.val
}

type binaryExpression struct {
	left     ast
	right    ast
	operator string
}

func (be binaryExpression) String() string {
	return "(" + be.left.String() + be.operator + be.right.String() + ")"
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
	return result
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
		result += st.String() + ";"
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
	return "return " + re.value.String()
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

func parse(src []token) {
	i := 0
	t := src[i]

	next := func() {
		i++
		if i < len(src) {
			t = src[i]
		}
	}

	accept := func(tType tokenType) bool {
		if t.tType == tType {
			next()
			return true
		}
		return false
	}

	expect := func(tType tokenType) bool {
		if accept(tType) {
			return true
		}
		panic(fmt.Sprintf("\nExpected %s, got %s\n", tType, t.tType))
	}

	lexeme := func() string {
		return src[i-1].lexeme
	}

	type grammar func() ast

	var (
		getStatement,
		getAtom,
		getNewOrMember,
		getFunctionCall,
		getPostfixUnary,
		getPrefixUnary,
		getExponent,
		getMultDiv,
		getAddSub,
		getShift,
		getCompare,
		getEquals,
		getBitAnd,
		getBitXor,
		getBitOr,
		getAnd,
		getOr,
		getCond,
		getAssign,
		getYield,
		getSpread,
		getExpression grammar
	)

	makeUnaryExp := func(ops []tokenType, leftType grammar, isPostfix bool) grammar {
		var expFunc grammar

		expFunc = func() ast {
			var result ast

			var value ast
			if isPostfix {
				value = leftType()
			}
			for _, op := range ops {
				if accept(op) {
					unExpr := unaryExpression{}
					unExpr.operator = lexeme()
					if isPostfix {
						unExpr.value = value
					} else {
						unExpr.value = leftType()
					}
					unExpr.isPostfix = isPostfix
					result = unExpr
				}
			}
			if result == nil && isPostfix {
				result = value
			} else if result == nil {
				result = leftType()
			}

			return result
		}

		return expFunc
	}

	makeBinExp := func(ops []tokenType, leftType grammar) grammar {
		var expFunc grammar

		expFunc = func() ast {
			var result ast

			left := leftType()
			for _, op := range ops {
				if accept(op) {
					be := binaryExpression{}
					be.left = left
					be.operator = lexeme()
					be.right = expFunc()
					result = be
				}
			}
			if result == nil {
				result = left
			}

			return result
		}

		return expFunc
	}

	getFunctionExpression := func(isDeclaration bool) functionExpression {
		funcExpr := functionExpression{}
		funcExpr.isDeclaration = isDeclaration
		funcExpr.args = make([]string, 0)
		if isDeclaration {
			expect(tNAME)
			funcExpr.name = lexeme()
		} else if accept(tNAME) {
			funcExpr.name = lexeme()
		} else {
			funcExpr.name = "anonymous"
		}
		expect(tPAREN_LEFT)
		for !accept(tPAREN_RIGHT) {
			if accept(tNAME) {
				funcExpr.args = append(funcExpr.args, lexeme())
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

	getObjectLiteral := func() objectLiteral {
		objLit := objectLiteral{}
		objLit.properties = make([]objectProperty, 0)
		for !accept(tCURLY_RIGHT) {
			prop := objectProperty{}
			if accept(tBRACKET_LEFT) {
				prop.isCalculated = true
				prop.key = getExpression()
				expect(tBRACKET_RIGHT)
				if src[i+1].tType != tCOLON {
					panic("Expected tCOLON, got " + src[i+1].tType.String())
				}
			} else if accept(tNAME) {
				prop.key = atomValue{lexeme()}
			} else if isValidPropertyName(t.lexeme) || t.tType == tNUMBER || t.tType == tSTRING {
				if src[i+1].tType != tCOLON {
					panic("Expected tCOLON, got " + src[i+1].tType.String())
				}
				next()
				prop.key = atomValue{lexeme()}
			}
			if accept(tCOLON) {
				prop.value = getSpread()
			}
			objLit.properties = append(objLit.properties, prop)
			if !accept(tCOMMA) {
				expect(tCURLY_RIGHT)
				break
			}
		}
		return objLit
	}

	getAtom = func() ast {
		var result ast
		switch {
		case accept(tNUMBER) || accept(tSTRING) || accept(tNAME) || accept(tTRUE) || accept(tFALSE):
			result = atomValue{lexeme()}

		case accept(tPAREN_LEFT):
			result = getExpression()
			expect(tPAREN_RIGHT)

		case accept(tCURLY_LEFT):
			result = getObjectLiteral()

		case accept(tFUNCTION):
			result = getFunctionExpression(false)

		default:
			panic("whoops")
		}

		return result
	}

	getNewOrMember = func() ast {
		var result ast

		switch {
		case accept(tNEW):
			fc := functionCall{}
			fc.name = getNewOrMember()
			fc.isConstructor = true
			fc.args = make([]ast, 0)
			if accept(tPAREN_LEFT) {
				for ok := true; ok; ok = accept(tCOMMA) {
					if accept(tNAME) {
						fc.args = append(fc.args, getSpread())
					} else {
						break
					}
				}
				expect(tPAREN_RIGHT)
			}
			result = fc

		default:
			object := getAtom()
			if accept(tDOT) {
				me := memberExpression{}
				me.object = object
				if !isValidPropertyName(t.lexeme) {
					panic(t.lexeme + " is not a valid property name")
				}
				next()
				me.property = atomValue{lexeme()}
				result = me
			} else if accept(tBRACKET_LEFT) {
				me := memberExpression{}
				me.object = object
				me.property = getExpression()
				me.isCalculated = true
				expect(tBRACKET_RIGHT)
				result = me
			} else {
				result = object
			}
		}

		return result
	}

	getFunctionCall = func() ast {
		var result ast

		funcName := getNewOrMember()
		switch {
		case accept(tPAREN_LEFT):
			fc := functionCall{}
			fc.name = funcName
			fc.args = make([]ast, 0)
			for !accept(tPAREN_RIGHT) {
				arg := getSpread()
				fc.args = append(fc.args, arg)
				if arg != nil {
					if accept(tCOMMA) {
						if t.tType == tCOMMA {
							panic("Expected tPAREN_RIGHT, got tCOMMA")
						}
					}
				}
			}
			result = fc

		default:
			result = funcName
		}

		return result
	}

	getPostfixUnary = makeUnaryExp([]tokenType{tINC, tDEC}, getFunctionCall, true)
	getPrefixUnary = makeUnaryExp([]tokenType{tNOT, tBITWISE_NOT, tPLUS, tMINUS, tINC, tDEC, tTYPEOF, tVOID, tDELETE}, getPostfixUnary, false)

	getExponent = makeBinExp([]tokenType{tEXP}, getPrefixUnary)
	getMultDiv = makeBinExp([]tokenType{tMULT, tDIV, tMOD}, getExponent)
	getAddSub = makeBinExp([]tokenType{tPLUS, tMINUS}, getMultDiv)
	getShift = makeBinExp([]tokenType{tBITWISE_SHIFT_LEFT, tBITWISE_SHIFT_RIGHT, tBITWISE_SHIFT_RIGHT_ZERO}, getAddSub)
	getCompare = makeBinExp([]tokenType{tLESS, tLESS_EQUALS, tGREATER, tGREATER_EQUALS, tIN, tINSTANCEOF}, getShift)
	getEquals = makeBinExp([]tokenType{tEQUALS, tNOT_EQUALS, tEQUALS_STRICT, tNOT_EQUALS_STRICT}, getCompare)
	getBitAnd = makeBinExp([]tokenType{tBITWISE_AND}, getEquals)
	getBitXor = makeBinExp([]tokenType{tBITWISE_XOR}, getBitAnd)
	getBitOr = makeBinExp([]tokenType{tBITWISE_OR}, getBitXor)
	getAnd = makeBinExp([]tokenType{tAND}, getBitOr)
	getOr = makeBinExp([]tokenType{tOR}, getAnd)

	getCond = func() ast {
		var result ast

		test := getOr()
		if accept(tQUESTION) {
			condExp := conditionalExpression{}
			condExp.test = test
			condExp.consequent = getOr()
			expect(tCOLON)
			condExp.alternate = getOr()
			result = condExp
		} else {
			result = test
		}

		return result
	}

	getAssign = makeBinExp([]tokenType{
		tASSIGN, tPLUS_ASSIGN, tMINUS_ASSIGN, tMULT_ASSIGN, tDIV_ASSIGN,
		tBITWISE_SHIFT_LEFT_ASSIGN, tBITWISE_SHIFT_RIGHT_ASSIGN, tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN,
		tBITWISE_AND_ASSIGN, tBITWISE_XOR_ASSIGN, tBITWISE_OR_ASSIGN,
	}, getCond)
	getYield = makeBinExp([]tokenType{tYIELD, tYIELD_STAR}, getAssign)
	getSpread = makeBinExp([]tokenType{tSPREAD}, getYield)

	getExpression = func() ast {
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
					decl := getFunctionExpression(false)
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
			wlSt.test = getExpression()
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
			ifSt.test = getExpression()
			expect(tPAREN_RIGHT)
			ifSt.consequent = getStatement()
			if accept(tELSE) {
				ifSt.alternate = getStatement()
			}
			result = ifSt

		// tVAR tNAME [tEQUALS add] {tCOMMA tNAME [tEQUALS add]} tSEMI
		case accept(tVAR) || accept(tCONST) || accept(tLET):
			keyword := lexeme()
			varSt := varStatement{make([]declaration, 0), keyword}
			for ok := true; ok; ok = accept(tCOMMA) {
				expect(tNAME)
				decl := declaration{}
				decl.name = lexeme()
				if accept(tASSIGN) {
					decl.value = getOr()
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
				defVar.pseudonym = lexeme()
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
					extVar.name = lexeme()
					if accept(tAS) {
						expect(tNAME)
						extVar.pseudonym = lexeme()
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
			impSt.path = lexeme()
			result = impSt
			expect(tSEMI)

		case accept(tFUNCTION):
			result = getFunctionExpression(true)

		case accept(tRETURN):
			rs := returnStatement{}
			rs.value = getExpression()
			result = rs
			expect(tSEMI)

		default:
			result = getExpression()
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
