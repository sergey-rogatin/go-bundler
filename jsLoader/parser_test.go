package jsLoader

import (
	"testing"
)

func setParser(text string) {
	toks := lex([]byte(text))
	sourceTokens = toks
	index = 0
	tok = sourceTokens[index]
}

func TestLambda(t *testing.T) {
	cases := []struct {
		src       string
		argsCount int
	}{
		{
			"() => {}",
			0,
		},
		{
			"foo => bar",
			1,
		},
		{
			"(foo = 3, bar,) => { foo; bar; }",
			2,
		},
	}

	for _, c := range cases {
		setParser(c.src)
		le, ok := getLambda()
		if !ok {
			t.Errorf("Lambda not parsed")
		}

		if len(le.args) != c.argsCount {
			t.Errorf("Wrong arguments")
		}
	}
}

func TestLambdaFalse(t *testing.T) {
	cases := []struct {
		src string
	}{
		{
			"(,) => {}",
		},
		{
			"(foo, bar) + baz",
		},
		{
			"foo = 3 => bar",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		_, ok := getLambda()
		if ok {
			t.Errorf("Lambda parsed incorrectly")
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
		ol := getStatement()
		ok := true
		if !ok {
			t.Errorf("Member expression not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
		}
	}
}

func TestObjectLiteral(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"{}",
			"{}",
		},
		{
			"{foo, bar}",
			"{foo,bar}",
		},
		{
			"{foo: 1+23, bar, 32: ttu}",
			"{foo:1+23,bar,[32]:ttu}",
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
	}

	for _, c := range cases {
		setParser(c.src)
		ol, ok := getObjectLiteral()
		if !ok {
			t.Errorf("Object literal not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
	}

	for _, c := range cases {
		setParser(c.src)
		ol, ok := getFunctionExpression(false)
		if !ok {
			t.Errorf("Function expression not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
		ol, ok := getArrayLiteral()
		if !ok {
			t.Errorf("Array literal not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
		ol, ok := getBlockStatement()
		if !ok {
			t.Errorf("Block statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
	}

	for _, c := range cases {
		setParser(c.src)
		ol, ok := getForStatement()
		if !ok {
			t.Errorf("For statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
		ol, ok := getWhileStatement()
		if !ok {
			t.Errorf("While statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
		ol, ok := getDoWhileStatement()
		if !ok {
			t.Errorf("Do-while statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
		ol, ok := getIfStatement()
		if !ok {
			t.Errorf("If statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
		ol, ok := getFunctionStatement()
		if !ok {
			t.Errorf("Function statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
			"import 'foo';",
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
			"import bar,{foo as bar}from'foo';",
		},
		{
			"import foo, {default as foo, bar, baz} from 'foo';",
			"import foo,{default as foo,bar,baz}from'foo';",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		ol, ok := getImportStatement()
		if !ok {
			t.Errorf("Import statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
		ol := getExpressionStatement()

		if c.res != ol.String() {
			t.Errorf("Expected %s, got %s", c.res, ol)
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
			"export default function foo() {};",
			"export default function foo(){};",
		},
		{
			"export var foo = 4, bar;",
			"export var foo=4,bar;",
		},
		{
			"export {};",
			"export {};",
		},
		{
			"export {foo as fee, bar as default, wee, };",
			"export {foo as fee,bar as default,wee};",
		},
		{
			"export * from 'foo';",
			"export * from'foo';",
		},
		{
			"export {} from 'foo';",
			"export {} from'foo';",
		},
		{
			"export function foo() {};",
			"export function foo(){};",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		ol, ok := getExportStatement()
		if !ok {
			t.Errorf("Export statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
			"{}=foo;",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		ol := getStatement()
		ok := true
		if !ok {
			t.Errorf("Destucturing statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
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
		ol, ok := getSwitchStatement()
		if !ok {
			t.Errorf("Switch statement not parsed")
		} else {
			if c.res != ol.String() {
				t.Errorf("Expected %s, got %s", c.res, ol)
			}
		}
	}
}
