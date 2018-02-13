package jsLoader

func printAst(n ast) string {
	switch n.t {
	case n_IDENTIFIER,
		n_NUMBER_LITERAL,
		n_STRING_LITERAL,
		n_CONTROL_WORD,
		n_IMPORT_ALIAS,
		n_IMPORT_NAME,
		n_EXPORT_ALIAS,
		n_EXPORT_NAME:
		return n.value

	case n_LAMBDA_EXPRESSION:
		params := n.children[0]
		body := n.children[1]
		return printAst(params) + "=>" + printAst(body)

	case n_FUNCTION_PARAMETERS, n_FUNCTION_ARGS:
		res := "("
		for i, param := range n.children {
			res += printAst(param)
			if i < len(n.children)-1 {
				res += ","
			}
		}
		res += ")"
		return res

	case n_BLOCK_STATEMENT, n_CLASS_BODY:
		res := "{"
		for _, st := range n.children {
			res += printAst(st)
		}
		res += "}"
		return res

	case n_BINARY_EXPRESSION, n_ASSIGNMENT_EXPRESSION:
		left := n.children[0]
		right := n.children[1]
		res := printAst(left)
		if n.value == "in" || n.value == "instanceof" {
			res += " " + n.value + " "
		} else {
			res += n.value
		}
		res += printAst(right)
		return res

	case n_PAREN_EXPRESSION:
		content := n.children[0]
		return "(" + printAst(content) + ")"

	case n_REGEXP_LITERAL:
		body := n.children[0]
		flags := n.children[1]
		return "/" + body.value + "/" + flags.value

	case n_OBJECT_LITERAL, n_OBJECT_PATTERN:
		res := "{"
		for i, prop := range n.children {
			res += printAst(prop)
			if i < len(n.children)-1 {
				res += ","
			}
		}
		res += "}"
		return res

	case n_OBJECT_PROPERTY:
		key := n.children[0]
		value := n.children[1]

		res := printAst(key)
		if value.t != n_EMPTY {
			res += ":" + printAst(value)
		}
		return res

	case n_OBJECT_METHOD, n_CLASS_METHOD:
		kind := n.value
		key := n.children[0]
		params := n.children[1]
		body := n.children[2]

		res := ""
		if kind != "" {
			res += kind + " "
		}
		res += printAst(key) + printAst(params) + printAst(body)
		return res

	case n_NON_IDENTIFIER_OBJECT_KEY,
		n_OBJECT_KEY:
		return n.value

	case n_CALCULATED_OBJECT_KEY:
		return "[" + printAst(n.children[0]) + "]"

	case n_EXPRESSION_STATEMENT:
		expr := n.children[0]
		return printAst(expr) + ";"

	case n_PROGRAM:
		res := ""
		for _, st := range n.children {
			res += printAst(st)
		}
		return res

	case n_MEMBER_EXPRESSION:
		object := n.children[0]
		property := n.children[1]
		res := printAst(object)
		if property.t != n_CALCULATED_PROPERTY_NAME {
			res += "."
		}
		res += printAst(property)
		return res

	case n_CALCULATED_PROPERTY_NAME:
		return "[" + printAst(n.children[0]) + "]"

	case n_CONSTRUCTOR_CALL:
		name := n.children[0]
		args := n.children[1]
		return "new " + printAst(name) + printAst(args)

	case n_FUNCTION_CALL:
		name := n.children[0]
		args := n.children[1]
		return printAst(name) + printAst(args)

	case n_PROPERTY_NAME:
		return n.value

	case n_SPREAD_ELEMENT:
		value := n.children[0]
		return "..." + printAst(value)

	case n_FUNCTION_DECLARATION, n_FUNCTION_EXPRESSION:
		name := n.children[0]
		params := n.children[1]
		body := n.children[2]
		res := "function"
		if name.value != "" {
			res += " " + printAst(name)
		}
		res += printAst(params) + printAst(body)
		return res

	case n_ASSIGNMENT_PATTERN:
		left := n.children[0]
		right := n.children[1]
		return printAst(left) + "=" + printAst(right)

	case n_REST_ELEMENT:
		return "..." + printAst(n.children[0])

	case n_ARRAY_LITERAL:
		res := "["
		for i, item := range n.children {
			res += printAst(item)
			if i < len(n.children)-1 {
				res += ","
			}
		}
		res += "]"
		return res

	case n_EMPTY_STATEMENT:
		return ";"

	case n_FOR_STATEMENT:
		init := n.children[0]
		test := n.children[1]
		final := n.children[2]
		body := n.children[3]

		res := "for("
		res += printAst(init) + ";"
		res += printAst(test) + ";"
		res += printAst(final) + ")"
		res += printAst(body)
		return res

	case n_VARIABLE_DECLARATION:
		kind := n.value
		res := kind + " "
		for i, decl := range n.children {
			res += printAst(decl)
			if i < len(n.children)-1 {
				res += ","
			}
		}
		return res

	case n_DECLARATOR:
		name := n.children[0]
		init := n.children[1]
		res := printAst(name)
		if init.t != n_EMPTY {
			res += "=" + printAst(init)
		}
		return res

	case n_PREFIX_UNARY_EXPRESSION:
		return n.value + printAst(n.children[0])

	case n_POSTFIX_UNARY_EXPRESSION:
		return printAst(n.children[0]) + n.value

	case n_FOR_IN_STATEMENT:
		left := n.children[0]
		right := n.children[1]
		body := n.children[2]
		return "for(" + printAst(left) +
			" in " + printAst(right) + ")" + printAst(body)

	case n_FOR_OF_STATEMENT:
		left := n.children[0]
		right := n.children[1]
		body := n.children[2]
		return "for(" + printAst(left) +
			" of " + printAst(right) + ")" + printAst(body)

	case n_WHILE_STATEMENT:
		cond := n.children[0]
		body := n.children[1]
		return "while(" + printAst(cond) + ")" +
			printAst(body)

	case n_DO_WHILE_STATEMENT:
		cond := n.children[0]
		body := n.children[1]
		return "do " + printAst(body) + "while(" +
			printAst(cond) + ")"

	case n_IF_STATEMENT:
		cond := n.children[0]
		body := n.children[1]
		alt := n.children[2]
		res := "if(" + printAst(cond) + ")" + printAst(body)
		if alt.t != n_EMPTY {
			res += "else " + printAst(alt)
		}
		return res

	case n_SEQUENCE_EXPRESSION:
		res := ""
		for i, item := range n.children {
			res += printAst(item)
			if i < len(n.children)-1 {
				res += ","
			}
		}
		return res

	case n_IMPORT_PATH:
		return n.value

	case n_IMPORT_ALL:
		if n.value != "" {
			return "*as " + n.value
		}
		return ""

	case n_IMPORT_VAR:
		name := n.children[0]
		alias := n.children[1]
		res := printAst(name)
		if alias.value != "" {
			res += " as " + printAst(alias)
		}
		return res

	case n_IMPORT_STATEMENT:
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

	case n_EXPORT_STATEMENT:
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

	case n_EXPORT_PATH:
		if n.value != "" {
			return " from" + n.value
		}
		return ""

	case n_EXPORT_VAR:
		result := printAst(n.children[0]) // name
		if len(n.children) > 1 {
			result += " as " + printAst(n.children[1]) // alias
		}
		return result

	case n_DECLARATION_STATEMENT:
		return printAst(n.children[0]) + ";"

	case n_RETURN_STATEMENT:
		res := "return"
		if len(n.children) > 0 {
			res += " " + printAst(n.children[0])
		}
		return res

	case n_CONDITIONAL_EXPRESSION:
		test := n.children[0]
		conseq := n.children[1]
		alt := n.children[2]
		return printAst(test) + "?" + printAst(conseq) + ":" + printAst(alt)

	case n_CLASS_STATEMENT, n_CLASS_EXPRESSION:
		name := n.children[0]
		extends := n.children[1]
		body := n.children[2]
		res := "class"
		if name.t != n_EMPTY {
			res += " " + printAst(name)
		}
		if extends.t != n_EMPTY {
			res += " extends " + printAst(extends)
		}
		res += printAst(body)
		return res

	case n_CLASS_PROPERTY:
		key := n.children[0]
		value := n.children[1]
		res := printAst(key)
		if value.t != n_EMPTY {
			res += "=" + printAst(value)
		}
		res += ";"
		return res

	case n_CLASS_STATIC_PROPERTY:
		return "static " + printAst(n.children[0])
	}

	return ""
}
