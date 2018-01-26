package main

type tokenType int32

var keywords = map[string]tokenType{
	"typeof":   tTYPEOF,
	"for":      tFOR,
	"as":       tAS,
	"return":   tRETURN,
	"function": tFUNCTION,
	"from":     tFROM,
	"default":  tDEFAULT,
	"export":   tEXPORT,
	"import":   tIMPORT,
	"var":      tVAR,
	"const":    tCONST,
	"let":      tLET,
	"of":       tOF,
	"else":     tELSE,
	"if":       tIF,
	"new":      tNEW,
	"try":      tTRY,
	"catch":    tCATCH,
	"while":    tWHILE,
	"do":       tDO,
}

var operators = map[string]tokenType{
	"~":   tBITWISE_NOT,
	"&":   tBITWISE_AND,
	"|":   tBITWISE_OR,
	"*":   tMULT,
	"**":  tEXP,
	"-":   tMINUS,
	"+":   tPLUS,
	"/":   tDIV,
	"%":   tMOD,
	"+=":  tPLUS_ASSIGN,
	"*=":  tMULT_ASSIGN,
	"-=":  tMINUS_ASSIGN,
	"/=":  tDIV_ASSIGN,
	"++":  tINC,
	"--":  tDEC,
	"!":   tNOT,
	"&&":  tAND,
	"||":  tOR,
	">":   tGREATER,
	"<":   tLESS,
	">=":  tGREATER_EQUALS,
	"<=":  tLESS_EQUALS,
	"==":  tEQUALS,
	"===": tEQUALS_STRICT,
	"!=":  tNOT_EQUALS,
	"!==": tNOT_EQUALS_STRICT,
	"(":   tPAREN_LEFT,
	")":   tPAREN_RIGHT,
	"[":   tBRACKET_LEFT,
	"]":   tBRACKET_RIGHT,
	"{":   tCURLY_LEFT,
	"}":   tCURLY_RIGHT,
	".":   tDOT,
	",":   tCOMMA,
	"?":   tQUESTION,
	";":   tSEMI,
	":":   tCOLON,
	"=":   tASSIGN,
	"=>":  tLAMBDA,
}

const (
	tUNDEFINED tokenType = iota

	tTYPEOF
	tFOR
	tAS
	tRETURN
	tFUNCTION
	tFROM
	tDEFAULT
	tEXPORT
	tIMPORT
	tVAR
	tCONST
	tLET
	tOF
	tELSE
	tIF
	tNEW
	tTRY
	tCATCH
	tWHILE
	tDO

	tBITWISE_NOT
	tBITWISE_AND
	tBITWISE_OR
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

	tNUMBER
	tSTRING
	tNAME
)

var tokenTypeNames = []string{
	"tUNDEFINED",
	"tTYPEOF",
	"tFOR",
	"tAS",
	"tRETURN",
	"tFUNCTION",
	"tFROM",
	"tDEFAULT",
	"tEXPORT",
	"tIMPORT",
	"tVAR",
	"tCONST",
	"tLET",
	"tOF",
	"tELSE",
	"tIF",
	"tNEW",
	"tTRY",
	"tCATCH",
	"tWHILE",
	"tDO",
	"tBITWISE_NOT",
	"tBITWISE_AND",
	"tBITWISE_OR",
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
	"tNUMBER",
	"tSTRING",
	"tNAME",
}

func (t tokenType) String() string {
	return tokenTypeNames[int32(t)]
}
