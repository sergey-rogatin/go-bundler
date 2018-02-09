package jsLoader

func printAst(n astNode) string {
	switch n.t {
	case g_IMPORT_ALIAS:
		return n.value

	case g_DECLARATION_STATEMENT:
		return printAst(n.children[0]) + ";"

	case g_DECLARATION_EXPRESSION:
		result := n.value + " "
		for i, d := range n.children {
			result += printAst(d) // declarator
			if i < len(n.children)-1 {
				result += ","
			}
		}
		return result

	case g_DECLARATOR:
		result := printAst(n.children[0]) // var name
		if len(n.children) > 1 {
			result += "=" + printAst(n.children[1]) // var value
		}
		return result

	case g_NAME, g_NUMBER_LITERAL, g_NULL, g_STRING_LITERAL,
		g_BOOL_LITERAL, g_UNDEFINED, g_THIS:
		return n.value

	case g_FUNCTION_DECLARATION:
		result := "function " + n.value
		result += printAst(n.children[0]) // function params
		result += printAst(n.children[1]) // function body

		return result

	case g_FUNCTION_PARAMETERS, g_FUNCTION_ARGS:
		result := "("
		for i, param := range n.children {
			result += printAst(param)
			if i < len(n.children)-1 {
				result += ","
			}
		}
		result += ")"

		return result

	case g_FUNCTION_PARAMETER:
		result := ""
		if n.flags&f_FUNCTION_PARAM_REST != 0 {
			result += "..."
		}
		result += printAst(n.children[0]) // param name
		if len(n.children) > 1 {
			result += "=" + printAst(n.children[1]) // param default value
		}

		return result

	case g_BLOCK_STATEMENT:
		result := "{"
		for _, st := range n.children {
			result += printAst(st)
		}
		result += "}"

		return result

	case g_EXPRESSION_STATEMENT:
		return printAst(n.children[0]) + ";"

	case g_IF_STATEMENT:
		result := "if("
		result += printAst(n.children[0]) + ")" // condition
		result += printAst(n.children[1])       // body
		if len(n.children) > 2 {
			result += " else " + printAst(n.children[2]) // alternate
		}

		return result

	case g_EXPRESSION:
		result := printAst(n.children[0]) // left
		if n.value == "in" || n.value == "instanceof" {
			result += " " + n.value + " "
		} else {
			result += n.value // operator
		}
		result += printAst(n.children[1]) // right

		return result

	case g_OBJECT_PROPERTY:
		result := ""
		if n.value != "" {
			result += n.value + " " // set or get
		}
		result += printAst(n.children[0]) // key

		if len(n.children) > 1 {
			value := n.children[1]
			if value.t != g_MEMBER_FUNCTION {
				result += ":"
			}

			result += printAst(value) // value
		}
		return result

	case g_OBJECT_PATTERN, g_OBJECT_LITERAL:
		result := "{"
		for i, prop := range n.children {
			result += printAst(prop)

			if i < len(n.children)-1 {
				result += ","
			}
		}
		result += "}"

		return result

	case g_ASSIGNMENT_PATTERN:
		result := printAst(n.children[0]) // left
		result += "="
		result += printAst(n.children[1]) // right

		return result

	case g_UNARY_POSTFIX_EXPRESSION:
		result := printAst(n.children[0]) // value
		result += n.value

		return result

	case g_UNARY_PREFIX_EXPRESSION:
		result := n.value
		if n.value == "typeof" || n.value == "void" || n.value == "delete" {
			result += " " + printAst(n.children[0]) + " "
		} else {
			result += printAst(n.children[0]) // value
		}

		return result

	case g_ARRAY_LITERAL, g_ARRAY_PATTERN:
		result := "["
		for i, item := range n.children {
			result += printAst(item)
			if i < len(n.children)-1 {
				result += ","
			}
		}
		result += "]"

		return result

	case g_LAMBDA_EXPRESSION:
		result := ""
		result += printAst(n.children[0]) // params
		result += "=>"
		result += printAst(n.children[1]) // body

		return result

	case g_FUNCTION_CALL:
		result := printAst(n.children[0]) // func name
		if len(n.children) > 1 {
			result += printAst(n.children[1]) // func args
		} else {
			result += "()"
		}

		return result

	case g_MEMBER_EXPRESSION:
		result := printAst(n.children[0]) // object

		prop := n.children[1]
		if prop.t == g_CALCULATED_PROPERTY_NAME {
			result += printAst(prop)
		} else {
			result += "."
			result += printAst(prop) // property
		}

		return result

	case g_CALCULATED_PROPERTY_NAME:
		return "[" + printAst(n.children[0]) + "]"

	case g_CONSTRUCTOR_CALL:
		result := "new " + printAst(n.children[0]) // func name
		result += printAst(n.children[1])          // args

		return result

	case g_VALID_PROPERTY_NAME:
		return n.value

	case g_MEMBER_FUNCTION:
		result := printAst(n.children[0]) // params
		result += printAst(n.children[1]) // body

		return result

	case g_FUNCTION_EXPRESSION:
		result := "function"
		if n.value != "" {
			result += " " + n.value
		}
		result += printAst(n.children[0]) // params
		result += printAst(n.children[1]) // body
		return result

	case g_PARENS_EXPRESSION:
		result := "(" + printAst(n.children[0]) + ")"
		return result

	case g_EMPTY_STATEMENT:
		return ";"

	case g_FOR_STATEMENT:
		result := "for("
		result += printAst(n.children[0]) + ";" // init
		result += printAst(n.children[1]) + ";" // test
		result += printAst(n.children[2]) + ")" // final
		result += printAst(n.children[3])       // body

		return result

	case g_FOR_OF_STATEMENT:
		result := "for("
		result += printAst(n.children[0]) + " of " // init
		result += printAst(n.children[1]) + ")"    // test
		result += printAst(n.children[2])          // body

		return result

	case g_FOR_IN_STATEMENT:
		result := "for("
		result += printAst(n.children[0]) + " in " // init
		result += printAst(n.children[1]) + ")"    // test
		result += printAst(n.children[2])          // body

		return result

	case g_WHILE_STATEMENT:
		result := "while("
		result += printAst(n.children[0]) + ")" // condition
		result += printAst(n.children[1])       // body

		return result

	case g_SEQUENCE_EXPRESSION:
		result := ""
		for i, c := range n.children {
			result += printAst(c)
			if i < len(n.children)-1 {
				result += ","
			}
		}

		return result

	case g_DO_WHILE_STATEMENT:
		result := "do "
		result += printAst(n.children[1])                   // body
		result += "while(" + printAst(n.children[0]) + ");" // condition

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

		all := printAst(n.children[1]) // import all
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
					result += printAst(vars[i])
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
		result += printAst(n.children[2]) // import path

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
		result += printAst(declaration)

		result += "export"
		vars := n.children[0].children
		if len(vars) == 1 {
			v := vars[0]
			alias := v.children[1]
			name := v.children[0]
			if alias.value == "default" {
				result += " default " + printAst(name)
				noDefault = false
			}
		}

		if noDefault {
			result += "{"
			for i, v := range vars {
				result += printAst(v)
				if i < len(vars)-1 {
					result += ","
				}
			}
			result += "}"
		}

		result += printAst(n.children[2]) // export path

		result += ";"
		return result

	case g_EXPORT_PATH:
		if n.value != "" {
			return " from" + n.value
		}
		return ""

	case g_EXPORT_VAR:
		result := printAst(n.children[0]) // name
		if len(n.children) > 1 {
			result += " as " + printAst(n.children[1]) // alias
		}
		return result

	case g_EXPORT_ALIAS:
		return n.value

	case g_EXPORT_NAME:
		return n.value

	case g_SWITCH_STATEMENT:
		result := "switch("
		result += printAst(n.children[0]) + ")" // test
		result += "{"
		for _, c := range n.children[1].children {
			result += printAst(c)
		}

		if len(n.children) > 2 {
			result += printAst(n.children[2]) // default
		}
		result += "}"

		return result

	case g_SWITCH_DEFAULT:
		result := "default:"
		for _, st := range n.children {
			result += printAst(st)
		}
		return result

	case g_SWITCH_CASE:
		result := "case "
		result += printAst(n.children[0]) + ":" // test
		for _, st := range n.children[1].children {
			result += printAst(st)
		}
		return result

	case g_SWITCH_CASE_TEST:
		result := printAst(n.children[0])
		return result

	case g_PROGRAM_STATEMENT:
		result := ""
		for _, st := range n.children {
			result += printAst(st)
		}
		return result

	case g_RETURN_STATEMENT:
		result := "return"
		if len(n.children) > 0 {
			result += " " + printAst(n.children[0])
		}
		result += ";"

		return result

	case g_IMPORT_NAME:
		return n.value

	case g_MULTISTATEMENT:
		result := ""
		for _, st := range n.children {
			result += printAst(st)
		}
		return result

	case g_HEX_LITERAL:
		return n.value

	case g_REGEXP_LITERAL:
		return n.value

	case g_TRY_CATCH_STATEMENT:
		result := "try"
		result += printAst(n.children[0])
		catch := printAst(n.children[2])
		finally := printAst(n.children[3])
		if catch != "" {
			result += "catch(" + printAst(n.children[1]) + ")"
			result += catch
		}
		if finally != "" {
			result += "finally"
			result += finally
		}

		return result

	case g_THROW_STATEMENT:
		return "throw " + printAst(n.children[0])

	case g_CONDITIONAL_EXPRESSION:
		return printAst(n.children[0]) + "?" + printAst(n.children[1]) + ":" + printAst(n.children[2])

	case g_MARKER:
		return n.value + ":"
	}

	return ""
}
