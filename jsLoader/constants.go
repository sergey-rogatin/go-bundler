package jsLoader

type tokenType int32

var keywords = map[string]tokenType{
	"break":      tBREAK,
	"case":       tCASE,
	"catch":      tCATCH,
	"class":      tCLASS,
	"const":      tCONST,
	"continue":   tCONTINUE,
	"debugger":   tDEBUGGER,
	"default":    tDEFAULT,
	"delete":     tDELETE,
	"do":         tDO,
	"else":       tELSE,
	"export":     tEXPORT,
	"extends":    tEXTENDS,
	"finally":    tFINALLY,
	"for":        tFOR,
	"function":   tFUNCTION,
	"if":         tIF,
	"import":     tIMPORT,
	"in":         tIN,
	"instanceof": tINSTANCEOF,
	"new":        tNEW,
	"return":     tRETURN,
	// "super":      tSUPER,
	"switch":     tSWITCH,
	"this":       tTHIS,
	"throw":      tTHROW,
	"try":        tTRY,
	"typeof":     tTYPEOF,
	"var":        tVAR,
	"void":       tVOID,
	"while":      tWHILE,
	"with":       tWITH,
	"yield":      tYIELD,
	"yield*":     tYIELD_STAR,
	"enum":       tENUM,
	"implements": tIMPLEMENTS,
	"interface":  tINTERFACE,
	"let":        tLET,
	"package":    tPACKAGE,
	"private":    tPRIVATE,
	"protected":  tPROTECTED,
	"public":     tPUBLIC,
	"static":     tSTATIC,
	"await":      tAWAIT,
	"async":      tASYNC,
	"null":       tNULL,
	"true":       tTRUE,
	"false":      tFALSE,
	// "as":         tAS,
	// "from":       tFROM,
	"of": tOF,
}

var operators = map[string]tokenType{
	"~":    tBITWISE_NOT,
	"&":    tBITWISE_AND,
	"|":    tBITWISE_OR,
	"^":    tBITWISE_XOR,
	"*":    tMULT,
	"**":   tEXP,
	"-":    tMINUS,
	"+":    tPLUS,
	"/":    tDIV,
	"%":    tMOD,
	"+=":   tPLUS_ASSIGN,
	"*=":   tMULT_ASSIGN,
	"-=":   tMINUS_ASSIGN,
	"/=":   tDIV_ASSIGN,
	"%=":   tMOD_ASSIGN,
	"<<=":  tBITWISE_SHIFT_LEFT_ASSIGN,
	">>=":  tBITWISE_SHIFT_RIGHT_ASSIGN,
	">>>=": tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN,
	"&=":   tBITWISE_AND_ASSIGN,
	"^=":   tBITWISE_XOR_ASSIGN,
	"|=":   tBITWISE_OR_ASSIGN,
	"++":   tINC,
	"--":   tDEC,
	"!":    tNOT,
	"&&":   tAND,
	"||":   tOR,
	">":    tGREATER,
	"<":    tLESS,
	">=":   tGREATER_EQUALS,
	"<=":   tLESS_EQUALS,
	"==":   tEQUALS,
	"===":  tEQUALS_STRICT,
	"!=":   tNOT_EQUALS,
	"!==":  tNOT_EQUALS_STRICT,
	"(":    tPAREN_LEFT,
	")":    tPAREN_RIGHT,
	"[":    tBRACKET_LEFT,
	"]":    tBRACKET_RIGHT,
	"{":    tCURLY_LEFT,
	"}":    tCURLY_RIGHT,
	".":    tDOT,
	",":    tCOMMA,
	"?":    tQUESTION,
	";":    tSEMI,
	":":    tCOLON,
	"=":    tASSIGN,
	"=>":   tLAMBDA,
	">>":   tBITWISE_SHIFT_RIGHT,
	"<<":   tBITWISE_SHIFT_LEFT,
	">>>":  tBITWISE_SHIFT_RIGHT_ZERO,
	"...":  tSPREAD,
	"\\":   tESCAPE,
}

const (
	tUNDEFINED tokenType = iota
	tEND_OF_INPUT

	tNEWLINE

	tBREAK
	tCASE
	tCATCH
	tCLASS
	tCONST
	tCONTINUE
	tDEBUGGER
	tDEFAULT
	tDELETE
	tDO
	tELSE
	tEXPORT
	tEXTENDS
	tFINALLY
	tFOR
	tFUNCTION
	tIF
	tIMPORT
	tIN
	tINSTANCEOF
	tNEW
	tRETURN
	tSUPER
	tSWITCH
	tTHIS
	tTHROW
	tTRY
	tTYPEOF
	tVAR
	tVOID
	tWHILE
	tWITH
	tYIELD
	tYIELD_STAR
	tENUM
	tIMPLEMENTS
	tINTERFACE
	tLET
	tPACKAGE
	tPRIVATE
	tPROTECTED
	tPUBLIC
	tSTATIC
	tAWAIT
	tASYNC
	tNULL
	tTRUE
	tFALSE
	// tAS
	// tFROM
	tOF

	tBITWISE_NOT
	tBITWISE_AND
	tBITWISE_OR
	tBITWISE_XOR
	tMULT
	tEXP
	tMINUS
	tPLUS
	tDIV
	tMOD
	tPLUS_ASSIGN
	tMULT_ASSIGN
	tMINUS_ASSIGN
	tDIV_ASSIGN
	tMOD_ASSIGN
	tBITWISE_SHIFT_LEFT_ASSIGN
	tBITWISE_SHIFT_RIGHT_ASSIGN
	tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN
	tBITWISE_AND_ASSIGN
	tBITWISE_XOR_ASSIGN
	tBITWISE_OR_ASSIGN
	tINC
	tDEC
	tNOT
	tAND
	tOR
	tGREATER
	tLESS
	tGREATER_EQUALS
	tLESS_EQUALS
	tEQUALS
	tEQUALS_STRICT
	tNOT_EQUALS
	tNOT_EQUALS_STRICT
	tPAREN_LEFT
	tPAREN_RIGHT
	tBRACKET_LEFT
	tBRACKET_RIGHT
	tCURLY_LEFT
	tCURLY_RIGHT
	tDOT
	tCOMMA
	tQUESTION
	tSEMI
	tCOLON
	tASSIGN
	tLAMBDA
	tBITWISE_SHIFT_RIGHT
	tBITWISE_SHIFT_LEFT
	tBITWISE_SHIFT_RIGHT_ZERO
	tSPREAD

	tNUMBER
	tSTRING
	tNAME
	tHEX
	tREGEXP
	tESCAPE
	tSPACE
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
}

func (t tokenType) String() string {
	return tokenTypeNames[int32(t)]
}

type operatorInfo struct {
	precedence         int
	isRightAssociative bool
}

var operatorTable = map[tokenType]operatorInfo{
	tASSIGN:                          operatorInfo{3, true},
	tPLUS_ASSIGN:                     operatorInfo{3, true},
	tMINUS_ASSIGN:                    operatorInfo{3, true},
	tMULT_ASSIGN:                     operatorInfo{3, true},
	tDIV_ASSIGN:                      operatorInfo{3, true},
	tMOD_ASSIGN:                      operatorInfo{3, true},
	tBITWISE_SHIFT_LEFT_ASSIGN:       operatorInfo{3, true},
	tBITWISE_SHIFT_RIGHT_ASSIGN:      operatorInfo{3, true},
	tBITWISE_SHIFT_RIGHT_ZERO_ASSIGN: operatorInfo{3, true},
	tBITWISE_AND_ASSIGN:              operatorInfo{3, true},
	tBITWISE_OR_ASSIGN:               operatorInfo{3, true},
	tBITWISE_XOR_ASSIGN:              operatorInfo{3, true},
	tOR:                              operatorInfo{5, false},
	tAND:                             operatorInfo{6, false},
	tBITWISE_OR:                      operatorInfo{7, false},
	tBITWISE_XOR:                     operatorInfo{8, false},
	tBITWISE_AND:                     operatorInfo{9, false},
	tEQUALS:                          operatorInfo{10, false},
	tNOT_EQUALS:                      operatorInfo{10, false},
	tEQUALS_STRICT:                   operatorInfo{10, false},
	tNOT_EQUALS_STRICT:               operatorInfo{10, false},
	tLESS:                            operatorInfo{11, false},
	tLESS_EQUALS:                     operatorInfo{11, false},
	tGREATER:                         operatorInfo{11, false},
	tGREATER_EQUALS:                  operatorInfo{11, false},
	tINSTANCEOF:                      operatorInfo{11, false},
	tIN:                              operatorInfo{11, false},
	tBITWISE_SHIFT_LEFT:              operatorInfo{12, false},
	tBITWISE_SHIFT_RIGHT:             operatorInfo{12, false},
	tBITWISE_SHIFT_RIGHT_ZERO:        operatorInfo{12, false},
	tPLUS:  operatorInfo{13, false},
	tMINUS: operatorInfo{13, false},
	tMULT:  operatorInfo{14, false},
	tDIV:   operatorInfo{14, false},
	tMOD:   operatorInfo{14, false},
	tEXP:   operatorInfo{15, true},
}

type nodeType int

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
)
