package jsLoader

func generateJsCode(n ast) string {
	switch n.t {
	case g_IMPORT_ALIAS:
		return n.value

	case g_DECLARATION_STATEMENT:
		return generateJsCode(n.children[0]) + ";"

	case g_DECLARATION_EXPRESSION:
		result := n.value + " "
		for i, d := range n.children {
			result += generateJsCode(d) // declarator
			if i < len(n.children)-1 {
				result += ","
			}
		}
		return result

	case g_DECLARATOR:
		result := generateJsCode(n.children[0]) // var name
		if len(n.children) > 1 {
			result += "=" + generateJsCode(n.children[1]) // var value
		}
		return result

	case g_NAME, g_NUMBER_LITERAL, g_NULL, g_STRING_LITERAL,
		g_BOOL_LITERAL, g_UNDEFINED, g_THIS:
		return n.value

	case g_FUNCTION_DECLARATION:
		result := "function " + n.value
		result += generateJsCode(n.children[0]) // function params
		result += generateJsCode(n.children[1]) // function body

		return result

	case g_FUNCTION_PARAMETERS, g_FUNCTION_ARGS:
		result := "("
		for i, param := range n.children {
			result += generateJsCode(param)
			if i < len(n.children)-1 {
				result += ","
			}
		}
		result += ")"

		return result

	case g_BLOCK_STATEMENT:
		result := "{"
		for _, st := range n.children {
			result += generateJsCode(st)
		}
		result += "}"

		return result

	case g_EXPRESSION_STATEMENT:
		return generateJsCode(n.children[0]) + ";"

	case g_IF_STATEMENT:
		result := "if("
		result += generateJsCode(n.children[0]) + ")" // condition
		result += generateJsCode(n.children[1])       // body
		if len(n.children) > 2 {
			result += " else " + generateJsCode(n.children[2]) // alternate
		}

		return result

	case g_EXPRESSION:
		result := generateJsCode(n.children[0]) // left
		if n.value == "in" || n.value == "instanceof" {
			result += " " + n.value + " "
		} else {
			result += n.value // operator
		}
		result += generateJsCode(n.children[1]) // right

		return result

	case g_OBJECT_PROPERTY:
		result := ""
		if n.value != "" {
			result += n.value + " " // set or get
		}
		result += generateJsCode(n.children[0]) // key

		if len(n.children) > 1 {
			value := n.children[1]
			if value.t != g_MEMBER_FUNCTION {
				result += ":"
			}

			result += generateJsCode(value) // value
		}
		return result

	case g_OBJECT_PATTERN, g_OBJECT_LITERAL:
		result := "{"
		for i, prop := range n.children {
			result += generateJsCode(prop)

			if i < len(n.children)-1 {
				result += ","
			}
		}
		result += "}"

		return result

	case g_ASSIGNMENT_PATTERN:
		result := generateJsCode(n.children[0]) // left
		result += "="
		result += generateJsCode(n.children[1]) // right

		return result

	case g_UNARY_POSTFIX_EXPRESSION:
		result := generateJsCode(n.children[0]) // value
		result += n.value

		return result

	case g_UNARY_PREFIX_EXPRESSION:
		result := n.value
		if n.value == "typeof" || n.value == "void" || n.value == "delete" {
			result += " " + generateJsCode(n.children[0]) + " "
		} else {
			result += generateJsCode(n.children[0]) // value
		}

		return result

	case g_ARRAY_LITERAL, g_ARRAY_PATTERN:
		result := "["
		for i, item := range n.children {
			result += generateJsCode(item)
			if i < len(n.children)-1 {
				result += ","
			}
		}
		result += "]"

		return result

	case g_LAMBDA_EXPRESSION:
		result := ""
		result += generateJsCode(n.children[0]) // params
		result += "=>"
		result += generateJsCode(n.children[1]) // body

		return result

	case g_FUNCTION_CALL:
		result := generateJsCode(n.children[0]) // func name
		if len(n.children) > 1 {
			result += generateJsCode(n.children[1]) // func args
		} else {
			result += "()"
		}

		return result

	case g_MEMBER_EXPRESSION:
		result := generateJsCode(n.children[0]) // object

		prop := n.children[1]
		if prop.t == g_CALCULATED_PROPERTY_NAME {
			result += generateJsCode(prop)
		} else {
			result += "."
			result += generateJsCode(prop) // property
		}

		return result

	case g_CALCULATED_PROPERTY_NAME:
		return "[" + generateJsCode(n.children[0]) + "]"

	case g_CONSTRUCTOR_CALL:
		result := "new " + generateJsCode(n.children[0]) // func name
		result += generateJsCode(n.children[1])          // args

		return result

	case g_VALID_PROPERTY_NAME:
		return n.value

	case g_MEMBER_FUNCTION:
		result := generateJsCode(n.children[0]) // params
		result += generateJsCode(n.children[1]) // body

		return result

	case g_FUNCTION_EXPRESSION:
		result := "function"
		if n.value != "" {
			result += " " + n.value
		}
		result += generateJsCode(n.children[0]) // params
		result += generateJsCode(n.children[1]) // body
		return result

	case g_PARENS_EXPRESSION:
		result := "(" + generateJsCode(n.children[0]) + ")"
		return result

	case g_EMPTY_STATEMENT:
		return ";"

	case g_FOR_STATEMENT:
		result := "for("
		result += generateJsCode(n.children[0]) + ";" // init
		result += generateJsCode(n.children[1]) + ";" // test
		result += generateJsCode(n.children[2]) + ")" // final
		result += generateJsCode(n.children[3])       // body

		return result

	case g_FOR_OF_STATEMENT:
		result := "for("
		result += generateJsCode(n.children[0]) + " of " // init
		result += generateJsCode(n.children[1]) + ")"    // test
		result += generateJsCode(n.children[2])          // body

		return result

	case g_FOR_IN_STATEMENT:
		result := "for("
		result += generateJsCode(n.children[0]) + " in " // init
		result += generateJsCode(n.children[1]) + ")"    // test
		result += generateJsCode(n.children[2])          // body

		return result

	case g_WHILE_STATEMENT:
		result := "while("
		result += generateJsCode(n.children[0]) + ")" // condition
		result += generateJsCode(n.children[1])       // body

		return result

	case g_SEQUENCE_EXPRESSION:
		result := ""
		for i, c := range n.children {
			result += generateJsCode(c)
			if i < len(n.children)-1 {
				result += ","
			}
		}

		return result

	case g_DO_WHILE_STATEMENT:
		result := "do "
		result += generateJsCode(n.children[1])                   // body
		result += "while(" + generateJsCode(n.children[0]) + ");" // condition

		return result

	case g_IMPORT_STATEMENT:
		result := "import"

		vars := n.children[0].children

		alreadyPrintedVar := -1

		// print default first
		if len(vars) > 0 {
			for i, v := range vars {
				name := v.children[0]
				alias := v.children[1]

				if name.value == "default" {
					result += " " + alias.value
					alreadyPrintedVar = i
					break
				}
			}
		}

		all := generateJsCode(n.children[1]) // import all
		if all != "" {
			if result != "import" {
				result += ","
			}
			result += all
		}

		if !(len(vars) == 1 && alreadyPrintedVar == 0 || len(vars) == 0) {
			if result != "import" {
				result += ","
			}
			result += "{"
			for i := 0; i < len(vars); i++ {
				if i != alreadyPrintedVar {
					result += generateJsCode(vars[i])
					if i < len(vars)-1 {
						result += ","
					}
				}
			}
			result += "}"
		}

		if result != "import" {
			result += " from"
		}
		result += generateJsCode(n.children[2]) // import path

		result += ";"
		return result

	case g_IMPORT_VAR:
		name := n.children[0]

		result := name.value
		if len(n.children) > 1 {
			result += " as " + n.children[1].value
		}

		return result

	case g_IMPORT_PATH:
		return n.value

	case g_IMPORT_ALL:
		if n.value != "" {
			return "*as " + n.value
		}

	case g_BREAK_STATEMENT:
		return "break;"

	case g_CONTINUE_STATEMENT:
		return "continue;"

	case g_DEBUGGER_STATEMENT:
		return "debugger;"

	case g_EXPORT_STATEMENT:
		result := ""
		noDefault := true

		declaration := n.children[1]
		result += generateJsCode(declaration)

		result += "export"
		vars := n.children[0].children
		if len(vars) == 1 {
			v := vars[0]
			alias := v.children[1]
			name := v.children[0]
			if alias.value == "default" {
				result += " default " + generateJsCode(name)
				noDefault = false
			}
		}

		if noDefault {
			result += "{"
			for i, v := range vars {
				result += generateJsCode(v)
				if i < len(vars)-1 {
					result += ","
				}
			}
			result += "}"
		}

		result += generateJsCode(n.children[2]) // export path

		result += ";"
		return result

	case g_EXPORT_PATH:
		if n.value != "" {
			return " from" + n.value
		}
		return ""

	case g_EXPORT_VAR:
		result := generateJsCode(n.children[0]) // name
		if len(n.children) > 1 {
			result += " as " + generateJsCode(n.children[1]) // alias
		}
		return result

	case g_EXPORT_ALIAS:
		return n.value

	case g_EXPORT_NAME:
		return n.value

	case g_SWITCH_STATEMENT:
		result := "switch("
		result += generateJsCode(n.children[0]) + ")" // test
		result += "{"
		for _, c := range n.children[1].children {
			result += generateJsCode(c)
		}

		if len(n.children) > 2 {
			result += generateJsCode(n.children[2]) // default
		}
		result += "}"

		return result

	case g_SWITCH_DEFAULT:
		result := "default:"
		for _, st := range n.children {
			result += generateJsCode(st)
		}
		return result

	case g_SWITCH_CASE:
		result := "case "
		result += generateJsCode(n.children[0]) + ":" // test
		for _, st := range n.children[1].children {
			result += generateJsCode(st)
		}
		return result

	case g_SWITCH_CASE_TEST:
		result := generateJsCode(n.children[0])
		return result

	case g_PROGRAM_STATEMENT:
		result := ""
		for _, st := range n.children {
			result += generateJsCode(st)
		}
		return result

	case g_RETURN_STATEMENT:
		result := "return"
		if len(n.children) > 0 {
			result += " " + generateJsCode(n.children[0])
		}
		result += ";"

		return result

	case g_IMPORT_NAME:
		return n.value

	case g_MULTISTATEMENT:
		result := ""
		for _, st := range n.children {
			result += generateJsCode(st)
		}
		return result

	case g_HEX_LITERAL:
		return n.value

	case g_REGEXP_LITERAL:
		return n.value

	case g_TRY_CATCH_STATEMENT:
		result := "try"
		result += generateJsCode(n.children[0])
		catch := generateJsCode(n.children[2])
		finally := generateJsCode(n.children[3])
		if catch != "" {
			result += "catch(" + generateJsCode(n.children[1]) + ")"
			result += catch
		}
		if finally != "" {
			result += "finally"
			result += finally
		}

		return result

	case g_THROW_STATEMENT:
		return "throw " + generateJsCode(n.children[0])

	case g_CONDITIONAL_EXPRESSION:
		return generateJsCode(n.children[0]) + "?" + generateJsCode(n.children[1]) + ":" + generateJsCode(n.children[2])

	case g_MARKER:
		return n.value + ":"

	case g_REST_EXPRESSION, g_SPREAD_EXPRESSION:
		return "..." + generateJsCode(n.children[0])

	case g_CLASS_STATEMENT:
		return generateJsCode(n.children[0])

	case g_CLASS_EXPRESSION:
		result := "class"
		name := n.children[0]
		if name.value != "" {
			result += " " + generateJsCode(n.children[0])
		}
		extends := n.children[1]
		if extends.t != g_UNDEFINED {
			result += " extends " + generateJsCode(extends)
		}

		body := n.children[2]
		result += generateJsCode(body)

		return result

	case g_CLASS_METHOD:
		result := ""
		result += generateJsCode(n.children[0]) // name
		result += generateJsCode(n.children[1]) // args and body

		return result

	case g_CLASS_PROPERTY:
		result := generateJsCode(n.children[0])
		if len(n.children) > 1 {
			result += "=" + generateJsCode(n.children[1])
		}

		return result + ";"

	case g_CLASS_STATIC_PROP:
		return "static " + generateJsCode(n.children[0])

	}

	return ""
}
