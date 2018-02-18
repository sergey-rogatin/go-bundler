package jsLoader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/lvl5hm/go-bundler/util"
)

type context struct {
	fileName      string
	fileImports   []string
	hasES6Exports bool
	env           map[string]string
}

func LoadFile(fileName string, env map[string]string) ([]byte, []string, error) {
	src, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}

	tokens := lex(src)
	initialProgram, parseErr := parseTokens(tokens)
	if parseErr != nil {
		return nil, nil, parseErr
	}

	resultProgram, fileImports := transformIntoModule(initialProgram, fileName, env)
	resultBytes := []byte(printAst(resultProgram))
	return resultBytes, fileImports, nil
}

func ParseAndPrint(src []byte) ([]byte, error) {
	tokens := lex(src)
	initialProgram, parseErr := parseTokens(tokens)
	if parseErr != nil {
		return nil, parseErr
	}
	resultBytes := []byte(printAst(initialProgram))

	return resultBytes, parseErr
}

func GetJsBundleFileHead() []byte {
	start := `function requireES6(module, impName) {
		if (module.hasES6Exports) {
			if (impName === '*') {
				return module.es6;
			}
			return module.es6[impName];
		}
		if (impName == '*' || impName === 'default') {
			return module.exports;
		}
		return module.exports[impName];
	}
	
	function require(module) {
		if (module.hasES6Exports) {
			return module.es6;
		}
		return module.exports;
	}
	var moduleFns={},modules={};
	var process={env:{}};`

	res, _ := ParseAndPrint([]byte(start))
	// res := []byte(start)
	return res
}

type cyclicDependencyWarning struct {
	path []string
}

func (c cyclicDependencyWarning) Error() string {
	return fmt.Sprintf(
		"Warning: circular dependency detected:\n%s\n",
		strings.Join(c.path, " -> "),
	)
}

func GetJsBundleFileTail(
	entryFileName string,
	importMap map[string][]string,
) ([]byte, []cyclicDependencyWarning) {
	moduleOrder := []string{}
	warnings := []cyclicDependencyWarning{}

	var createImportTree func(string, []string)
	createImportTree = func(fileName string, path []string) {
		if i := util.IndexOf(path, fileName); i >= 0 {
			warnings = append(warnings, cyclicDependencyWarning{append(path[i:], fileName)})
			return
		}

		fileImports := importMap[fileName]
		for _, importPath := range fileImports {
			createImportTree(importPath, append(path, fileName))
		}

		moduleName := "'" + CreateVarNameFromPath(fileName) + "'"
		if util.IndexOf(moduleOrder, moduleName) < 0 {
			moduleOrder = append(moduleOrder, moduleName)
		}
	}

	createImportTree(entryFileName, []string{})
	jsModuleOrder := fmt.Sprintf("var moduleOrder = [%s];", strings.Join(moduleOrder, ","))
	tail := []byte(jsModuleOrder + "moduleOrder.forEach((moduleName)=>modules[moduleName]=moduleFns[moduleName]())")

	result, _ := ParseAndPrint(tail)
	// result := tail
	return result, warnings
}

func modifyAst(n ast, ctx *context) ast {
	switch n.t {

	case n_BINARY_EXPRESSION:
		return modifyBinaryExpression(n, ctx)

	case n_MEMBER_EXPRESSION:
		return modifyMemberExpression(n, ctx)

	case n_IF_STATEMENT:
		return modifyIfStatement(n, ctx)

	case n_FUNCTION_CALL:
		return modifyFunctionCall(n, ctx)

	case n_EXPORT_STATEMENT:
		return modifyExport(n, ctx)

	case n_PROGRAM:
		return modifyProgram(n, ctx)

	case n_IMPORT_STATEMENT:
		return modifyImport(n, ctx)

	default:
		res := n
		res.children = []ast{}
		for _, c := range n.children {
			res.children = append(res.children, modifyAst(c, ctx))
		}
		return res
	}
}

func modifyBinaryExpression(n ast, ctx *context) ast {
	operator := n.value
	left := modifyAst(n.children[0], ctx)
	right := modifyAst(n.children[1], ctx)

	if operator == "===" || operator == "==" {
		if left.t == n_STRING_LITERAL && left.t == right.t {
			if printAst(left) == printAst(right) {
				return makeNode(n_BOOL_LITERAL, "true")
			}
			return makeNode(n_BOOL_LITERAL, "false")
		}
	}

	return makeNode(n_BINARY_EXPRESSION, operator, left, right)
}

func modifyMemberExpression(n ast, ctx *context) ast {
	object := n.children[0]
	prop := n.children[1]

	if object.t == n_MEMBER_EXPRESSION {
		innerObject := object.children[0]
		innerProp := object.children[1]

		if innerObject.t == n_IDENTIFIER &&
			innerObject.value == "process" &&
			innerProp.t == n_IDENTIFIER &&
			innerProp.value == "env" {

			if envVal, ok := ctx.env[prop.value]; ok {
				return makeNode(n_STRING_LITERAL, envVal)
			}
		}
	}

	children := []ast{}
	for _, c := range n.children {
		children = append(children, modifyAst(c, ctx))
	}
	n.children = children
	return n
}

func modifyIfStatement(n ast, ctx *context) ast {
	cond := modifyAst(n.children[0], ctx)
	conseq := n.children[1]
	altern := n.children[2]

	if cond.t == n_BOOL_LITERAL {
		if cond.value == "true" {
			return modifyAst(conseq, ctx)
		}
		return modifyAst(altern, ctx)
	}

	res := n
	res.children = []ast{}
	for _, c := range n.children {
		res.children = append(res.children, modifyAst(c, ctx))
	}
	return res
}

func modifyProgram(n ast, ctx *context) ast {
	children := []ast{}
	for _, c := range n.children {
		children = append(children, modifyAst(c, ctx))
	}
	n.children = children

	moduleName := CreateVarNameFromPath(ctx.fileName)

	hasES6ExportsStr := "false"
	if ctx.hasES6Exports {
		hasES6ExportsStr = "true"
	}

	beforeSt := []ast{
		makeNode(
			n_DECLARATION_STATEMENT, "",
			makeNode(
				n_VARIABLE_DECLARATION, "var",
				makeNode(
					n_DECLARATOR, "",
					makeNode(n_IDENTIFIER, "module"),
					makeNode(
						n_OBJECT_LITERAL, "",
						makeNode(
							n_OBJECT_PROPERTY, "",
							makeNode(n_OBJECT_KEY, "exports"),
							makeNode(n_OBJECT_LITERAL, ""),
						),
						makeNode(
							n_OBJECT_PROPERTY, "",
							makeNode(n_OBJECT_KEY, "es6"),
							makeNode(n_OBJECT_LITERAL, ""),
						),
						makeNode(
							n_OBJECT_PROPERTY, "",
							makeNode(n_OBJECT_KEY, "hasES6Exports"),
							makeNode(n_BOOL_LITERAL, hasES6ExportsStr),
						),
					),
				),
				makeNode(
					n_DECLARATOR, "",
					makeNode(n_IDENTIFIER, "exports"),
					makeNode(
						n_MEMBER_EXPRESSION, "",
						makeNode(n_IDENTIFIER, "module"),
						makeNode(n_IDENTIFIER, "exports"),
					),
				),
			),
		),
	}

	afterSt := []ast{
		makeNode(
			n_RETURN_STATEMENT, "",
			makeNode(n_IDENTIFIER, "module"),
		),
	}

	allSt := append(
		beforeSt,
		append(n.children, afterSt...)...,
	)

	funcExpr := makeNode(
		n_FUNCTION_EXPRESSION, "",
		makeNode(n_EMPTY, ""),
		makeNode(n_FUNCTION_PARAMETERS, ""),
		makeNode(n_BLOCK_STATEMENT, "", allSt...),
	)

	decl := makeNode(
		n_EXPRESSION_STATEMENT, "",
		makeNode(
			n_ASSIGNMENT_EXPRESSION, "=",
			makeNode(
				n_MEMBER_EXPRESSION, "",
				makeNode(n_IDENTIFIER, "moduleFns"),
				makeNode(n_IDENTIFIER, moduleName),
			),
			funcExpr,
		),
	)

	return decl
}

func getImport(resolvedPath, eType, name string) ast {
	moduleName := CreateVarNameFromPath(resolvedPath)

	funcName := eType

	args := []ast{
		makeNode(
			n_MEMBER_EXPRESSION, "",
			makeNode(n_IDENTIFIER, "modules"),
			makeNode(n_IDENTIFIER, moduleName),
		),
	}

	if name != "" {
		args = append(args, makeNode(n_STRING_LITERAL, name))
	}

	return makeNode(
		n_FUNCTION_CALL, "",
		makeNode(n_IDENTIFIER, funcName),
		makeNode(n_FUNCTION_ARGS, "", args...),
	)
}

func modifyFunctionCall(n ast, ctx *context) ast {
	children := []ast{}
	for _, c := range n.children {
		children = append(children, modifyAst(c, ctx))
	}
	n.children = children

	name := n.children[0]

	// modify require
	if name.value == "require" {
		args := n.children[1].children

		path := args[0].value
		resolvedPath := resolveES6ImportPath(path, ctx.fileName)
		ctx.fileImports = append(ctx.fileImports, resolvedPath)
		moduleName := CreateVarNameFromPath(resolvedPath)

		args[0] = makeNode(
			n_MEMBER_EXPRESSION, "",
			makeNode(n_IDENTIFIER, "modules"),
			makeNode(n_IDENTIFIER, moduleName),
		)
		return n
	}

	return n
}

func modifyImport(n ast, ctx *context) ast {
	children := []ast{}
	for _, c := range n.children {
		children = append(children, modifyAst(c, ctx))
	}
	n.children = children

	vars := n.children[0].children
	all := n.children[1].value
	path := n.children[2].value
	resolvedPath := resolveES6ImportPath(path, ctx.fileName)
	ctx.fileImports = append(ctx.fileImports, resolvedPath)

	decls := []ast{}

	if all != "" {
		alias := makeNode(n_IDENTIFIER, all)
		allDecl := makeNode(
			n_DECLARATOR, "", alias, getImport(
				resolvedPath, "requireES6", "*",
			),
		)
		decls = append(decls, allDecl)
	}

	for _, v := range vars {
		name := v.children[0]
		alias := v.children[1]
		d := makeNode(
			n_DECLARATOR, "", alias, getImport(
				resolvedPath, "requireES6", name.value,
			),
		)
		decls = append(decls, d)
	}

	if len(decls) > 0 {
		return makeNode(
			n_DECLARATION_STATEMENT, "",
			makeNode(n_VARIABLE_DECLARATION, "var", decls...),
		)
	}
	return makeNode(n_EMPTY, "")
}

func modifyExport(n ast, ctx *context) ast {
	ctx.hasES6Exports = true

	children := []ast{}
	for _, c := range n.children {
		children = append(children, modifyAst(c, ctx))
	}
	n.children = children

	vars := n.children[0].children
	declaration := n.children[1]
	path := n.children[2].value

	es6Exports := makeNode(
		n_MEMBER_EXPRESSION, "",
		makeNode(n_IDENTIFIER, "module"),
		makeNode(n_IDENTIFIER, "es6"),
	)

	if path != "" {
		resolvedPath := resolveES6ImportPath(path, ctx.fileName)
		ctx.fileImports = append(ctx.fileImports, resolvedPath)

		if n.children[0].t == n_EXPORT_ALL {
			return makeNode(
				n_EXPRESSION_STATEMENT, "",
				makeNode(
					n_FUNCTION_CALL, "",
					makeNode(
						n_MEMBER_EXPRESSION, "",
						makeNode(n_IDENTIFIER, "Object"),
						makeNode(n_IDENTIFIER, "assign"),
					),
					makeNode(
						n_FUNCTION_ARGS, "",
						es6Exports,
						getImport(resolvedPath, "requireES6", "*"),
					),
				),
			)
		}
	}

	assignments := []ast{}
	for _, v := range vars {
		name := v.children[0]
		alias := v.children[1].value

		left := makeNode(
			n_MEMBER_EXPRESSION, "",
			es6Exports,
			makeNode(n_IDENTIFIER, alias),
		)

		var right ast
		if path != "" {
			resolvedPath := resolveES6ImportPath(path, ctx.fileName)
			right = getImport(resolvedPath, "requireES6", name.value)
		} else {
			right = name
		}

		ass := makeNode(n_ASSIGNMENT_EXPRESSION, "=", left, right)
		assignments = append(assignments, ass)
	}

	seq := makeNode(n_SEQUENCE_EXPRESSION, "", assignments...)
	exprSt := makeNode(n_EXPRESSION_STATEMENT, "", seq)

	multi := makeNode(n_MULTISTATEMENT, "", declaration, exprSt)

	return multi
}

func transformIntoModule(src ast, fileName string, env map[string]string) (ast, []string) {
	ctx := context{}
	ctx.fileImports = []string{}
	ctx.fileName = fileName
	ctx.env = env

	res := modifyAst(src, &ctx)
	// res := src

	return res, ctx.fileImports
}

func CreateVarNameFromPath(path string) string {
	newName := strings.Replace(path, "/", "_", -1)
	newName = strings.Replace(newName, ".", "_", -1)
	newName = strings.Replace(newName, "-", "_", -1)
	return newName
}

func makeToken(text string) token {
	res := lex([]byte(text))
	return res[0]
}

type packageJSON struct {
	Main string
}

func resolveES6ImportPath(importPath, currentFileName string) string {
	pathParts := strings.Split(importPath, "/")

	locationParts := strings.Split(currentFileName, "/")
	locationParts = locationParts[:len(locationParts)-1]

	// import from node_modules
	isNode := false
	if len(pathParts) > 0 {
		if pathParts[0] != "." && pathParts[0] != ".." {
			locationParts = []string{"node_modules"}
			isNode = true
		}
	}

	for _, part := range pathParts {
		if part == ".." {
			locationParts = locationParts[:len(locationParts)-1]
			pathParts = pathParts[1:]
		}
		if part == "." {
			pathParts = pathParts[1:]
		}
	}

	fullFileName := strings.Join(append(locationParts, pathParts...), "/")
	if isNode && len(pathParts) == 1 {
		packageFile, err := ioutil.ReadFile(fullFileName + "/package.json")
		if err != nil {
			// fmt.Println("Cannot find package.json!")
			fullFileName += "/" + "index.js"
		} else {
			packJSON := packageJSON{}
			json.Unmarshal(packageFile, &packJSON)
			if packJSON.Main != "" {
				fullFileName += "/" + packJSON.Main
			} else {
				fullFileName += "/" + "index.js"
			}
		}
	}

	ext := filepath.Ext(fullFileName)
	if ext == "" {
		ext = ".js"
	} else {
		ext = ""
	}

	result := fullFileName + ext
	return result
}
