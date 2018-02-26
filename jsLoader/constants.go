package jsLoader

type tokenType int32

var keywords = map[string]tokenType{
	"break":      t_BREAK,
	"case":       t_CASE,
	"catch":      t_CATCH,
	"class":      t_CLASS,
	"const":      t_CONST,
	"continue":   t_CONTINUE,
	"debugger":   t_DEBUGGER,
	"default":    t_DEFAULT,
	"delete":     t_DELETE,
	"do":         t_DO,
	"else":       t_ELSE,
	"export":     t_EXPORT,
	"extends":    t_EXTENDS,
	"finally":    t_FINALLY,
	"for":        t_FOR,
	"function":   t_FUNCTION,
	"if":         t_IF,
	"import":     t_IMPORT,
	"in":         t_IN,
	"instanceof": t_INSTANCEOF,
	"new":        t_NEW,
	"return":     t_RETURN,
	// "super":      tSUPER,
	"switch":     t_SWITCH,
	"this":       t_THIS,
	"throw":      t_THROW,
	"try":        t_TRY,
	"typeof":     t_TYPEOF,
	"var":        t_VAR,
	"void":       t_VOID,
	"while":      t_WHILE,
	"with":       t_WITH,
	"yield":      t_YIELD,
	"yield*":     t_YIELD_STAR,
	"enum":       t_ENUM,
	"implements": t_IMPLEMENTS,
	"interface":  t_INTERFACE,
	"let":        t_LET,
	"package":    t_PACKAGE,
	"private":    t_PRIVATE,
	"protected":  t_PROTECTED,
	"public":     t_PUBLIC,
	"static":     t_STATIC,
	"await":      t_AWAIT,
	"async":      t_ASYNC,
	"null":       t_NULL,
	"true":       t_TRUE,
	"false":      t_FALSE,
	// "as":         tAS,
	// "from":       tFROM,
	"of": t_OF,
}

var operators = map[string]tokenType{
	"~":    t_BITWISE_NOT,
	"&":    t_BITWISE_AND,
	"|":    t_BITWISE_OR,
	"^":    t_BITWISE_XOR,
	"*":    t_MULT,
	"**":   t_EXP,
	"-":    t_MINUS,
	"+":    t_PLUS,
	"/":    t_DIV,
	"%":    t_MOD,
	"+=":   t_PLUS_ASSIGN,
	"*=":   t_MULT_ASSIGN,
	"-=":   t_MINUS_ASSIGN,
	"/=":   t_DIV_ASSIGN,
	"%=":   t_MOD_ASSIGN,
	"<<=":  t_BITWISE_SHIFT_LEFT_ASSIGN,
	">>=":  t_BITWISE_SHIFT_RIGHT_ASSIGN,
	">>>=": t_BITWISE_SHIFT_RIGHT_ZERO_ASSIGN,
	"&=":   t_BITWISE_AND_ASSIGN,
	"^=":   t_BITWISE_XOR_ASSIGN,
	"|=":   t_BITWISE_OR_ASSIGN,
	"++":   t_INC,
	"--":   t_DEC,
	"!":    t_NOT,
	"&&":   t_AND,
	"||":   t_OR,
	">":    t_GREATER,
	"<":    t_LESS,
	">=":   t_GREATER_EQUALS,
	"<=":   t_LESS_EQUALS,
	"==":   t_EQUALS,
	"===":  t_EQUALS_STRICT,
	"!=":   t_NOT_EQUALS,
	"!==":  t_NOT_EQUALS_STRICT,
	"(":    t_PAREN_LEFT,
	")":    t_PAREN_RIGHT,
	"[":    t_BRACKET_LEFT,
	"]":    t_BRACKET_RIGHT,
	"{":    t_CURLY_LEFT,
	"}":    t_CURLY_RIGHT,
	".":    t_DOT,
	",":    t_COMMA,
	"?":    t_QUESTION,
	";":    t_SEMI,
	":":    t_COLON,
	"=":    t_ASSIGN,
	"=>":   t_LAMBDA,
	">>":   t_BITWISE_SHIFT_RIGHT,
	"<<":   t_BITWISE_SHIFT_LEFT,
	">>>":  t_BITWISE_SHIFT_RIGHT_ZERO,
	"...":  t_SPREAD,
	"\\":   t_ESCAPE,
}

const (
	t_UNDEFINED tokenType = iota
	t_ANY
	t_END_OF_INPUT

	t_NEWLINE

	t_BREAK
	t_CASE
	t_CATCH
	t_CLASS
	t_CONST
	t_CONTINUE
	t_DEBUGGER
	t_DEFAULT
	t_DELETE
	t_DO
	t_ELSE
	t_EXPORT
	t_EXTENDS
	t_FINALLY
	t_FOR
	t_FUNCTION
	t_IF
	t_IMPORT
	t_IN
	t_INSTANCEOF
	t_NEW
	t_RETURN
	t_SUPER
	t_SWITCH
	t_THIS
	t_THROW
	t_TRY
	t_TYPEOF
	t_VAR
	t_VOID
	t_WHILE
	t_WITH
	t_YIELD
	t_YIELD_STAR
	t_ENUM
	t_IMPLEMENTS
	t_INTERFACE
	t_LET
	t_PACKAGE
	t_PRIVATE
	t_PROTECTED
	t_PUBLIC
	t_STATIC
	t_AWAIT
	t_ASYNC
	t_NULL
	t_TRUE
	t_FALSE
	// tAS
	// tFROM
	t_OF

	t_BITWISE_NOT
	t_BITWISE_AND
	t_BITWISE_OR
	t_BITWISE_XOR
	t_MULT
	t_EXP
	t_MINUS
	t_PLUS
	t_DIV
	t_MOD
	t_PLUS_ASSIGN
	t_MULT_ASSIGN
	t_MINUS_ASSIGN
	t_DIV_ASSIGN
	t_MOD_ASSIGN
	t_BITWISE_SHIFT_LEFT_ASSIGN
	t_BITWISE_SHIFT_RIGHT_ASSIGN
	t_BITWISE_SHIFT_RIGHT_ZERO_ASSIGN
	t_BITWISE_AND_ASSIGN
	t_BITWISE_XOR_ASSIGN
	t_BITWISE_OR_ASSIGN
	t_INC
	t_DEC
	t_NOT
	t_AND
	t_OR
	t_GREATER
	t_LESS
	t_GREATER_EQUALS
	t_LESS_EQUALS
	t_EQUALS
	t_EQUALS_STRICT
	t_NOT_EQUALS
	t_NOT_EQUALS_STRICT
	t_PAREN_LEFT
	t_PAREN_RIGHT
	t_BRACKET_LEFT
	t_BRACKET_RIGHT
	t_CURLY_LEFT
	t_CURLY_RIGHT
	t_DOT
	t_COMMA
	t_QUESTION
	t_SEMI
	t_COLON
	t_ASSIGN
	t_LAMBDA
	t_BITWISE_SHIFT_RIGHT
	t_BITWISE_SHIFT_LEFT
	t_BITWISE_SHIFT_RIGHT_ZERO
	t_SPREAD

	t_NUMBER
	t_STRING_QUOTE
	t_NAME
	t_HEX
	t_REGEXP
	t_ESCAPE
	t_SPACE
	t_TEMPLATE_LITERAL_QUOTE
)

var tokenTypeNames = []string{
	"tUNDEFINED",
	"tEND_OF_INPUT",
	"tNEWLINE",
	"tBREAK",
	"tCASE",
	"tCATCH",
	"tCLASS",
	"tCONST",
	"tCONTINUE",
	"tDEBUGGER",
	"tDEFAULT",
	"tDELETE",
	"tDO",
	"tELSE",
	"tEXPORT",
	"tEXTENDS",
	"tFINALLY",
	"tFOR",
	"tFUNCTION",
	"tIF",
	"tIMPORT",
	"tIN",
	"tINSTANCEOF",
	"tNEW",
	"tRETURN",
	"tSUPER",
	"tSWITCH",
	"tTHIS",
	"tTHROW",
	"tTRY",
	"tTYPEOF",
	"tVAR",
	"tVOID",
	"tWHILE",
	"tWITH",
	"tYIELD",
	"tYIELD_STAR",
	"tENUM",
	"tIMPLEMENTS",
	"tINTERFACE",
	"tLET",
	"tPACKAGE",
	"tPRIVATE",
	"tPROTECTED",
	"tPUBLIC",
	"tSTATIC",
	"tAWAIT",
	"tASYNC",
	"tNULL",
	"tTRUE",
	"tFALSE",
	// "tAS",
	// "tFROM",
	"tOF",
	"tBITWISE_NOT",
	"tBITWISE_AND",
	"tBITWISE_OR",
	"tBITWISE_XOR",
	"tMULT",
	"tEXP",
	"tMINUS",
	"tPLUS",
	"tDIV",
	"tMOD",
	"tPLUS_ASSIGN",
	"tMULT_ASSIGN",
	"tMINUS_ASSIGN",
	"tDIV_ASSIGN",
	"tMOD_ASSIGN",
	"tBITWISE_SHIFT_LEFT_ASSIGN",
	"tBITWISE_SHIFT_RIGHT_ASSIGN",
	"tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN",
	"tBITWISE_AND_ASSIGN",
	"tBITWISE_XOR_ASSIGN",
	"tBITWISE_OR_ASSIGN",
	"tINC",
	"tDEC",
	"tNOT",
	"tAND",
	"tOR",
	"tGREATER",
	"tLESS",
	"tGREATER_EQUALS",
	"tLESS_EQUALS",
	"tEQUALS",
	"tEQUALS_STRICT",
	"tNOT_EQUALS",
	"tNOT_EQUALS_STRICT",
	"tPAREN_LEFT",
	"tPAREN_RIGHT",
	"tBRACKET_LEFT",
	"tBRACKET_RIGHT",
	"tCURLY_LEFT",
	"tCURLY_RIGHT",
	"tDOT",
	"tCOMMA",
	"tQUESTION",
	"tSEMI",
	"tCOLON",
	"tASSIGN",
	"tLAMBDA",
	"tBITWISE_SHIFT_RIGHT",
	"tBITWISE_SHIFT_LEFT",
	"tBITWISE_SHIFT_RIGHT_ZERO",
	"tSPREAD",
	"tNUMBER",
	"tSTRING",
	"tNAME",
	"tHEX",
	"tREGEXP",
	"tESCAPE",
	"tSPACE",
	"tTEMPLATE_LITERAL",
}

func (t tokenType) String() string {
	return tokenTypeNames[int32(t)]
}

type operatorInfo struct {
	precedence         int
	isRightAssociative bool
}

var operatorTable = map[tokenType]operatorInfo{
	t_ASSIGN:                          operatorInfo{3, true},
	t_PLUS_ASSIGN:                     operatorInfo{3, true},
	t_MINUS_ASSIGN:                    operatorInfo{3, true},
	t_MULT_ASSIGN:                     operatorInfo{3, true},
	t_DIV_ASSIGN:                      operatorInfo{3, true},
	t_MOD_ASSIGN:                      operatorInfo{3, true},
	t_BITWISE_SHIFT_LEFT_ASSIGN:       operatorInfo{3, true},
	t_BITWISE_SHIFT_RIGHT_ASSIGN:      operatorInfo{3, true},
	t_BITWISE_SHIFT_RIGHT_ZERO_ASSIGN: operatorInfo{3, true},
	t_BITWISE_AND_ASSIGN:              operatorInfo{3, true},
	t_BITWISE_OR_ASSIGN:               operatorInfo{3, true},
	t_BITWISE_XOR_ASSIGN:              operatorInfo{3, true},
	t_OR:                              operatorInfo{5, false},
	t_AND:                             operatorInfo{6, false},
	t_BITWISE_OR:                      operatorInfo{7, false},
	t_BITWISE_XOR:                     operatorInfo{8, false},
	t_BITWISE_AND:                     operatorInfo{9, false},
	t_EQUALS:                          operatorInfo{10, false},
	t_NOT_EQUALS:                      operatorInfo{10, false},
	t_EQUALS_STRICT:                   operatorInfo{10, false},
	t_NOT_EQUALS_STRICT:               operatorInfo{10, false},
	t_LESS:                            operatorInfo{11, false},
	t_LESS_EQUALS:                     operatorInfo{11, false},
	t_GREATER:                         operatorInfo{11, false},
	t_GREATER_EQUALS:                  operatorInfo{11, false},
	t_INSTANCEOF:                      operatorInfo{11, false},
	t_IN:                              operatorInfo{11, false},
	t_BITWISE_SHIFT_LEFT:              operatorInfo{12, false},
	t_BITWISE_SHIFT_RIGHT:             operatorInfo{12, false},
	t_BITWISE_SHIFT_RIGHT_ZERO:        operatorInfo{12, false},
	t_PLUS:  operatorInfo{13, false},
	t_MINUS: operatorInfo{13, false},
	t_MULT:  operatorInfo{14, false},
	t_DIV:   operatorInfo{14, false},
	t_MOD:   operatorInfo{14, false},
	t_EXP:   operatorInfo{15, true},
}

type nodeType int32

var nodeToString = []string{
	"n_EMPTY",
	"n_INVALID",
	"n_PROGRAM",
	"n_IDENTIFIER",
	"n_EXPRESSION_STATEMENT",
	"n_FUNCTION_DECLARATION",
	"n_BLOCK_STATEMENT",
	"n_FUNCTION_EXPRESSION",
	"n_FUNCTION_PARAMETERS",
	"n_VARIABLE_DECLARATION",
	"n_DECLARATOR",
	"n_STRING_LITERAL",
	"n_DOUBLE_QUOTE_STRING_LITERAL",
	"n_NUMBER_LITERAL",
	"n_BOOL_LITERAL",
	"n_UNDEFINED",
	"n_NULL",
	"n_OBJECT_LITERAL",
	"n_ARRAY_LITERAL",
	"n_REGEXP_LITERAL",
	"n_REGEXP_BODY",
	"n_REGEXP_FLAGS",
	"n_ASSIGNMENT_PATTERN",
	"n_OBJECT_PATTERN",
	"n_ARRAY_PATTERN",
	"n_OBJECT_PROPERTY",
	"n_REST_ELEMENT",
	"n_LAMBDA_EXPRESSION",
	"n_SEQUENCE_EXPRESSION",
	"n_ASSIGNMENT_EXPRESSION",
	"n_CONDITIONAL_EXPRESSION",
	"n_BINARY_EXPRESSION",
	"n_PREFIX_UNARY_EXPRESSION",
	"n_POSTFIX_UNARY_EXPRESSION",
	"n_MEMBER_EXPRESSION",
	"n_PROPERTY_NAME",
	"n_CALCULATED_PROPERTY_NAME",
	"n_PAREN_EXPRESSION",
	"n_FUNCTION_CALL",
	"n_FUNCTION_ARGS",
	"n_CONSTRUCTOR_CALL",
	"n_IMPORT_STATEMENT",
	"n_IMPORT_ALL",
	"n_IMPORT_NAME",
	"n_IMPORT_ALIAS",
	"n_IMPORT_VAR",
	"n_IMPORT_PATH",
	"n_EXPORT_STATEMENT",
	"n_EXPORT_VARS",
	"n_EXPORT_VAR",
	"n_EXPORT_NAME",
	"n_EXPORT_ALIAS",
	"n_EXPORT_PATH",
	"n_EXPORT_DECLARATION",
	"n_EXPORT_ALL",
	"n_IF_STATEMENT",
	"n_WHILE_STATEMENT",
	"n_DO_WHILE_STATEMENT",
	"n_FOR_STATEMENT",
	"n_FOR_IN_STATEMENT",
	"n_FOR_OF_STATEMENT",
	"n_EMPTY_EXPRESSION",
	"n_EMPTY_STATEMENT",
	"n_OBJECT_METHOD",
	"n_OBJECT_KEY",
	"n_NON_IDENTIFIER_OBJECT_KEY",
	"n_CALCULATED_OBJECT_KEY",
	"n_SPREAD_ELEMENT",
	"n_CONTROL_WORD",
	"n_RETURN_STATEMENT",
	"n_DECLARATION_STATEMENT",
	"n_CLASS_STATEMENT",
	"n_CLASS_EXPRESSION",
	"n_CLASS_BODY",
	"n_CLASS_METHOD",
	"n_CLASS_PROPERTY",
	"n_CLASS_STATIC_PROPERTY",
	"n_LABEL",
	"n_MULTISTATEMENT",
	"n_TRY_CATCH_STATEMENT",
	"n_SWITCH_STATEMENT",
	"n_SWITCH_CASE",
	"n_SWITCH_DEFAULT_CASE",
	"n_TAGGED_LITERAL",
	"n_TEMPLATE_LITERAL",
	"n_THROW_STATEMENT",
	"n_CONTROL_STATEMENT",
	"n_YIELD_EXPRESSION",
	"n_AWAIT_EXPRESSION",
	"n_ASYNC_FUNCTION",
	"n_WITH_STATEMENT",
	"n_LINE_COMMENT",
	"n_BLOCK_COMMENT",
}

func (n nodeType) String() string {
	return nodeToString[int(n)]
}

const (
	n_EMPTY nodeType = iota
	n_INVALID
	n_PROGRAM
	n_IDENTIFIER
	n_EXPRESSION_STATEMENT
	n_FUNCTION_DECLARATION
	n_BLOCK_STATEMENT

	n_FUNCTION_EXPRESSION
	n_FUNCTION_PARAMETERS

	n_VARIABLE_DECLARATION
	n_DECLARATOR

	n_STRING_LITERAL
	n_DOUBLE_QUOTE_STRING_LITERAL
	n_NUMBER_LITERAL
	n_BOOL_LITERAL
	n_UNDEFINED
	n_NULL
	n_OBJECT_LITERAL
	n_ARRAY_LITERAL
	n_REGEXP_LITERAL

	n_REGEXP_BODY
	n_REGEXP_FLAGS

	n_ASSIGNMENT_PATTERN
	n_OBJECT_PATTERN
	n_ARRAY_PATTERN
	n_OBJECT_PROPERTY

	n_REST_ELEMENT

	n_LAMBDA_EXPRESSION
	n_SEQUENCE_EXPRESSION
	n_ASSIGNMENT_EXPRESSION
	n_CONDITIONAL_EXPRESSION
	n_BINARY_EXPRESSION
	n_PREFIX_UNARY_EXPRESSION
	n_POSTFIX_UNARY_EXPRESSION
	n_MEMBER_EXPRESSION
	n_PROPERTY_NAME
	n_CALCULATED_PROPERTY_NAME
	n_PAREN_EXPRESSION

	n_FUNCTION_CALL
	n_FUNCTION_ARGS

	n_CONSTRUCTOR_CALL

	n_IMPORT_STATEMENT
	n_IMPORT_ALL
	n_IMPORT_NAME
	n_IMPORT_ALIAS
	n_IMPORT_VAR
	n_IMPORT_PATH

	n_EXPORT_STATEMENT
	n_EXPORT_VARS
	n_EXPORT_VAR
	n_EXPORT_NAME
	n_EXPORT_ALIAS
	n_EXPORT_PATH
	n_EXPORT_DECLARATION
	n_EXPORT_ALL

	n_IF_STATEMENT
	n_WHILE_STATEMENT
	n_DO_WHILE_STATEMENT
	n_FOR_STATEMENT
	n_FOR_IN_STATEMENT
	n_FOR_OF_STATEMENT

	n_EMPTY_EXPRESSION
	n_EMPTY_STATEMENT

	n_OBJECT_METHOD
	n_OBJECT_KEY
	n_NON_IDENTIFIER_OBJECT_KEY
	n_CALCULATED_OBJECT_KEY
	n_SPREAD_ELEMENT
	n_CONTROL_WORD
	n_RETURN_STATEMENT
	n_DECLARATION_STATEMENT

	n_CLASS_STATEMENT
	n_CLASS_EXPRESSION
	n_CLASS_BODY
	n_CLASS_METHOD
	n_CLASS_PROPERTY
	n_CLASS_STATIC_PROPERTY

	n_LABEL
	n_MULTISTATEMENT
	n_TRY_CATCH_STATEMENT

	n_SWITCH_STATEMENT
	n_SWITCH_CASE
	n_SWITCH_DEFAULT_CASE

	n_TAGGED_LITERAL
	n_TEMPLATE_LITERAL
	n_THROW_STATEMENT
	n_CONTROL_STATEMENT
	n_YIELD_EXPRESSION
	n_AWAIT_EXPRESSION
	n_ASYNC_FUNCTION
	n_WITH_STATEMENT
	n_LINE_COMMENT
	n_BLOCK_COMMENT
	n_OBJECT_MEMBER
)
