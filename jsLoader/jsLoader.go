package jsLoader

import (
	"fmt"
	"path/filepath"
	"strings"
)

type LoaderError struct {
	err      error
	fileName string
}

func (le LoaderError) Error() string {
	return fmt.Sprintf("Error loading file %s:\n %s", le.fileName, le.err)
}

func LoadFile(src []byte, filePath string) ([]byte, []string, error) {
	tokens := lex(src)
	initialProgram, parseErr := parse(tokens)
	if parseErr != nil {
		loaderErr := LoaderError{}
		loaderErr.err = parseErr
		loaderErr.fileName = filePath

		return nil, nil, loaderErr
	}

	resultProgram, fileImports := transformIntoModule(initialProgram, filePath)
	resultBytes := []byte(resultProgram.String())
	return resultBytes, fileImports, nil
}

func CreateVarNameFromPath(path string) string {
	newName := strings.Replace(path, "/", "_", -1)
	newName = strings.Replace(newName, ".", "_", -1)
	return newName
}

func transformIntoModule(programAst program, fileName string) (program, []string) {
	result := program{}
	result.statements = []ast{}
	fileImports := []string{}

	expObj := CreateVarNameFromPath(fileName)

	bs := blockStatement{}
	bs.statements = []ast{}

	returnExport := objectLiteral{}
	returnExport.properties = []objectProperty{}

	for _, item := range programAst.statements {
		switch st := item.(type) {

		case exportStatement:
			for _, v := range st.vars {
				prop := objectProperty{}

				if v.name != nil {
					prop.value = v.name
				}
				prop.key.value = v.alias
				returnExport.properties = append(returnExport.properties, prop)
			}

		case importStatement:
			resolvedPath := resolveES6ImportPath(st.path.val.lexeme, fileName)
			exportObjName := CreateVarNameFromPath(resolvedPath)

			fileImports = append(fileImports, resolvedPath)
			ext := filepath.Ext(resolvedPath)

			if len(st.vars) > 0 {
				vd := varDeclaration{}
				vd.keyword = atom{makeToken("var")}
				vd.declarations = []declarator{}

				for _, impVar := range st.vars {
					decl := declarator{}
					decl.left = impVar.alias

					if ext == ".js" {
						me := memberExpression{}
						me.object = atom{makeToken(exportObjName)}
						me.property.value = impVar.name

						decl.value = me
					} else {
						decl.value = atom{makeToken("\"" + exportObjName + ext + "\"")}
					}
					vd.declarations = append(vd.declarations, decl)
				}

				declStatement := expressionStatement{}
				declStatement.expr = vd
				bs.statements = append(bs.statements, declStatement)
			}

		default:
			bs.statements = append(bs.statements, st)
		}
	}

	rs := returnStatement{}
	rs.value = returnExport

	bs.statements = append(bs.statements, rs)

	fe := functionExpression{}
	fe.body = bs

	fc := functionCall{}
	fc.name = fe

	decl := declarator{}
	decl.left = atom{makeToken(expObj)}
	decl.value = fc

	vd := varDeclaration{}
	vd.declarations = []declarator{decl}
	vd.keyword = atom{makeToken("var")}

	es := expressionStatement{}
	es.expr = vd
	result.statements = append(result.statements, es)

	return result, fileImports
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
