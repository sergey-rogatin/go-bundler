package main

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
	"super":      tSUPER,
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
	"as":         tAS,
	"from":       tFROM,
	"of":         tOF,
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
	tAS
	tFROM
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
	"tAS",
	"tFROM",
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
}

func (t tokenType) String() string {
	return tokenTypeNames[int32(t)]
}
