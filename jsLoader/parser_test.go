package jsLoader

import (
	"testing"
)

var ps parserState

func setParser(text string) {
	toks := lex([]byte(text))
	ps = parserState{
		sourceTokens: toks,
		index:        0,
		tok:          toks[0],
	}
}

func TestLambda(t *testing.T) {
	cases := []struct {
		src string
		exp string
	}{
		{
			"() => {}",
			"()=>{}",
		},
		{
			"foo => bar",
			"(foo)=>bar",
		},
		{
			"(foo = 3, bar,) => { foo; bar; }",
			"(foo=3,bar)=>{foo;bar;}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseExpression(&ps)

		res := printAst(le)
		if res != c.exp {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.exp, printAst(le))
		}
	}
}

func TestMemberExpression(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"foo.a[b].c().d;",
			"foo.a[b].c().d;",
		},
		{
			"new a.v().sd;",
			"new a.v().sd;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestObjectLiteral(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"{32: foo}",
			"{32:foo}",
		},
		{
			"{foo, bar}",
			"{foo,bar}",
		},
		{
			"{foo: 1+23, bar, 32: ttu}",
			"{foo:1+23,bar,32:ttu}",
		},
		{
			"{[312 + foo]: bar}",
			"{[312+foo]:bar}",
		},
		{
			"{foo() {}}",
			"{foo(){}}",
		},
		{
			"{['foo' + 32]() {}}",
			"{['foo'+32](){}}",
		},
		{
			"{set foo() {}, get foo() {}}",
			"{set foo(){},get foo(){}}",
		},
		{
			`{
					get _enabled () { return _enabled;},
				};`,
			"{get _enabled(){return _enabled;;}}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseExpression(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestFunctionExpression(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"function foo() {}",
			"function foo(){}",
		},
		{
			"function(foo = ee = 321, bar) {}",
			"function(foo=ee=321,bar){}",
		},
		{
			"function(a,b,...c) {}",
			"function(a,b,...c){}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseExpression(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"[foo, bar, 213*(21+3), () => foo,]",
			"[foo,bar,213*(21+3),()=>foo]",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseExpression(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestBlockStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"{foo; bar = 321;;}",
			"{foo;bar=321;;}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestMarkerStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"{foo:let i=0;}",
			"{foo:let i=0;}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestForStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"for (;;);",
			"for(;;);",
		},
		{
			"for(var i=0;i<10;i++);",
			"for(var i=0;i<10;i++);",
		},
		{
			"for(;i<10;i++);",
			"for(;i<10;i++);",
		},
		{
			"for(i;;i++);",
			"for(i;;i++);",
		},
		{
			"for(i;i<23;);",
			"for(i;i<23;);",
		},
		{
			"for(i;i<23;) {}",
			"for(i;i<23;){}",
		},
		{
			"for(i;i<23;) foo = 3;",
			"for(i;i<23;)foo=3;",
		},
		{
			"for(foo of bar()) foo = 3;",
			"for(foo of bar())foo=3;",
		},
		{
			"for(const foo in bar) {foo = 3;}",
			"for(const foo in bar){foo=3;}",
		},
		{
			"for(a in b; i < 21; i++) {foo = 3;}",
			"for(a in b;i<21;i++){foo=3;}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestInAndInstanceof(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"for(a in b; i < 21; i++) {foo = 3;}",
			"for(a in b;i<21;i++){foo=3;}",
		},
		{
			"a instanceof b;",
			"a instanceof b;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestExpressions(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"!!(a+b);",
			"!!(a+b);",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestStringLiterals(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"'foo \\' + fsbds'",
			"'foo \\' + fsbds'",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestWhileStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"while (foo);",
			"while(foo);",
		},
		{
			"while (foo*bar < 3) {}",
			"while(foo*bar<3){}",
		},
		{
			"while (foo, bar += 3) bar();",
			"while(foo,bar+=3)bar();",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestDoWhileStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"do {} while(foo);",
			"do {}while(foo);",
		},
		{
			"do ; while (foo*bar < 3);",
			"do ;while(foo*bar<3);",
		},
		{
			"do bar();while(foo,bar+=3);",
			"do bar();while(foo,bar+=3);",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestIfStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"if (foo) bar;",
			"if(foo)bar;",
		},
		{
			"if(foo){}",
			"if(foo){}",
		},
		{
			"if(foo, bar = 3){foo();}",
			"if(foo,bar=3){foo();}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestFunctionStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"function foo() {}",
			"function foo(){}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestImportStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"import 'foo';",
			"import'foo';",
		},
		{
			"import * as foo from 'foo';",
			"import*as foo from'foo';",
		},
		{
			"import bar, * as foo from 'foo';",
			"import bar,*as foo from'foo';",
		},
		{
			"import bar, {foo as bar} from 'foo';",
			"import bar,{foo as bar} from'foo';",
		},
		{
			"import foo, {default as foo, bar, baz} from 'foo';",
			"import foo,{default as foo,bar as bar,baz as baz} from'foo';",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestExpressionStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"var foo = 3, bar;",
			"var foo=3,bar;",
		},
		{
			"break;",
			"break;",
		},
		{
			"continue;",
			"continue;",
		},
		{
			"debugger;",
			"debugger;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestExportStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"export default foo;",
			"export default foo;",
		},
		{
			"export default foo + 231;",
			"export default foo+231;",
		},
		{
			"export default function() {};",
			"export default function(){};",
		},
		{
			"export default function foo() {}",
			"function foo(){}export default foo;",
		},
		{
			"export var foo = 4, bar;",
			"var foo=4,bar;export{foo as foo,bar as bar};",
		},
		{
			"export {};",
			"export{};",
		},
		{
			"export {foo as fee, bar as default, wee, };",
			"export{foo as fee,bar as default,wee as wee};",
		},
		// {
		// 	"export * from 'foo';",
		// 	"export * from'foo';",
		// },
		{
			"export {} from 'foo';",
			"export{} from'foo';",
		},
		{
			"export function foo() {};",
			"function foo(){}export{foo as foo};",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestObjectDestructuring(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"var {} = foo;",
			"var {}=foo;",
		},
		{
			"var {a:b} = foo;",
			"var {a:b}=foo;",
		},
		{
			"var {a:b=32} = foo;",
			"var {a:b=32}=foo;",
		},
		{
			"var {a:b={c:f=5,d:e}=3} = foo;",
			"var {a:b={c:f=5,d:e}=3}=foo;",
		},
		{
			"({}) = foo;",
			"({})=foo;",
		},
		// {
		// 	"var {foo,...bar=2} = doo;",
		// 	"var {foo,...bar=2}=doo;",
		// },
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestSwitchStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"switch(foo){}",
			"switch(foo){}",
		},
		{
			"switch(foo+23){case a: b;c;d; default: e;f;g;}",
			"switch(foo+23){case a:b;c;d;default:e;f;g;}",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestNewlineAndSemi(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"var\n foo\n",
			"var foo;",
		},
		{
			"{foo}",
			"{foo;}",
		},
		{
			"var a = {foo}",
			"var a={foo};",
		},
		{
			`const {
				addEntityType,

				addEntity,
			} = engine`,
			"const {addEntityType,addEntity}=engine;",
		},
		{
			`for(
				i
				;
				i<321;
				i++
				);`,
			"for(i;i<321;i++);",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestImportTransform(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"import foo from './bar';",
			"var foo=bar_js.default;",
		},
		{
			"import foo, {bar as baz, default as fooz} from './bar'",
			"var foo=bar_js.default,baz=bar_js.bar,fooz=bar_js.default;",
		},
		{
			"import './foo';",
			"",
		},
		{
			"import a, * as b from './foo';",
			"var b=foo_js,a=foo_js.default;",
		},
		{
			"import * as b from './foo';",
			"var b=foo_js;",
		},
		{
			"import {} from './foo';",
			"",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		ast, _ := parseTokens(ps.sourceTokens)
		transAst, _ := transformIntoModule(ast, "a.js")

		str := printAst(transAst)
		cutStr := str[35 : len(str)-19]

		if cutStr != c.res {
			t.Errorf("%v", ast)
			t.Errorf("Expected %s, got %s", c.res, cutStr)
		}
	}
}

func TestExportTransform(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"export default foo;",
			"exports.default=foo;",
		},
		{
			"export default function foo(){};",
			"function foo(){}exports.default=foo;",
		},
		{
			"export default function foo(){};",
			"function foo(){}exports.default=foo;",
		},
		{
			"export default function(){};",
			"exports.default=function(){};",
		},
		{
			"export var foo=3, bar;",
			"var foo=3,bar;exports.foo=foo,exports.bar=bar;",
		},
		{
			"export {a,b as default};",
			"exports.a=a,exports.default=b;",
		},
		{
			"export function foo() {}",
			"function foo(){}exports.foo=foo;",
		},
		{
			"export {foo as bar};",
			"exports.bar=foo;",
		},
		{
			"export {foo as bar,a} from './bar';",
			"exports.bar=bar_js.foo,exports.a=bar_js.a;",
		},
		{
			"export * from './bar';",
			"Object.assign(exports,bar_js);",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		ast, _ := parseTokens(ps.sourceTokens)
		transAst, _ := transformIntoModule(ast, "a.js")

		str := printAst(transAst)
		cutStr := str[35 : len(str)-19]

		if cutStr != c.res {
			t.Errorf("%v", ast)
			t.Errorf("Expected %s, got %s", c.res, cutStr)
		}
	}
}

func TestRequireTransform(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"var a = require('./foo');",
			"var a=foo_js;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		ast, _ := parseTokens(ps.sourceTokens)
		transAst, _ := transformIntoModule(ast, "a.js")

		str := printAst(transAst)
		cutStr := str[35 : len(str)-19]

		if cutStr != c.res {
			t.Errorf("%v", ast)
			t.Errorf("Expected %s, got %s", c.res, cutStr)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			`return {
				result: mapResult,
				keyPrefix: keyPrefix,
				func: mapFunction,
				context: mapContext,
				count: 0
			};`,
			"return {result:mapResult,keyPrefix:keyPrefix," +
				"func:mapFunction,context:mapContext,count:0};",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestConditional(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"foo?bar:baz;",
			"foo?bar:baz;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

func TestRegexp(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"var foo = a * /[a-zA-Z]/gi;",
			"var foo=a*/[a-zA-Z]/gi;",
		},
		{
			"var foo = a * /[a-zA-Z]/ + 3;",
			"var foo=a*/[a-zA-Z]/+3;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le := parseStatement(&ps)

		act := printAst(le)
		if act != c.res {
			t.Errorf("%v", le)
			t.Errorf("Expected %s, got %s", c.res, printAst(le))
		}
	}
}

// func TestArrayPattern(t *testing.T) {
// 	cases := []struct {
// 		src string
// 		res string
// 	}{
// 		{
// 			"[foo,bar,]",
// 			"[foo,bar,]",
// 		},
// 		{
// 			"[foo,,]",
// 			"[foo,,]",
// 		},
// 		{
// 			"[a = 23, foo]",
// 			"[a=23,foo]",
// 		},
// 		{
// 			"[{foo:bar = 23} = 23, foo]",
// 			"[{foo:bar=23}=23,foo]",
// 		},
// 		// {
// 		// 	"[a, ...b]",
// 		// 	"[a,...b]",
// 		// },
// 	}

// 	for _, c := range cases {
// 		setParser(c.src)
// 		le := parseStatement()

// 		//t.Errorf("%v", le)
// 		act := printAst(le)
// 		if act != c.res {
// 			t.Errorf("Expected %s, got %s", c.res, printAst(le))
// 		}
// 	}
// }
