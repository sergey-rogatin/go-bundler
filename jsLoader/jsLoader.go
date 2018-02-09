package jsLoader

import (
	"fmt"
	"path/filepath"
	"strings"
)

// "strings"

type LoaderError struct {
	err      error
	fileName string
}

func (le LoaderError) Error() string {
	return fmt.Sprintf("Error loading file %s:\n %s", le.fileName, le.err)
}

func LoadFile(src []byte, filePath string) ([]byte, []string, error) {
	tokens := lex(src)
	initialProgram, parseErr := parseTokens(tokens)
	if parseErr != nil {
		loaderErr := LoaderError{}
		loaderErr.err = parseErr
		loaderErr.fileName = filePath

		return nil, nil, loaderErr
	}

	resultProgram, fileImports := transformIntoModule(initialProgram, filePath)
	resultBytes := []byte(printAst(resultProgram))
	return resultBytes, fileImports, nil
}

func transformIntoModule(src astNode, fileName string) (astNode, []string) {
	fileImports := []string{}

	var modifyAst,
		modifyProgram,
		modifyImport,
		modifyExport,
		modifyFunctionCall func(astNode) astNode

	modifyAst = func(n astNode) astNode {
		switch n.t {

		case g_FUNCTION_CALL:
			return modifyFunctionCall(n)

		case g_EXPORT_STATEMENT:
			return modifyExport(n)

		case g_PROGRAM_STATEMENT:
			return modifyProgram(n)

		case g_IMPORT_STATEMENT:
			return modifyImport(n)

		default:
			res := n
			res.children = []astNode{}
			for _, c := range n.children {
				res.children = append(res.children, modifyAst(c))
			}
			return res
		}
	}

	modifyFunctionCall = func(n astNode) astNode {
		children := []astNode{}
		for _, c := range n.children {
			children = append(children, modifyAst(c))
		}
		n.children = children

		nameNode := children[0]
		args := children[1].children

		if nameNode.value == "require" {
			if len(args) == 1 && args[0].t == g_STRING_LITERAL {
				path := args[0].value
				resolvedPath := resolveES6ImportPath(path, fileName)
				fileImports = append(fileImports, resolvedPath)

				objectName := CreateVarNameFromPath(resolvedPath)
				object := makeNode(g_NAME, objectName)
				return object
			}

			fmt.Printf("Wrong arguments to require function")
			return n
		}

		return n
	}

	modifyProgram = func(n astNode) astNode {
		statements := []astNode{}

		// add var exports = {}
		exportsObj := makeNode(g_NAME, "exports")
		{
			right := makeNode(g_OBJECT_LITERAL, "")
			decl := makeNode(g_DECLARATOR, "", exportsObj, right)
			declExpr := makeNode(g_DECLARATION_EXPRESSION, "var", decl)
			declSt := makeNode(g_DECLARATION_STATEMENT, "", declExpr)
			statements = append(statements, declSt)
		}

		// add all other statements
		for _, st := range n.children {
			statements = append(statements, modifyAst(st))
		}

		// add return exports
		ret := makeNode(g_RETURN_STATEMENT, "", exportsObj)
		statements = append(statements, ret)

		params := makeNode(g_FUNCTION_PARAMETERS, "")
		blockSt := makeNode(g_BLOCK_STATEMENT, "", statements...)
		funcExpr := makeNode(g_FUNCTION_EXPRESSION, "", params, blockSt)
		funcCall := makeNode(g_FUNCTION_CALL, "", funcExpr)

		objectName := CreateVarNameFromPath(fileName)
		object := makeNode(g_NAME, objectName)

		decl := makeNode(g_DECLARATOR, "", object, funcCall)

		declExpr := makeNode(g_DECLARATION_EXPRESSION, "var", decl)
		declSt := makeNode(g_DECLARATION_STATEMENT, "", declExpr)

		return declSt
	}

	modifyImport = func(n astNode) astNode {
		vars := n.children[0].children
		importAll := n.children[1].value
		importPath := n.children[2].value

		resolvedPath := resolveES6ImportPath(importPath, fileName)
		fileImports = append(fileImports, resolvedPath)

		ext := filepath.Ext(resolvedPath)

		objectName := CreateVarNameFromPath(resolvedPath)
		object := makeNode(g_NAME, objectName)

		declarators := []astNode{}

		if importAll != "" {
			alias := makeNode(g_NAME, importAll)
			d := makeNode(g_DECLARATOR, "", alias, object)
			declarators = append(declarators, d)
		}

		for _, v := range vars {
			alias := v.children[1]

			var d astNode

			if ext == ".js" {
				property := v.children[0]
				member := makeNode(g_MEMBER_EXPRESSION, "", object, property)

				d = makeNode(g_DECLARATOR, "", alias, member)
			} else {
				filePath := "'" + objectName + ext + "'"
				fileURL := makeNode(g_STRING_LITERAL, filePath)

				d = makeNode(g_DECLARATOR, "", alias, fileURL)
			}

			declarators = append(declarators, d)
		}

		if len(declarators) > 0 {
			declExpr := makeNode(
				g_DECLARATION_EXPRESSION, "var", declarators...,
			)
			declSt := makeNode(g_DECLARATION_STATEMENT, "", declExpr)

			return declSt
		}

		return makeNode(g_EMPTY_EXPRESSION, "")
	}

	modifyExport = func(n astNode) astNode {
		vars := n.children[0].children
		exportsObj := makeNode(g_NAME, "exports")

		var importObj astNode
		pathNode := n.children[2]
		if pathNode.value != "" {
			resolvedPath := resolveES6ImportPath(pathNode.value, fileName)
			fileImports = append(fileImports, resolvedPath)
			objectName := CreateVarNameFromPath(resolvedPath)
			importObj = makeNode(g_NAME, objectName)
		}

		if !(n.flags&f_EXPORT_ALL != 0) {
			assignments := []astNode{}
			for _, v := range vars {
				exportedName := v.children[1]
				left := makeNode(g_MEMBER_EXPRESSION, "", exportsObj, exportedName)
				var right astNode

				if pathNode.value != "" {
					property := v.children[0]
					right = makeNode(g_MEMBER_EXPRESSION, "", importObj, property)
				} else {
					right = v.children[0]
				}

				d := makeNode(g_EXPRESSION, "=", left, right)
				assignments = append(assignments, d)
			}
			seqExpr := makeNode(g_SEQUENCE_EXPRESSION, "=", assignments...)
			exprSt := makeNode(g_EXPRESSION_STATEMENT, "", seqExpr)

			decl := n.children[1]

			multiSt := makeNode(g_MULTISTATEMENT, "", decl, exprSt)

			return multiSt
		}

		obj := makeNode(g_NAME, "Object")
		assign := makeNode(g_NAME, "assign")
		funcName := makeNode(g_MEMBER_EXPRESSION, "", obj, assign)

		args := []astNode{
			exportsObj,
			importObj,
		}
		argsNode := makeNode(g_FUNCTION_ARGS, "", args...)
		objectAssignCall := makeNode(g_FUNCTION_CALL, "", funcName, argsNode)

		exprSt := makeNode(g_EXPRESSION_STATEMENT, "", objectAssignCall)
		return exprSt
	}

	res := modifyAst(src)

	return res, fileImports
}

func CreateVarNameFromPath(path string) string {
	newName := strings.Replace(path, "/", "_", -1)
	newName = strings.Replace(newName, ".", "_", -1)
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
