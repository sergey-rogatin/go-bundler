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
	// "of":         tOF,
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
	// tOF

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
	// "tOF",
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
}

func (t tokenType) String() string {
	return tokenTypeNames[int32(t)]
}

type grammarType int

var grammarTypeToString = []string{
	"g_UNDEFINED",
	"g_INVALID_GRAMMAR",
	"g_PROGRAM_STATEMENT",
	"g_RETURN_STATEMENT",
	"g_EXPRESSION_STATEMENT",
	"g_BLOCK_STATEMENT",
	"g_IMPORT_STATEMENT",
	"g_IMPORT_VARS",
	"g_IMPORT_VAR",
	"g_IMPORT_NAME",
	"g_IMPORT_ALIAS",
	"g_IMPORT_PATH",
	"g_IMPORT_ALL",
	"g_EXPRESSION",
	"g_FUNCTION_CALL",
	"g_CONSTRUCTOR_CALL",
	"g_VARIABLE_NAME",
	"g_OBJECT_LITERAL_PROPERTY_NAME",
	"g_CALCULATED_PROPERTY_NAME",
	"g_VALID_PROPERTY_NAME",
	"g_LITERAL",
	"g_DECLARATOR",
	"g_FUNCTION_PARAMETER",
	"g_FUNCTION_PARAMETERS",
	"g_FUNCTION_EXPRESSION",
	"g_LAMBDA_EXPRESSION",
	"g_MEMBER_FUNCTION",
	"g_ARRAY_LITERAL",
	"g_OBJECT_LITERAL",
	"g_OBJECT_PROPERTY",
	"g_OBJECT_PATTERN",
	"g_ASSIGNMENT_PATTERN",
	"g_ARRAY_PATTERN",
	"g_FUNCTION_DECLARATION",
	"g_EMPTY_STATEMENT",
	"g_EMPTY_EXPRESSION",
	"g_FOR_STATEMENT",
	"g_FOR_IN_STATEMENT",
	"g_FOR_OF_STATEMENT",
	"g_IF_STATEMENT",
	"g_DECLARATION_STATEMENT",
	"g_DECLARATION_EXPRESSION",
	"g_WHILE_STATEMENT",
	"g_DO_WHILE_STATEMENT",
	"g_SWITCH_CASES",
	"g_SWITCH_CASE_STATEMENTS",
	"g_SWITCH_STATEMENT",
	"g_SWITCH_CASE",
	"g_SWITCH_CASE_TEST",
	"g_SWITCH_DEFAULT",
	"g_EXPORT_STATEMENT",
	"g_EXPORT_VARS",
	"g_EXPORT_VAR",
	"g_EXPORT_NAME",
	"g_EXPORT_ALIAS",
	"g_EXPORT_DECLARATION",
	"g_EXPORT_PATH",
	"g_EXPORT_ALL",
	"g_NUMBER_LITERAL",
	"g_STRING_LITERAL",
	"g_NULL",
	"g_BOOL_LITERAL",
	"g_THIS",
	"g_UNARY_PREFIX_EXPRESSION",
	"g_UNARY_POSTFIX_EXPRESSION",
	"g_MEMBER_EXPRESSION",
	"g_FUNCTION_ARGS",
	"g_PARENS_EXPRESSION",
	"g_SEQUENCE_EXPRESSION",
	"g_BREAK_STATEMENT",
	"g_CONTINUE_STATEMENT",
	"g_DEBUGGER_STATEMENT",
	"g_MULTISTATEMENT",
	"g_HEX_LITERAL",
	"g_REGEXP_LITERAL",
	"g_TRY_CATCH_STATEMENT",
	"g_THROW_STATEMENT",
	"g_CONDITIONAL_EXPRESSION",
	"g_MARKER",
	"g_REST_EXPRESSION",
	"g_SPREAD_EXPRESSION",
	"g_CLASS_EXPRESSION",
	"g_CLASS_MEMBER",
	"g_CLASS_METHOD",
	"g_CLASS_STATIC_PROP",
	"g_CLASS_PROPERTY",
	"g_CLASS_STATEMENT",
}

func (g grammarType) String() string {
	return grammarTypeToString[int(g)]
}

const (
	g_UNDEFINED grammarType = iota
	g_INVALID_GRAMMAR
	g_PROGRAM_STATEMENT
	g_RETURN_STATEMENT
	g_EXPRESSION_STATEMENT
	g_BLOCK_STATEMENT
	g_IMPORT_STATEMENT
	g_IMPORT_VARS
	g_IMPORT_VAR
	g_IMPORT_NAME
	g_IMPORT_ALIAS
	g_IMPORT_PATH
	g_IMPORT_ALL
	g_EXPRESSION
	g_FUNCTION_CALL
	g_CONSTRUCTOR_CALL
	g_VARIABLE_NAME
	g_OBJECT_LITERAL_PROPERTY_NAME
	g_CALCULATED_PROPERTY_NAME
	g_VALID_PROPERTY_NAME
	g_LITERAL
	g_DECLARATOR
	g_FUNCTION_PARAMETER
	g_FUNCTION_PARAMETERS
	g_FUNCTION_EXPRESSION
	g_LAMBDA_EXPRESSION
	g_MEMBER_FUNCTION
	g_ARRAY_LITERAL
	g_OBJECT_LITERAL
	g_OBJECT_PROPERTY
	g_OBJECT_PATTERN
	g_ASSIGNMENT_PATTERN
	g_ARRAY_PATTERN
	g_FUNCTION_DECLARATION
	g_EMPTY_STATEMENT
	g_EMPTY_EXPRESSION
	g_FOR_STATEMENT
	g_FOR_IN_STATEMENT
	g_FOR_OF_STATEMENT
	g_IF_STATEMENT
	g_DECLARATION_STATEMENT
	g_DECLARATION_EXPRESSION
	g_WHILE_STATEMENT
	g_DO_WHILE_STATEMENT
	g_SWITCH_CASES
	g_SWITCH_CASE_STATEMENTS
	g_SWITCH_STATEMENT
	g_SWITCH_CASE
	g_SWITCH_CASE_TEST
	g_SWITCH_DEFAULT
	g_EXPORT_STATEMENT
	g_EXPORT_VARS
	g_EXPORT_VAR
	g_EXPORT_NAME
	g_EXPORT_ALIAS
	g_EXPORT_DECLARATION
	g_EXPORT_PATH
	g_EXPORT_ALL
	g_NUMBER_LITERAL
	g_STRING_LITERAL
	g_NULL
	g_BOOL_LITERAL
	g_THIS
	g_UNARY_PREFIX_EXPRESSION
	g_UNARY_POSTFIX_EXPRESSION
	g_MEMBER_EXPRESSION
	g_FUNCTION_ARGS
	g_PARENS_EXPRESSION
	g_SEQUENCE_EXPRESSION
	g_BREAK_STATEMENT
	g_CONTINUE_STATEMENT
	g_DEBUGGER_STATEMENT
	g_MULTISTATEMENT
	g_HEX_LITERAL
	g_REGEXP_LITERAL
	g_TRY_CATCH_STATEMENT
	g_THROW_STATEMENT
	g_CONDITIONAL_EXPRESSION
	g_MARKER
	g_REST_EXPRESSION
	g_SPREAD_EXPRESSION
	g_CLASS_EXPRESSION
	g_CLASS_MEMBER
	g_CLASS_METHOD
	g_CLASS_STATIC_PROP
	g_CLASS_PROPERTY
	g_CLASS_STATEMENT
)

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
