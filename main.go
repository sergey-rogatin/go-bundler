package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type safeFile struct {
	file *os.File
	lock sync.RWMutex
}

func newSafeFile(fileName string) *safeFile {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	return &safeFile{file, sync.RWMutex{}}
}

func (sf *safeFile) close() {
	sf.file.Close()
}

func containsStr(arr []string, c string) bool {
	for _, b := range arr {
		if b == c {
			return true
		}
	}
	return false
}

func trimQuotesFromString(jsString string) string {
	return jsString[1 : len(jsString)-1]
}

func getExtension(resolvedImportPath string) string {
	parts := strings.Split(resolvedImportPath, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "js"
}

func createVarNameFromPath(resolvedImportPath string) string {
	newName := strings.Replace(resolvedImportPath, "/", "_", -1)
	newName = strings.Replace(newName, ".", "_", -1)
	return newName
}

func isKeyword(t token) bool {
	_, ok := keywords[t.lexeme]
	return ok && t.tType != tNAME
}

func copyFile(dst, src string) {
	from, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer to.Close()
	io.Copy(to, from)
}

type cachedFile struct {
	data        []byte
	lastModTime time.Time
	imports     []string
}

type bundleCache struct {
	files map[string]cachedFile
	lock  sync.RWMutex
}

func (c *bundleCache) read(fileName string) (cachedFile, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	file, ok := c.files[fileName]
	return file, ok
}

func (c *bundleCache) write(fileName string, data cachedFile) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.files[fileName] = data
}

func addFileToBundle(
	resolvedPath string,
	bundleSf *safeFile,
	finishedImportsCh chan string,
	cache *bundleCache,
) {
	ext := getExtension(resolvedPath)
	fileStats, _ := os.Stat(resolvedPath)

	var data []byte
	var fileImports []string

	switch ext {
	case "js":
		cachedData, ok := cache.read(resolvedPath)
		if ok && cachedData.lastModTime == fileStats.ModTime() {
			data = cachedData.data
			fileImports = cachedData.imports
		} else {
			src, err := ioutil.ReadFile(resolvedPath)
			if err != nil {
				panic(err)
			}

			data, fileImports = loadJsFile(src, resolvedPath)
		}

	default:
		bundleDir := filepath.Dir(bundleSf.file.Name())
		dstFileName := bundleDir + "/" + createVarNameFromPath(resolvedPath) + "." + ext
		copyFile(dstFileName, resolvedPath)
	}

	cache.write(resolvedPath, cachedFile{
		data:        data,
		lastModTime: fileStats.ModTime(),
		imports:     fileImports,
	})

	addFilesToBundle(fileImports, bundleSf, cache)
	writeToSafeFile(data, bundleSf)
	finishedImportsCh <- resolvedPath
}

func addFilesToBundle(
	files []string,
	bundleSf *safeFile,
	cache *bundleCache,
) {
	filesImportedCh := make(chan string, len(files))

	for _, unbundledFile := range files {
		go addFileToBundle(unbundledFile, bundleSf, filesImportedCh, cache)
	}

	for counter := 0; counter < len(files); counter++ {
		<-filesImportedCh
	}
}

func writeToSafeFile(data []byte, sf *safeFile) {
	sf.lock.Lock()
	defer sf.lock.Unlock()

	sf.file.Write(data)
}

func createBundle(entryFileName, bundleFileName string, cache *bundleCache) {
	buildStartTime := time.Now()

	os.MkdirAll(filepath.Dir(bundleFileName), 0666)
	os.Remove(bundleFileName)
	sf := newSafeFile(bundleFileName)
	defer sf.close()

	addFilesToBundle([]string{entryFileName}, sf, cache)

	fmt.Printf("Build finished in %s\n", time.Since(buildStartTime))
}

type configJSON struct {
	Entry        string
	TemplateHTML string
	BundleDir    string
	WatchFiles   bool
	DevServer    struct {
		Enable bool
		Port   int
	}
}

// type astNode struct {
// 	t        token
// 	children []astNode
// }

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
	args          []string
	isConstructor bool
}

func (fc functionCall) String() string {
	result := ""
	if fc.isConstructor {
		result += "new "
	}
	result += fc.name.String() + "("
	for i, arg := range fc.args {
		result += arg
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
		result += ":" + prop.value.String()
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
	exportedName string
	pseudonym    string
}

type importStatement struct {
	path string
	vars []importedVar
}

func (is importStatement) String() string {
	result := "import{"
	for i, v := range is.vars {
		result += v.exportedName
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

type program struct {
	statements []ast
}

func (ps program) String() string {
	result := ""
	for _, s := range ps.statements {
		result += s.String() + ";"
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

	getAtom = func() ast {
		var result ast
		if accept(tNUMBER) || accept(tSTRING) || accept(tNAME) || accept(tTRUE) || accept(tFALSE) {
			result = atomValue{lexeme()}
		} else if accept(tPAREN_LEFT) {
			result = getExpression()
			expect(tPAREN_RIGHT)
		} else if accept(tCURLY_LEFT) {
			objLit := objectLiteral{}
			objLit.properties = make([]objectProperty, 0)
			for !accept(tCURLY_RIGHT) {
				prop := objectProperty{}
				if accept(tBRACKET_LEFT) {
					prop.key = getExpression()
					prop.isCalculated = true
					expect(tBRACKET_RIGHT)
				} else if accept(tSTRING) || accept(tNAME) {
					prop.key = atomValue{lexeme()}
				}
				expect(tCOLON)
				prop.value = getSpread()
				objLit.properties = append(objLit.properties, prop)
				if !accept(tCOMMA) {
					break
				}
			}
			expect(tCURLY_RIGHT)
			result = objLit
		} else {
			panic("whoops")
		}

		return result
	}

	getNewOrMember = func() ast {
		var result ast

		if accept(tNEW) {
			fc := functionCall{}
			fc.name = getNewOrMember()
			fc.isConstructor = true
			fc.args = make([]string, 0)
			expect(tPAREN_LEFT)
			for ok := true; ok; ok = accept(tCOMMA) {
				if accept(tNAME) {
					fc.args = append(fc.args, lexeme())
				} else {
					break
				}
			}
			expect(tPAREN_RIGHT)
			result = fc
		} else {
			object := getAtom()
			if accept(tDOT) {
				me := memberExpression{}
				me.object = object
				expect(tNAME)
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
		if accept(tPAREN_LEFT) {
			fc := functionCall{}
			fc.name = funcName
			fc.args = make([]string, 0)
			for ok := true; ok; ok = accept(tCOMMA) {
				if accept(tNAME) {
					fc.args = append(fc.args, lexeme())
				} else {
					break
				}
			}
			expect(tPAREN_RIGHT)
			result = fc
		} else {
			result = funcName
		}
		// TODO new without params somehow?

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
		case accept(tVAR):
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
				defVar.exportedName = "default"

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
					extVar.exportedName = lexeme()
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

func main() {
	src := []byte("foo = { bar: 3213, baz: 123 };")

	tokens := lex(src)
	//fmt.Println(tokens)
	parse(tokens)

	// config := configJSON{}
	// config.TemplateHTML = "test/template.html"
	// config.WatchFiles = false
	// config.DevServer.Enable = false

	// if len(os.Args) > 1 {
	// 	configFileName := os.Args[1]

	// 	configFile, err := ioutil.ReadFile(configFileName)
	// 	if err != nil {
	// 		fmt.Println("Unable to load config file!")
	// 		config = configJSON{}
	// 	}
	// 	json.Unmarshal(configFile, &config)
	// }

	// // config defaults
	// if config.Entry == "" {
	// 	config.Entry = "test/index.js"
	// }
	// if config.BundleDir == "" {
	// 	config.BundleDir = "test/build"
	// }
	// if config.DevServer.Port == 0 {
	// 	config.DevServer.Port = 8080
	// }

	// entryName := config.Entry
	// bundleName := filepath.Join(config.BundleDir, "bundle.js")

	// cache := bundleCache{}
	// cache.files = make(map[string]cachedFile)
	// createBundle(entryName, bundleName, &cache)

	// if config.TemplateHTML != "" {
	// 	bundleHTMLTemplate(config.TemplateHTML, bundleName)
	// }

	// // dev server and watching files
	// if config.DevServer.Enable {
	// 	if config.WatchFiles {
	// 		go watchBundledFiles(&cache, entryName, bundleName)
	// 	}
	// 	fmt.Printf("Dev server listening at port %v\n", config.DevServer.Port)
	// 	server := http.FileServer(http.Dir(config.BundleDir))
	// 	err := http.ListenAndServe(fmt.Sprintf(":%v", config.DevServer.Port), server)
	// 	log.Fatal(err)
	// } else if config.WatchFiles {
	// 	watchBundledFiles(&cache, entryName, bundleName)
	// }
}
