package jsLoader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/lvl5hm/go-bundler/loaders"

	"github.com/lvl5hm/go-bundler/util"
)

type environment struct {
	fileName      string
	fileImports   []string
	hasES6Exports bool
	env           map[string]string
}

type context struct {
	declaredVars map[string]ast
	parent       *context

	isInsideDeclaration bool
}

func getReplaceImport(ctx *context, varName string) (ast, bool) {
	if replace, ok := ctx.declaredVars[varName]; ok {
		if replace.t != n_EMPTY {
			return replace, true
		}
		return replace, false
	}
	if ctx.parent != nil {
		return getReplaceImport(ctx.parent, varName)
	}
	return ast{}, false
}

func setDeclaredVar(ctx *context, varName string, value ast) {
	ctx.declaredVars[varName] = value
}

var Loader jsLoader

type jsLoader struct{}

func (j jsLoader) BeforeBuild(fileName string, config *loaders.ConfigJSON) {}

func (j jsLoader) LoadAndTransformFile(
	fileName string,
	config *loaders.ConfigJSON,
) ([]byte, []string, error) {
	src, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}
	return j.TransformFile(fileName, src, config)
}

func (j jsLoader) TransformFile(
	fileName string,
	src []byte,
	config *loaders.ConfigJSON,
) ([]byte, []string, error) {
	tokens := lex(src)
	initialProgram, parseErr := parseTokens(tokens)
	if parseErr != nil {
		return nil, nil, parseErr
	}

	resultProgram, fileImports := transformIntoModule(initialProgram, fileName, config.Env)
	resultBytes := []byte(printAst(resultProgram))
	return resultBytes, fileImports, nil
}

func LoadFile(fileName string, config *loaders.ConfigJSON) ([]byte, []string, error) {
	src, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}

	tokens := lex(src)
	initialProgram, parseErr := parseTokens(tokens)
	if parseErr != nil {
		return nil, nil, parseErr
	}

	resultProgram, fileImports := transformIntoModule(initialProgram, fileName, config.Env)
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
    if (!module) {
      return undefined;
    }
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

		moduleName := "'" + loaders.CreateVarNameFromPath(fileName) + "'"
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

func modifyChildren(n ast, e *environment, ctx *context) ast {
	res := n
	res.children = []ast{}
	for _, c := range n.children {
		res.children = append(res.children, modifyAst(c, e, ctx))
	}
	return res
}

func modifyAst(n ast, e *environment, ctx *context) ast {
	switch n.t {

	case n_DECLARATOR:
		ctx.isInsideDeclaration = true
		n.children[0] = modifyAst(n.children[0], e, ctx)
		ctx.isInsideDeclaration = false

		n.children[1] = modifyAst(n.children[1], e, ctx)
		return n

	case n_FUNCTION_PARAMETERS:
		ctx.isInsideDeclaration = true
		res := modifyChildren(n, e, ctx)
		ctx.isInsideDeclaration = false
		return res

	case n_FUNCTION_EXPRESSION,
		n_FUNCTION_DECLARATION,
		n_LAMBDA_EXPRESSION,
		n_BLOCK_STATEMENT,
		n_FOR_IN_STATEMENT,
		n_FOR_OF_STATEMENT,
		n_FOR_STATEMENT:

		childCtx := context{}
		childCtx.parent = ctx
		childCtx.isInsideDeclaration = ctx.isInsideDeclaration
		childCtx.declaredVars = map[string]ast{}

		return modifyChildren(n, e, &childCtx)

	case n_CALCULATED_OBJECT_KEY,
		n_OBJECT_KEY,
		n_NON_IDENTIFIER_OBJECT_KEY:

		if ctx.isInsideDeclaration {
			setDeclaredVar(ctx, n.value, ast{})
		}
		return n

	case n_BINARY_EXPRESSION:
		return modifyBinaryExpression(n, e, ctx)

	case n_MEMBER_EXPRESSION:
		return modifyMemberExpression(n, e, ctx)

	case n_IF_STATEMENT:
		return modifyIfStatement(n, e, ctx)

	case n_FUNCTION_CALL:
		return modifyFunctionCall(n, e, ctx)

	case n_EXPORT_STATEMENT:
		return modifyExport(n, e, ctx)

	case n_PROGRAM:
		return modifyProgram(n, e, ctx)

	case n_IMPORT_STATEMENT:
		return modifyImport(n, e, ctx)

	case n_IDENTIFIER:
		return modifyIdentifier(n, e, ctx)

	default:
		return modifyChildren(n, e, ctx)
	}
}

func modifyIdentifier(n ast, e *environment, ctx *context) ast {
	if ctx.isInsideDeclaration {
		setDeclaredVar(ctx, n.value, ast{})
	} else {
		if replace, ok := getReplaceImport(ctx, n.value); ok {
			return replace
		}
	}
	return n
}

func modifyBinaryExpression(n ast, e *environment, ctx *context) ast {
	operator := n.value
	left := modifyAst(n.children[0], e, ctx)
	right := modifyAst(n.children[1], e, ctx)

	if operator == "===" || operator == "==" {
		if left.t == n_STRING_LITERAL && left.t == right.t {
			if printAst(left) == printAst(right) {
				return newNode(n_BOOL_LITERAL, "true")
			}
			return newNode(n_BOOL_LITERAL, "false")
		}
	} else if operator == "!==" || operator == "!=" {
		if left.t == n_STRING_LITERAL && left.t == right.t {
			if printAst(left) != printAst(right) {
				return newNode(n_BOOL_LITERAL, "true")
			}
			return newNode(n_BOOL_LITERAL, "false")
		}
	}

	return newNode(n_BINARY_EXPRESSION, operator, left, right)
}

func modifyMemberExpression(n ast, e *environment, ctx *context) ast {
	object := n.children[0]
	prop := n.children[1]

	if object.t == n_MEMBER_EXPRESSION {
		innerObject := object.children[0]
		innerProp := object.children[1]

		if innerObject.t == n_IDENTIFIER &&
			innerObject.value == "process" &&
			innerProp.t == n_IDENTIFIER &&
			innerProp.value == "env" {

			if envVal, ok := e.env[prop.value]; ok {
				return newNode(n_STRING_LITERAL, envVal)
			}
		}
	}

	return modifyChildren(n, e, ctx)
}

func modifyIfStatement(n ast, e *environment, ctx *context) ast {
	cond := modifyAst(n.children[0], e, ctx)
	conseq := n.children[1]
	altern := n.children[2]

	if cond.t == n_BOOL_LITERAL {
		if cond.value == "true" {
			return modifyAst(conseq, e, ctx)
		}
		return modifyAst(altern, e, ctx)
	}

	return modifyChildren(n, e, ctx)
}

func modifyProgram(n ast, e *environment, ctx *context) ast {
	children := []ast{}
	for _, c := range n.children {
		children = append(children, modifyAst(c, e, ctx))
	}
	n.children = children

	moduleName := loaders.CreateVarNameFromPath(e.fileName)

	hasES6ExportsStr := "false"
	if e.hasES6Exports {
		hasES6ExportsStr = "true"
	}

	beforeSt := []ast{
		newNode(
			n_DECLARATION_STATEMENT, "",
			newNode(
				n_VARIABLE_DECLARATION, "var",
				newNode(
					n_DECLARATOR, "",
					newNode(n_IDENTIFIER, "module"),
					newNode(
						n_OBJECT_LITERAL, "",
						newNode(
							n_OBJECT_PROPERTY, "",
							newNode(n_OBJECT_KEY, "exports"),
							newNode(n_OBJECT_LITERAL, ""),
						),
						newNode(
							n_OBJECT_PROPERTY, "",
							newNode(n_OBJECT_KEY, "es6"),
							newNode(n_OBJECT_LITERAL, ""),
						),
						newNode(
							n_OBJECT_PROPERTY, "",
							newNode(n_OBJECT_KEY, "hasES6Exports"),
							newNode(n_BOOL_LITERAL, hasES6ExportsStr),
						),
					),
				),
				newNode(
					n_DECLARATOR, "",
					newNode(n_IDENTIFIER, "exports"),
					newNode(
						n_MEMBER_EXPRESSION, "",
						newNode(n_IDENTIFIER, "module"),
						newNode(n_IDENTIFIER, "exports"),
					),
				),
			),
		),
	}

	afterSt := []ast{
		newNode(
			n_RETURN_STATEMENT, "",
			newNode(n_IDENTIFIER, "module"),
		),
	}

	allSt := append(
		beforeSt,
		append(n.children, afterSt...)...,
	)

	funcExpr := newNode(
		n_FUNCTION_EXPRESSION, "",
		newNode(n_EMPTY, ""),
		newNode(n_FUNCTION_PARAMETERS, ""),
		newNode(n_BLOCK_STATEMENT, "", allSt...),
	)

	decl := newNode(
		n_EXPRESSION_STATEMENT, "",
		newNode(
			n_ASSIGNMENT_EXPRESSION, "=",
			newNode(
				n_MEMBER_EXPRESSION, "",
				newNode(n_IDENTIFIER, "moduleFns"),
				newNode(n_IDENTIFIER, moduleName),
			),
			funcExpr,
		),
	)

	return decl
}

func getImport(resolvedPath, eType, name string) ast {
	moduleName := loaders.CreateVarNameFromPath(resolvedPath)

	funcName := eType

	args := []ast{
		newNode(
			n_MEMBER_EXPRESSION, "",
			newNode(n_IDENTIFIER, "modules"),
			newNode(n_IDENTIFIER, moduleName),
		),
	}

	if name != "" {
		args = append(args, newNode(n_STRING_LITERAL, name))
	}

	return newNode(
		n_FUNCTION_CALL, "",
		newNode(n_IDENTIFIER, funcName),
		newNode(n_FUNCTION_ARGS, "", args...),
	)
}

func modifyFunctionCall(n ast, e *environment, ctx *context) ast {
	children := []ast{}
	for _, c := range n.children {
		children = append(children, modifyAst(c, e, ctx))
	}
	n.children = children

	name := n.children[0]

	// modify require
	if name.value == "require" {
		args := n.children[1].children

		path := args[0].value
		resolvedPath := resolveES6ImportPath(path, e.fileName)
		e.fileImports = append(e.fileImports, resolvedPath)
		moduleName := loaders.CreateVarNameFromPath(resolvedPath)

		args[0] = newNode(
			n_MEMBER_EXPRESSION, "",
			newNode(n_IDENTIFIER, "modules"),
			newNode(n_IDENTIFIER, moduleName),
		)
		return n
	}

	return n
}

func modifyImport(n ast, e *environment, ctx *context) ast {
	n = modifyChildren(n, e, ctx)

	vars := n.children[0].children
	all := n.children[1].value
	path := n.children[2].value
	resolvedPath := resolveES6ImportPath(path, e.fileName)
	e.fileImports = append(e.fileImports, resolvedPath)

	decls := []ast{}

	if all != "" {
		alias := newNode(n_IDENTIFIER, all)
		importReplace := getImport(resolvedPath, "requireES6", "*")

		allDecl := newNode(
			n_DECLARATOR, "", alias, importReplace,
		)
		decls = append(decls, allDecl)
		setDeclaredVar(ctx, all, importReplace)
	}

	for _, v := range vars {
		name := v.children[0]
		alias := v.children[1]
		importReplace := getImport(resolvedPath, "requireES6", name.value)

		d := newNode(
			n_DECLARATOR, "", alias, importReplace,
		)
		decls = append(decls, d)
		setDeclaredVar(ctx, alias.value, importReplace)
	}

	if len(decls) > 0 {
		return newNode(
			n_DECLARATION_STATEMENT, "",
			newNode(n_VARIABLE_DECLARATION, "var", decls...),
		)
	}
	return newNode(n_EMPTY, "")
}

func modifyExport(n ast, e *environment, ctx *context) ast {
	e.hasES6Exports = true

	n = modifyChildren(n, e, ctx)

	vars := n.children[0].children
	declaration := n.children[1]
	path := n.children[2].value

	es6Exports := newNode(
		n_MEMBER_EXPRESSION, "",
		newNode(n_IDENTIFIER, "module"),
		newNode(n_IDENTIFIER, "es6"),
	)

	if path != "" {
		resolvedPath := resolveES6ImportPath(path, e.fileName)
		e.fileImports = append(e.fileImports, resolvedPath)

		if n.children[0].t == n_EXPORT_ALL {
			return newNode(
				n_EXPRESSION_STATEMENT, "",
				newNode(
					n_FUNCTION_CALL, "",
					newNode(
						n_MEMBER_EXPRESSION, "",
						newNode(n_IDENTIFIER, "Object"),
						newNode(n_IDENTIFIER, "assign"),
					),
					newNode(
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

		left := newNode(
			n_MEMBER_EXPRESSION, "",
			es6Exports,
			newNode(n_IDENTIFIER, alias),
		)

		var right ast
		if path != "" {
			resolvedPath := resolveES6ImportPath(path, e.fileName)
			right = getImport(resolvedPath, "requireES6", name.value)
		} else {
			right = name
		}

		ass := newNode(n_ASSIGNMENT_EXPRESSION, "=", left, right)
		assignments = append(assignments, ass)
	}

	seq := newNode(n_SEQUENCE_EXPRESSION, "", assignments...)
	exprSt := newNode(n_EXPRESSION_STATEMENT, "", seq)

	multi := newNode(n_MULTISTATEMENT, "", declaration, exprSt)

	return multi
}

func transformIntoModule(src ast, fileName string, env map[string]string) (ast, []string) {
	e := environment{}
	e.fileImports = []string{}
	e.fileName = fileName
	e.env = env

	ctx := context{}
	ctx.declaredVars = map[string]ast{}

	res := modifyAst(src, &e, &ctx)

	return res, e.fileImports
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
