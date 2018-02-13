package jsLoader

// import (
// 	"testing"
// )

// var ps parserState

// func setParser(text string) {
// 	toks := lex([]byte(text))
// 	ps = parserState{
// 		sourceTokens: toks,
// 		index:        0,
// 		tok:          toks[0],
// 	}
// }

// func TestMarkerStatement(t *testing.T) {
// 	cases := []struct {
// 		src string
// 		res string
// 	}{
// 		{
// 			"{foo:let i=0;}",
// 			"{foo:let i=0;}",
// 		},
// 	}

// 	for _, c := range cases {
// 		setParser(c.src)
// 		le := parseStatement(&ps)

// 		act := generateJsCode(le)
// 		if act != c.res {
// 			t.Errorf("%v", le)
// 			t.Errorf("Expected %s, got %s", c.res, generateJsCode(le))
// 		}
// 	}
// }

// func TestSwitchStatement(t *testing.T) {
// 	cases := []struct {
// 		src string
// 		res string
// 	}{
// 		{
// 			"switch(foo){}",
// 			"switch(foo){}",
// 		},
// 		{
// 			"switch(foo+23){case a: b;c;d; default: e;f;g;}",
// 			"switch(foo+23){case a:b;c;d;default:e;f;g;}",
// 		},
// 	}

// 	for _, c := range cases {
// 		setParser(c.src)
// 		le := parseStatement(&ps)

// 		act := generateJsCode(le)
// 		if act != c.res {
// 			t.Errorf("%v", le)
// 			t.Errorf("Expected %s, got %s", c.res, generateJsCode(le))
// 		}
// 	}
// }

// func TestNewlineAndSemi(t *testing.T) {
// 	cases := []struct {
// 		src string
// 		res string
// 	}{
// 		{
// 			"var\n foo\n",
// 			"var foo;",
// 		},
// 		{
// 			"{foo}",
// 			"{foo;}",
// 		},
// 		{
// 			"var a = {foo}",
// 			"var a={foo};",
// 		},
// 		{
// 			`const {
// 				addEntityType,

// 				addEntity,
// 			} = engine`,
// 			"const {addEntityType,addEntity}=engine;",
// 		},
// 		{
// 			`for(
// 				i
// 				;
// 				i<321;
// 				i++
// 				);`,
// 			"for(i;i<321;i++);",
// 		},
// 	}

// 	for _, c := range cases {
// 		setParser(c.src)
// 		le := parseStatement(&ps)

// 		act := generateJsCode(le)
// 		if act != c.res {
// 			t.Errorf("%v", le)
// 			t.Errorf("Expected %s, got %s", c.res, generateJsCode(le))
// 		}
// 	}
// }

// func TestImportTransform(t *testing.T) {
// 	cases := []struct {
// 		src string
// 		res string
// 	}{
// 	// {
// 	// 	"import foo from './bar';",
// 	// 	"var foo=bar_js.default;",
// 	// },
// 	// {
// 	// 	"import foo, {bar as baz, default as fooz} from './bar'",
// 	// 	"var foo=bar_js.default,baz=bar_js.bar,fooz=bar_js.default;",
// 	// },
// 	// {
// 	// 	"import './foo';",
// 	// 	"",
// 	// },
// 	// {
// 	// 	"import a, * as b from './foo';",
// 	// 	"var b=foo_js,a=foo_js.default;",
// 	// },
// 	// {
// 	// 	"import * as b from './foo';",
// 	// 	"var b=foo_js;",
// 	// },
// 	// {
// 	// 	"import {} from './foo';",
// 	// 	"",
// 	// },
// 	}

// 	for _, c := range cases {
// 		setParser(c.src)
// 		ast, _ := parseTokens(ps.sourceTokens)
// 		transAst, _ := transformIntoModule(ast, "a.js")

// 		str := generateJsCode(transAst)
// 		cutStr := str[20 : len(str)-10]

// 		if cutStr != c.res {
// 			t.Errorf("%v", ast)
// 			t.Errorf("Expected %s, got %s", c.res, cutStr)
// 		}
// 	}
// }

// func TestExportTransform(t *testing.T) {
// 	cases := []struct {
// 		src string
// 		res string
// 	}{
// 		{
// 			"export default foo;",
// 			"exports.default=foo;",
// 		},
// 		{
// 			"export default function foo(){};",
// 			"function foo(){}exports.default=foo;",
// 		},
// 		{
// 			"export default function foo(){};",
// 			"function foo(){}exports.default=foo;",
// 		},
// 		{
// 			"export default function(){};",
// 			"exports.default=function(){};",
// 		},
// 		{
// 			"export var foo=3, bar;",
// 			"var foo=3,bar;exports.foo=foo,exports.bar=bar;",
// 		},
// 		{
// 			"export {a,b as default};",
// 			"exports.a=a,exports.default=b;",
// 		},
// 		{
// 			"export function foo() {}",
// 			"function foo(){}exports.foo=foo;",
// 		},
// 		{
// 			"export {foo as bar};",
// 			"exports.bar=foo;",
// 		},
// 		{
// 			"export {foo as bar,a} from './bar';",
// 			"exports.bar=modules.bar_js.foo,exports.a=modules.bar_js.a;",
// 		},
// 		{
// 			"export * from './bar';",
// 			"Object.assign(exports,modules.bar_js);",
// 		},
// 	}

// 	for _, c := range cases {
// 		setParser(c.src)
// 		ast, _ := parseTokens(ps.sourceTokens)
// 		transAst, _ := transformIntoModule(ast, "a.js")

// 		str := generateJsCode(transAst)
// 		cutStr := str[41 : len(str)-17]

// 		if cutStr != c.res {
// 			t.Errorf("%v", ast)
// 			t.Errorf("Expected %s, got %s", c.res, cutStr)
// 		}
// 	}
// }

// func TestRequireTransform(t *testing.T) {
// 	cases := []struct {
// 		src string
// 		res string
// 	}{
// 		{
// 			"var a = require('./foo');",
// 			"var a=modules.foo_js.default;",
// 		},
// 	}

// 	for _, c := range cases {
// 		setParser(c.src)
// 		ast, _ := parseTokens(ps.sourceTokens)
// 		transAst, _ := transformIntoModule(ast, "a.js")

// 		str := generateJsCode(transAst)
// 		cutStr := str[41 : len(str)-17]

// 		if cutStr != c.res {
// 			t.Errorf("%v", ast)
// 			t.Errorf("Expected %s, got %s", c.res, cutStr)
// 		}
// 	}
// }

// 	for _, c := range cases {
// 		setParser(c.src)
// 		le := parseStatement(&ps)

// 		act := generateJsCode(le)
// 		if act != c.res {
// 			t.Errorf("%v", le)
// 			t.Errorf("Expected %s, got %s", c.res, generateJsCode(le))
// 		}
// 	}
// }

// func TestClass(t *testing.T) {
// 	cases := []struct {
// 		src string
// 		res string
// 	}{
// 		{
// 			"class foo extends bar {}",
// 			"class foo extends bar{}",
// 		},
// 		{
// 			"class foo{bar(){}}",
// 			"class foo{bar(){}}",
// 		},
// 		{
// 			"class foo{a = fsa; b;static c;}",
// 			"class foo{a=fsa;b;static c;}",
// 		},
// 		{
// 			"class foo{static foo = () => {}}",
// 			"class foo{static foo=()=>{};}",
// 		},
// 		{
// 			"var foo = class foo{};",
// 			"var foo=class foo{};",
// 		},
// 		{
// 			"var foo = class{};",
// 			"var foo=class{};",
// 		},
// 	}

// 	for _, c := range cases {
// 		setParser(c.src)
// 		le := parseStatement(&ps)

// 		act := generateJsCode(le)
// 		if act != c.res {
// 			t.Errorf("%v", le)
// 			t.Errorf("Expected %s, got %s", c.res, generateJsCode(le))
// 		}
// 	}
// }
