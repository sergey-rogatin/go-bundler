package jsLoader

import (
	"testing"
)

func TestImportTransform(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"import foo from './bar';",
			"var foo=requireES6(modules.bar_js,'default');",
		},
		{
			"import foo, {bar as baz, default as fooz} from './bar'",
			"var foo=requireES6(modules.bar_js,'default'),baz=requireES6(modules.bar_js,'bar'),fooz=requireES6(modules.bar_js,'default');",
		},
		{
			"import './foo';",
			"",
		},
		{
			"import a, * as b from './foo';",
			"var b=requireES6(modules.foo_js,'*'),a=requireES6(modules.foo_js,'default');",
		},
		{
			"import * as b from './foo';",
			"var b=requireES6(modules.foo_js,'*');",
		},
		{
			"import {} from './foo';",
			"",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		a, _ := parseTokens(ps.tokens)

		ctx := context{}
		ctx.declaredVars = map[string]ast{}
		e := environment{"a.js", []string{}, false, nil}
		transAst := modifyImport(a.children[0], &e, &ctx)

		str := printAst(transAst)
		cutStr := str

		if cutStr != c.res {
			t.Errorf("%v", a)
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
			"var a=require(modules.foo_js);",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		a, _ := parseTokens(ps.tokens)
		ctx := context{}
		ctx.declaredVars = map[string]ast{}
		e := environment{"a.js", []string{}, false, nil}
		transAst := modifyAst(a.children[0], &e, &ctx)

		str := printAst(transAst)
		cutStr := str

		if cutStr != c.res {
			t.Errorf("%v", a)
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
			"module.es6.default=foo;",
		},
		{
			"export default function foo(){};",
			"function foo(){}module.es6.default=foo;",
		},
		{
			"export default function foo(){};",
			"function foo(){}module.es6.default=foo;",
		},
		{
			"export default function(){};",
			"module.es6.default=function(){};",
		},
		{
			"export var foo=3, bar;",
			"var foo=3,bar;module.es6.foo=foo,module.es6.bar=bar;",
		},
		{
			"export {a,b as default};",
			"module.es6.a=a,module.es6.default=b;",
		},
		{
			"export function foo() {}",
			"function foo(){}module.es6.foo=foo;",
		},
		{
			"export {foo as bar};",
			"module.es6.bar=foo;",
		},
		{
			"export {foo as bar,a} from './bar';",
			"module.es6.bar=requireES6(modules.bar_js,'foo'),module.es6.a=requireES6(modules.bar_js,'a');",
		},
		{
			"export * from './bar';",
			"Object.assign(module.es6,requireES6(modules.bar_js,'*'));",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		a, _ := parseTokens(ps.tokens)
		ctx := context{}
		ctx.declaredVars = map[string]ast{}
		e := environment{"a.js", []string{}, false, nil}
		transAst := modifyExport(a.children[0], &e, &ctx)

		str := printAst(transAst)
		cutStr := str

		if cutStr != c.res {
			t.Errorf("%v", a)
			t.Errorf("Expected %s, got %s", c.res, cutStr)
		}
	}
}

func TestProgramTransform(t *testing.T) {
	cases := []struct {
		src string
		res string
	}{
		{
			"foo;",
			"moduleFns.a_js=function(){var module={exports:{},es6:{},hasES6Exports:false},exports=module.exports;foo;return module;};",
		},
	}

	for _, c := range cases {
		setParser(c.src)
		a, _ := parseTokens(ps.tokens)
		ctx := context{}
		ctx.declaredVars = map[string]ast{}
		e := environment{"a.js", []string{}, false, nil}
		transAst := modifyAst(a, &e, &ctx)

		str := printAst(transAst)
		cutStr := str

		if cutStr != c.res {
			t.Errorf("%v", a)
			t.Errorf("%v", transAst)
			t.Errorf("Expected %s, got %s", c.res, cutStr)
		}
	}
}
