package jsLoader

import (
	"fmt"
	"strings"
)

type LoaderError struct {
	err      string
	fileName string
}

func (le LoaderError) Error() string {
	return fmt.Sprintf("Error loading file %s:\n%s", le.fileName, le.err)
}

func getNiceError(pe *parsingError, src []byte) string {
	start := pe.t.charIndex
	for src[start] != '\n' {
		start--
	}
	start++

	end := pe.t.charIndex
	for src[end] != '\n' {
		end++
	}

	resLine0 := fmt.Sprintf("Unexpected token '%s' at %v:%v\n", pe.t.lexeme, pe.t.line, pe.t.column)
	resStart := fmt.Sprintf("%04d", pe.t.line) + " "
	resLine1 := resStart + string(src[start:end]) + "\n"
	resLine2 := ""
	for i := 0; i <= pe.t.column+len(resStart); i++ {
		resLine1 += " "
	}
	resLine2 += "^"

	return resLine0 + resLine1 + resLine2
}

func LoadFile(src []byte, filePath string) ([]byte, []string, error) {
	tokens := lex(src)
	initialProgram, parseErr := parseTokens(tokens)
	if parseErr != nil {
		loaderErr := LoaderError{}

		loaderErr.err = getNiceError(parseErr, src)
		loaderErr.fileName = filePath

		return nil, nil, loaderErr
	}

	resultProgram, fileImports := transformIntoModule(initialProgram, filePath)
	resultBytes := []byte(printAst(resultProgram))
	return resultBytes, fileImports, nil
}

type importedVar struct {
}

type context struct {
	fileName      string
	fileImports   []string
	hasES6Exports bool
}

func modifyAst(n ast, ctx *context) ast {
	switch n.t {

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
	varName := "'" + name + "'"

	args := []ast{
		makeNode(
			n_MEMBER_EXPRESSION, "",
			makeNode(n_IDENTIFIER, "modules"),
			makeNode(n_IDENTIFIER, moduleName),
		),
	}

	if name != "" {
		args = append(args, makeNode(n_IDENTIFIER, varName))
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

func transformIntoModule(src ast, fileName string) (ast, []string) {
	ctx := context{}
	ctx.fileImports = []string{}
	ctx.fileName = fileName

	res := modifyAst(src, &ctx)

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

func resolveES6ImportPath(importPath, currentFileName string) string {
	importPath = trimQuotesFromString(importPath)
	pathParts := strings.Split(importPath, "/")

	locationParts := strings.Split(currentFileName, "/")
	locationParts = locationParts[:len(locationParts)-1]

	// import from node_modules
	if len(pathParts) > 0 {
		if pathParts[0] != "." && pathParts[0] != ".." {
			locationParts = []string{"node_modules"}
			if len(pathParts) == 1 {
				pathParts = append(pathParts, "index.js")
			}
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

	ext := ""
	if strings.Index(pathParts[len(pathParts)-1], ".") < 0 {
		ext = ".js"
	}

	result := fullFileName + ext
	return result
}
